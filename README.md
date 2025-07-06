# ansiose

A terminal UI (TUI) tool for demoing, exploring, and (eventually) composing ANSI escape codes meant for the end-user.

## Features
- (Planned) Interactive Bubble Tea UI for visualizing ANSI SGR and control codes
- Supports Kitty text sizing protocol (in kitty, better detection planned)


## Usage
Run with Go:

```sh
go run ansiose.go
```

Or build a binary:

```sh
go build -o ansiose ansiose.go
./ansiose
```

- Press `q` or `Ctrl+C` to quit.
- Requires a terminal that supports ANSI codes for full effect.

## Status
- Early-stage demo and playground for ANSI/SGR codes
- UI and features are experimental and subject to change

---
