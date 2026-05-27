"""PyInstaller / python 统一入口。"""
from app.pyinstaller_bootstrap import bootstrap_pyinstaller

bootstrap_pyinstaller()

from app.main import run

if __name__ == "__main__":
    run()
