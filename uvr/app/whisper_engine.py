"""faster-whisper 推理（默认 small + CUDA，懒加载）。"""
from __future__ import annotations

import logging
import os
import threading
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from app.config import settings

logger = logging.getLogger("uvr-api.whisper")


def _ensure_hf_endpoint() -> None:
    endpoint = (settings.hf_endpoint or "").strip()
    if endpoint and not os.environ.get("HF_ENDPOINT"):
        os.environ["HF_ENDPOINT"] = endpoint.rstrip("/")
        logger.info("HF_ENDPOINT=%s", os.environ["HF_ENDPOINT"])
    # hf-xet / CAS 在国内镜像下常 401，强制走普通下载
    os.environ.setdefault("HF_HUB_DISABLE_XET", "1")
    os.environ.setdefault("HF_HUB_ENABLE_HF_TRANSFER", "0")

_lock = threading.Lock()
_model: Any = None
_resolved_device: str = ""
_resolved_compute: str = ""
_last_error: str | None = None


@dataclass
class WhisperStatus:
    ok: bool
    model: str
    device: str
    compute_type: str
    loaded: bool
    error: str | None = None


@dataclass
class WhisperResult:
    text: str
    language: str | None
    duration: float | None
    elapsed_sec: float
    segments: list[dict[str, Any]]


def _prefer_cuda() -> bool:
    want = (settings.whisper_device or "cuda").strip().lower()
    if want == "cpu":
        return False
    if want in ("cuda", "auto", ""):
        try:
            import ctranslate2

            return ctranslate2.get_cuda_device_count() > 0
        except Exception:
            return False
    return False


def resolve_device_and_compute() -> tuple[str, str]:
    if _prefer_cuda():
        compute = (settings.whisper_compute_type or "float16").strip() or "float16"
        return "cuda", compute
    return "cpu", "int8"


def get_status() -> WhisperStatus:
    global _last_error
    model_name = settings.whisper_model or "small"
    try:
        import faster_whisper  # noqa: F401
        import ctranslate2  # noqa: F401

        device, compute = resolve_device_and_compute()
        return WhisperStatus(
            ok=True,
            model=model_name,
            device=device,
            compute_type=compute,
            loaded=_model is not None,
            error=_last_error,
        )
    except Exception as exc:
        _last_error = str(exc)
        device, compute = "cpu", "int8"
        return WhisperStatus(
            ok=False,
            model=model_name,
            device=device,
            compute_type=compute,
            loaded=False,
            error=str(exc),
        )


def get_model() -> Any:
    global _model, _resolved_device, _resolved_compute, _last_error
    if _model is not None:
        return _model

    with _lock:
        if _model is not None:
            return _model

        from faster_whisper import WhisperModel

        _ensure_hf_endpoint()
        device, compute = resolve_device_and_compute()
        download_root = settings.whisper_download_root
        download_root.mkdir(parents=True, exist_ok=True)
        model_ref = _resolve_model_ref(settings.whisper_model, download_root)

        logger.info(
            "加载 Whisper model=%s device=%s compute=%s root=%s",
            model_ref,
            device,
            compute,
            download_root,
        )
        try:
            _model = WhisperModel(
                model_ref,
                device=device,
                compute_type=compute,
                download_root=str(download_root),
            )
            _resolved_device = device
            _resolved_compute = compute
            _last_error = None
        except Exception as exc:
            _last_error = str(exc)
            if device == "cuda":
                logger.warning("CUDA 加载 Whisper 失败，回退 CPU: %s", exc)
                _model = WhisperModel(
                    model_ref,
                    device="cpu",
                    compute_type="int8",
                    download_root=str(download_root),
                )
                _resolved_device = "cpu"
                _resolved_compute = "int8"
                _last_error = f"cuda_fallback: {exc}"
            else:
                raise
        return _model


def _resolve_model_ref(model: str, download_root: Path) -> str:
    """优先使用本机已下载目录，避免启动时再走 Hub。"""
    raw = (model or "small").strip()
    as_path = Path(raw)
    if as_path.is_dir() and (as_path / "model.bin").is_file():
        return str(as_path.resolve())

    local_named = download_root / f"faster-whisper-{raw}"
    if local_named.is_dir() and (local_named / "model.bin").is_file():
        return str(local_named.resolve())

    # snapshot_download 默认缓存布局
    for child in download_root.glob("models--Systran--faster-whisper-*"):
        snaps = list((child / "snapshots").glob("*")) if (child / "snapshots").is_dir() else []
        for snap in snaps:
            if (snap / "model.bin").is_file():
                if raw in child.name or child.name.endswith(raw):
                    return str(snap.resolve())
    return raw


def transcribe_file(
    path: Path,
    *,
    language: str | None = "zh",
    vad_filter: bool = True,
) -> WhisperResult:
    started = time.perf_counter()
    model = get_model()
    lang = language if language else None
    segments_iter, info = model.transcribe(
        str(path),
        language=lang,
        vad_filter=vad_filter,
        beam_size=5,
    )
    segments: list[dict[str, Any]] = []
    texts: list[str] = []
    for seg in segments_iter:
        text = (seg.text or "").strip()
        if not text:
            continue
        texts.append(text)
        segments.append(
            {
                "start": float(seg.start),
                "end": float(seg.end),
                "text": text,
            }
        )
    elapsed = time.perf_counter() - started
    return WhisperResult(
        text="".join(texts).strip() or " ".join(texts).strip(),
        language=getattr(info, "language", None),
        duration=float(getattr(info, "duration", 0.0) or 0.0) or None,
        elapsed_sec=round(elapsed, 3),
        segments=segments,
    )
