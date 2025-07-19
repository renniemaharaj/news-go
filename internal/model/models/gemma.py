# To run this code you need to install the following dependencies:
# pip install google-genai

import base64
import os
import io
import sys
from google import genai
from google.genai import types
# from dotenv import load_dotenv

# Load environment variables from .env file
# load_dotenv()

# Force UTF-8 output for Windows terminals
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# Read input from command-line arguments
if len(sys.argv) < 2:
    print("Error: Missing API Key", file=sys.stderr)
    sys.exit(1)
API_KEY = sys.argv[1].strip()

# Read input from command-line arguments
if len(sys.argv) < 3:
    print("Error: Missing prompt", file=sys.stderr)
    sys.exit(1)
PROMPT = sys.argv[2].strip()

print(os.getenv("GEMINI_API_KEY"))
def generate(prompt):
    client = genai.Client(
        api_key=API_KEY,
    )

    model = "gemma-3-4b-it"
    contents = [
        types.Content(
            role="user",
            parts=[
                types.Part.from_text(text=PROMPT),
            ],
        ),
    ]
    generate_content_config = types.GenerateContentConfig(
        response_mime_type="text/plain",
    )

    for chunk in client.models.generate_content_stream(
        model=model,
        contents=contents,
        config=generate_content_config,
    ):
        print(chunk.text, end="")

if __name__ == "__main__":
    generate(PROMPT)
