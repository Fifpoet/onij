# 阻止 PyInstaller 收集 onnx.reference（打包阶段 import 会 access violation）
excludedimports = [
    "onnx.reference",
    "onnx.reference.op_run",
    "onnx.reference.reference_evaluator",
    "onnx.tools",
    "onnx.backend",
]
