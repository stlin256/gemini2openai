# GeminiAI to OpenAI API Adapter

---
[中文](#-中文) | [English](#-english)
---

## 🇨🇳 中文

本项目是一个轻量级的 Go 应用，其功能是将 OpenAI 聊天 API 格式的请求转换为 Google Gemini API 格式。它允许您使用兼容 OpenAI 的客户端来调用 Gemini 后端，包括第三方的 Gemini 服务提供商。

### ✨ 功能特性

-   将 OpenAI `/v1/chat/completions` 请求转换为 Gemini API 格式。
-   将 Gemini API 的响应转换回 OpenAI 格式。
-   通过 `config.yaml` 文件进行灵活配置。
-   通过 `.gitignore` 保护您的敏感配置信息。
-   **请求日志**：可将每个请求的详细信息（包括IP、请求体、响应体等）记录到本地JSON Lines文件中。
-   支持可选的 Bearer Token 认证。
-   基于 Go 的高性能后端，为高并发场景设计。
-   提供 `Dockerfile`，便于容器化部署。

### 🚀 快速开始

#### 环境要求

-   Go 1.22 或更高版本
-   Docker (如果需要容器化部署)

#### ⚙️ 配置

1.  **复制配置文件模板**:
    -   将 `config.yaml.example` 复制为 `config.yaml`。
    -   `.gitignore` 文件已配置好，会忽略 `config.yaml`，避免您的密钥被意外提交。
2.  **编辑 `config.yaml`**:
    ```yaml
    server:
      port: 8080

    log:
      # 是否启用请求日志记录
      enabled: true
      # 日志文件路径
      path: "requests.log"

    gemini:
      # 你的第三方Gemini API基础URL
      base_url: "your genai base_url" 
      # 你的第三方Gemini API密钥
      api_key: "your api key"

    auth:
      # 用于保护此代理服务的API密钥，客户端需要在Authorization头中提供
      # 留空 ("") 表示禁用认证
      openai_api_key: "your openai compatible api key"
    ```

#### 💻 本地运行 (Linux / macOS)

1.  安装依赖:
    ```bash
    go mod tidy
    ```
2.  运行应用:
    ```bash
    go run main.go
    ```
    服务将启动在 `config.yaml` 中指定的端口上。

#### 💻 本地运行 (Windows)

1.  **安装 Go**:
    -   访问 [Go 官方下载页面](https://golang.org/dl/) 下载并安装适用于 Windows 的 Go。
    -   安装过程中，请确保勾选 "Add Go to your PATH" 选项。
2.  **打开终端**:
    -   你可以使用 `CMD` 或者 `PowerShell`。
3.  **下载项目代码**:
    -   如果你安装了 Git，可以使用 `git clone`。
    -   或者，直接下载项目的 ZIP 压缩包并解压。
4.  **编译并运行**:
    -   在项目的根目录下，打开终端并运行以下命令来编译一个 `.exe` 可执行文件：
      ```powershell
      go build -o gemini2openai.exe .
      ```
    -   编译成功后，直接运行该程序：
      ```powershell
      .\gemini2openai.exe
      ```
    服务将启动在 `config.yaml` 中指定的端口上。

#### 🐳 使用 Docker 运行

1.  构建 Docker 镜像:
    ```bash
    docker build -t gemini2openai .
    ```
2.  运行 Docker 容器:
    ```bash
    # Linux / macOS
    docker run -p 8080:8080 -v $(pwd)/config.yaml:/root/config.yaml gemini2openai

    # Windows (CMD)
    docker run -p 8080:8080 -v "%cd%\config.yaml:/root/config.yaml" gemini2openai
    ```
    此命令会将主机的 8080 端口映射到容器的 8080 端口，并将本地的 `config.yaml` 挂载到容器中。

### 🔌 API 调用示例

向 `http://localhost:8080/v1/chat/completions` 发送一个标准的 OpenAI 聊天 API 格式的 POST 请求。

**使用 cURL 的示例:**

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer you_api_key" \
  -d '{
    "model": "gemini-2.5-flash",
    "messages": [
      {
        "role": "user",
        "content": "你好，你是谁？"
      }
    ]
  }'
```

---

## 🇬🇧 English

This project is a lightweight Go application that acts as an adapter to convert API requests from the OpenAI Chat Completions format to the Google Gemini format. It allows you to use OpenAI-compatible clients with a Gemini backend, including third-party Gemini providers.

### ✨ Features

-   Converts OpenAI `/v1/chat/completions` requests to Gemini API format.
-   Converts Gemini API responses back to OpenAI format.
-   Flexible configuration via a `config.yaml` file.
-   Protects your sensitive configuration with `.gitignore`.
-   **Request Logging**: Can log detailed information for each request (IP, request body, response body, etc.) to a local JSON Lines file.
-   Supports optional Bearer Token authentication.
-   High-performance Go backend, designed for concurrent requests.
-   Provides a `Dockerfile` for easy containerization.

### 🚀 Getting Started

#### Prerequisites

-   Go 1.22 or later
-   Docker (for containerized deployment)

#### ⚙️ Configuration

1.  **Copy the config template**:
    -   Copy `config.yaml.example` to `config.yaml`.
    -   The `.gitignore` file is already configured to ignore `config.yaml`, preventing your keys from being accidentally committed.
2.  **Edit `config.yaml`**:
    ```yaml
    server:
      port: 8080

    log:
      # Enable or disable request logging
      enabled: true
      # Path to the log file
      path: "requests.log"

    gemini:
      # Your third-party Gemini API base URL
      base_url: "your genai base_url" 
      # Your third-party Gemini API key
      api_key: "your api key"

    auth:
      # API key to protect this proxy service. Clients must provide it in the Authorization header.
      # Leave it empty ("") to disable authentication.
      openai_api_key: "your openai compatible api key"
    ```

#### 💻 Running Locally (Linux / macOS)

1.  Install dependencies:
    ```bash
    go mod tidy
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on the port specified in your `config.yaml`.

#### 💻 Running Locally (Windows)

1.  **Install Go**:
    -   Visit the [official Go downloads page](https://golang.org/dl/) to download and install Go for Windows.
    -   During installation, ensure the "Add Go to your PATH" option is checked.
2.  **Open a terminal**:
    -   You can use `CMD` or `PowerShell`.
3.  **Get the project code**:
    -   If you have Git, use `git clone`.
    -   Alternatively, download the project as a ZIP file and extract it.
4.  **Build and run**:
    -   In the project's root directory, open a terminal and run the following command to build an `.exe` executable:
      ```powershell
      go build -o gemini2openai.exe .
      ```
    -   Once built, run the program:
      ```powershell
      .\gemini2openai.exe
      ```
    The server will start on the port specified in `config.yaml`.

#### 🐳 Running with Docker

1.  Build the Docker image:
    ```bash
    docker build -t gemini2openai .
    ```
2.  Run the Docker container:
    ```bash
    # Linux / macOS
    docker run -p 8080:8080 -v $(pwd)/config.yaml:/root/config.yaml gemini2openai

    # Windows (CMD)
    docker run -p 8080:8080 -v "%cd%\config.yaml:/root/config.yaml" gemini2openai
    ```
    This command maps port 8080 on your host to port 8080 in the container and mounts your local `config.yaml` into it.

### 🔌 API Usage Example

Send a standard OpenAI Chat Completions POST request to `http://localhost:8080/v1/chat/completions`.

**Example with cURL:**

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer you_api_key" \
  -d '{
    "model": "gemini-2.5-flash",
    "messages": [
      {
        "role": "user",
        "content": "Hello, who are you?"
      }
    ]
  }'