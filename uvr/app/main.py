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
)

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger("uvr-api")

app = FastAPI(
    title="ONIJ UVR API",
    description="将本机 Ultimate Vocal Remover 模型通过 HTTP 暴露为伴奏提取服务",
    version="0.1.0",
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
    return HealthResponse(
        ok=settings.uvr_exe.is_file() and settings.models_root.is_dir(),
        uvr_exe=str(settings.uvr_exe),
        uvr_exe_exists=settings.uvr_exe.is_file(),
        models_root=str(settings.models_root),
        ffmpeg=str(ffmpeg),
        ffmpeg_exists=ffmpeg.is_file(),
    )


@app.get("/api/v1/models", response_model=list[ModelInfo])
def models() -> list[ModelInfo]:
    return [ModelInfo(**item) for item in list_local_models()]


@app.post("/api/v1/separate/instrumental", response_model=SeparateResponse)
async def separate_instrumental_api(
    file: UploadFile = File(..., description="原始音频文件"),
    architecture: str = Form("mdx"),
    model: str | None = Form(None),
    output_format: str = Form("WAV"),
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


def run() -> None:
    import sys

    kwargs = {
        "host": settings.host,
        "port": settings.port,
        "reload": False,
    }
    if getattr(sys, "frozen", False):
        uvicorn.run(app, **kwargs)
    else:
        uvicorn.run("app.main:app", **kwargs)


if __name__ == "__main__":
    run()
