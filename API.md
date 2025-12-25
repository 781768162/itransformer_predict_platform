# API 接口文档

本文档定义了新能源预测平台的 API 接口。所有接口都通过 API 网关进行访问。

## 基础 URL

所有 API 的基础 URL 为：`/api/v1`

## 认证

需要认证的接口必须在 HTTP 请求头中包含 `Authorization` 字段，其值为 `Bearer <token>`。`token` 在用户登录后获取。

---

## 1. 用户认证模块

由 API 网关直接处理。

### 1.1 用户注册

*   **Endpoint**: `POST /user/register`
*   **描述**: 创建一个新用户账户。
*   **请求体**:
    ```json
    {
      "username": "newuser",
      "password": "securepassword123"
    }
    ```
*   **成功响应 (201 Created)**:
    ```json
    {
      "status": "success",
      "message": "User registered successfully",
      "data": {
        "user_id": "uuid-for-newuser",
        "username": "newuser"
      }
    }
    ```
*   **失败响应 (400 Bad Request)**:
    ```json
    {
      "status": "error",
      "message": "Username already exists"
    }
    ```

### 1.2 用户登录

*   **Endpoint**: `POST /user/login`
*   **描述**: 用户登录以获取认证 Token。
*   **请求体**:
    ```json
    {
      "username": "newuser",
      "password": "securepassword123"
    }
    ```
*   **成功响应 (200 OK)**:
    ```json
    {
      "status": "success",
      "data": {
        "token": "your_jwt_token_here",
        "expires_in": 3600
      }
    }
    ```
*   **失败响应 (401 Unauthorized)**:
    ```json
    {
      "status": "error",
      "message": "Invalid username or password"
    }
    ```

### 1.3 用户信息

*   **Endpoint**: `GET /user/userinfo`
*   **认证**: 需要
*   **描述**: 返回当前用户信息。
*   **请求体**: 无
*   **成功响应 (200 OK)**:
    ```json
    {
      "status": "success",
      "message": "User info successfully"
    }
    ```
*   **失败响应 (401 Unauthorized)**:
    ```json
    {
      "status": "error",
      "message": "Authentication failed"
    }
    ```

---

## 2. 预测模块

这些请求将由 API 网关路由到相应的 Python 算法服务。

### 2.1 发电量预测


### 2.2 用电量预测


