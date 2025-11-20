# Subtitle Generator

A web application that enables you to upload a video and automatically generate subtitles (transcription) using AI-powered models (OpenAI Whisper or local Whisper CLI).

## Features

- **Simple web UI** to upload a video file and generate subtitles
- **HTMX-powered interactions** for async UI updates (upload progress, transcript display)
- **Automatic audio extraction** from video using ffmpeg
- **Dual transcription options**: Local Whisper CLI or OpenAI API
- **Real-time transcript display** in the browser
- **Minimal dependencies** - built with Go standard library and HTMX

## How It Works

1. **User uploads video file** via the web interface (`/upload`)
2. **Backend saves the file** in `static/uploads/` and renders a video player
3. **User clicks "Generate Subtitles"**, sending a POST request to `/transcribe`
4. **Audio is extracted** from the uploaded video using ffmpeg
5. **Audio file is transcribed** using Whisper (local CLI or OpenAI API)
6. **Transcript is displayed** in the browser via HTMX

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Browser (Client)                          │
│                                                                   │
│  ┌──────────────┐        ┌──────────────┐    ┌──────────────┐  │
│  │  Upload Form │───────▶│ Video Player │───▶│  Transcript   │  │
│  │   (HTMX)     │        │   (HTMX)     │    │    (HTMX)     │  │
│  └──────────────┘        └──────────────┘    └──────────────┘  │
│         │                        │                    ▲          │
└─────────┼────────────────────────┼────────────────────┼─────────┘
          │                        │                    │
          │ POST /upload           │ POST /transcribe   │ HTML Response
          │                        │                    │
          ▼                        ▼                    │
┌─────────────────────────────────────────────────────────────────┐
│                      Go Web Server (main.go)                     │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                      HTTP Handlers                        │   │
│  │  ┌────────────┐  ┌────────────┐  ┌──────────────────┐   │   │
│  │  │   Home     │  │   Upload   │  │   Transcribe     │   │   │
│  │  │  Handler   │  │  Handler   │  │    Handler       │   │   │
│  │  └────────────┘  └────────────┘  └──────────────────┘   │   │
│  │                                            │               │   │
│  └────────────────────────────────────────────┼──────────────┘   │
│                                               │                   │
│  ┌────────────────────────────────────────────┼──────────────┐   │
│  │                     Services               │               │   │
│  │  ┌──────────────────┐    ┌─────────────────▼──────────┐  │   │
│  │  │  Audio Service   │    │  Transcription Service     │  │   │
│  │  │  (audio.go)      │◀───│  - Local Whisper           │  │   │
│  │  │  - ffmpeg        │    │  - OpenAI API              │  │   │
│  │  └──────────────────┘    └────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
└───────────────────────────────────────────────────────────────┬─┘
                                                                 │
                                                                 ▼
                                         ┌────────────────────────────┐
                                         │    File System             │
                                         │  - static/uploads/         │
                                         │  - templates/              │
                                         └────────────────────────────┘
```

## Sequence Diagram

```
User          Browser         Server          Audio Service    Whisper
 │               │               │                   │             │
 │ Select Video  │               │                   │             │
 ├──────────────▶│               │                   │             │
 │               │               │                   │             │
 │ Click Upload  │               │                   │             │
 ├──────────────▶│               │                   │             │
 │               │ POST /upload  │                   │             │
 │               ├──────────────▶│                   │             │
 │               │               │ Save video        │             │
 │               │               │ to static/uploads/│             │
 │               │               ├──────────┐        │             │
 │               │               │          │        │             │
 │               │               │◀─────────┘        │             │
 │               │ HTML (player) │                   │             │
 │               │◀──────────────┤                   │             │
 │◀──────────────┤               │                   │             │
 │               │               │                   │             │
 │ View Video    │               │                   │             │
 │               │               │                   │             │
 │ Generate Subs │               │                   │             │
 ├──────────────▶│               │                   │             │
 │               │POST /transcribe│                  │             │
 │               ├──────────────▶│                   │             │
 │               │               │ Extract Audio     │             │
 │               │               ├──────────────────▶│             │
 │               │               │                   │ ffmpeg      │
 │               │               │                   ├────────┐    │
 │               │               │                   │        │    │
 │               │               │                   │◀───────┘    │
 │               │               │ audio.mp3         │             │
 │               │               │◀──────────────────┤             │
 │               │               │                   │             │
 │               │               │ Transcribe Audio  │             │
 │               │               ├─────────────────────────────────▶│
 │               │               │                   │  whisper CLI│
 │               │               │                   │  or OpenAI  │
 │               │               │                   │             │
 │               │               │ Transcript Text   │             │
 │               │               │◀─────────────────────────────────┤
 │               │               │                   │             │
 │               │ HTML (transcript)                 │             │
 │               │◀──────────────┤                   │             │
 │◀──────────────┤               │                   │             │
 │               │               │                   │             │
 │ Read Transcript                                   │             │
 │               │               │                   │             │
