import json
import io
import os
import sys
import requests
from dotenv import load_dotenv

# Force UTF-8 output for Windows terminals
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

load_dotenv()
# Read input from command-line arguments
if len(sys.argv) < 2:
    print("Error: Missing input text argument", file=sys.stderr)
    sys.exit(1)

input_text = sys.argv[1].strip()

# Construct Groq API request
url = "https://api.groq.com/openai/v1/chat/completions"
headers = {
    "Content-Type": "application/json",
    "Authorization": "Bearer "+ os.getenv("GROQ_API_KEY")
}
payload = {
    "model": "meta-llama/llama-4-scout-17b-16e-instruct",
    "messages": [
        {
            "role": "user",
            "content": input_text
        }
    ]
}

# Send request to Groq API
response = requests.post(url, headers=headers, json=payload)

# Error handling
if response.status_code != 200:
    print(f"Error: {response.status_code} - {response.text}", file=sys.stderr)
    sys.exit(1)

# Parse and print the output
result = response.json()
raw_output = result["choices"][0]["message"]["content"].strip()
print(raw_output)
