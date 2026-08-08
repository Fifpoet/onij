#!/usr/bin/env python3
import os
import subprocess
import sys

import paramiko

HOST = "117.72.33.136"
USER = "root"


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


def main() -> int:
    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not password:
        print("sshpwd 未配置", file=sys.stderr)
        return 1

    sh = """
cd ~/workbench/go/onij
git pull origin feat/onij
sudo rsync -av --delete release/web/ /var/www/onij/web/
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {} \\;
sudo find /var/www/onij/web -type f -exec chmod 644 {} \\;
# assets 目录若缺 +x 会导致 nginx 403
sudo chmod 755 /var/www/onij/web /var/www/onij/web/assets 2>/dev/null || true
echo SYNC_OK
"""
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=password, timeout=30)
    _, stdout, _ = c.exec_command(sh, get_pty=True, timeout=180)
    out = stdout.read().decode(errors="replace")
    print(out)
    c.close()
    return 0 if "SYNC_OK" in out else 1


if __name__ == "__main__":
    raise SystemExit(main())
