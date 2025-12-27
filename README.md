# GameMeow 🎮

A Game Boy emulator written in Go, built for personal learning of emulation development.

## Overview

GameMeow is a personal learning project to understand how emulation works by implementing the core components of a Game Boy system in Go. This project covers CPU instruction handling, memory management, graphics rendering, input processing, and sound synthesis.

## Project Structure

- **cpu/** - CPU implementation including instruction sets (arithmetic, load operations, registers)
- **memory/** - Memory management unit (MMU) for address translation and memory access
- **graphics/** - Graphics rendering system
- **input/** - Input handling for game controls
- **sound/** - Sound synthesis and audio processing
- **cartridge/** - Game cartridge loading and management
- **emulator/** - Main emulator orchestration

## Getting Started

### Prerequisites

- Go 1.13 or later

### Building

```bash
go build
```

### Running

```bash
go run .
```

## Learning Goals

This project is a deep dive into:
- CPU instruction decoding and execution
- Memory-mapped I/O and address translation
- Real-time audio and graphics synchronization
- Game Boy hardware architecture and how it all fits together

## License

Educational use

---

**GameMeow** - Learn emulation development with Go! 🐱
