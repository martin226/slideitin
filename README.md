# Slide it In

A powerful presentation generator that transforms documents into beautiful slide decks using AI.

![Slide it In](docs/hero.png)

Try it now: https://justslideitin.com/

## Overview

Slide it In is a web application that uses AI to automatically generate professional presentations from your documents. Simply upload your PDF, Markdown, or TXT files, select a visual theme, customize your settings, and let the AI create a beautifully formatted slide deck in seconds.

## Features

- 🤖 **AI-Powered Presentation Generation**: Uses Gemini 1.5 Flash to analyze your documents and create compelling slide content
- 📊 **Multiple Visual Themes**: Choose from various professionally designed themes (Default, Beam, Rose Pine, Gaia, Uncover, Graph Paper)
- 🎯 **Audience Targeting**: Tailor your presentation for different audiences (General, Academic, Technical, Professional, Executive)
- ⚙️ **Detail Level Control**: Customize the amount of detail extracted from your documents (Minimal, Medium, Detailed)
- 🔄 **Real-time Progress Updates**: Monitor the generation process with live status updates
- 📱 **Responsive Design**: Works seamlessly on desktop and mobile devices
- 📤 **PDF Export**: Download your generated presentations as PDF files

## Tech Stack

### Frontend
- Next.js 14 (React framework)
- TypeScript
- Tailwind CSS for styling
- Framer Motion for animations
- React-PDF for PDF preview

### Backend
- Go (Gin web framework)
- Google Cloud Firestore for job storage and status tracking
- Gemini 1.5 Flash API for AI content generation
- Marp for converting markdown to presentation slides
- Server-Sent Events (SSE) for real-time status updates

## Project Structure

```
slideitin/
├── frontend/               # Next.js frontend application
│   ├── app/                # Next.js app directory
│   ├── components/         # React components
│   │   ├── flow/           # Presentation generation flow components
│   │   └── ...
│   ├── lib/                # Utility functions and API client
│   └── ...
│
└── backend/                # Go backend application
    ├── controllers/        # API controllers
    ├── models/             # Data models
    ├── services/           # Business logic
    │   ├── prompts/        # AI prompt templates
    │   ├── queue/          # Job queue management
    │   ├── slides/         # Slide generation service
    │   │   └── themes/     # Presentation theme files
    │   └── ...
    └── main.go             # Application entry point
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go 1.20+
- Google Cloud account with Firestore enabled
- Gemini API key

### Frontend Setup

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

4. The frontend will be available at http://localhost:3000

### Backend Setup

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Copy the example environment file:
   ```
   cp .env.example .env
   ```

3. Edit the `.env` file with your Google Cloud project ID and Gemini API key

4. Install Go dependencies:
   ```
   go mod download
   ```

5. Run the backend server:
   ```
   go run main.go
   ```

6. The API will be available at http://localhost:8080

## Usage Flow

1. **Start**: Visit the homepage and click "Upload Document"
2. **Theme Selection**: Choose a visual theme for your presentation
3. **File Upload**: Upload your PDF, Markdown, or TXT files
4. **Settings**: Configure audience type and detail level
5. **Processing**: Wait while the AI generates your presentation
6. **Download**: Preview and download your generated presentation

## Environment Variables

### Backend

```
PORT=8080
GEMINI_API_KEY=your_gemini_api_key
GOOGLE_CLOUD_PROJECT=your_gcp_project_id
GOOGLE_APPLICATION_CREDENTIALS=./service-account.json
```

## License

[MIT License](LICENSE)

## Acknowledgements

- [Marp](https://marp.app/) for presentation generation
- [Google Gemini](https://ai.google.dev/gemini-api) for AI content generation 