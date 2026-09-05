"""对人声 stem 抽 F0，写出 pitch.json 供 KTV 画横条。"""
from __future__ import annotations

import json
import logging
import math
import os
import threading
import time
import urllib.request
from pathlib import Path
from typing import Any

import numpy as np

from app.config import settings

logger = logging.getLogger("uvr-api.pitch")

HOP_MS = 10.0  # hop 160 @ 16kHz
MIN_NOTE_MS = 80.0
GAP_FILL_MS = 60.0
MERGE_GAP_MS = 100.0
STICKY_HOLD_MS = 100.0
PITCH_VERSION = 5
PITCH_FILENAME = "pitch.json"
# RMVPE 在 Live 人声上常锁到低五度（E3→A2）；pyin 更稳，RMVPE 仅作回退
PYIN_FMIN_NOTE = "C2"
PYIN_FMAX_NOTE = "G5"
_PERIOD_FACTORS = (1.0, 0.5, 2.0, 2.0 / 3.0, 1.5)

_lock = threading.RLock()
_rmvpe = None
_last_error: str | None = None


def pitch_json_path(job_dir: Path) -> Path:
    return job_dir / "out" / PITCH_FILENAME


def pitch_json_stale(dest: Path) -> bool:
    """version < 5 的 json 用 RMVPE，Live 开口常低五度且缺横条，需重抽。"""
    if not dest.is_file() or dest.stat().st_size <= 20:
        return True
    try:
        data = json.loads(dest.read_text(encoding="utf-8"))
        return int(data.get("version") or 0) < PITCH_VERSION
    except Exception:
        return True


def _model_path() -> Path:
    return settings.work_dir / "rmvpe" / "rmvpe.pt"


def _download_rmvpe(dest: Path) -> None:
    dest.parent.mkdir(parents=True, exist_ok=True)
    endpoint = (settings.hf_endpoint or "https://huggingface.co").rstrip("/")
    urls = [
        f"{endpoint}/lj1995/VoiceConversionWebUI/resolve/main/rmvpe.pt",
        "https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/rmvpe.pt",
    ]
    last_err: Exception | None = None
    for url in urls:
        try:
            logger.info("下载 RMVPE 权重: %s", url)
            req = urllib.request.Request(url, headers={"User-Agent": "ONIJ-UVR-API/0.3"})
            with urllib.request.urlopen(req, timeout=180) as resp:
                data = resp.read()
            if len(data) < 1_000_000:
                raise RuntimeError(f"rmvpe.pt 过小 ({len(data)} bytes)")
            tmp = dest.with_suffix(".pt.part")
            tmp.write_bytes(data)
            tmp.replace(dest)
            logger.info("RMVPE 权重已保存 %s (%d bytes)", dest, len(data))
            return
        except Exception as exc:
            last_err = exc
            logger.warning("下载失败 %s: %s", url, exc)
    raise RuntimeError(f"无法下载 rmvpe.pt: {last_err}")


def _ensure_hf_endpoint() -> None:
    endpoint = (settings.hf_endpoint or "").strip()
    if endpoint and not os.environ.get("HF_ENDPOINT"):
        os.environ["HF_ENDPOINT"] = endpoint.rstrip("/")


def _get_rmvpe():
    global _rmvpe, _last_error
    if _rmvpe is not None:
        return _rmvpe
    with _lock:
        if _rmvpe is not None:
            return _rmvpe
        _ensure_hf_endpoint()
        path = _model_path()
        if not path.is_file():
            _download_rmvpe(path)
        from app.pyinstaller_bootstrap import bootstrap_pyinstaller

        bootstrap_pyinstaller()
        import torch
        from app.rmvpe import RMVPE

        device = "cuda" if torch.cuda.is_available() else "cpu"
        is_half = device.startswith("cuda")
        logger.info("加载 RMVPE device=%s half=%s path=%s", device, is_half, path)
        _rmvpe = RMVPE(str(path), is_half=is_half, device=device)
        _last_error = None
        return _rmvpe


def _hz_to_midi(hz: float) -> float:
    return 69.0 + 12.0 * math.log2(hz / 440.0)


def _fill_short_gaps(midi: np.ndarray, hop_ms: float, max_gap_ms: float) -> np.ndarray:
    out = midi.copy()
    max_frames = max(1, int(round(max_gap_ms / hop_ms)))
    n = len(out)
    i = 0
    while i < n:
        if out[i] > 0:
            i += 1
            continue
        j = i
        while j < n and out[j] <= 0:
            j += 1
        gap = j - i
        left = out[i - 1] if i > 0 else 0.0
        right = out[j] if j < n else 0.0
        if 0 < gap <= max_frames and left > 0 and (right <= 0 or abs(right - left) <= 1):
            out[i:j] = left
        i = j
    return out


