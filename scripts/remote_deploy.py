#!/usr/bin/env python3
"""远程部署：同步静态资源并重启 onij-server。"""
import io
import os
import subprocess
import sys

import paramiko

HOST = "117.72.33.136"
USER = "root"
BRANCH = "feat/onij"
ENV_KEYS = ("ip", "mypwd", "sk", "QINIU_SK")


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


def write_remote_env(sftp: paramiko.SFTPClient, path: str) -> None:
    lines = []
    for key in ENV_KEYS:
        val = user_env(key)
        if val:
            lines.append(f"{key}={val}")
    if not user_env("ip") or not user_env("mypwd"):
        raise RuntimeError("本机 User 环境变量 ip / mypwd 未配置，无法启动后端")
    content = "\n".join(lines) + "\n"
    with sftp.file(path, "w") as f:
        f.write(content)
    sftp.chmod(path, 0o600)


def main() -> int:
    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not password:
        print("DEPLOY_SSH_PWD / sshpwd 未配置", file=sys.stderr)
        return 1

    deploy_sh = f"""set -a
source /tmp/onij-deploy.env
set +a
rm -f /tmp/onij-deploy.env
unset QINIU_PUBLIC_DOMAIN
chmod +x ~/workbench/go/onij/release/bin/onij-server
cd ~/workbench/go/onij/release/bin
nohup ./onij-server >> nohup.out 2>&1 &
sleep 4
curl -sf http://127.0.0.1:8889/ping && echo ping_ok || (echo '--- nohup tail ---'; tail -30 nohup.out; echo ping_failed)
sudo rsync -av --delete ~/workbench/go/onij/release/web/ /var/www/onij/web/
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {{}} \\;
sudo find /var/www/onij/web -type f -exec chmod 644 {{}} \\;
nginx -t && systemctl reload nginx
grep -o 'index-[^\"]*\\.js' /var/www/onij/web/index.html | head -1
echo DEPLOY_OK
"""

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(HOST, username=USER, password=password, timeout=30)
        client.exec_command(f"cd ~/workbench/go/onij && git pull origin {BRANCH}", timeout=60)[1].read()
        write_remote_env(client.open_sftp(), "/tmp/onij-deploy.env")
        _, stdout, stderr = client.exec_command(deploy_sh, get_pty=True, timeout=180)
        out = stdout.read().decode(errors="replace")
        err = stderr.read().decode(errors="replace")
        if out:
            print(out)
        if err:
            print(err, file=sys.stderr)
        ok = "DEPLOY_OK" in out and "ping_ok" in out
        return 0 if ok else 1
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
