# -*- mode: python ; coding: utf-8 -*-
# 在 uvr 目录执行: scripts\build_exe.bat
#
# 仅排除 onnx.reference（打包分析阶段会 access violation）。
# 勿排除 onnx.helper / external_data_helper 等，onnx2torch 运行时需要。

from PyInstaller.utils.hooks import collect_all, collect_data_files, collect_submodules

block_cipher = None

# 只排除确定不需要、且会引发打包/体积问题的模块
EXCLUDES = [
    "onnx.reference",
    "onnx.tools",
    "onnx.backend",
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

# onnx2torch / audio_separator 运行时需要的 onnx 子模块（勿放进 EXCLUDES）
ONNX_RUNTIME_IMPORTS = [
    "onnx",
    "onnx.defs",
    "onnx.helper",
    "onnx.numpy_helper",
    "onnx.external_data_helper",
    "onnx.checker",
    "onnx.shape_inference",
    "onnx.version_converter",
    "onnx.serialization",
    "onnx.onnx_pb",
    "onnx.onnx_ml_pb2",
    "onnx.onnx_operators_ml_pb2",
    "onnx.onnx_data_pb2",
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
    "onnx2torch",
    "onnx2torch.utils",
    "onnx2torch.node_converters",
] + ONNX_RUNTIME_IMPORTS

hiddenimports += collect_submodules("app")
hiddenimports += collect_submodules("audio_separator")
hiddenimports += collect_submodules("onnx2torch")
datas += collect_data_files("audio_separator")

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
    runtime_hooks=["scripts/hooks/pyi_rth_meipass.py"],
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
    console=False,
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
