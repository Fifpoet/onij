from __future__ import annotations

import sys
from pathlib import Path

DEFAULT_WORK_DIR = Path(r"C:\Users\onij\Downloads\Workbench\data\uvr")


def app_root() -> Path:
    """源码目录或 PyInstaller 解压后的 exe 所在目录。"""
    if getattr(sys, "frozen", False):
        return Path(sys.executable).resolve().parent
    return Path(__file__).resolve().parent.parent


def default_work_dir() -> Path:
    return DEFAULT_WORK_DIR


def default_output_dir() -> Path:
    return DEFAULT_WORK_DIR / "_output"
