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
set -x
ls -la /var/www/onij/web/index.html
ls -la /var/www/onij/web/assets/index-B5Rs1Dfg.js 2>&1 || true
namei -l /var/www/onij/web/assets/index-B5Rs1Dfg.js 2>&1 || true
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {} \\;
sudo find /var/www/onij/web -type f -exec chmod 644 {} \\;
ls -la /var/www/onij/web/assets/index-B5Rs1Dfg.js
curl -sI http://127.0.0.1/assets/index-B5Rs1Dfg.js | head -5
curl -sI http://onij.fun/assets/index-B5Rs1Dfg.js | head -5
echo FIX_OK
"""
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=password, timeout=30)
    _, stdout, stderr = c.exec_command(sh, get_pty=True, timeout=120)
    out = stdout.read().decode(errors="replace")
    err = stderr.read().decode(errors="replace")
    print(out)
    if err:
        print(err, file=sys.stderr)
    c.close()
    return 0 if "FIX_OK" in out else 1


if __name__ == "__main__":
    raise SystemExit(main())
