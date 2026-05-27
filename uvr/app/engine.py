from __future__ import annotations

import logging
import os
import re
import subprocess
import time
import uuid
from pathlib import Path
from urllib.parse import unquote, urlparse

from app.config import settings
from app.schemas import Architecture, SeparateInstrumentalParams

logger = logging.getLogger(__name__)

ARCH_DIRS: dict[Architecture, str] = {
    "mdx": "MDX_Net_Models",
    "vr": "VR_Models",
    "demucs": "Demucs_Models",
}

_AUDIO_EXTS = {".mp3", ".flac", ".wav", ".m4a", ".ogg", ".opus", ".aac", ".wma", ".webm"}
_CONTENT_TYPE_EXT = {
    "audio/mpeg": ".mp3",
    "audio/mp3": ".mp3",
    "audio/flac": ".flac",
    "audio/wav": ".wav",
    "audio/x-wav": ".wav",
    "audio/mp4": ".m4a",
    "audio/aac": ".aac",
    "audio/ogg": ".ogg",
    "audio/opus": ".opus",
}


def setup_uvr_runtime() -> None:
    """把 UVR 安装目录下的 ffmpeg 加入 PATH，并校验 PyTorch 可用。"""
    uvr_root = settings.uvr_root
    if not uvr_root.is_dir():
        raise FileNotFoundError(f"UVR 安装目录不存在: {uvr_root}")

    ffmpeg = uvr_root / "ffmpeg.exe"
    if ffmpeg.is_file():
        root = str(uvr_root.resolve())
        current = os.environ.get("PATH", "")
        if root not in current.split(os.pathsep):
            os.environ["PATH"] = root + os.pathsep + current

    _ensure_torch_ready()
    _ensure_onnxruntime_ready()


def resolve_model_path(architecture: Architecture, model: str | None) -> tuple[Path, str]:
    model_name = model or settings.default_model
    model_dir = settings.models_root / ARCH_DIRS[architecture]
    if not model_dir.is_dir():
        raise FileNotFoundError(f"模型目录不存在: {model_dir}")

    direct = model_dir / model_name
    if direct.is_file():
        return model_dir, model_name

    # 允许只传不含扩展名的模型名
    for candidate in model_dir.iterdir():
        if candidate.is_file() and candidate.stem == Path(model_name).stem:
            return model_dir, candidate.name

    raise FileNotFoundError(
        f"未找到模型 {model_name!r}，请检查 {model_dir} 或先通过 UVR.exe 下载模型",
    )


def list_local_models() -> list[dict[str, str]]:
    items: list[dict[str, str]] = []
    for arch, folder in ARCH_DIRS.items():
        model_dir = settings.models_root / folder
        if not model_dir.is_dir():
            continue
        for path in sorted(model_dir.iterdir()):
            if not path.is_file():
                continue
            if path.suffix.lower() not in {".onnx", ".pth", ".ckpt", ".yaml", ".th"}:
                continue
            if path.suffix.lower() == ".yaml":
                continue
            items.append(
                {
                    "architecture": arch,
                    "filename": path.name,
                    "path": str(path),
                },
            )
    return items


def _ensure_torch_ready() -> None:
    try:
        import torch
    except OSError as exc:
        raise RuntimeError(
            "PyTorch 无法加载（常见原因：安装了不兼容的 torch 版本）。"
            "请在 uvr 目录执行: .venv\\Scripts\\pip install -r requirements.txt",
        ) from exc
    logger.info("PyTorch %s cuda=%s", torch.__version__, torch.cuda.is_available())


def _ensure_onnxruntime_ready() -> None:
    try:
        import onnxruntime as ort
    except (OSError, ImportError) as exc:
        raise RuntimeError(
            "onnxruntime 无法加载。请在 uvr 目录执行: .venv\\Scripts\\pip install -r requirements.txt",
        ) from exc
    providers = ort.get_available_providers()
    logger.info("onnxruntime %s providers=%s", ort.__version__, providers)
    if "CUDAExecutionProvider" not in providers:
        logger.warning(
            "ONNX Runtime 未启用 CUDA，MDX 分离将主要占用 CPU。"
            "有 NVIDIA 显卡请安装 requirements-gpu.txt 中的 onnxruntime-gpu",
        )


def _pick_instrumental(outputs: list[str], single_stem: str | None) -> str:
    if not outputs:
        raise RuntimeError("分离完成但未产生输出文件")

    stem = (single_stem or "Instrumental").lower()
    ranked: list[tuple[int, str]] = []
    for item in outputs:
        name = Path(item).name.lower()
        score = 0
        if stem in name:
            score += 10
        if "instrumental" in name or "inst" in name:
            score += 8
        if "no vocal" in name or "novocal" in name:
            score += 6
        if "vocal" in name and "instrumental" not in name:
            score -= 5
        if "karaoke" in name:
            score += 2
        ranked.append((score, item))

    ranked.sort(key=lambda x: x[0], reverse=True)
    best_score, best_path = ranked[0]
    if best_score <= 0 and len(outputs) == 2:
        # 常见输出为 [vocals, instrumental]，取第二个
        return outputs[-1]
    return best_path


