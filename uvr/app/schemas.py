from typing import Literal

from pydantic import BaseModel, Field


Architecture = Literal["mdx", "vr", "demucs"]
OutputFormat = Literal["WAV", "MP3", "FLAC", "M4A"]


class SeparateInstrumentalParams(BaseModel):
    """与 UVR5 常见选项对齐的可选参数（通过表单字段传入）。"""

    architecture: Architecture = Field(
        default="mdx",
        description="分离架构：mdx / vr / demucs",
    )
    model: str | None = Field(
        default=None,
        description="模型文件名，如 UVR-MDX-NET-Inst_HQ_3.onnx；默认使用配置中的 default_model",
    )
    output_format: OutputFormat = Field(default="MP3", description="输出音频格式")
    output_name: str | None = Field(
        default=None,
        description="输出文件基名（不含扩展名）",
    )

    use_gpu: bool = Field(default=True, description="是否使用 GPU（CUDA）")
    device_id: int = Field(default=0, ge=0, description="GPU 设备编号")

    # MDX / MDX-C
    overlap: float = Field(default=0.25, ge=0.0, le=0.99, description="MDX overlap")
    overlap_mdxc: int = Field(default=2, ge=2, le=50, description="MDX-C overlap")
    segment_size: int = Field(default=256, ge=32, description="MDX segment size")
    batch_size: int = Field(default=1, ge=1, description="batch size")

    # VR
    window_size: int = Field(default=512, description="VR window size")
    aggression: int = Field(default=5, ge=0, le=20, description="VR aggression")
    vr_enable_tta: bool = Field(default=False, description="VR TTA")
    vr_high_end_process: bool = Field(default=False, description="VR high end process")
    vr_post_process_threshold: float = Field(default=0.2, description="VR post process threshold")
    vr_model_param: str | None = Field(
        default=None,
        description="VR 模型参数名，如 4band_v3",
    )

    # Demucs
    demucs_stems: str | None = Field(
        default=None,
        description="Demucs stems，如 Vocals / Instrumental / All Stems",
    )
    demucs_shifts: int = Field(default=1, ge=0, description="Demucs shifts")
    demucs_overlap: float = Field(default=0.25, ge=0.0, le=0.99, description="Demucs overlap")

    # 通用
    pitch_shift: int = Field(default=0, description="音高偏移（半音，若底层支持）")
    normalize: bool = Field(default=False, description="输出归一化")
    invert_spect: bool = Field(default=False, description="反转频谱（部分模型）")
    single_stem: str | None = Field(
        default="Instrumental",
        description="仅导出指定 stem，伴奏场景默认 Instrumental",
    )


class HealthResponse(BaseModel):
    ok: bool
    uvr_exe: str
    uvr_exe_exists: bool
    models_root: str
    ffmpeg: str
    ffmpeg_exists: bool
    whisper_ok: bool = False
    whisper_model: str = "small"
    whisper_device: str = ""
    whisper_compute_type: str = ""
    whisper_loaded: bool = False
    whisper_error: str | None = None


class WhisperHealthResponse(BaseModel):
    ok: bool
    model: str
    device: str
    compute_type: str
    loaded: bool
    error: str | None = None


class WhisperSegment(BaseModel):
    start: float
    end: float
    text: str


class WhisperTranscribeResponse(BaseModel):
    text: str
    language: str | None = None
    duration: float | None = None
    elapsed_sec: float
    segments: list[WhisperSegment] = Field(default_factory=list)


class ModelInfo(BaseModel):
    architecture: Architecture
    filename: str
    path: str


class SeparateResponse(BaseModel):
    job_id: str
    instrumental_path: str
    instrumental_filename: str
    all_outputs: list[str]
    elapsed_sec: float


class SeparateFromUrlRequest(BaseModel):
    id: str | None = Field(
        default=None,
        description="业务 ID（如网易音乐 song id），用作 data 子目录名，相同 ID 复用已有伴奏",
    )
    source_url: str = Field(description="原音频可访问 URL")
    architecture: Architecture = "mdx"
    model: str | None = None
    output_format: OutputFormat = "MP3"
    output_name: str | None = None
    use_gpu: bool = True
    device_id: int = 0
    overlap: float = 0.25
    overlap_mdxc: int = 2
    segment_size: int = 256
    batch_size: int = 1
    window_size: int = 512
    aggression: int = 5
    vr_enable_tta: bool = False
    vr_high_end_process: bool = False
    vr_post_process_threshold: float = 0.2
    vr_model_param: str | None = None
    demucs_stems: str | None = None
    demucs_shifts: int = 1
    demucs_overlap: float = 0.25
    pitch_shift: int = 0
    normalize: bool = False
    invert_spect: bool = False
    single_stem: str | None = "Instrumental"
