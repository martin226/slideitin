# SlideItIn Backend

This is the backend service for SlideItIn, a tool that generates presentation slides from documents.

## Features

- Secure file upload handling
- Presentation generation
- Asynchronous job queue for processing slide generation jobs
- RESTful API for frontend communication
- Placeholders for Gemini API integration
- WebSocket for real-time job status updates

## Requirements

- Go 1.18+
- Gemini API key (optional for development)

## Setup and Installation

1. Clone the repository
2. Navigate to the `backend` directory
3. Set environment variables:

```bash
export PORT=8080 # Optional, defaults to 8080
export GEMINI_API_KEY=your_gemini_api_key
```

4. Build and run:

```bash
go build -o slideitin
./slideitin
```

## API Endpoints

### Generate Slides

Creates a new slide generation job based on uploaded files.

**Endpoint**: `POST /api/slides`

**Request**:
- Content-Type: `multipart/form-data`
- Body:
  - `files`: One or more files (PDF, DOCX, TXT, etc.)
  - `theme`: Theme for the slides (string)
  - `settings`: JSON string containing slide settings

**Response**:
```json
{
  "jobID": "unique-job-id",
  "status": "queued"
}
```

### Get Job Status

Retrieves the current status of a slide generation job.

**Endpoint**: `GET /api/jobs/:jobID`

**Response**:
```json
{
  "jobID": "unique-job-id",
  "status": "queued|processing|completed|failed",
  "message": "Status message",
  "result": "URL to the generated slides (if completed)"
}
```

## WebSocket

Real-time job status updates are available via WebSocket connection.

**Endpoint**: `ws://localhost:8080/ws/jobs/:jobID`

**Messages**:
```json
{
  "status": "processing",
  "message": "Generating slides...",
  "progress": 50
}
```

## Implementation Notes

This is a placeholder implementation. The actual Gemini integration and slide generation will be implemented in the future.

## Environment Variables

- `PORT`: The port to run the server on (default: 8080)
- `GEMINI_API_KEY`: Your Gemini API key
- `MAX_UPLOAD_SIZE`: Maximum upload size in bytes (default: 10MB)

## License

MIT 