# termtrace
A reproducible terminal workflow recorder for debugging and knowledge sharing.

termtrace captures commands, outputs, and execution context so sessions can be replayed step by step. It provides a deterministic, machine-readable trace of terminal workflows, sitting between shell history and screen recording.

## Why
Debugging and operational workflows are hard to reproduce because they live in shell history and human memory. termtrace turns terminal activity into a replayable artifact.

## Features (MVP)
- Record terminal sessions
- Replay sessions step by step
- Capture commands, outputs, timestamps, and exit codes

## Usage

```shell
termtrace record
termtrace stop
termtrace replay session.wf
```
## Status
Early development. APIs and file formats may change.

## License
MIT
