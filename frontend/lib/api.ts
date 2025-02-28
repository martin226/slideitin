// API service for communicating with the backend

// Get API URL from environment variable or fall back to localhost for development
export const API_BASE_URL = 
  process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/v1';

// Types
export interface SlideRequest {
  theme: string;
  settings: {
    slideDetail: string;
    audience: string;
  };
}

export interface SlideResponse {
  id: string;
  status: string;
  message: string;
  createdAt: number;
  updatedAt: number;
}

export interface SlideUpdate {
  id: string;
  status: string;
  message: string;
  resultUrl: string;
  updatedAt: number;
}

// Generate slides by sending data and files to the backend
export async function generateSlides(
  data: SlideRequest,
  files: File[]
): Promise<SlideResponse> {
  const formData = new FormData();
  
  // Add JSON data
  formData.append('data', JSON.stringify(data));
  
  // Add files
  files.forEach(file => {
    formData.append('files', file);
  });
  
  try {
    const response = await fetch(`${API_BASE_URL}/generate`, {
      method: 'POST',
      body: formData,
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to generate slides');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error generating slides:', error);
    throw error;
  }
}

// Create an EventSource for server-sent events to get status updates
export function subscribeToSlideUpdates(
  slideId: string,
  onUpdate: (update: SlideUpdate) => void,
  onError: (error: Error) => void
): () => void {
  // Create EventSource for SSE connection
  const eventSource = new EventSource(`${API_BASE_URL}/slides/${slideId}`);
  
  // Handle normal update events
  eventSource.addEventListener('update', (event) => {
    try {
      const data = JSON.parse(event.data);
      onUpdate(data);
      
      // If the status is completed or failed, prepare for stream to end
      if (data.status === 'completed' || data.status === 'failed') {
        console.log(`Job ${data.status}. Stream will close soon.`);
      }
    } catch (error) {
      console.error('Error parsing SSE data:', error);
      onError(new Error('Failed to parse server update'));
    }
  });
  
  // Handle explicit close events from server
  eventSource.addEventListener('close', () => {
    console.log('Server indicated stream is closing normally');
    eventSource.close();
  });
  
  // Handle connection errors
  eventSource.addEventListener('error', (event) => {
    // Don't report an error if we've already received a completed/failed status
    // or if the connection was closed cleanly
    if (eventSource.readyState === EventSource.CLOSED) {
      console.log('EventSource connection closed');
    } else {
      console.error('SSE connection error:', event);
      onError(new Error('Connection to server lost'));
    }
    eventSource.close();
  });
  
  // Return a cleanup function
  return () => {
    eventSource.close();
  };
} 