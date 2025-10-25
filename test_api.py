import requests
import json

# API endpoint URL
url = "http://localhost:8080/v1/chat/completions"

# Headers
# IMPORTANT: Replace "YOUR_OPENAI_COMPATIBLE_API_KEY" with the key from your config.yaml
headers = {
    "Content-Type": "application/json",
    "Authorization": "Bearer YOUR_OPENAI_COMPATIBLE_API_KEY"
}

# Request payload
data = {
    "model": "gemini-2.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "Hello, who are you?"
        }
    ]
}

try:
    # Make the POST request
    response = requests.post(url, headers=headers, data=json.dumps(data))

    # Check the response
    if response.status_code == 200:
        print("Request successful!")
        print("Response:")
        print(response.json())
    else:
        print(f"Request failed with status code: {response.status_code}")
        print("Response:")
        print(response.text)

except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")
