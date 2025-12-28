import json
import sys
import requests

BASE_URL = "http://localhost:8080/api/v1"
LOGIN_PATH = "/user/login"
REGISTER_PATH = "/user/register"

# 配置测试账号
USERNAME = "test_user"
PASSWORD = "test_password"


def register(username: str, password: str) -> None:
    url = f"{BASE_URL}{REGISTER_PATH}"
    payload = {"user_name": username, "password": password}
    resp = requests.post(url, json=payload, timeout=10)
    print(f"[register] HTTP {resp.status_code}")
    try:
        data = resp.json()
    except Exception:
        print("响应非 JSON:", resp.text[:500])
        return
    print(json.dumps(data, ensure_ascii=False, indent=2))


def login(username: str, password: str) -> None:
    url = f"{BASE_URL}{LOGIN_PATH}"
    payload = {"user_name": username, "password": password}
    resp = requests.post(url, json=payload, timeout=10)
    print(f"HTTP {resp.status_code}")
    try:
        data = resp.json()
    except Exception:
        print("响应非 JSON:", resp.text[:500])
        return
    print("响应 JSON:")
    print(json.dumps(data, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    u = sys.argv[1] if len(sys.argv) > 1 else USERNAME
    p = sys.argv[2] if len(sys.argv) > 2 else PASSWORD
    # 先注册再登录，重复注册时可能返回已存在的错误，忽略即可
    register(u, p)
    login(u, p)