def separate_instrumental(
    input_path: Path,
    params: SeparateInstrumentalParams,
    job_output_dir: Path | None = None,
) -> tuple[Path, list[Path], float]:
    setup_uvr_runtime()
    input_path = _ensure_wav_input(input_path)
    model_dir, model_filename = resolve_model_path(params.architecture, params.model)
    out_dir = job_output_dir or (settings.output_dir / uuid.uuid4().hex)
    out_dir.mkdir(parents=True, exist_ok=True)

    started = time.perf_counter()
    outputs = _run_separator(input_path, out_dir, model_dir, model_filename, params)
    elapsed = time.perf_counter() - started

    output_paths = [Path(p) for p in outputs]
    inst = Path(_pick_instrumental([str(p) for p in output_paths], params.single_stem))
    if params.output_name and inst.exists():
        target = out_dir / f"{params.output_name}{inst.suffix}"
        if target != inst:
            inst.replace(target)
            inst = target

    return inst, output_paths, elapsed


def _guess_audio_extension(url: str, content_type: str | None) -> str:
    path = unquote(urlparse(url).path)
    ext = Path(path).suffix.lower()
    if ext in _AUDIO_EXTS:
        return ext
    if content_type:
        mime = content_type.split(";", 1)[0].strip().lower()
        if mime in _CONTENT_TYPE_EXT:
            return _CONTENT_TYPE_EXT[mime]
    return ".mp3"


def download_url_to_file(url: str, dest_dir: Path, timeout: int = 180) -> Path:
    import urllib.request

    dest_dir.mkdir(parents=True, exist_ok=True)
    req = urllib.request.Request(url, headers={"User-Agent": "ONIJ-UVR-API/0.1"})
    with urllib.request.urlopen(req, timeout=timeout) as resp:
        content_type = resp.headers.get("Content-Type")
        data = resp.read()
    ext = _guess_audio_extension(url, content_type)
    dest = dest_dir / f"input{ext}"
    dest.write_bytes(data)
    logger.info("已下载音频 %s (%d bytes, %s)", dest.name, len(data), content_type or "unknown")
    return dest


def _ffmpeg_path() -> Path:
    ffmpeg = settings.uvr_root / "ffmpeg.exe"
    if not ffmpeg.is_file():
        raise FileNotFoundError(f"未找到 ffmpeg: {ffmpeg}")
    return ffmpeg


def _ensure_wav_input(input_path: Path) -> Path:
    """MP3 等格式先转 WAV，避免 audio-separator 对未知 subtype 处理失败。"""
    suffix = input_path.suffix.lower()
    if suffix in {".wav", ".flac"}:
        return input_path.resolve()

    wav_path = input_path.with_suffix(".wav")
    if wav_path.is_file() and wav_path.stat().st_size > 0:
        return wav_path.resolve()

    setup_uvr_runtime()
    cmd = [
        str(_ffmpeg_path()),
        "-y",
        "-loglevel",
        "error",
        "-i",
        str(input_path),
        "-acodec",
        "pcm_s16le",
        "-ar",
        "44100",
        str(wav_path),
    ]
    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0:
        detail = (proc.stderr or proc.stdout or "").strip()
        raise RuntimeError(f"ffmpeg 转 WAV 失败: {detail or proc.returncode}")
    if not wav_path.is_file() or wav_path.stat().st_size == 0:
        raise RuntimeError(f"ffmpeg 未生成有效 WAV: {wav_path}")

    logger.info("输入已转为 WAV: %s -> %s", input_path.name, wav_path.name)
    return wav_path.resolve()


def _collect_output_files(output_dir: Path) -> list[str]:
    exts = {".wav", ".flac", ".mp3", ".m4a", ".ogg"}
    return sorted(
        str(p)
        for p in output_dir.iterdir()
        if p.is_file() and p.suffix.lower() in exts
    )


def _sanitize_job_id(music_id: str | None) -> str:
    """将业务 ID（如网易 music id）转为安全的目录名。"""
    if not music_id:
        return uuid.uuid4().hex
    safe = re.sub(r"[^\w\-]", "", str(music_id).strip())
    if not safe:
        return uuid.uuid4().hex
    return safe


def _find_cached_instrumental(
    job_dir: Path,
    single_stem: str | None,
) -> tuple[Path, list[Path]] | None:
    """若该 ID 目录已有伴奏输出，直接复用。"""
    out_dir = job_dir / "out"
    if not out_dir.is_dir():
        return None
    outputs = _collect_output_files(out_dir)
    if not outputs:
        return None
    try:
        inst = Path(_pick_instrumental(outputs, single_stem))
    except RuntimeError:
        return None
    if not inst.is_file():
        return None
    return inst, [Path(p) for p in outputs]


