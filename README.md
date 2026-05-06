# NASA Name Generator Bot

A Telegram bot that generates stylized name images using letter assets with 3 random variants per letter.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Project Structure

```
nasa/
├── cmd/bot/main.go           # Telegram bot entry point
├── internal/nasa/
│   ├── image.go              # Image generation with caching
│   └── image_test.go         # Unit tests
├── Images/                   # Letter assets (a1.jpg - z3.jpg)
├── .env                      # Bot token (not committed)
├── .gitignore
├── go.mod
└── README.md
```

## Setup

1. Create a `.env` file in the project root:
   ```
   BOT_TOKEN=your_telegram_bot_token
   ```

2. Ensure the `Images/` folder contains all letter variants:
   - `a1.jpg`, `a2.jpg`, `a3.jpg` through `z1.jpg`, `z2.jpg`, `z3.jpg`
   - 3 variants per letter for random selection

## Running

Build and run:
```bash
GOCACHE=/tmp/go-build go build -o nasa-bot ./cmd/bot
./nasa-bot
```

Or run directly:
```bash
GOCACHE=/tmp/go-build go run ./cmd/bot
```

## Bot Commands

- `/start` - Show intro message
- `/help` - Show usage instructions
- `/nasa <name>` - Generate image for a name
- Any text - Generate image from the text

## How It Works

1. Bot receives a name or text message
2. For each letter, randomly selects one of 3 variants (e.g., for 'k' picks k1, k2, or k3)
3. Combines all letter images horizontally
4. Caches decoded images in memory for fast subsequent requests
5. Sends the generated image back to the Telegram chat

## Features

- Letter images cached in memory at startup for fast response
- Random variant selection creates unique images each time
- Pre-warms cache on initialization
- Clean error handling with user-friendly messages