# Discord Canary Updater for Arch Linux

A simple automation tool to keep Discord Canary up to date on Arch Linux systems.

[Русский](README_ru.md) | [中文](README_zh.md)

## Problem

Arch Linux users who use the official Discord Canary version (tar.gz) face an inconvenient issue: the application may require downloading an update every time it launches. Manually downloading the archive, extracting it, and replacing files daily is tedious and cumbersome.

## Solution

This simple Go script automates the process. It performs the following actions:

- Downloads the latest discord-canary.tar.gz archive from the official Discord server.
- Extracts the archive.
- Replaces the old version with the new one.

## Installation

1. Download the latest version from [Releases](https://github.com/Lazlo105/Discord-canary-updater-for-Arch-linux/releases/latest).
2. Grant execution permissions:
   ```bash
   chmod +x discord-updater
   ```
3. Run the updater:
   ```bash
   ./discord-updater
   ```

## Requirements

- Go 1.20 or higher (for building from source)
- Make (optional, for using Makefile targets)

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/Lazlo105/Discord-canary-updater-for-Arch-linux.git
   cd Discord-canary-updater-for-Arch-linux
   ```

2. Build the project:
   ```bash
   make build
   ```

3. Run the application:
   ```bash
   ./bin/updater
   ```

## Development

- Format code: `make fmt`
- Run tests: `make test`
- Run linter: `make lint`
- View coverage: `make test-coverage`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTORS.md) before submitting a pull request.

## Security

If you discover a security vulnerability, please refer to our [Security Policy](SECURITY.md).