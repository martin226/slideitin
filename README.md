# Slide it In

A powerful presentation generator that transforms documents into beautiful slide decks using AI.

![Slide it In](docs/hero.png)

Try it now: https://justslideitin.com/

Case study: https://martinsit.ca/writing/shipping-presentation-generation-3-days

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


### Usage Flow

1. **Start**: Visit the homepage and click "Upload Document"
2. **Theme Selection**: Choose a visual theme for your presentation
3. **File Upload**: Upload your PDF, Markdown, or TXT files
4. **Settings**: Configure audience type and detail level
5. **Processing**: Wait while the AI generates your presentation
6. **Download**: Preview and download your generated presentation

## Tech Stack

### Frontend
- Next.js 14 (React framework)
- TypeScript
- Tailwind CSS for styling
- Framer Motion for animations

### Backend
- Go (Gin web framework)
- Google Cloud Firestore for job storage and status tracking
- Google Cloud Tasks for asynchronous job processing
- Server-Sent Events (SSE) for real-time status updates

### Slides Service
- Go microservice
- Gemini 1.5 Flash API for AI content generation
- Marp for converting markdown to presentation slides

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
├── backend/                # Go backend services
│   ├── api/                # Main API application
│   │   ├── controllers/    # API controllers
│   │   ├── models/         # Data models
│   │   ├── services/       # Business logic
│   │   │   ├── queue/      # Cloud Tasks queue management
│   │   │   └── ...
│   │   └── main.go         # Application entry point
│   │
│   └── slides-service/     # Microservice for slide generation
│       ├── controllers/    # API controllers
│       ├── models/         # Data models
│       ├── services/       # Business logic
│       │   ├── prompts/    # AI prompt templates
│       │   ├── slides/     # Slide generation service
│       │   │   └── themes/ # Presentation theme files
│       │   └── ...
│       └── main.go         # Service entry point
```

## Getting Started

This application is designed to run on Google Cloud Platform and cannot be run locally due to its GCP-specific dependencies.

### Prerequisites

- Google Cloud Platform account with billing enabled
- Google Cloud SDK installed
- Permissions to create and manage:
  - Cloud Run services
  - Cloud Tasks
  - Firestore database
  - Service accounts

The `cloudbuild.yaml` file handles the building and deployment of all services:
- Builds Docker images for frontend, backend, and slides-service
- Deploys each service to Cloud Run
- Configures service-to-service communication

## License

[MIT License](LICENSE)

## Acknowledgements

- [Marp](https://marp.app/) for presentation generation
- [Google Gemini](https://ai.google.dev/gemini-api) for AI content generation 