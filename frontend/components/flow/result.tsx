"use client"

import { Download, Edit, X, RefreshCw } from "lucide-react"
import { useState } from "react"
import { API_BASE_URL } from "@/lib/api"

interface ResultProps {
  onRestart?: () => void;
  resultUrl: string;
}

const Result = ({ onRestart, resultUrl }: ResultProps) => {
  const [tutorialOpen, setTutorialOpen] = useState(false)
  
  console.log("Result component rendered with resultUrl:", resultUrl);

  // Custom Modal component
  const TutorialModal = () => {
    if (!tutorialOpen) return null;
    
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div className="bg-white rounded-lg max-w-md w-full p-6 relative">
          <button 
            onClick={() => setTutorialOpen(false)}
            className="absolute right-4 top-4 text-gray-500 hover:text-gray-700"
          >
            <X size={20} />
          </button>
          
          <h2 className="text-2xl font-bold mb-4">How to Edit Your Slides</h2>
          
          <div className="space-y-6 py-2">
            <div className="flex items-start gap-4">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-amber-500 flex items-center justify-center text-white font-bold">1</div>
              <div>
                <h3 className="font-semibold text-lg">Download the PDF</h3>
                <p className="text-gray-600">First, download your slides as a PDF using the "Download as PDF" button.</p>
              </div>
            </div>
            
            <div className="flex items-start gap-4">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-amber-500 flex items-center justify-center text-white font-bold">2</div>
              <div>
                <h3 className="font-semibold text-lg">Convert to PowerPoint</h3>
                <p className="text-gray-600">Go to <a href="https://www.adobe.com/acrobat/online/pdf-to-ppt.html" target="_blank" rel="noopener noreferrer" className="text-blue-500 underline">Adobe's PDF to PPT converter</a>, convert your file, then edit the PowerPoint slides.</p>
              </div>
            </div>

            <p className="text-gray-600 text-xs">
              (unfortunately, PDFs are proprietary and we cannot convert them ourselves, but Adobe's converter works almost perfectly)
            </p>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="w-full max-w-7xl mx-auto p-6">
      <div className="flex flex-col gap-6">
        {/* Create Another button positioned at the top right */}
        <div className="flex justify-end w-full">
          <button 
            onClick={() => onRestart?.()}
            className="text-sm flex items-center gap-1 text-gray-600 hover:text-amber-600 transition-colors focus:outline-none"
          >
            <RefreshCw size={14} />
            <span>Create Another Presentation</span>
          </button>
        </div>
        
        {/* Slides viewer with embedded PDF */}
        <div className="w-full bg-gray-100 rounded-lg shadow-md">
          {/* Using aspect ratio container for 16:9 */}
          <div className="relative" style={{ paddingBottom: "56.25%" }}>
            <iframe 
              src={API_BASE_URL + resultUrl}
              className="absolute inset-0 w-full h-full border-0 rounded-lg"
              title="Slide Presentation"
              allowFullScreen
            />
          </div>
        </div>
        
        {/* Buttons in a row underneath */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center mt-4">
          <button 
            className="py-3 px-6 rounded-lg bg-amber-500 hover:bg-amber-600 transition-colors flex items-center justify-center gap-2 text-white font-medium"
            onClick={() => window.open(API_BASE_URL + resultUrl + "?download=true", '_blank')}
          >
            <Download size={18} />
            Download as PDF
          </button>
          
          <button 
            onClick={() => setTutorialOpen(true)}
            className="py-3 px-6 rounded-lg bg-amber-400 hover:bg-amber-500 transition-colors flex items-center justify-center gap-2 text-white font-medium"
          >
            <Edit size={18} />
            Edit
          </button>
        </div>
      </div>
      
      {/* Tutorial Modal */}
      <TutorialModal />
    </div>
  )
}

export default Result 