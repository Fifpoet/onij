# -*- mode: python ; coding: utf-8 -*-
# 在 uvr 目录执行: scripts\build_exe.bat
#
# 注意: collect_all('fastapi'/'onnx' 等) 会让 PyInstaller 在分析阶段 import
# onnx.reference 并触发 access violation。仅收集 torch / onnxruntime 二进制。

from PyInstaller.utils.hooks import collect_all, collect_submodules

block_cipher = None

# 运行时不需要、且会在打包阶段导致崩溃的模块
EXCLUDES = [
    "onnx.reference",
    "onnx.tools",
    "onnx.backend",
    "onnx.checker",
    "onnx.helper",
    "onnx.shape_inference",
    "onnx.version_converter",
    "onnx.numpy_helper",
    "onnx.external_data_helper",
    "sklearn",
    "torch.testing",
    "torch.distributed",
    "torch.utils.tensorboard",
    "matplotlib",
    "IPython",
    "notebook",
    "pytest",
    "sympy.testing",
    "numba.tests",
    "scipy.tests",
]

datas = []
binaries = []
hiddenimports = [
    "app",
    "app.main",
    "app.engine",
    "app.config",
    "app.paths",
    "app.schemas",
    "uvicorn.logging",
    "uvicorn.loops",
    "uvicorn.loops.auto",
    "uvicorn.protocols",
    "uvicorn.protocols.http",
    "uvicorn.protocols.http.auto",
    "uvicorn.protocols.websockets",
    "uvicorn.protocols.websockets.auto",
    "uvicorn.lifespan",
    "uvicorn.lifespan.on",
    "multipart",
    "pydantic_settings",
    "onnx",
    "onnx.defs",
    "onnx2torch",
    "onnx2torch.utils",
    "onnx2torch.node_converters",
]

hiddenimports += collect_submodules("app")
hiddenimports += collect_submodules("audio_separator")

for pkg in ("torch", "onnxruntime"):
    try:
        tmp = collect_all(pkg)
        datas += tmp[0]
        binaries += tmp[1]
        hiddenimports += tmp[2]
    except Exception:
        pass

a = Analysis(
    ["launcher.py"],
    pathex=["."],
    binaries=binaries,
    datas=datas,
    hiddenimports=hiddenimports,
    hookspath=["scripts/hooks"],
    hooksconfig={},
    runtime_hooks=[],
    excludes=EXCLUDES,
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
    module_collection_mode={
        "onnx.reference": "ignore",
        "onnx.tools": "ignore",
        "onnx.backend": "ignore",
        "sklearn": "ignore",
        "torch.testing": "ignore",
        "torch.distributed": "ignore",
    },
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name="uvr-api",
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=False,
    console=True,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=False,
    upx_exclude=[],
    name="uvr-api",
)