def _merge_close_notes(notes: list[dict[str, float]], max_gap_s: float) -> list[dict[str, float]]:
    if not notes:
        return notes
    out: list[dict[str, float]] = [dict(notes[0])]
    for n in notes[1:]:
        last = out[-1]
        if n["midi"] == last["midi"] and n["t0"] - last["t1"] <= max_gap_s:
            last["t1"] = max(last["t1"], n["t1"])
        else:
            out.append(dict(n))
    min_s = MIN_NOTE_MS / 1000.0
    return [n for n in out if n["t1"] - n["t0"] >= min_s]


def _sticky_merge(notes: list[dict[str, float]]) -> list[dict[str, float]]:
    """只抹 1–2 半音的短抖动；跨三度以上的短音保留，避免 E3 被前一个 A2 吞掉。"""
    if not notes:
        return notes
    hold_s = STICKY_HOLD_MS / 1000.0
    cur = notes[0]["midi"]
    for n in notes:
        if n["midi"] == cur:
            continue
        if abs(n["midi"] - cur) >= 3 or n["t1"] - n["t0"] >= hold_s:
            cur = n["midi"]
        else:
            n["midi"] = cur
    return _merge_close_notes(notes, MERGE_GAP_MS / 1000.0)


def _frame_rms(audio: np.ndarray, hop: int, n_frames: int, win: int = 320) -> np.ndarray:
    rms = np.zeros(n_frames, dtype=np.float64)
    n = len(audio)
    for i in range(n_frames):
        a = i * hop
        b = min(n, a + win)
        if b > a:
            sl = audio[a:b]
            rms[i] = float(np.sqrt(np.mean(sl * sl)))
    return rms


