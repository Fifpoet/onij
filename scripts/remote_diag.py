#!/usr/bin/env python3
import os
import paramiko

password = os.environ.get("DEPLOY_SSH_PWD", "")
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("117.72.33.136", username="root", password=password, timeout=30)

cmds = [
    "find /root -maxdepth 4 -type f -name '*.env' 2>/dev/null",
    "find /root -maxdepth 4 -type f -name 'start*.sh' 2>/dev/null",
    "ls -la /root/env* /root/*.sh 2>/dev/null || true",
    "grep -n nohup ~/.bash_history | tail -10",
]
for cmd in cmds:
    print(">>>", cmd)
    _, o, _ = c.exec_command(cmd, timeout=30)
    print(o.read().decode(errors="replace") or "(empty)\n")
c.close()
