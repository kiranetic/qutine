# qutine

A secure, minimal container runtime that runs containers ephemerally with encrypted filesystems.

## Features
- User-specific, rootless installation
- Password-protected commands
- Encrypted in-memory filesystem (tmpfs)
- Minimal subcommands: `run`, `stop`, `list`
- Secure, non-exportable image format

## Installation
```bash
./install.sh

## Usage
qtn run <encrypted-image> <command>
