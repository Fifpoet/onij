from __future__ import annotations

import logging
import shutil
import uuid
from pathlib import Path

import uvicorn
from fastapi import FastAPI, File, Form, HTTPException, UploadFile
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse

from app.config import settings
from app.engine import (
    _sanitize_job_id,
    ensure_job_pitch,
    list_local_models,
    separate_instrumental,
    separate_instrumental_from_url,
    setup_uvr_runtime,
)
from app.schemas import (
    HealthResponse,
    ModelInfo,
    SeparateFromUrlRequest,
    SeparateInstrumentalParams,
    SeparateResponse,
    WhisperHealthResponse,
    WhisperSegment,
    WhisperTranscribeResponse,
)
from app.whisper_engine import get_status as whisper_status
from app.whisper_engine import transcribe_file

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger("uvr-api")

app = FastAPI(
    title="ONIJ UVR API",
    description="本机 UVR 伴奏提取 + Whisper 语音识别 + RMVPE 音高",
    version="0.3.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.on_event("startup")
def _startup() -> None:
    try:
        setup_uvr_runtime()
    except FileNotFoundError as exc:
        logger.warning("UVR 环境未就绪: %s", exc)


@app.get("/health", response_model=HealthResponse)
def health() -> HealthResponse:
    ffmpeg = settings.uvr_root / "ffmpeg.exe"
    wh = whisper_status()
    uvr_ok = settings.uvr_exe.is_file() and settings.models_root.is_dir()
    return HealthResponse(
        ok=uvr_ok,
        uvr_exe=str(settings.uvr_exe),
        uvr_exe_exists=settings.uvr_exe.is_file(),
        models_root=str(settings.models_root),
        ffmpeg=str(ffmpeg),
        ffmpeg_exists=ffmpeg.is_file(),
        whisper_ok=wh.ok,
        whisper_model=wh.model,
        whisper_device=wh.device,
        whisper_compute_type=wh.compute_type,
        whisper_loaded=wh.loaded,
        whisper_error=wh.error,
    )


@app.get("/api/v1/whisper/health", response_model=WhisperHealthResponse)
def whisper_health() -> WhisperHealthResponse:
    wh = whisper_status()
    return WhisperHealthResponse(
        ok=wh.ok,
        model=wh.model,
        device=wh.device,
        compute_type=wh.compute_type,
        loaded=wh.loaded,
        error=wh.error,
    )


@app.post("/api/v1/whisper/transcribe", response_model=WhisperTranscribeResponse)
async def whisper_transcribe(
    file: UploadFile = File(..., description="音频文件，如 wav/mp3/webm/m4a"),
    language: str = Form("zh"),
    vad_filter: bool = Form(True),
) -> WhisperTranscribeResponse:
    if not file.filename:
        raise HTTPException(status_code=400, detail="缺少文件名")

    job_id = uuid.uuid4().hex
    job_dir = settings.work_dir / "whisper" / job_id
    job_dir.mkdir(parents=True, exist_ok=True)
    suffix = Path(file.filename).suffix or ".wav"
    input_path = job_dir / f"input{suffix}"

    try:
        with input_path.open("wb") as f:
            shutil.copyfileobj(file.file, f)
        result = transcribe_file(
            input_path,
            language=language or None,
            vad_filter=vad_filter,
        )
    except FileNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except Exception as exc:
        logger.exception("Whisper 转写失败 job=%s", job_id)
        raise HTTPException(status_code=500, detail=str(exc)) from exc
    finally:
        try:
            shutil.rmtree(job_dir, ignore_errors=True)
        except Exception:
            pass

    return WhisperTranscribeResponse(
        text=result.text,
        language=result.language,
        duration=result.duration,
        elapsed_sec=result.elapsed_sec,
        segments=[WhisperSegment(**s) for s in result.segments],
    )


@app.get("/api/v1/models", response_model=list[ModelInfo])
def models() -> list[ModelInfo]:
    return [ModelInfo(**item) for item in list_local_models()]


@app.post("/api/v1/separate/instrumental", response_model=SeparateResponse)
async def separate_instrumental_api(
    file: UploadFile = File(..., description="原始音频文件"),
    architecture: str = Form("mdx"),
    model: str | None = Form(None),
    output_format: str = Form("MP3"),
    output_name: str | None = Form(None),
    use_gpu: bool = Form(True),
    device_id: int = Form(0),
    overlap: float = Form(0.25),
    overlap_mdxc: int = Form(2),
    segment_size: int = Form(256),
    batch_size: int = Form(1),
    window_size: int = Form(512),
    aggression: int = Form(5),
    vr_enable_tta: bool = Form(False),
    vr_high_end_process: bool = Form(False),
    vr_post_process_threshold: float = Form(0.2),
    vr_model_param: str | None = Form(None),
    demucs_stems: str | None = Form(None),
    demucs_shifts: int = Form(1),
    demucs_overlap: float = Form(0.25),
    pitch_shift: int = Form(0),
    normalize: bool = Form(False),
    invert_spect: bool = Form(False),
    single_stem: str | None = Form("Instrumental"),
) -> SeparateResponse:
    if not file.filename:
        raise HTTPException(status_code=400, detail="缺少文件名")

    params = SeparateInstrumentalParams(
        architecture=architecture,  # type: ignore[arg-type]
        model=model,
        output_format=output_format,  # type: ignore[arg-type]
        output_name=output_name,
        use_gpu=use_gpu,
        device_id=device_id,
        overlap=overlap,
        overlap_mdxc=overlap_mdxc,
        segment_size=segment_size,
        batch_size=batch_size,
        window_size=window_size,
        aggression=aggression,
        vr_enable_tta=vr_enable_tta,
        vr_high_end_process=vr_high_end_process,
        vr_post_process_threshold=vr_post_process_threshold,
        vr_model_param=vr_model_param,
        demucs_stems=demucs_stems,
        demucs_shifts=demucs_shifts,
        demucs_overlap=demucs_overlap,
        pitch_shift=pitch_shift,
        normalize=normalize,
        invert_spect=invert_spect,
        single_stem=single_stem,
    )

    job_id = uuid.uuid4().hex
    job_dir = settings.work_dir / job_id
    job_dir.mkdir(parents=True, exist_ok=True)
    suffix = Path(file.filename).suffix or ".wav"
    input_path = job_dir / f"input{suffix}"

    try:
        with input_path.open("wb") as f:
            shutil.copyfileobj(file.file, f)

        inst_path, all_paths, elapsed = separate_instrumental(
            input_path,
            params,
            job_output_dir=job_dir / "out",
        )
    except FileNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except Exception as exc:
        logger.exception("分离失败 job=%s", job_id)
        raise HTTPException(status_code=500, detail=str(exc)) from exc

    return SeparateResponse(
        job_id=job_id,
        instrumental_path=str(inst_path.resolve()),
        instrumental_filename=inst_path.name,
        all_outputs=[str(p.resolve()) for p in all_paths],
        elapsed_sec=round(elapsed, 3),
    )


@app.post("/api/v1/separate/instrumental/from-url", response_model=SeparateResponse)
async def separate_instrumental_from_url_api(body: SeparateFromUrlRequest) -> SeparateResponse:
    params = SeparateInstrumentalParams(
        architecture=body.architecture,
        model=body.model,
        output_format=body.output_format,
        output_name=body.output_name,
        use_gpu=body.use_gpu,
        device_id=body.device_id,
        overlap=body.overlap,
        overlap_mdxc=body.overlap_mdxc,
        segment_size=body.segment_size,
        batch_size=body.batch_size,
        window_size=body.window_size,
        aggression=body.aggression,
        vr_enable_tta=body.vr_enable_tta,
        vr_high_end_process=body.vr_high_end_process,
        vr_post_process_threshold=body.vr_post_process_threshold,
        vr_model_param=body.vr_model_param,
        demucs_stems=body.demucs_stems,
        demucs_shifts=body.demucs_shifts,
        demucs_overlap=body.demucs_overlap,
        pitch_shift=body.pitch_shift,
        normalize=body.normalize,
        invert_spect=body.invert_spect,
        single_stem=body.single_stem,
    )
    try:
        inst_path, all_paths, elapsed, job_id = separate_instrumental_from_url(
            body.source_url,
            params,
            music_id=body.id,
        )
    except FileNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except Exception as exc:
        logger.exception("URL 分离失败")
        raise HTTPException(status_code=500, detail=str(exc)) from exc

    return SeparateResponse(
        job_id=job_id,
        instrumental_path=str(inst_path.resolve()),
        instrumental_filename=inst_path.name,
        all_outputs=[str(p.resolve()) for p in all_paths],
        elapsed_sec=round(elapsed, 3),
    )


@app.get("/api/v1/jobs/{job_id}/instrumental")
def download_instrumental(job_id: str) -> FileResponse:
    job_dir = settings.work_dir / job_id / "out"
    if not job_dir.is_dir():
        raise HTTPException(status_code=404, detail="任务不存在")

    candidates = list(job_dir.glob("*"))
    inst_files = [
        p
        for p in candidates
        if p.is_file() and ("inst" in p.name.lower() or "instrumental" in p.name.lower())
    ]
    if not inst_files and candidates:
        inst_files = [p for p in candidates if p.is_file()]

    if not inst_files:
        raise HTTPException(status_code=404, detail="未找到伴奏文件")

    target = inst_files[0]
    return FileResponse(path=target, filename=target.name, media_type="application/octet-stream")


@app.get("/api/v1/jobs/{job_id}/pitch")
def download_pitch(job_id: str):
    safe_id = _sanitize_job_id(job_id)
    dest = ensure_job_pitch(safe_id)
    if dest is None or not dest.is_file():
        raise HTTPException(status_code=404, detail="未找到音高数据（需要人声模型 UVR-MDX-NET-Voc_FT.onnx）")
    return FileResponse(path=dest, filename="pitch.json", media_type="application/json")


def run() -> None:
    import copy
    import sys

    from uvicorn.config import LOGGING_CONFIG

    kwargs: dict = {
        "host": settings.host,
        "port": settings.port,
        "reload": False,
    }
    if getattr(sys, "frozen", False):
        # console=False 打包时 stdout/stderr 为 None，uvicorn 默认日志配置会崩溃
        import os

        if sys.stdout is None:
            sys.stdout = open(os.devnull, "w", encoding="utf-8")  # type: ignore[assignment]
        if sys.stderr is None:
            sys.stderr = open(os.devnull, "w", encoding="utf-8")  # type: ignore[assignment]
        log_config = copy.deepcopy(LOGGING_CONFIG)
        log_config["formatters"]["default"]["use_colors"] = False
        log_config["formatters"]["access"]["use_colors"] = False
        kwargs["log_config"] = log_config
        kwargs["use_colors"] = False
        uvicorn.run(app, **kwargs)
    else:
        uvicorn.run("app.main:app", **kwargs)


if __name__ == "__main__":
    run()
