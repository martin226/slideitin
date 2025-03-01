"use client"

import { useState } from "react"
import { ChevronRight, Upload, ArrowLeft, FileText } from "lucide-react"
import { motion, AnimatePresence } from "framer-motion"

// File Upload Component
const FileUpload = ({ onNext, onBack, initialFiles }: { 
  onNext: (files: File[]) => void; 
  onBack: () => void;
  initialFiles: File[];
}) => {
  const [files, setFiles] = useState<File[]>(initialFiles)
  const [isDragging, setIsDragging] = useState(false)
  const [textContent, setTextContent] = useState("")

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
    
    // allow PDF, MD, TXT - check both MIME type and file extension for .md files
    const droppedFiles = Array.from(e.dataTransfer.files).filter(
      (file) => 
        file.type === "application/pdf" || 
        file.type === "text/markdown" || 
        file.type === "text/plain" ||
        file.name.toLowerCase().endsWith('.md')
    )
    
    if (droppedFiles.length > 0) {
      setFiles((prev) => [...prev, ...droppedFiles])
    }
  }

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files?.length) {
      const selectedFiles = Array.from(e.target.files).filter(
        (file) => 
          file.type === "application/pdf" || 
          file.type === "text/markdown" || 
          file.type === "text/plain" ||
          file.name.toLowerCase().endsWith('.md')
      )
      setFiles((prev) => [...prev, ...selectedFiles])
    }
  }

  const removeFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index))
  }

  const handleTextContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setTextContent(e.target.value)
  }

  const handleSubmit = () => {
    // If we have text content, convert it to a File object
    if (textContent.trim()) {
      const textBlob = new Blob([textContent], { type: 'text/plain' })
      const textFile = new File([textBlob], 'content.txt', { type: 'text/plain' })
      
      onNext([...files, textFile])
    } else {
      onNext(files)
    }
  }

  const hasContent = files.length > 0 || textContent.trim().length > 0

  return (
    <div className="w-full max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Upload Documents</h2>
      
      {/* Split container for desktop, column for mobile */}
      <div className="flex flex-col md:flex-row gap-4 md:gap-6 mb-8">
        {/* Left side - File upload */}
        <div className="w-full md:w-1/2">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className={`border-2 border-dashed rounded-lg p-6 md:p-8 text-center h-full ${
              isDragging ? "border-amber-500 bg-amber-50" : "border-gray-300"
            } transition-colors relative overflow-hidden`}
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
              className="relative z-10 h-full flex flex-col justify-center items-center"
            >
              <Upload className="mx-auto h-12 w-12 md:h-16 md:w-16 text-amber-500 mb-4" />
              <p className="text-lg md:text-xl text-gray-700 mb-2 font-medium">
                Drag and drop your files here
              </p>
              <p className="text-sm md:text-md text-gray-500 mb-4 md:mb-6">or click to browse your files</p>
              <input
                type="file"
                id="file-upload"
                className="hidden"
                multiple
                accept=".pdf, .md, .txt"
                onChange={handleFileInput}
              />
              <label
                htmlFor="file-upload"
                className="px-4 py-2 md:px-6 md:py-3 border border-gray-300 rounded-md text-gray-700 font-medium cursor-pointer hover:bg-gray-50 transition-colors inline-flex items-center gap-2"
              >
                <span>Browse Files</span>
              </label>
              <p className="text-xs md:text-sm text-gray-500 mt-4 md:mt-6">Supports PDF, MD, and TXT files</p>
            </motion.div>
          </motion.div>
        </div>

        {/* Middle divider - Only visible on desktop */}
        <div className="hidden md:flex flex-col items-center justify-center">
          <div className="h-36 w-px bg-gray-300"></div>
          <div className="py-3 px-4 rounded-full bg-amber-100 text-amber-800 font-medium my-2">or</div>
          <div className="h-36 w-px bg-gray-300"></div>
        </div>

        {/* Mobile divider */}
        <div className="md:hidden flex items-center justify-center mb-2">
          <div className="flex-1 h-px bg-gray-300"></div>
          <div className="py-3 px-4 rounded-full bg-amber-100 text-amber-800 font-medium mx-4">or</div>
          <div className="flex-1 h-px bg-gray-300"></div>
        </div>

        {/* Right side - Text content */}
        <div className="w-full md:w-1/2 flex flex-col">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="border-2 border-gray-300 rounded-lg flex flex-col flex-grow h-full"
          >
            <div className="p-3 flex items-center border-b border-gray-200 flex-shrink-0">
              <FileText className="h-5 w-5 text-amber-500 mr-2" />
              <h3 className="text-sm font-medium text-gray-700">Paste your content here</h3>
            </div>
            <textarea
              className="w-full flex-grow p-4 resize-none focus:outline-none focus:ring-1 focus:ring-amber-400 rounded-b-lg min-h-[200px]"
              placeholder="This will be used in addition to any files you upload."
              autoFocus
              value={textContent}
              onChange={handleTextContentChange}
            ></textarea>
          </motion.div>
        </div>
      </div>

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
          onClick={handleSubmit}
          disabled={!hasContent}
          className={`px-6 py-3 rounded-md flex items-center gap-2 ${
            !hasContent 
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