#
#
# llama_local.py - Run LLaMA locally using llama-cpp-python bindings
#
# This script replaces remote LLM APIs like Groq with a locally loaded LLaMA model.
# You must install the model, build dependencies, and the llama-cpp-python package.
#   - pip install llama-cpp-python
#
# -------------
# SETUP INSTRUCTIONS (Windows)
# -------------
#
# 1. Install Python 3.10+ (preferably from https://www.python.org)
#
# 2. Install C++ Build Tools
#       - Required to compile the llama-cpp backend. You only need to do this once.
#
#   a. Download and install Visual Studio Build Tools:
#      https://visualstudio.microsoft.com/visual-cpp-build-tools/
#
#   b. During installation, select:
#       - "Desktop development with C++"
#       - Make sure the following are checked:
#           - MSVC v143 or newer
#           - Windows 10 or 11 SDK
#           - CMake tools for Windows
#           - C++ CMake tools for Windows
#
# 3. (Optional but recommended) Install Windows Developer Tools via terminal:
#
#
#   winget install --id Microsoft.WindowsTerminal -e
#   winget install Microsoft.VisualStudio.2022.BuildTools

# Finally, find and download models from https://huggingface.co/
#    

import sys
import io
from llama_cpp import Llama

# Force UTF-8 output
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# Check for input argument
if len(sys.argv) < 2:
    print("Error: Missing input text argument", file=sys.stderr)
    sys.exit(1)

input_text = sys.argv[1].strip() # Pass your prompt as arguments

# Load your local model (adjust path to your model file)
model_path = "models/llama-2-7b-chat.Q4_K_M.gguf"  # Change this to your GGUF file

llm = Llama(
    model_path=model_path,
    n_ctx=2048, # Set based on requirements and model capability
    n_threads=8,     # Set based on your CPU
    verbose=False
)

# Create the prompt
messages = [{"role": "user", "content": input_text}]

# Call the model
response = llm.create_chat_completion(
    messages=messages,
    temperature=0.7,
    top_p=0.95,
    max_tokens=1024,
    stop=["</s>", "<|endoftext|>"]
)

# Print result
print(response['choices'][0]['message']['content'].strip())
