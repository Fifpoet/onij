#!/usr/bin/env python3
"""线上申请 Let's Encrypt 证书并写入 nginx HTTPS 配置。"""
import os
import subprocess
import sys
from pathlib import Path

import paramiko

HOST = "117.72.33.136"
USER = "root"
REMOTE_NGINX = "/etc/nginx/conf.d/web-vue.conf"
ROOT = Path(__file__).resolve().parents[1]
LOCAL_NGINX = ROOT / "deploy" / "nginx" / "web-vue.conf"


def user_env(name: str) -> str:
    val = os.environ.get(name, "").strip()
    if val:
        return val
    if sys.platform == "win32":
        r = subprocess.run(
            ["powershell", "-NoProfile", "-Command", f"[Environment]::GetEnvironmentVariable('{name}','User')"],
            capture_output=True,
            text=True,
            check=False,
        )
        return (r.stdout or "").strip()
    return ""


def run(client: paramiko.SSHClient, sh: str, timeout: int = 180) -> str:
    _, stdout, stderr = client.exec_command(sh, timeout=timeout)
    out = stdout.read().decode(errors="replace")
    err = stderr.read().decode(errors="replace")
    text = out + (("\n" + err) if err else "")
    print(text)
    return text


def main() -> int:
    password = os.environ.get("DEPLOY_SSH_PWD", "").strip() or user_env("sshpwd")
    if not password:
        print("DEPLOY_SSH_PWD / sshpwd 未配置", file=sys.stderr)
        return 1
    nginx_body = LOCAL_NGINX.read_text(encoding="utf-8")

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(HOST, username=USER, password=password, timeout=30)
        sftp = client.open_sftp()

        # 1) 先保证 HTTP 能响应 ACME（证书还不存在时不能听 443）
        http_only = """
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name onij.fun www.onij.fun _;
    root /var/www/onij/web;
    index index.html;
    location ^~ /.well-known/acme-challenge/ {
        default_type "text/plain";
    }
    location / {
        try_files $uri $uri/ /index.html;
    }
}
"""
        with sftp.file(REMOTE_NGINX, "w") as f:
            f.write(http_only)
        t = run(client, "nginx -t && systemctl reload nginx && echo NGINX_HTTP_OK")
        if "NGINX_HTTP_OK" not in t:
            return 1

        t = run(
            client,
            """
set -e
export DEBIAN_FRONTEND=noninteractive
if ! command -v certbot >/dev/null 2>&1; then
  apt-get update -qq
  apt-get install -y certbot
fi
certbot --version
certbot certonly --webroot -w /var/www/onij/web \
  -d onij.fun -d www.onij.fun \
  --non-interactive --agree-tos --register-unsafely-without-email \
  --keep-until-expiring
test -f /etc/letsencrypt/live/onij.fun/fullchain.pem
echo CERT_OK
""",
            timeout=300,
        )
        if "CERT_OK" not in t:
            print("证书申请失败：看上面 certbot 日志。常见原因是 80 未对公网开放或域名未解析到本机。")
            return 1

        with sftp.file(REMOTE_NGINX, "w") as f:
            f.write(nginx_body)
        t = run(client, "nginx -t && systemctl reload nginx && echo NGINX_SSL_OK")
        if "NGINX_SSL_OK" not in t:
            return 1

        t = run(
            client,
            r"""
curl -sf --max-time 5 http://127.0.0.1/.well-known/acme-challenge/ >/dev/null || true
curl -skf --max-time 5 https://127.0.0.1/go/ping -H 'Host: onij.fun' && echo GO_PROXY_OK
curl -skI --max-time 5 https://127.0.0.1/ -H 'Host: onij.fun' | head -8
ss -lntp | grep -E ':443|:80' || true
""",
        )
        if "GO_PROXY_OK" not in t:
            print("本机 443 反代 /go/ping 失败")
            return 1
        print("HTTPS_SETUP_OK")
        return 0
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
