import asyncio
import aiohttp
import time
import json

# --- Configuration ---
URL = "http://localhost:8080/v1/chat/completions"
# IMPORTANT: Replace "YOUR_OPENAI_COMPATIBLE_API_KEY" with the key from your config.yaml
API_KEY = "YOUR_OPENAI_COMPATIBLE_API_KEY"
CONCURRENT_REQUESTS = 20  # Number of requests to send at the same time
# ---------------------

HEADERS = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {API_KEY}"
}

DATA = {
    "model": "gemini-2.5-flash",
    "messages": [
        {"role": "user", "content": "Tell me a short story about a robot."}
    ]
}

async def send_request(session, request_num):
    """Sends a single request and prints the result."""
    start_time = time.time()
    try:
        async with session.post(URL, headers=HEADERS, data=json.dumps(DATA)) as response:
            duration = time.time() - start_time
            if response.status == 200:
                # response_json = await response.json() # Uncomment to see full response
                print(f"Request {request_num:2d}: Success ({response.status}) in {duration:.2f}s")
                return True
            else:
                response_text = await response.text()
                print(f"Request {request_num:2d}: Failed  ({response.status}) in {duration:.2f}s - {response_text}")
                return False
    except Exception as e:
        duration = time.time() - start_time
        print(f"Request {request_num:2d}: Error in {duration:.2f}s - {e}")
        return False

async def main():
    """Runs the concurrency test."""
    print(f"Starting concurrency test with {CONCURRENT_REQUESTS} parallel requests...")
    start_total_time = time.time()
    
    async with aiohttp.ClientSession() as session:
        tasks = [send_request(session, i + 1) for i in range(CONCURRENT_REQUESTS)]
        results = await asyncio.gather(*tasks)
    
    total_duration = time.time() - start_total_time
    success_count = sum(1 for r in results if r)
    
    print("\n--- Test Summary ---")
    print(f"Total requests: {CONCURRENT_REQUESTS}")
    print(f"Successful:     {success_count}")
    print(f"Failed:         {CONCURRENT_REQUESTS - success_count}")
    print(f"Total time:     {total_duration:.2f}s")
    if total_duration > 0:
        print(f"Requests/sec:   {CONCURRENT_REQUESTS / total_duration:.2f}")
    print("--------------------")

if __name__ == "__main__":
    asyncio.run(main())