#!/usr/bin/env python3
"""部署 onij 到服务器。

服务器配置低（无 Go / pnpm），所以不在远端编译：
本地编译好 linux 二进制和前端产物，SFTP 上传，远端只做「替换 + 重启」。

流程：
  1. 远端 git pull（只要源码，二进制/静态资源已被 .gitignore 排除）
  2. 上传 release/bin/onij-server 与 release/web/
  3. 重启进程，等 /ping
  4. 同步 /var/www/onij/web，reload nginx

用法：
    python scripts/remote_deploy.py              # 完整部署
    python scripts/remote_deploy.py --skip-web   # 只更新后端
    python scripts/remote_deploy.py --skip-build # 用已有产物，不重新编译
"""
from __future__ import annotations

import os
import subprocess
import sys
import time
from pathlib import Path

import paramiko

HOST = "117.72.33.136"
USER = "root"
BRANCH = "feat/onij"
REMOTE_REPO = "/root/workbench/go/onij"
ENV_KEYS = ("ip", "mypwd", "sk", "QINIU_SK", "dskey")

ROOT = Path(__file__).resolve().parent.parent
LOCAL_BIN = ROOT / "release" / "bin" / "onij-server"
LOCAL_WEB = ROOT / "release" / "web"


def user_env(name: str) -> str:
    val = os.environ.get(name, "").strip()
    if val:
        return val
    if sys.platform == "win32":
        ps = f"[Environment]::GetEnvironmentVariable('{name}','User')"
        r = subprocess.run(
            ["powershell", "-NoProfile", "-Command", ps],
            capture_output=True, text=True, check=False,
        )
        return (r.stdout or "").strip()
    return ""


def run(client: paramiko.SSHClient, sh: str, timeout: int = 180) -> str:
    _, stdout, stderr = client.exec_command(sh, timeout=timeout)
    out = stdout.read().decode(errors="replace")
    err = stderr.read().decode(errors="replace")
    text = out + (("\n" + err) if err else "")
    try:
        print(text)
    except UnicodeEncodeError:
        enc = getattr(sys.stdout, "encoding", None) or "utf-8"
        print(text.encode(enc, errors="replace").decode(enc, errors="replace"))
    return text


def write_remote_env(sftp: paramiko.SFTPClient, path: str, password: str) -> None:
    lines = []
    for key in ENV_KEYS:
        val = user_env(key)
        if val:
            lines.append(f"{key}={val}")
    if not user_env("ip") or not user_env("mypwd"):
        raise RuntimeError("本机 User 环境变量 ip / mypwd 未配置，无法启动后端")
    with sftp.file(path, "w") as f:
        f.write("\n".join(lines) + "\n")
    sftp.chmod(path, 0o600)


def build(skip_build: bool, skip_web: bool) -> None:
    if skip_build:
        print("[build] 跳过，使用已有产物")
        return
    npx = "pnpm" if subprocess.run(
        ["pnpm", "-v"], capture_output=True, shell=True
    ).returncode == 0 else "npm"
    print(f"[build] 前端（{npx}）")
    subprocess.run(f"{npx} run build:vite", cwd=ROOT, shell=True, check=True)
    print("[build] 后端（linux/amd64）")
    env = {**os.environ, "CGO_ENABLED": "0", "GOOS": "linux", "GOARCH": "amd64"}
    subprocess.run(
        ["go", "build", "-trimpath", "-ldflags=-s -w", "-o", str(LOCAL_BIN), "."],
        cwd=ROOT / "server", env=env, check=True,
    )


def upload(sftp: paramiko.SFTPClient, skip_web: bool) -> None:
    """上传二进制与前端产物。"""
    if not LOCAL_BIN.exists():
        raise RuntimeError(f"缺少二进制：{LOCAL_BIN}")
    remote_bin = f"{REMOTE_REPO}/release/bin/onij-server"
    print(f"[upload] {LOCAL_BIN.name} ({LOCAL_BIN.stat().st_size // 1048576} MB)")
    sftp.put(str(LOCAL_BIN), remote_bin + ".new")
    sftp.chmod(remote_bin + ".new", 0o755)

    if skip_web:
        return
    remote_web = f"{REMOTE_REPO}/release/web"
    sftp_mkdirs(sftp, remote_web + "/assets")
    print("[upload] index.html")
    sftp.put(str(LOCAL_WEB / "index.html"), remote_web + "/index.html")
    files = sorted((LOCAL_WEB / "assets").iterdir())
    print(f"[upload] assets {len(files)} 个文件")
    for i, f in enumerate(files, 1):
        if f.is_file():
            sftp.put(str(f), f"{remote_web}/assets/{f.name}")
        if i % 10 == 0:
            print(f"  {i}/{len(files)}")


