#!/usr/bin/env python3
"""把 release/web 静态资源上传到公开空间 onijweb，并确保 CDN 域名 static.onij.fun 存在。"""
from __future__ import annotations

import hmac
import hashlib
import json
import mimetypes
import os
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

AK = "lKrcOB9iUaKvKhspGh1hgo4-Dy0yFlH0mgRSgPRY"
BUCKET = "onijweb"
DOMAIN = "static.onij.fun"
UP_HOST = "https://up-z0.qiniup.com"
SKIP_SUFFIX = {".gz", ".map"}


def user_env(name: str) -> str:
    val = os.environ.get(name, "").strip()
    if val:
        return val
    if sys.platform == "win32":
        ps = f"[Environment]::GetEnvironmentVariable('{name}','User')"
        r = subprocess.run(
            ["powershell", "-NoProfile", "-Command", ps],
            capture_output=True,
            text=True,
            check=False,
        )
        return (r.stdout or "").strip()
    return ""


def get_sk() -> str:
    for key in ("QINIU_SK", "sk", "SK"):
        val = user_env(key)
        if val:
            return val
    return ""


def qbox(sk: str, path_query: str, body: bytes = b"") -> str:
    data = path_query.encode() + b"\n" + (body or b"")
    sig = hmac.new(sk.encode(), data, hashlib.sha1).digest()
    return "QBox " + AK + ":" + base64_url(sig)


def base64_url(raw: bytes) -> str:
    return __import__("base64").urlsafe_b64encode(raw).decode("ascii")


def call_qbox(sk: str, method: str, host: str, path: str, body: bytes | None = None, content_type: str | None = None):
    data = body
    # JSON 不参与 QBox 签名；仅 form-urlencoded 把 body 算进去
    sign_body = data if content_type == "application/x-www-form-urlencoded" else b""
    headers = {"Authorization": qbox(sk, path, sign_body or b"")}
    if content_type:
        headers["Content-Type"] = content_type
    req = urllib.request.Request("https://" + host + path, data=data, method=method, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            raw = resp.read().decode("utf-8", "replace")
            return resp.status, json.loads(raw) if raw else {}
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8", "replace")
        try:
            parsed = json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            parsed = {"error": raw}
        return e.code, parsed


def upload_token(sk: str, key: str) -> str:
    deadline = int(time.time()) + 3600
    policy = json.dumps({"scope": f"{BUCKET}:{key}", "deadline": deadline}, separators=(",", ":"))
    encoded = base64_url(policy.encode())
    sign = hmac.new(sk.encode(), encoded.encode(), hashlib.sha1).digest()
    return f"{AK}:{base64_url(sign)}:{encoded}"


def upload_file(sk: str, local: Path, key: str) -> None:
    boundary = "----onijWebCdn" + str(int(time.time() * 1000))
    mime = mimetypes.guess_type(local.name)[0] or "application/octet-stream"
    file_bytes = local.read_bytes()
    token = upload_token(sk, key)
    parts = [
        f"--{boundary}\r\nContent-Disposition: form-data; name=\"token\"\r\n\r\n{token}\r\n".encode(),
        f"--{boundary}\r\nContent-Disposition: form-data; name=\"key\"\r\n\r\n{key}\r\n".encode(),
        (
            f"--{boundary}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"{local.name}\"\r\n"
            f"Content-Type: {mime}\r\n\r\n"
        ).encode()
        + file_bytes
        + b"\r\n"
        + f"--{boundary}--\r\n".encode(),
    ]
    body = b"".join(parts)
    req = urllib.request.Request(
        UP_HOST,
        data=body,
        method="POST",
        headers={"Content-Type": f"multipart/form-data; boundary={boundary}"},
    )
    with urllib.request.urlopen(req, timeout=120) as resp:
        if resp.status >= 300:
            raise RuntimeError(f"upload {key} http {resp.status}")


def collect_files(web_root: Path) -> list[tuple[Path, str]]:
    out: list[tuple[Path, str]] = []
    for p in web_root.rglob("*"):
        if not p.is_file():
            continue
        if p.suffix.lower() in SKIP_SUFFIX:
            continue
        key = p.relative_to(web_root).as_posix()
        out.append((p, key))
    return out


def main() -> int:
    if os.environ.get("SKIP_QINIU_WEB_CDN", "").strip():
        print("skip qiniu web cdn (SKIP_QINIU_WEB_CDN)")
        return 0
    sk = get_sk()
    if not sk:
        print("skip qiniu web cdn (no sk / QINIU_SK)")
        return 0
    root = Path(__file__).resolve().parents[1]
    web = root / "release" / "web"
    if not web.is_dir():
        raise SystemExit(f"missing {web}, run build-release first")

    files = collect_files(web)
    print(f"upload {len(files)} files -> {BUCKET}")
    for i, (path, key) in enumerate(files, 1):
        upload_file(sk, path, key)
        if i == 1 or i % 10 == 0 or i == len(files):
            print(f"  {i}/{len(files)} {key}")

    code, body = call_qbox(sk, "GET", "api.qiniu.com", f"/domain/{DOMAIN}")
    if code == 200 and body.get("name") == DOMAIN:
        print("domain exists", body.get("cname"), body.get("operatingState"), body.get("protocol"))
        return 0

    create = {
        "type": "normal",
        "platform": "web",
        "geoCover": "china",
        "protocol": "http",
        "ipTypes": 1,
        "source": {
            "sourceType": "qiniuBucket",
            "sourceQiniuBucket": BUCKET,
            "testURLPath": "/tecent.png",
        },
        "cache": {
            "cacheControls": [
                {"time": 30, "timeunit": 3, "type": "all", "rule": "*"},
            ]
        },
    }
    raw = json.dumps(create, separators=(",", ":")).encode()
    code, body = call_qbox(
        sk,
        "POST",
        "api.qiniu.com",
        f"/domain/{DOMAIN}",
        raw,
        "application/json",
    )
    print("create domain", code, body)
    if code not in (200, 400):
        return 1
    code, body = call_qbox(sk, "GET", "api.qiniu.com", f"/domain/{DOMAIN}")
    print("cname", body.get("cname"), "state", body.get("operatingState"))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
