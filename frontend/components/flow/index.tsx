"use client"

import { useState } from "react"
import { CheckCircle } from "lucide-react"
import { motion, AnimatePresence } from "framer-motion"

// Import individual components
import ThemeSelection from "./theme-selection"
import FileUpload from "./file-upload"
import Settings from "./settings"
import Success from "./success"
import Result from "./result"

// Main Upload Flow Component
export default function UploadFlow({ onBack }: { onBack?: () => void }) {
  const [step, setStep] = useState(1)
  const [data, setData] = useState({
    theme: "",
    files: [] as File[],
    settings: {
      slideDetail: "medium",
      audience: "general"
    }
  })
  const [resultUrl, setResultUrl] = useState<string | null>(null)

  const nextStep = () => {
    setStep((prev) => prev + 1)
  }

  const prevStep = () => {
    setStep((prev) => prev - 1)
  }

  const restartFlow = () => {
    // Reset data and go back to step 1
    setData({
      theme: "",
      files: [],
      settings: {
        slideDetail: "medium",
        audience: "general"
      }
    })
    setResultUrl(null)
    setStep(1)
  }

  const handleThemeSelect = (theme: string) => {
    setData((prev) => ({ ...prev, theme }))
    nextStep()
  }

  const handleFilesSelect = (files: File[]) => {
    setData((prev) => ({ ...prev, files }))
    nextStep()
  }

  const handleSettingsSubmit = (settings: { slideDetail: string; audience: string }) => {
    setData((prev) => ({ ...prev, settings }))
    nextStep()
  }

  // Update settings and go back
  const handleSettingsBack = (settings: { slideDetail: string; audience: string }) => {
    setData((prev) => ({ ...prev, settings }))
    prevStep()
  }

  // Handle job completion - called when job status is "completed"
  const handleJobCompletion = (jobData: { resultUrl: string }) => {
    console.log("Job completion handler called with:", jobData);
    console.log("Setting resultUrl to:", jobData.resultUrl);
    setResultUrl(jobData.resultUrl)
    console.log("Transitioning to result step");
    nextStep()
  }

  const steps = [
    { title: "Theme", completed: step > 1 },
    { title: "Upload Files", completed: step > 2 },
    { title: "Settings", completed: step > 3 },
    { title: "Processing", completed: step > 4 },
    { title: "Result", completed: false }
  ]

  return (
    <div className="h-full bg-amber-50 flex flex-col">
      <div className="container mx-auto px-4 flex flex-col h-full">
        {/* Progress tracker - at the top */}
        <motion.div 
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="pt-8"
        >
          {/* Step circles and lines */}
          <div className="flex justify-center items-center space-x-3 md:space-x-6">
            {steps.map((s, i) => (
              <div key={i} className="flex flex-col items-center">
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  className={`w-12 h-12 md:w-14 md:h-14 rounded-full flex items-center justify-center ${
                    step > i + 1 ? "bg-amber-400 text-white" : 
                    step === i + 1 ? "bg-amber-500 text-white" : "bg-gray-200 text-gray-500"
                  } transition-colors shadow-sm relative cursor-pointer`}
                >
                  {s.completed ? (
                    <CheckCircle className="w-6 h-6" />
                  ) : (
                    <span className="text-lg font-bold">{i + 1}</span>
                  )}
                  
                  {/* Line connector */}
                  {i < steps.length - 1 && (
                    <div className="absolute left-full top-1/2 -translate-y-1/2 h-0.5 bg-gray-200 w-6 md:w-12">
                      <div 
                        className={`h-full ${
                          step > i + 1 ? "bg-amber-400" : "bg-gray-200"
                        }`}
                        style={{ width: `${step > i + 1 ? '100%' : '0%'}` }}
                      ></div>
                    </div>
                  )}
                </motion.div>
                
                <span className={`mt-2 text-xs md:text-sm text-center whitespace-nowrap font-medium ${
                  step === i + 1 ? "text-amber-700" : "text-gray-600"
                }`}>
                  {s.title}
                </span>
              </div>
            ))}
          </div>
        </motion.div>

        {/* Form steps - positioned slightly higher than center */}
        <div className="flex-1 flex items-start justify-center mt-10 mb-16">
          <AnimatePresence mode="wait">
            <motion.div
              key={step}
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
              className="max-w-4xl w-full"
            >
              {step === 1 && <ThemeSelection 
                onNext={handleThemeSelect} 
                onBack={() => onBack?.()} 
                initialTheme={data.theme}
              />}
              {step === 2 && <FileUpload 
                onNext={handleFilesSelect} 
                onBack={() => prevStep()} 
                initialFiles={data.files}
              />}
              {step === 3 && <Settings 
                onNext={handleSettingsSubmit} 
                onBack={handleSettingsBack} 
                initialSettings={data.settings}
              />}
              {step === 4 && <Success 
                data={data}
                onComplete={(jobData: { resultUrl: string }) => handleJobCompletion(jobData)}
              />}
              {step === 5 && <Result 
                onRestart={restartFlow} 
                resultUrl={resultUrl!} 
              />}
            </motion.div>
          </AnimatePresence>
        </div>
      </div>
    </div>
  )
} 