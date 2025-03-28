# Gonyx Framework

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Blocktunium/gonyx)](https://github.com/Blocktunium/gonyx/releases/latest)
[![License](https://img.shields.io/github/license/Blocktunium/gonyx)](LICENSE)
[![Website](https://img.shields.io/badge/website-gonyx.io-blue)](https://gonyx.io)

Gonyx is a modern, high-performance framework designed to streamline application development with a focus on simplicity and efficiency. It provides a robust foundation for building scalable applications with built-in support for modern development practices.

Visit our official website at [https://gonyx.io](https://gonyx.io) for more information.

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://golang.org/dl/) version 1.23 or higher

To verify your Go installation:
```bash
go version
```

## Quick Start

### Download and Installation

Visit our download page at [https://gonyx.io/download](https://gonyx.io/download) for the latest versions and installation instructions.

Alternatively, you can download the binaries directly for your platform:

- [Linux (x64)](https://github.com/Blocktunium/gonyx/releases/download/v0.3.0/gonyx_linux_amd64.zip)
- [Linux (arm64)](https://github.com/Blocktunium/gonyx/releases/download/v0.3.0/gonyx_linux_arm64.zip)
- [macOS (x64)](https://github.com/Blocktunium/gonyx/releases/download/v0.3.0/gonyx_macos_amd64.zip)
- [macOS (arm64)](https://github.com/Blocktunium/gonyx/releases/download/v0.3.0/gonyx_macos_arm64.zip)
- [Windows (x64)](https://github.com/Blocktunium/gonyx/releases/download/v0.3.0/gonyx_windows_amd64.zip)

After downloading, extract the archive and use the binary inside to create new projects:

```bash
# Linux/macOS
unzip gonyx_<platform>_x64.zip
chmod +x gonyx
./gonyx init hello_world --path .

# Windows
# Extract the zip file and run:
gonyx.exe init hello_world --path .
```

## Features

- Fast and efficient core architecture
- Built-in development tools
- Cross-platform support
- Modern development workflow
- Extensible plugin system

## Contrib Packages

The following contributed packages extend Gonyx's functionality:

- **gormkit**: A wrapper around GORM package to handle rational databases (sqlite, mysql, postgresql)
- **mongokit**: A wrapper around the MongoDB driver package

## Documentation

Comprehensive documentation is available at [https://gonyx.io/docs/0.3.0/intro](https://gonyx.io/docs/0.3.0/intro).

## Contributing

We welcome contributions from the community! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contributors

<a href="https://github.com/Blocktunium/gonyx/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Blocktunium/gonyx" />
</a>

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

