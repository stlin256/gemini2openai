import openai

# Configure the OpenAI client to point to your local server
client = openai.OpenAI(
    base_url="your base url",
    api_key="your api key"  # This key should match the one in your config.yaml
)

def test_streaming_chat():
    """
    Interactive chat session with streaming responses.
    """
    print("Starting interactive chat session (press Ctrl+C to exit).")
    while True:
        try:
            # Get user input for the question
            user_prompt = input("\nEnter your question: ")
            if not user_prompt:
                print("Question cannot be empty.")
                continue

            # Get user input for the model, with a default value
            model_name = input("Enter model ID (default: gemini-2.5-pro): ")
            if not model_name:
                model_name = "gemini-2.5-pro"

            print(f"\n--- Assistant's Response (Model: {model_name}) ---")

            # Create a chat completion request with streaming enabled
            stream = client.chat.completions.create(
                model=model_name,
                messages=[{"role": "user", "content": user_prompt}],
                stream=True,
            )

            # Print the response chunks as they arrive
            for chunk in stream:
                content = chunk.choices[0].delta.content
                if content:
                    print(content, end='', flush=True)
            
            print("\n--------------------------------------------------")

        except KeyboardInterrupt:
            print("\nExiting chat session.")
            break
        except Exception as e:
            print(f"\nAn error occurred: {e}")
            break

if __name__ == "__main__":
    test_streaming_chat()