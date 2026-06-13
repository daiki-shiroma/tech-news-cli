# Tech News CLI Instructions

You can use this high-performance Go-based CLI to fetch tech news (Hacker News) and show them to the user.

## Available Commands

- `./tech-news-cli list`: Lists stories.
  - Options:
    - `-n <number>`: Number of stories (default: 10).
    - `-t <type>`: Type (`top`, `new`, `best`).
    - `--json`: Output as raw JSON.
- `./tech-news-cli view <id>`: Shows details and comments.
  - Options:
    - `--json`: Output as raw JSON.

## Performance Note
This implementation uses Go routines to fetch story details in parallel, making it significantly faster than sequential approaches.

## Workflow for AI Agents
1. Execute `./tech-news-cli list -n 20 --json` to get a broad overview.
2. Filter the JSON locally to find relevant topics.
3. Use `./tech-news-cli view <id> --json` to deep-dive into interesting stories or comments.
