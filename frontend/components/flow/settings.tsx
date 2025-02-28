"use client"

import { useState } from "react"
import { ArrowLeft, Wand2 } from "lucide-react"
import { motion } from "framer-motion"

interface SettingsProps {
  onNext: (settings: { slideDetail: string; audience: string }) => void;
  onBack: (settings: { slideDetail: string; audience: string }) => void;
  initialSettings: {
    slideDetail: string;
    audience: string;
  };
}

const Settings = ({ onNext, onBack, initialSettings }: SettingsProps) => {
  const [slideDetail, setSlideDetail] = useState(initialSettings.slideDetail || "medium")
  const [audience, setAudience] = useState(initialSettings.audience || "general")

  const handleSubmit = () => {
    onNext({
      slideDetail,
      audience
    })
  }

  // Update parent state when navigating back
  const handleBack = () => {
    // Save current settings before going back
    onBack({
      slideDetail,
      audience
    })
  }

  return (
    <div className="w-full max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Additional Settings</h2>
      
      <div className="space-y-10">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <h3 className="text-lg font-medium text-gray-700 mb-4">Slide Detail Level</h3>
          <div className="space-y-3">
            {["minimal", "medium", "detailed"].map((option) => (
              <motion.label 
                key={option} 
                className={`flex items-center gap-4 p-4 border rounded-lg cursor-pointer hover:bg-amber-50/50 transition-colors ${
                  slideDetail === option ? "border-amber-500 bg-amber-50/30" : "border-gray-200"
                }`}
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.99 }}
              >
                <input
                  type="radio"
                  name="slideDetail"
                  value={option}
                  checked={slideDetail === option}
                  onChange={() => setSlideDetail(option)}
                  className="h-5 w-5 text-amber-500"
                />
                <div>
                  <span className="block font-medium text-gray-700 capitalize text-lg">{option}</span>
                  <span className="text-sm text-gray-500">
                    {option === "minimal" && "Less text, more visuals - perfect for high-level overviews"}
                    {option === "medium" && "Balanced text and visuals - suitable for most presentations"}
                    {option === "detailed" && "Comprehensive content with details - ideal for in-depth analysis"}
                  </span>
                </div>
              </motion.label>
            ))}
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <h3 className="text-lg font-medium text-gray-700 mb-4">Target Audience</h3>
          <select
            value={audience}
            onChange={(e) => setAudience(e.target.value)}
            className="w-full p-4 border rounded-lg focus:border-amber-500 focus:ring focus:ring-amber-200 focus:ring-opacity-50 transition-all"
          >
            <option value="general">General</option>
            <option value="academic">Academic</option>
            <option value="technical">Technical</option>
            <option value="professional">Business</option>
            <option value="executive">Executive</option>
          </select>
          <p className="text-sm text-gray-500 mt-2 ml-1">
            Tailors the presentation style and language to your specific audience
          </p>
        </motion.div>
      </div>

      <div className="mt-10 flex justify-between">
        <button
          onClick={handleBack}
          className="px-4 py-3 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center gap-2 transition-colors"
        >
          <ArrowLeft className="h-4 w-4" /> Back
        </button>
          
        <motion.button
          whileHover={{ scale: 1.03 }}
          whileTap={{ scale: 0.98 }}
          onClick={handleSubmit}
          className="font-bold px-8 py-3 bg-amber-500 hover:bg-amber-600 text-white rounded-lg flex items-center gap-2 shadow-md hover:shadow-lg transition-all"
        >
          <Wand2 className="h-5 w-5" />
          Slide it In
        </motion.button>
      </div>
    </div>
  )
}

export default Settings 