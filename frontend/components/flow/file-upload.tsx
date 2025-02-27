"use client"

import { useState } from "react"
import { ChevronRight, Upload, ArrowLeft } from "lucide-react"
import { motion, AnimatePresence } from "framer-motion"

// File Upload Component
const FileUpload = ({ onNext, onBack, initialFiles }: { 
  onNext: (files: File[]) => void; 
  onBack: () => void;
  initialFiles: File[];
}) => {
  const [files, setFiles] = useState<File[]>(initialFiles)
  const [isDragging, setIsDragging] = useState(false)

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }

  const handleDragLeave = () => {
    setIsDragging(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    
    // allow PDF, MD, TXT
    const droppedFiles = Array.from(e.dataTransfer.files).filter(
      (file) => file.type === "application/pdf" || file.type === "text/markdown" || file.type === "text/plain"
    )
    
    if (droppedFiles.length > 0) {
      setFiles((prev) => [...prev, ...droppedFiles])
    }
  }

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.length) {
      const selectedFiles = Array.from(e.target.files).filter(
        (file) => file.type === "application/pdf" || file.type === "text/markdown" || file.type === "text/plain"
      )
      setFiles((prev) => [...prev, ...selectedFiles])
    }
  }

  const removeFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index))
  }

  return (
    <div className="w-full max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Upload Documents</h2>
      
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className={`border-2 border-dashed rounded-lg p-10 text-center ${
          isDragging ? "border-amber-500 bg-amber-50" : "border-gray-300"
        } transition-colors mb-8 relative overflow-hidden`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {isDragging && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 0.15 }}
            className="absolute inset-0 bg-amber-200 z-0"
          />
        )}
        <motion.div
          initial={{ scale: 0.9 }}
          animate={{ scale: 1 }}
          transition={{ type: "spring", stiffness: 300, damping: 20 }}
          className="relative z-10"
        >
          <Upload className="mx-auto h-16 w-16 text-amber-500 mb-4" />
          <p className="text-xl text-gray-700 mb-2 font-medium">
            Drag and drop your files here
          </p>
          <p className="text-md text-gray-500 mb-6">or click to browse your files</p>
          <input
            type="file"
            id="file-upload"
            className="hidden"
            multiple
            accept=".pdf"
            onChange={handleFileInput}
          />
          <label
            htmlFor="file-upload"
            className="px-6 py-3 border border-gray-300 rounded-md text-gray-700 font-medium cursor-pointer hover:bg-gray-50 transition-colors inline-flex items-center gap-2"
          >
            <span>Browse Files</span>
          </label>
          <p className="text-sm text-gray-500 mt-6">Supports PDF, MD, and TXT files</p>
        </motion.div>
      </motion.div>

      <AnimatePresence>
        {files.length > 0 && (
          <motion.div 
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="mb-8 overflow-hidden"
          >
            <h3 className="font-medium text-gray-700 mb-3">Selected Files ({files.length})</h3>
            <div className="space-y-2">
              {files.map((file, index) => (
                <motion.div 
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                  transition={{ delay: index * 0.05 }}
                  className="flex items-center justify-between border border-gray-200 p-4 rounded-md hover:border-amber-200 hover:bg-amber-50/30 transition-colors"
                >
                  <span className="truncate text-gray-600">{file.name}</span>
                  <button
                    onClick={() => removeFile(index)}
                    className="text-gray-400 hover:text-red-500 transition-colors"
                  >
                    Remove
                  </button>
                </motion.div>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="flex justify-between">
        <button
          onClick={onBack}
          className="px-4 py-3 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center gap-2 transition-colors"
        >
          <ArrowLeft className="h-4 w-4" /> Back
        </button>
          
        <motion.button
          whileHover={{ scale: 1.03 }}
          whileTap={{ scale: 0.98 }}
          onClick={() => onNext(files)}
          disabled={files.length === 0}
          className={`px-6 py-3 rounded-md flex items-center gap-2 ${
            files.length === 0 
              ? "bg-gray-200 text-gray-500 cursor-not-allowed" 
              : "bg-amber-500 hover:bg-amber-600 text-white shadow-md hover:shadow-lg"
          } transition-all`}
        >
          Next <ChevronRight className="h-4 w-4" />
        </motion.button>
      </div>
    </div>
  )
}

export default FileUpload 