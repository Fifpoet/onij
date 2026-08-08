"""run.bat 启动前依赖自检。"""
import sys

import onnxruntime as ort
import torch
from audio_separator.separator import Separator  # noqa: F401
from faster_whisper import WhisperModel  # noqa: F401
import ctranslate2

print(f"torch={torch.__version__} cuda={torch.cuda.is_available()}", file=sys.stderr)
if torch.cuda.is_available():
    print(f"gpu={torch.cuda.get_device_name(0)}", file=sys.stderr)
print(f"onnxruntime={ort.__version__}", file=sys.stderr)
print(f"providers={ort.get_available_providers()}", file=sys.stderr)
print(f"ctranslate2_cuda_devices={ctranslate2.get_cuda_device_count()}", file=sys.stderr)
