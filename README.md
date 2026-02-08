### Claude code
A minimal claude code like cli tool.

### Setup .env
```
OPENROUTER_API_KEY=
BASE_URL=https://api.moonshot.ai/v1
MODEL=kimi-k2-0905-preview
```

### Example

```
➜  claude-code git:(main) ✗ sh run.sh

╔════════════════════════════════════════╗
║      Claude Code CLI - Go Edition      ║
╚════════════════════════════════════════╝

Type your message and press Enter.
Press Ctrl+C to exit.

> hello
Hi there! How can I help you today?
> what is the last commit message 
⟳ Calling tool: Bash
The last commit message is **"add cli styling"** (commit `181497f`).
> 
```