def separate_instrumental_from_url(
    source_url: str,
    params: SeparateInstrumentalParams,
    job_output_dir: Path | None = None,
    music_id: str | None = None,
) -> tuple[Path, list[Path], float, str]:
    """返回 (伴奏路径, 全部输出, 耗时, job_id)"""
    job_id = _sanitize_job_id(music_id)
    job_dir = settings.work_dir / job_id
    job_dir.mkdir(parents=True, exist_ok=True)

    cached = _find_cached_instrumental(job_dir, params.single_stem)
    if cached is not None:
        inst_path, all_paths = cached
        logger.info("复用已有伴奏 job_id=%s path=%s", job_id, inst_path)
        return inst_path, all_paths, 0.0, job_id

    raw_path = download_url_to_file(source_url, job_dir)
    inst_path, all_paths, elapsed = separate_instrumental(
        raw_path,
        params,
        job_output_dir=job_output_dir or job_dir / "out",
    )
    return inst_path, all_paths, elapsed, job_id


# 与 audio_separator.Separator 默认 arch 参数对齐；勿只传部分字段以免覆盖后缺 hop_length 等
_DEFAULT_MDX_PARAMS = {
    "hop_length": 1024,
    "segment_size": 256,
    "overlap": 0.25,
    "batch_size": 1,
    "enable_denoise": False,
}
_DEFAULT_VR_PARAMS = {
    "batch_size": 1,
    "window_size": 512,
    "aggression": 5,
    "enable_tta": False,
    "enable_post_process": False,
    "post_process_threshold": 0.2,
    "high_end_process": False,
}
_DEFAULT_DEMUCS_PARAMS = {
    "segment_size": "Default",
    "shifts": 2,
    "overlap": 0.25,
    "segments_enabled": True,
}


def _build_separator_kwargs(
    model_dir: Path,
    output_dir: Path,
    params: SeparateInstrumentalParams,
    use_autocast: bool,
) -> dict:
    kwargs: dict = {
        "log_level": logging.INFO,
        "model_file_dir": str(model_dir),
        "output_dir": str(output_dir),
        "output_format": params.output_format,
        "use_autocast": use_autocast,
        "invert_using_spec": params.invert_spect,
    }
    if params.single_stem:
        kwargs["output_single_stem"] = params.single_stem

    if params.architecture == "mdx":
        kwargs["mdx_params"] = {
            **_DEFAULT_MDX_PARAMS,
            "overlap": params.overlap,
            "segment_size": params.segment_size,
            "batch_size": params.batch_size,
        }
    elif params.architecture == "vr":
        kwargs["vr_params"] = {
            **_DEFAULT_VR_PARAMS,
            "window_size": params.window_size,
            "aggression": params.aggression,
            "enable_tta": params.vr_enable_tta,
            "high_end_process": params.vr_high_end_process,
            "post_process_threshold": params.vr_post_process_threshold,
        }
    elif params.architecture == "demucs":
        kwargs["demucs_params"] = {
            **_DEFAULT_DEMUCS_PARAMS,
            "shifts": params.demucs_shifts,
            "overlap": params.demucs_overlap,
        }
    return kwargs


def _run_separator(
    input_path: Path,
    output_dir: Path,
    model_dir: Path,
    model_filename: str,
    params: SeparateInstrumentalParams,
) -> list[str]:
    from app.pyinstaller_bootstrap import bootstrap_pyinstaller

    bootstrap_pyinstaller()
    try:
        import torch
        from audio_separator.separator import Separator
    except ImportError as exc:
        raise RuntimeError(
            "缺少 audio-separator，请在 uvr 目录执行: run.bat",
        ) from exc

    cuda_ok = torch.cuda.is_available()
    use_autocast = bool(params.use_gpu and cuda_ok)
    if params.use_gpu and not cuda_ok:
        logger.warning(
            "请求 use_gpu=true 但当前 PyTorch 无 CUDA。"
            "有 NVIDIA 显卡请重启 run.bat 以安装 requirements-gpu.txt",
        )

    separator = Separator(
        **_build_separator_kwargs(model_dir, output_dir, params, use_autocast),
    )
    separator.load_model(model_filename)
    outputs = separator.separate(str(input_path))

    if isinstance(outputs, str):
        outputs = [outputs]
    else:
        outputs = list(outputs)

    if not outputs:
        outputs = _collect_output_files(output_dir)
    if not outputs:
        raise RuntimeError(
            "分离未产生输出文件，请查看日志中 audio-separator 的 ERROR",
        )
    return outputs
