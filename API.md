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

*   **路径**: `POST /user/register`
*   **描述**: 创建一个新用户账户。
*   **请求体**:
    ```json
    {
      "user_name": string,
      "password": string
    }
    ```
*   **响应体**:
    ```json
    {
      "message": string
    }
    ```

### 1.2 用户登录

*   **路径**: `POST /user/login`
*   **描述**: 用户登录以获取认证 Token。
*   **请求体**:
    ```json
    {
      "user_name": string,
      "password": string
    }
    ```
*   **响应体**:
    ```json
    {
      "token": string,
      "expire_time": string,
      "message": string
    }
    ```

---

## 2. 预测模块

这些请求将由 API 网关路由到相应的 Python 算法服务。

### 2.1 投递任务

*   **路径**: `POST /service/create_task`
*   **描述**: 投递预测任务，header需要token。
*   **请求体（multipart/form-data）**:
    * Content-Type: `multipart/form-data`
    * 字段：
      * `file`：CSV 文件
      * `date`：预测开始时间点, string
*   **响应体**:
    ```json
    {
      "task_id": int,
      "message": string
    }
    ```

---

### 2.2 查询任务

*   **路径**: `POST /service/get_task`
*   **描述**: 查询任务状态，header需要token。
*   **请求体**:
    ```json
    {
      "task_id": int
    }
    ```
*   **响应体**:
    ```json
    {
      "message": string,
      "status": string,
      "result": [float, float, ..., float] // 长度 24 的数组，形如 [v0, v1, ... v23]
    }
    ```

---