def _gate_f0_by_energy(f0: np.ndarray, audio: np.ndarray, hop: int = 160) -> np.ndarray:
    """按约 4s 窗口的局部能量门限，避免整曲后段大声把 Live 轻开口抹掉。"""
    n = len(f0)
    rms = _frame_rms(audio, hop, n)
    voiced = f0 > 0
    if not np.any(voiced):
        return f0
    out = f0.copy()
    win = max(50, int(round(4000.0 / HOP_MS)))
    step = max(1, win // 2)
    for i0 in range(0, n, step):
        i1 = min(n, i0 + win)
        loc = voiced[i0:i1]
        if not np.any(loc):
            continue
        vrms = rms[i0:i1][loc]
        floor = float(np.percentile(vrms, 30)) * 0.28
        weak = rms[i0:i1] < max(floor, 1e-5)
        sl = out[i0:i1]
        sl[weak] = 0.0
        out[i0:i1] = sl
    return out


def _correct_period_errors(f0: np.ndarray) -> np.ndarray:
    """把明显的低八度 / 低五度帧扳回主体音域（RMVPE 回退用）。"""
    voiced = f0 > 0
    if int(np.count_nonzero(voiced)) < 40:
        return f0
    midi = np.zeros(len(f0), dtype=np.float64)
    midi[voiced] = 69.0 + 12.0 * np.log2(np.maximum(f0[voiced], 1e-6) / 440.0)
    core = midi[(midi >= 50.0) & (midi <= 72.0)]
    if core.size < 40:
        return f0
    ref = float(np.median(core))
    out = f0.copy()
    for i in np.flatnonzero(voiced):
        hz = float(f0[i])
        best_hz = hz
        best_d = abs(_hz_to_midi(hz) - ref)
        for fac in _PERIOD_FACTORS:
            cand = hz * fac
            if cand < 65.0 or cand > 900.0:
                continue
            d = abs(_hz_to_midi(cand) - ref)
            if d + 0.45 < best_d:
                best_d = d
                best_hz = cand
        out[i] = best_hz
    return out


def _infer_f0(audio: np.ndarray) -> tuple[np.ndarray, str]:
    try:
        import librosa

        f0, _voiced, _prob = librosa.pyin(
            audio,
            fmin=float(librosa.note_to_hz(PYIN_FMIN_NOTE)),
            fmax=float(librosa.note_to_hz(PYIN_FMAX_NOTE)),
            sr=16000,
            hop_length=160,
        )
        return np.where(np.isfinite(f0), f0.astype(np.float64), 0.0), "pyin"
    except Exception:
        logger.exception("pyin 失败，回退 RMVPE")
        raw = np.asarray(_get_rmvpe().infer_from_audio(audio, thred=0.08), dtype=np.float64)
        return _correct_period_errors(raw), "rmvpe"


def _run_lengths(midi: np.ndarray) -> np.ndarray:
    runs = np.zeros(len(midi), dtype=np.int32)
    n = len(midi)
    i = 0
    while i < n:
        j = i + 1
        while j < n and midi[j] == midi[i]:
            j += 1
        runs[i:j] = j - i
        i = j
    return runs


def _suppress_speech(midi: np.ndarray, hop_ms: float) -> np.ndarray:
    """讲话段音高又碎又跳，没有稳定横条，整段抹掉。"""
    out = midi.copy()
    runs = _run_lengths(out)
    win = max(3, int(round(240.0 / hop_ms)))
    min_sing = max(2, int(round(80.0 / hop_ms)))
    n = len(out)
    for i in range(n):
        if out[i] <= 0 or runs[i] >= min_sing:
            continue
        a = max(0, i - win // 2)
        b = min(n, i + win // 2)
        seg = out[a:b]
        voiced = seg[seg > 0]
        if voiced.size < 6 or float(np.std(voiced)) >= 1.4:
            out[i] = 0.0
    return out


def _clip_to_singing_range(midi: np.ndarray) -> np.ndarray:
    """主体音域外的短高音多半是观众尖叫。"""
    voiced = midi[midi > 0]
    if voiced.size < 30:
        return midi
    runs = _run_lengths(midi)
    core = midi[(midi > 0) & (runs >= 8)]
    if core.size < 20:
        core = voiced
    lo, hi = np.percentile(core, [10, 90])
    out = midi.copy()
    out[(out > 0) & ((out < lo - 5) | (out > hi + 5))] = 0.0
    return out


def _notes_from_f0(f0: np.ndarray, hop_ms: float) -> list[dict[str, float]]:
    midi = np.zeros(len(f0), dtype=np.float64)
    voiced = f0 > 0
    if np.any(voiced):
        midi[voiced] = np.round(69.0 + 12.0 * np.log2(np.maximum(f0[voiced], 1e-6) / 440.0))
    midi = _suppress_speech(midi, hop_ms)
    midi = _clip_to_singing_range(midi)
    midi = _fill_short_gaps(midi, hop_ms, GAP_FILL_MS)

    notes: list[dict[str, float]] = []
    i = 0
    n = len(midi)
    min_frames = max(1, int(round(MIN_NOTE_MS / hop_ms)))
    while i < n:
        if midi[i] <= 0:
            i += 1
            continue
        j = i + 1
        while j < n and midi[j] == midi[i]:
            j += 1
        if (j - i) >= min_frames:
            notes.append(
                {
                    "t0": round(i * hop_ms / 1000.0, 3),
                    "t1": round(j * hop_ms / 1000.0, 3),
                    "midi": float(midi[i]),
                }
            )
        i = j
    notes = _merge_close_notes(notes, MERGE_GAP_MS / 1000.0)
    return _sticky_merge(notes)


def _load_mono_16k(path: Path) -> np.ndarray:
    import librosa

    audio, _sr = librosa.load(str(path), sr=16000, mono=True)
    if audio.size == 0:
        raise RuntimeError(f"空音频: {path}")
    peak = float(np.max(np.abs(audio)))
    if peak > 1e-3:
        audio = audio / peak * 0.95
    return audio.astype(np.float32)


def extract_pitch_dict(vocals_path: Path, song_id: str | None = None) -> dict[str, Any]:
    started = time.perf_counter()
    audio = _load_mono_16k(vocals_path)
    f0, algo = _infer_f0(audio)
    f0 = _gate_f0_by_energy(f0, audio, hop=160)
    hop_ms = HOP_MS
    duration_ms = int(round(len(audio) / 16000 * 1000))
    notes = _notes_from_f0(f0, hop_ms)
    voiced_midi = [_hz_to_midi(float(h)) for h in f0 if h > 0]
    ref_midi = float(np.median(voiced_midi)) if voiced_midi else 60.0
    elapsed = time.perf_counter() - started
    logger.info(
        "F0 完成 song=%s algo=%s notes=%d voiced=%d elapsed=%.2fs src=%s",
        song_id,
        algo,
        len(notes),
        len(voiced_midi),
        elapsed,
        vocals_path.name,
    )
    return {
        "song_id": song_id or "",
        "version": PITCH_VERSION,
        "algo": algo,
        "hop_ms": hop_ms,
        "duration_ms": duration_ms,
        "ref_midi": round(ref_midi, 2),
        "notes": notes,
        "elapsed_sec": round(elapsed, 3),
        "vocals_file": vocals_path.name,
    }


def write_pitch_json(job_dir: Path, vocals_path: Path, song_id: str | None = None) -> Path:
    out = pitch_json_path(job_dir)
    out.parent.mkdir(parents=True, exist_ok=True)
    payload = extract_pitch_dict(vocals_path, song_id=song_id)
    out.write_text(json.dumps(payload, ensure_ascii=False), encoding="utf-8")
    return out


def ensure_pitch_json(
    job_dir: Path,
    vocals_path: Path | None,
    song_id: str | None = None,
    overwrite: bool = False,
) -> Path | None:
    dest = pitch_json_path(job_dir)
    if not overwrite and dest.is_file() and dest.stat().st_size > 20:
        return dest
    if vocals_path is None or not vocals_path.is_file():
        return None
    with _lock:
        if not overwrite and dest.is_file() and dest.stat().st_size > 20:
            return dest
        try:
            return write_pitch_json(job_dir, vocals_path, song_id=song_id)
        except Exception:
            logger.exception("抽 F0 失败 job=%s", song_id)
            return None
