#!/usr/bin/env python3
"""仅同步 release/web 到 /var/www/onij/web/"""
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
sudo rsync -av --delete ~/workbench/go/onij/release/web/ /var/www/onij/web/
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {} \\;
sudo find /var/www/onij/web -type f -exec chmod 644 {} \\;
grep -o 'index-[^\"]*\\.js' /var/www/onij/web/index.html | head -1
curl -sf http://onij.fun/assets/index-B5Rs1Dfg.js | head -c 120 | grep -q onij.fun && echo bundle_ok || echo bundle_check
echo SYNC_OK
"""
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        c.connect(HOST, username=USER, password=password, timeout=30)
        # 先 scp 本地 release/web 到服务器临时目录再 rsync？服务器上 release/web 是旧的
        # 需要上传本地 build 产物
        sftp = c.open_sftp()
        import pathlib
        local_root = pathlib.Path(__file__).resolve().parents[1] / "release" / "web"
        remote_tmp = "/tmp/onij-web-deploy"
        c.exec_command(f"rm -rf {remote_tmp} && mkdir -p {remote_tmp}")[1].read()

        def upload_dir(local: pathlib.Path, remote: str) -> None:
            for item in local.iterdir():
                rpath = f"{remote}/{item.name}"
                if item.is_dir():
                    try:
                        sftp.mkdir(rpath)
                    except OSError:
                        pass
                    upload_dir(item, rpath)
                else:
                    sftp.put(str(item), rpath)

        print("uploading release/web ...")
        upload_dir(local_root, remote_tmp)
        sftp.close()

        deploy = f"""
sudo rsync -av --delete {remote_tmp}/ /var/www/onij/web/
rm -rf {remote_tmp}
sudo chown -R www-data:www-data /var/www/onij/web
sudo find /var/www/onij/web -type d -exec chmod 755 {{}} \\;
sudo find /var/www/onij/web -type f -exec chmod 644 {{}} \\;
grep -o 'index-[^\"]*\\.js' /var/www/onij/web/index.html | head -1
echo SYNC_OK
"""
        _, stdout, _ = c.exec_command(deploy, get_pty=True, timeout=300)
        out = stdout.read().decode(errors="replace")
        print(out)
        return 0 if "SYNC_OK" in out else 1
    finally:
        c.close()


if __name__ == "__main__":
    raise SystemExit(main())
