# UVR HTTP 服务

在 Windows 本机将 **Ultimate Vocal Remover (UVR5)** 的模型与分离能力通过 HTTP 暴露，默认端口 **5555**。

> UVR.exe 本身为 GUI 程序，无官方 CLI。本服务通过 [audio-separator](https://github.com/karaokemicrophone/audio-separator) 调用**同一套 UVR 模型目录**与安装包内的 `ffmpeg.exe`，效果与 UVR 内分离一致。

## 前置条件

- Windows + 已安装 UVR5（默认路径见 `config.example.toml`）
- Python 3.10+
- 建议 NVIDIA GPU + CUDA（可在请求中关闭 `use_gpu`）

## 安装与启动

```bat
cd uvr
run.bat
```

或手动：

```bat
python -m venv .venv
.venv\Scripts\pip install -r requirements.txt
set UVR_ROOT=C:\Users\onij\Downloads\Application\Ultimate Vocal Remover
.venv\Scripts\python -m app.main
```

服务地址：`http://127.0.0.1:5555`

## 环境变量

| 变量 | 说明 | 默认 |
|------|------|------|
| `UVR_ROOT` | UVR 安装目录 | `C:\Users\onij\Downloads\Application\Ultimate Vocal Remover` |
| `UVR_EXE` | UVR.exe 路径 | `{UVR_ROOT}\UVR.exe` |
| `HOST` | 监听地址 | `0.0.0.0` |
| `PORT` | 端口 | `5555` |
| `DEFAULT_MODEL` | 默认 MDX 伴奏模型 | `UVR-MDX-NET-Inst_HQ_3.onnx` |

## API

### 健康检查

```http
GET /health
```

响应含 UVR 与 Whisper 状态字段（`whisper_ok` / `whisper_device` 等）。

### Whisper 健康检查

```http
GET /api/v1/whisper/health
```

### Whisper 转写

```http
POST /api/v1/whisper/transcribe
Content-Type: multipart/form-data
```

| 字段 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `file` | file | 必填 | 音频文件 |
| `language` | str | `zh` | 语言；空字符串则自动检测 |
| `vad_filter` | bool | `true` | VAD 过滤静音 |

默认模型：`faster-whisper` **small** + **CUDA**（`float16`）。模型缓存目录：`{WORK_DIR}/whisper-models/`。

### 列出本机模型

```http
GET /api/v1/models
```

### 提取伴奏（multipart）

```http
POST /api/v1/separate/instrumental
Content-Type: multipart/form-data
```

| 字段 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `file` | file | 必填 | 原始音频 |
| `architecture` | str | `mdx` | `mdx` / `vr` / `demucs` |
| `model` | str | 配置默认 | 模型文件名 |
| `output_format` | str | `WAV` | `WAV` / `MP3` / `FLAC` / `M4A` |
| `use_gpu` | bool | `true` | GPU 加速 |
| `device_id` | int | `0` | GPU 编号 |
| `overlap` | float | `0.25` | MDX overlap |
| `segment_size` | int | `256` | MDX segment |
| `batch_size` | int | `1` | batch |
| `window_size` | int | `512` | VR window |
| `aggression` | int | `5` | VR aggression |
| `vr_enable_tta` | bool | `false` | VR TTA |
| `single_stem` | str | `Instrumental` | 仅导出指定 stem |
| … | | | 其余 VR/Demucs 参数见 `app/schemas.py` |

响应示例：

```json
{
  "job_id": "abc123",
  "instrumental_path": "C:/.../uvr/data/abc123/out/song_(Instrumental)_....wav",
  "instrumental_filename": "song_(Instrumental)_....wav",
  "all_outputs": ["..."],
  "elapsed_sec": 12.345
}
```

下载伴奏文件：

```http
GET /api/v1/jobs/{job_id}/instrumental
```

### curl 示例

```bash
curl -X POST "http://127.0.0.1:5555/api/v1/separate/instrumental" \
  -F "file=@song.mp3" \
  -F "architecture=mdx" \
  -F "model=UVR-MDX-NET-Inst_HQ_3.onnx" \
  -F "use_gpu=true"
```

## 目录结构

```
uvr/
  app/
    main.py       # FastAPI 入口
    engine.py     # 分离引擎（UVR 模型目录 + audio-separator）
    schemas.py    # 请求/响应模型
    config.py     # 配置
  data/           # 上传与任务临时文件
  output/         # 可选输出目录
  run.bat
```

## 说明

- 首次使用某模型前，请先用 **UVR.exe** 下载对应模型到 `models/` 目录。
- 分离耗时取决于音频长度与 GPU；长任务建议客户端轮询或自行增加异步队列。
- KTV 页「伴」功能可对接 `POST /api/v1/separate/instrumental`。
