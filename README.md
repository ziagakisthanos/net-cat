# Net-Cat

A simple TCP-based chat application written in Go, originally known as `TCP-Chat`. Net-Cat allows multiple clients to connect to a server, exchange public and private messages, and change nicknames in real-time.

---

## Features

- **Public Chat**: Broadcast messages to all connected users.
- **Private Messaging**: Whisper (`-w`/`-whisper`) to a specific user.
- **Nickname Changes**: Change your display name on the fly with `-n`/`-name`.
- **Command Help**: Built-in help manual via `-h`/`-help`.
- **Message History**: New clients receive the recent chat history on join.
- **Server Logs**: All messages are logged to a file per server port.

---

## Prerequisites

- Go 1.18+ installed ([https://golang.org/dl/](https://golang.org/dl/))
- Git (optional, for cloning the repo)

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/net-cat.git
   cd net-cat
   ```
2. Initialize modules and download dependencies:
   ```bash
   go mod tidy
   ```
3. Build the server binary:
   ```bash
   go build -o net-cat ./cmd/net-cat
   ```

---

## Usage

```
# Run with default port (8989)
./net-cat

# Or specify a custom port
./net-cat 9000
```

By default, the server listens on `:<port>`. Clients can connect using any TCP client, e.g.:

```bash
nc localhost 8989
```

Once connected, follow on-screen prompts to enter your nickname and start chatting.

---

## Commands

| Command           | Description                                  |
| ----------------- | -------------------------------------------- |
| `<message>`       | Send a public message to everyone            |
| `-w <user> <msg>` | Send a private message (whisper) to `<user>` |
| `-n <newname>`    | Change your nickname                         |
| `-name <newname>` | Alternative syntax for name change           |
| `-h` / `-help`    | Show help manual                             |

---

## Project Structure

```
├── cmd/
│   └── net-cat/          # Main application entrypoint (main.go)
├── pkg/chat/            # Core chat server implementation
│   ├── client.go
│   ├── consts.go
│   ├── help.go
│   ├── history.go
│   ├── logger.go
│   ├── namechange.go
│   ├── server.go
│   ├── structs.go
│   ├── validation.go
│   └── whisper.go
├── go.mod
└── README.md
```

---

## Contributing

Contributions are welcome! Please fork the repository and open a pull request with your changes.

---

## License

Copyright © 2025 Maria Tzemanaki and Athanasios Ziagakis.