```

## Key Components

### Project Structure

```
subtitle-gen/
├── main.go                 # Entry point, HTTP server setup, route definitions
├── go.mod                  # Go module definition
├── handlers/               # HTTP request handlers
│   ├── home.go            # Renders the main page
│   ├── upload.go          # Handles video file uploads
│   └── transcribe.go      # Coordinates audio extraction and transcription
├── services/               # Business logic services
│   ├── audio.go           # Audio extraction using ffmpeg
│   ├── local_whisper.go   # Local Whisper CLI integration
│   └── openai.go          # OpenAI Whisper API integration
├── templates/              # HTML templates
│   ├── layout.html        # Base layout template
│   ├── index.html         # Main upload page
│   ├── player.html        # Video player fragment (HTMX response)
│   └── transcript.html    # Transcript display fragment (HTMX response)
├── static/                 # Static assets
│   ├── css/               # Stylesheets
│   └── uploads/           # Uploaded video files (created at runtime)
└── tools/                  # Additional tools
    └── whisper.cpp/       # Optional local Whisper implementation
```

### Key Files

- **`main.go`**: Sets up the HTTP server, defines routes (`/`, `/upload`, `/transcribe`), and serves static files
- **`handlers/upload.go`**: Receives multipart form data, saves the video file, and returns an HTMX fragment with the video player
- **`handlers/transcribe.go`**: Orchestrates the transcription workflow by calling audio extraction and transcription services
- **`services/audio.go`**: Uses ffmpeg to extract audio from video files as MP3
- **`services/local_whisper.go`**: Invokes the Whisper CLI tool for local transcription
- **`services/openai.go`**: Calls the OpenAI Whisper API for cloud-based transcription

## Requirements

### System Dependencies

- **Go 1.18+** (tested with Go 1.25.4)
- **ffmpeg** - for audio extraction from video files
  ```bash
  # macOS
  brew install ffmpeg
  
  # Ubuntu/Debian
  sudo apt-get install ffmpeg
  
  # Windows
  # Download from https://ffmpeg.org/download.html
  ```

### Transcription Options

Choose **one** of the following:

#### Option 1: Local Whisper CLI (Free, runs on your machine)

```bash
pip install openai-whisper
```

**Note**: Requires Python 3.8+ and downloads AI models (~150MB-1.5GB depending on model size)

#### Option 2: OpenAI API (Cloud-based, requires API key)

- Sign up for an OpenAI account at https://platform.openai.com/
- Generate an API key
- Set the `OPENAI_API_KEY` environment variable

## Quickstart

1. **Clone the repository**:
   ```bash
   git clone https://github.com/DhrubajitPC/subtitle-gen.git
   cd subtitle-gen
   ```

2. **Install dependencies**:
   ```bash
   # Install ffmpeg (see Requirements section above)
   
   # Install Whisper CLI (for local transcription)
   pip install openai-whisper
   ```

3. **Run the server**:
   ```bash
   go run main.go
   ```
   
   The server will start on `http://localhost:8080`

4. **Use the application**:
   - Open your browser to `http://localhost:8080`
   - Click "Select Video File" and choose a video
   - Click "Upload" to upload the video
   - Once the video player appears, click "Generate Subtitles"
   - Wait for the transcription to complete (may take 30s-2min depending on video length)
   - View the transcript in the right panel

### Environment Variables

- **`PORT`**: Server port (default: `8080`)
  ```bash
  PORT=3000 go run main.go
  ```

- **`OPENAI_API_KEY`**: Required only if using OpenAI API for transcription
  ```bash
  export OPENAI_API_KEY=your-api-key-here
  go run main.go
  ```

## Building for Production

```bash
# Build the binary
go build -o subtitle-generator

# Run the binary
./subtitle-generator
```

## Technology Stack

- **Backend**: Go (standard library + net/http)
- **Frontend**: HTML, HTMX for dynamic interactions
- **AI/ML**: OpenAI Whisper (local CLI or API)
- **Media Processing**: ffmpeg for audio extraction

## Contributing

Contributions are welcome! Please refer to the architecture and sequence diagrams above to understand the codebase structure before making changes.

### Development Guidelines

1. Keep changes minimal and focused
2. Follow the existing code structure (handlers → services pattern)
3. Test both local Whisper and OpenAI API modes if making changes to transcription logic
4. Ensure ffmpeg compatibility for audio extraction changes

### Submitting a Pull Request

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source. Please check the repository for license details.

## Troubleshooting

### "whisper CLI tool not found"
- Ensure you've installed openai-whisper: `pip install openai-whisper`
- Verify it's in your PATH: `which whisper` or `whisper --help`

### "ffmpeg failed"
- Verify ffmpeg is installed: `ffmpeg -version`
- Check that the video file format is supported by ffmpeg

### Upload fails or times out
- Check the file size limit (default: 100MB)
- Ensure the `static/uploads/` directory has write permissions

### Transcription is slow
- Local Whisper processes on your CPU/GPU - larger videos take longer
- Consider using a smaller Whisper model (e.g., `tiny` or `base` instead of `large`)
- For faster results, use the OpenAI API option

## Acknowledgments

- [OpenAI Whisper](https://github.com/openai/whisper) - For the amazing speech recognition model
- [HTMX](https://htmx.org/) - For simplifying dynamic web interactions
- [ffmpeg](https://ffmpeg.org/) - For powerful media processing capabilities
