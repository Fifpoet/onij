#!/usr/bin/env python3
"""探测服务器环境：是否装了 Go / Node 工具链，以及部署目录现状。"""
import os
import sys

import paramiko

HOST, USER = "117.72.33.136", "root"


def user_env(name: str) -> str:
    val = os.environ.get(name, "").strip()
    if val:
        return val
    if sys.platform == "win32":
        import subprocess

        ps = f"[Environment]::GetEnvironmentVariable('{name}','User')"
        r = subprocess.run(
            ["powershell", "-NoProfile", "-Command", ps],
            capture_output=True, text=True, check=False,
        )
        return (r.stdout or "").strip()
    return ""


def main() -> int:
    pwd = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not pwd:
        print("sshpwd 未配置", file=sys.stderr)
        return 1

    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=pwd, timeout=30)

    cmds = {
        "hostname": "hostname",
        "go": "go version 2>/dev/null || echo NO_GO",
        "node": "node -v 2>/dev/null || echo NO_NODE",
        "pnpm": "pnpm -v 2>/dev/null || echo NO_PNPM",
        "repo": "ls -d ~/workbench/go/onij 2>/dev/null || echo NO_REPO",
        "branch": "cd ~/workbench/go/onij && git rev-parse --abbrev-ref HEAD 2>/dev/null",
        "head": "cd ~/workbench/go/onij && git log --oneline -1 2>/dev/null",
        "gitdir": "du -sh ~/workbench/go/onij/.git 2>/dev/null",
        "arch": "uname -m",
        "ping": "curl -sf --max-time 2 http://127.0.0.1:8889/ping && echo PING_OK || echo PING_DOWN",
    }
    for label, cmd in cmds.items():
        _, out, err = c.exec_command(cmd, timeout=60)
        o = out.read().decode(errors="replace").strip()
        e = err.read().decode(errors="replace").strip()
        print(f"[{label}] {o or e}")
    c.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
