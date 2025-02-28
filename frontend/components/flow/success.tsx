"use client"

import { Wand2 } from "lucide-react"
import { motion } from "framer-motion"
import { useState, useEffect, useRef, useCallback } from "react"
import { generateSlides, subscribeToSlideUpdates, SlideRequest, SlideUpdate } from "../../lib/api"

interface SuccessProps {
  onComplete: (jobData: { resultUrl: string }) => void;
  data: {
    theme: string;
    files: File[];
    settings: {
      slideDetail: string;
      audience: string;
    };
  };
}

const Success = ({ onComplete, data }: SuccessProps) => {
  const [status, setStatus] = useState("Initializing...")
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const generationStarted = useRef(false)
  const jobCompleted = useRef(false)

  const handleStatusUpdate = useCallback((update: SlideUpdate) => {
    // Update status message
    setStatus(update.message);
    console.log("Received status update:", update);
    
    // Update progress based on status
    if (update.status === "queued") {
      setProgress(10);
    } else if (update.status === "processing") {
      // Determine progress based on specific status messages
      if (update.message.includes("Analyzing")) {
        setProgress(30);
      } else if (update.message.includes("Generating content")) {
        setProgress(50);
      } else if (update.message.includes("Creating presentation")) {
        setProgress(70);
      } else if (update.message.includes("Finalizing")) {
        setProgress(90);
      } else {
        setProgress(50); // Default for processing
      }
    } else if (update.status === "completed") {
      setProgress(100);
      console.log("Job completed with resultUrl:", update.resultUrl);
      
      // Mark job as completed
      jobCompleted.current = true;
      
      // Store reference to prevent cleanup issues
      const jobResultUrl = update.resultUrl;
      
      // Move to next step after a brief delay to show completion
      setTimeout(() => {
        console.log("Calling onComplete with:", { resultUrl: jobResultUrl });
        onComplete({ resultUrl: jobResultUrl });
      }, 1000);
    } else if (update.status === "failed") {
      setError(update.message || "Job failed");
      // Mark failed jobs as completed too
      jobCompleted.current = true;
    }
  }, [setProgress, setStatus, onComplete, setError]);

  useEffect(() => {
    let cleanup: (() => void) | null = null;

    const startSlideGeneration = async () => {
      // Only start the slide generation once
      if (generationStarted.current) {
        return;
      }
      generationStarted.current = true;

      try {
        // Prepare data for API call
        const slideRequest: SlideRequest = {
          theme: data.theme,
          settings: data.settings
        };

        // Call API to generate slides
        setStatus("Submitting job...")
        const response = await generateSlides(slideRequest, data.files);
        
        // No need to store job ID as state since it's never used
        setStatus(response.message || "Job submitted successfully");
        
        // Subscribe to status updates via SSE using the response ID directly
        cleanup = subscribeToSlideUpdates(
          response.id,
          handleStatusUpdate,
          (error) => {
            // Only set error if we haven't completed the job yet
            if (!jobCompleted.current) {
              console.error("SSE Error:", error);
              setError(error.message);
            } else {
              console.log("Ignoring error after job completion:", error.message);
            }
          }
        );
        
      } catch (error) {
        console.error("Failed to generate slides:", error);
        setError(error instanceof Error ? error.message : "Unknown error occurred");
      }
    };

    startSlideGeneration();

    // Cleanup function to close SSE connection
    return () => {
      if (cleanup) cleanup();
    };
  }, [data, onComplete, handleStatusUpdate]);

  return (
    <div className="w-full max-w-4xl mx-auto text-center py-12">
      <motion.div 
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ type: "spring", stiffness: 300, damping: 20 }}
        className="inline-flex items-center justify-center w-24 h-24 rounded-full bg-amber-100 mb-6"
      >
        <Wand2 className="h-12 w-12 text-amber-600" />
      </motion.div>
      <motion.h2 
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="text-3xl font-bold text-gray-800 mb-3"
      >
        Creating Your Slides
      </motion.h2>
      <motion.p 
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
        className="text-xl text-gray-600 mb-8"
      >
        {error ? (
          <span className="text-red-500">{error}</span>
        ) : (
          status
        )}
      </motion.p>
      <motion.div 
        initial={{ width: 0 }}
        animate={{ width: "100%" }}
        transition={{ delay: 0.4 }}
        className="max-w-md h-6 mx-auto bg-gray-200 rounded-full overflow-hidden shadow-inner"
      >
        <motion.div 
          initial={{ width: "0%" }}
          animate={{ width: `${progress}%` }}
          transition={{ 
            duration: 0.5, 
            ease: "easeInOut"
          }}
          className="h-full bg-gradient-to-r from-orange-400 via-orange-500 to-orange-600 rounded-full relative"
        >
          <motion.div
            animate={{ 
              opacity: [0.3, 0.6, 0.3],
              scale: [1, 1.02, 1]
            }}
            transition={{
              duration: 1.5,
              repeat: Infinity,
              repeatType: "loop"
            }}
            className="absolute inset-0 bg-white opacity-30 rounded-full"
          />
        </motion.div>
      </motion.div>
      
      <motion.p
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.6 }}
        className="text-sm text-gray-500 mt-6"
      >
        This usually takes about 15 seconds to 1 minute depending on document size
      </motion.p>
    </div>
  )
}

export default Success 