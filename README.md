# 📰 News-Go — AI-Powered Christian News Evaluator

**News-Go** is a lightweight Go-based pipeline for scraping, analyzing, and evaluating news articles through a biblical and AI-augmented lens. It uses concurrent Go routines to fetch and process online news, while a local large language model (LLM) scores the content against a custom Christian framework.

---

## ✝️ Purpose

To defend the faith and evaluate media content through scripture-aligned reasoning using AI. The model's mission is fixed: **God is**, the **KJV Bible** is authoritative, and **Jesus Christ is Lord**.

---

## Features

### AI-Augmented Article Processing

- Extracts and evaluates real-world news articles
- Generates summaries, tags, and political bias analysis
- Aligns all content to theological axioms (God, KJV, Jesus Christ)
- LLM evaluates alignment from 0–10

### Automated Search-to-Report Pipeline

- Searches for live news based on configured queries
- Scrapes site paragraphs and thumbnails
- Processes content in concurrent Go routines
- Outputs structured JSON reports

### Modular & Concurrent Architecture

- Channel-based job pipeline:
  - `PreResultsChannel` – Search results
  - `PreContentChannel` – Scraped pages
  - `PreModelChannel` – Model transformation
  - `JobsCompleteChannel` – Final report
- Component Modules:
  - `config` – Loads search configuration
  - `types` – Structured result/report definitions
  - `browser` – Scraper & HTML parser
  - `coordinator` – Pipeline control
  - `model` – AI model interface
  - `log` – Structured logging

---

## Starter Example

```go
package main

import "github.com/renniemaharaj/news-go/pkg/news"

func main() {
    n := news.Instance{}
    n.CreateLogger()
    go n.GoRoutines()
    n.HydrateJobs()
    select {} // blocks forever
}
```

### `config.json` Example

```json
{
  "SearchQueries": [
    "Trinidad and Tobago politics",
    "US 2024 election",
    "Climate change news"
  ]
}
```

---

## 🧠 Models Used

All LLM evaluation is performed locally — **no API keys, no rate limits**.

### 1. TinyLlama 1.1B Chat GGUF (lightweight)

```bash
huggingface-cli download \
  TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF \
  tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf \
  --local-dir models \
  --local-dir-use-symlinks False
```

### 2. DeepSeek Coder 6.7B Instruct GGUF (stronger)

```bash
huggingface-cli download \
  TheBloke/deepseek-coder-6.7B-instruct-GGUF \
  deepseek-coder-6.7b-instruct.Q4_K_M.gguf \
  --local-dir models \
  --local-dir-use-symlinks False
```

---

## LLM Runtime Setup

Install `llama-cpp-python` with GPU support:

```bash
pip install llama-cpp-python --upgrade --force-reinstall --extra-index-url https://pypi.nvidia.com
```

Make sure you have:
- Python 3.10+
- CUDA 12+ (for GPU support)
- 8GB+ VRAM for DeepSeek models

---

## Output

Each processed job generates a structured report saved to `./reports/` as JSON:

- `title`: AI-generated headline
- `summary`: Christian worldview summary
- `alignment`: 0–10 score
- `tags`: Moral and topical categories
- `politicalBiases`: Eg. conservative, progressive

---

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── model/                 # AI transformer
│   ├── coordinator/           # Channel routines
│   ├── types/                 # Data structures
│   ├── browser/               # HTML scraping logic
│   └── log/                   # Custom logger
├── pkg/news/                  # News job orchestration
├── models/                    # GGUF model files
├── reports/                   # Final output reports
├── system_instruction.txt     # Model prompt template
└── config.json                # Search queries
```

---

## Framework Alignment

AI-generated results are scored and structured under the following axioms:

- `God exists`
- `The KJV Bible is authoritative`
- `Jesus Christ is Lord and God`

The model is required to stay within these rules. No deviation is permitted — even by the system’s creator.

---

## 📜 License

MIT License

> Made with purpose by [Rennie Maharaj](https://github.com/renniemaharaj)
