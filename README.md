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
```

## Usage
```console
qtn run <encrypted-image> <command>
```

## License
This project is currently unlicensed. A license will be added in the future.
