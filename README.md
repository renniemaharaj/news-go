# News Scraper

A lightweight Go-based news scraper that automates search, scraping, and AI-assisted summarization of online news articles. Built to collect multiple viewpoints and tag results with framework-based alignment, bias detection, and summary generation.

## Features

### AI-Augmented Article Processing

- Extracts meaningful content from paragraphs and images on real-world news pages
- Automatically summarizes article content using AI
- Categorizes topics and tags political bias for each result

### Automated Search-to-Report Pipeline

- Uses search engine queries to find news links
- Scrapes each site for paragraph text and thumbnail images
- Hydrates channels for concurrent processing using Go routines

### Modular & Concurrent Architecture

- Built with Go routines to process results in parallel
- Cleanly separated modules:

  - `config`: Loads search query configuration
  - `types`: Structured result and report data models
  - `browser`: Handles HTTP requests and HTML parsing
  - `coordinator`: Controls job flow across multiple stages
  - `log`: Rich logger with success/info/error tracking

### Framework-Aware Result Structuring

- Supports Model/System/Editor fields for content control
- Fields include:

  - Title (AI-generated)
  - Summary
  - Tags & Political Biases
  - Raw Text Content & Images

## Usage

```go
package main

import "github.com/renniemaharaj/news-go/pkg/news"

func main() {
    n := news.Instance{}

    n.CreateLogger()
    go n.GoRoutines()
    n.HydrateJobs()
}
```

Ensure a `config.json` file is present with the following structure:

```json
{
  "SearchQueries": ["2024 US election", "AI regulation", "Climate change news"]
}
```

## Output

Each query generates a `Report`:

- List of `Result` entries
- Each with title, content, summary, tags, and bias evaluation
- Ready for post-processing or user-facing display

## License

MIT License
