#!/usr/bin/env python3
"""远程部署：git pull → 重启 onij-server → 同步静态资源。"""
import os
import subprocess
import sys
import time

import paramiko

HOST = "117.72.33.136"
USER = "root"
BRANCH = "feat/onij"
ENV_KEYS = ("ip", "mypwd", "sk", "QINIU_SK", "dskey")


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


def main() -> int:
    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not password:
        print("DEPLOY_SSH_PWD / sshpwd 未配置", file=sys.stderr)
        return 1

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(HOST, username=USER, password=password, timeout=30)
        write_remote_env(client.open_sftp(), "/tmp/onij-deploy.env")

        pull = run(
            client,
            f"""
set -e
cd ~/workbench/go/onij
git fetch origin {BRANCH}
git checkout {BRANCH}
git reset --hard origin/{BRANCH}
echo PULL_OK $(git rev-parse --short HEAD)
""",
            timeout=300,
        )
        if "PULL_OK" not in pull:
            return 1

        # 勿用 pkill -f：会误杀包含路径字符串的部署脚本本身
        run(client, "pkill -x onij-server || true; sleep 1")
        run(
            client,
            r"""
cat > /tmp/start-onij.sh <<'EOF'
#!/bin/bash
set -a
. /tmp/onij-deploy.env
set +a
unset QINIU_PUBLIC_DOMAIN
cd /root/workbench/go/onij/release/bin
chmod +x ./onij-server
setsid ./onij-server >> nohup.out 2>&1 < /dev/null &
echo STARTED $!
EOF
bash /tmp/start-onij.sh
rm -f /tmp/onij-deploy.env /tmp/start-onij.sh
""",
        )

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
            print(f"waiting ping... {i+1}")
        if not ok:
            run(client, "tail -40 /root/workbench/go/onij/release/bin/nohup.out")
            return 1

        run(
            client,
            r"""curl -sf -X POST http://127.0.0.1:8889/tran/list -H 'Content-Type: application/json' -d '{}' | head -c 200; echo; echo tran_ok""",
        )

        web = run(
            client,
            r"""
set -e
sudo rsync -av --delete ~/workbench/go/onij/release/web/ /var/www/onij/web/
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {} \;
sudo find /var/www/onij/web -type f -exec chmod 644 {} \;
nginx -t && systemctl reload nginx
grep -o 'index-[^"]*\.js' /var/www/onij/web/index.html | head -1
echo DEPLOY_OK
""",
            timeout=180,
        )
        return 0 if "DEPLOY_OK" in web else 1
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
