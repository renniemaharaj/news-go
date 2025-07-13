import json
import io
import sys
from pathlib import Path
from llama_cpp import Llama

# Force UTF-8 output for Windows terminals
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# Load system instruction template
with open("system_instruction.txt", "r", encoding="utf-8") as f:
    prompt_template = f.read()

# Read input from command-line arguments
if len(sys.argv) < 2:
    print("Error: Missing input text argument", file=sys.stderr)
    sys.exit(1)

input_text = sys.argv[1].strip()

# Inject into prompt template
prompt = prompt_template.replace("promptGoesHere", input_text)

# Load the model
llm = Llama(
    model_path="./models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
    n_ctx=2048,
    n_threads=8,
    verbose=False,
)

# Generate output
response = llm(prompt, max_tokens=1000)
raw_output = response["choices"][0]["text"].strip()

# Return raw output for Go to parse
print(raw_output)
