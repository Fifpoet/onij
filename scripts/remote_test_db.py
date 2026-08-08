#!/usr/bin/env python3
import os
import subprocess
import sys

import paramiko

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


def main() -> int:
    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    ip = user_env("ip")
    mypwd = user_env("mypwd")
    print(f"local ip configured: {'yes' if ip else 'no'}, mypwd configured: {'yes' if mypwd else 'no'}")

    lines = [f"{k}={user_env(k)}" for k in ENV_KEYS if user_env(k)]
    content = "\n".join(lines) + "\n"

    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect("117.72.33.136", username="root", password=password, timeout=30)
    sftp = c.open_sftp()
    with sftp.file("/tmp/onij-test.env", "w") as f:
        f.write(content)
    sftp.chmod("/tmp/onij-test.env", 0o600)
    sftp.close()

    test_sh = r"""
set -a
source /tmp/onij-test.env
set +a
echo "ip_len=${#ip} mypwd_len=${#mypwd} sk_len=${#sk}"
mysql -uroot -p"$mypwd" -h "$ip" -e 'select 1 as ok' 2>&1 | tail -3
rm -f /tmp/onij-test.env
pgrep -af onij-server || echo no_process
"""
    _, o, _ = c.exec_command(test_sh, get_pty=True, timeout=60)
    print(o.read().decode(errors="replace"))
    c.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
