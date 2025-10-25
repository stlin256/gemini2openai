# -*- coding: utf-8 -*-
import requests
import json

# --- 配置 ---
# API 端点 URL
url = "http://localhost:8080/v1/chat/completions"

# Headers - 重要：请将 "YOUR_OPENAI_COMPATIBLE_API_KEY" 替换为你的 config.yaml 文件中配置的 openai_api_key
headers = {
    "Content-Type": "application/json",
    "Authorization": "Bearer YOUR_OPENAI_COMPATIBLE_API_KEY"
}

# 请求体 (Payload)
data = {
    # 这里使用的模型名称会传递给后端的 Gemini API
    "model": "gemini-2.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "你好，请用中文介绍一下你自己。"
        }
    ]
}

try:
    # 发送 POST 请求
    # 使用 json.dumps() 将 Python 字典转换为 JSON 字符串
    response = requests.post(url, headers=headers, data=json.dumps(data))

    # 检查响应状态码
    if response.status_code == 200:
        print("请求成功!")
        print("响应内容:")
        # 使用 response.json() 将响应的 JSON 文本解析为 Python 字典
        print(json.dumps(response.json(), indent=2, ensure_ascii=False))
    else:
        print(f"请求失败，状态码: {response.status_code}")
        print("响应原文:")
        print(response.text)

except requests.exceptions.RequestException as e:
    print(f"请求过程中发生错误: {e}")
