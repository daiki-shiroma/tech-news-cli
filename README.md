# tech-news-cli

High-performance Hacker News CLI, powered by Go.

<br/>

<a href="https://github.com/shiromadaiki/tech-news-cli/actions"><img src="https://github.com/shiromadaiki/tech-news-cli/workflows/Build/badge.svg" /></a>
<a href="https://www.npmjs.com/package/tech-news-cli"><img src="https://img.shields.io/npm/v/tech-news-cli.svg" /></a>
<a href="https://www.npmjs.com/package/tech-news-cli"><img src="https://img.shields.io/npm/dm/tech-news-cli.svg" /></a>
<a href="https://github.com/shiromadaiki/tech-news-cli/blob/master/LICENSE"><img src="https://img.shields.io/npm/l/tech-news-cli.svg" /></a>

<br/>

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
  - [Human-friendly Output](#human-friendly-output)
  - [Machine-friendly Output (JSON)](#machine-friendly-output-json)
- [Features](#features)
- [AI Integration](#ai-integration)
- [License](#license)

## Introduction

`tech-news-cli` is a high-performance command-line interface for browsing Hacker News. It is built in Go to leverage its powerful concurrency model, allowing for rapid fetching of story details and comments through parallel requests.

Whether you are a human browsing for the latest tech trends or an AI agent like Claude Code or Cursor looking for structured data, `tech-news-cli` provides the speed and flexibility you need.

## Installation

Install globally via NPM:

```bash
npm install -g tech-news-cli
```

*Note: The NPM package includes pre-compiled binaries for major platforms (macOS, Linux, Windows).*

## Usage

### Human-friendly Output

The default output is optimized for readability in your terminal.

```bash
# List the top 10 stories
tech-news list

# List 20 of the newest stories
tech-news list --type new --number 20

# View details and top comments for a specific story
tech-news view 8863
```

### Machine-friendly Output (JSON)

Use the `--json` flag to get raw data, perfect for piping into other tools or processing with AI agents.

```bash
# Get top stories as JSON
tech-news list --json

# Get details for a specific item as JSON
tech-news view 8863 --json
```

## Features

- **🚀 Blazing Fast**: Leverages Go routines for parallel data fetching.
- **🤖 AI-Native**: First-class support for JSON output, making it easy for AI agents to parse and analyze.
- **🛠 Simple & Focused**: Does one thing well—bringing Hacker News to your terminal.
- **📦 Zero Dependencies**: The compiled binary is all you need.

## AI Integration

`tech-news-cli` is designed to be the eyes and ears of your AI development assistant. You can instruct your AI to:

1.  "Fetch the top 10 stories from Hacker News in JSON."
2.  "Find articles related to 'LLMs' or 'Rust' from that list."
3.  "Fetch the comments for the most relevant article and summarize the community sentiment."

## License

MIT © [Daiki Shiroma](https://github.com/shiromadaiki)
