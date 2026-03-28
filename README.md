# suapapa's tools

A collection of useful Go-based command-line utilities.

## Installation

```bash
go install github.com/suapapa/tools/cmd/...@latest
```

## Available Tools

- **adball**: Run an `adb` command on all attached Android devices simultaneously.
- **at**: Network test and benchmark client/server.
- **convunzip**: Unzip files with filename encoding conversion (euc-kr to utf-8).
- **dump2png**: Convert raw grayscale memory dumps to PNG images.
- **fileserver**: A simple, quick-start HTTP file server.
- **hcc**: HangulClock Controller for communicating over serial.
- **hilog**: A log highlighter with customizable schemes.
- **hn**: Human-readable number converter (Korean, English, Binary).
- **movconcat**: Batch concatenate video files using FFmpeg.
- **mung**: A single-file HTTP server with download count limits.
- **myage**: Calculate age based on birth year.
- **myip**: Display local IP and network interface information.
- **pushover**: Send push notifications via the Pushover API.
- **qr**: Generate QR code PNG images from text.
- **regmask**: Generate and visualize bitmasks from CSV definitions.
- **since**: Calculate the human-readable duration since a given date.
- **smi2srt**: Convert SAMI (.smi) subtitles to SubRip (.srt) format.
- **tb**: A basic Tomboy note manifest viewer.
- **thumb**: Batch generate image thumbnails.
- **wdate**: A world date/timezone utility with a local database.
- **weather**: Fetch and display current weather information.

## Development

This project uses Go 1.25+.

```bash
# Build all tools
go build ./...

# Run tests
go test ./...
```

## License

MIT
