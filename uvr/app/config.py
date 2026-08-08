from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict

from app.paths import app_root, default_output_dir, default_work_dir


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=app_root() / ".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    uvr_root: Path = Path(
        r"C:\Users\onij\Downloads\Application\Ultimate Vocal Remover",
    )
    uvr_exe: Path | None = None

    host: str = "0.0.0.0"
    port: int = 5555

    default_architecture: str = "mdx"
    default_model: str = "UVR-MDX-NET-Inst_HQ_3.onnx"

    work_dir: Path = default_work_dir()
    output_dir: Path = default_output_dir()

    # Whisper / faster-whisper
    whisper_model: str = "medium"
    whisper_device: str = "cuda"  # cuda | cpu | auto
    whisper_compute_type: str = "float16"
    whisper_download_root: Path | None = None
    # 国内可设 https://hf-mirror.com ，写入 HF_ENDPOINT
    hf_endpoint: str = "https://hf-mirror.com"

    def model_post_init(self, __context: object) -> None:
        if self.uvr_exe is None:
            self.uvr_exe = self.uvr_root / "UVR.exe"
        if self.whisper_download_root is None:
            self.whisper_download_root = self.work_dir / "whisper-models"
        self.work_dir.mkdir(parents=True, exist_ok=True)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self.whisper_download_root.mkdir(parents=True, exist_ok=True)

    @property
    def models_root(self) -> Path:
        return self.uvr_root / "models"

    @property
    def mdx_models_dir(self) -> Path:
        return self.models_root / "MDX_Net_Models"

    @property
    def vr_models_dir(self) -> Path:
        return self.models_root / "VR_Models"

    @property
    def demucs_models_dir(self) -> Path:
        return self.models_root / "Demucs_Models"


settings = Settings()
