"""PyInstaller onedir：在导入 torch/audio-separator 前补全 sys._MEIPASS。"""
from __future__ import annotations

import os
import sys


def bootstrap_pyinstaller() -> None:
    if not getattr(sys, "frozen", False):
        return
    if hasattr(sys, "_MEIPASS"):
        return
    base = os.path.dirname(os.path.abspath(sys.executable))
    internal = os.path.join(base, "_internal")
    sys._MEIPASS = internal if os.path.isdir(internal) else base