def sftp_mkdirs(sftp: paramiko.SFTPClient, path: str) -> None:
    parts, cur = path.strip("/").split("/"), ""
    for p in parts:
        cur += "/" + p
        try:
            sftp.stat(cur)
        except IOError:
            sftp.mkdir(cur)


def main() -> int:
    skip_web = "--skip-web" in sys.argv
    skip_build = "--skip-build" in sys.argv

    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not password:
        print("DEPLOY_SSH_PWD / sshpwd 未配置", file=sys.stderr)
        return 1

    build(skip_build, skip_web)

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(HOST, username=USER, password=password, timeout=30)
    try:
        sftp = client.open_sftp()
        sftp_mkdirs(sftp, f"{REMOTE_REPO}/release/bin")

        # 1. 只拉源码；二进制稍后覆盖上传
        pull = run(
            client,
            f"""
set -e
cd {REMOTE_REPO}
git fetch --depth=1 origin {BRANCH} 2>/dev/null || git fetch origin {BRANCH}
git checkout -f {BRANCH}
git reset --hard origin/{BRANCH}
echo PULL_OK $(git rev-parse --short HEAD)
""",
            timeout=300,
        )
        if "PULL_OK" not in pull:
            print("拉取失败", file=sys.stderr)
            return 1

        # 2. 上传产物
        upload(sftp, skip_web)

        # 3. 就位并重启（勿用 pkill -f，会误杀含路径的脚本）
        run(client, "pkill -x onij-server || true; sleep 1")
        write_remote_env(sftp, "/tmp/onij-deploy.env", password)
        run(
            client,
            f"""
set -e
mv -f {REMOTE_REPO}/release/bin/onij-server.new {REMOTE_REPO}/release/bin/onij-server
chmod +x {REMOTE_REPO}/release/bin/onij-server
cat > /tmp/start-onij.sh <<'EOF'
#!/bin/bash
set -a
. /tmp/onij-deploy.env
set +a
unset QINIU_PUBLIC_DOMAIN
cd /root/workbench/go/onij/release/bin
setsid ./onij-server >> nohup.out 2>&1 < /dev/null &
echo STARTED $!
EOF
bash /tmp/start-onij.sh
rm -f /tmp/onij-deploy.env /tmp/start-onij.sh
""",
        )

        # 4. 等健康检查
        ok = False
        for i in range(30):
            time.sleep(2)
            r = run(
                client,
                "curl -sf --max-time 2 http://127.0.0.1:8889/ping && echo ping_ok || echo wait",
                timeout=30,
            )
            if "ping_ok" in r:
                ok = True
                break
            print(f"waiting ping... {i + 1}")
        if not ok:
            run(client, f"tail -40 {REMOTE_REPO}/release/bin/nohup.out")
            return 1

        if skip_web:
            print("DEPLOY_OK (backend only)")
            return 0

        # 5. 静态资源落到 nginx（hashed 文件名会变，assets 不 --delete）
        web = run(
            client,
            f"""
set -e
cp {REMOTE_REPO}/deploy/nginx/web-vue.conf /etc/nginx/conf.d/web-vue.conf
mkdir -p /var/www/onij/web/assets
rsync -a --delete --exclude 'assets/' {REMOTE_REPO}/release/web/ /var/www/onij/web/
rsync -a {REMOTE_REPO}/release/web/assets/ /var/www/onij/web/assets/
mkdir -p /var/www/onij/web/assets
chown -R www-data:www-data /var/www/onij/web
find /var/www/onij/web -type d -exec chmod 755 {{}} \\;
find /var/www/onij/web -type f -exec chmod 644 {{}} \\;
nginx -t && systemctl reload nginx
echo "入口: $(grep -o 'index-[^\\"]*\\.js' /var/www/onij/web/index.html | head -1)"
echo DEPLOY_OK
""",
            timeout=180,
        )
        return 0 if "DEPLOY_OK" in web else 1
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
