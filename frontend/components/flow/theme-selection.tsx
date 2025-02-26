"use client"

import { useState } from "react"
import { ChevronRight, ArrowLeft } from "lucide-react"
import { motion } from "framer-motion"

// Theme Selection Component
const ThemeSelection = ({ onNext, onBack, initialTheme }: { 
  onNext: (theme: string) => void; 
  onBack: () => void;
  initialTheme: string;
}) => {
  const [selectedTheme, setSelectedTheme] = useState(initialTheme)
  const themes = [
    { id: "minimal", name: "Minimal", color: "bg-gray-50" },
    { id: "corporate", name: "Corporate", color: "bg-blue-50" },
    { id: "creative", name: "Creative", color: "bg-purple-50" },
    { id: "academic", name: "Academic", color: "bg-amber-50/70" },
    { id: "modern", name: "Modern", color: "bg-pink-50" },
    { id: "elegant", name: "Elegant", color: "bg-orange-50" },
  ]

  return (
    <div className="w-full max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Choose a Presentation Theme</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
        {themes.map((theme, index) => (
          <motion.div
            key={theme.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ 
              opacity: 1, 
              y: 0,
              transition: { delay: index * 0.1 }
            }}
            whileHover={{ 
              scale: 1.03,
              boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.05), 0 10px 10px -5px rgba(0, 0, 0, 0.02)"
            }}
            onClick={() => {
              setSelectedTheme(theme.id)
              // Don't automatically proceed to next step
            }}
            className={`border-2 ${selectedTheme === theme.id ? 'border-amber-500' : 'border-gray-200'} hover:border-amber-500 rounded-lg p-5 cursor-pointer transition-all overflow-hidden`}
          >
            <div className={`h-32 ${theme.color} rounded-md mb-4 flex items-center justify-center relative`}>
              <div className="absolute inset-0 bg-gradient-to-br from-transparent to-white/20"></div>
              <span className="text-gray-700 font-medium relative z-10">{theme.name} Preview</span>
            </div>
            <p className="font-medium text-gray-700 text-center">{theme.name}</p>
          </motion.div>
        ))}
      </div>
      
      <div className="mt-8 flex justify-between">
        <button
          onClick={onBack}
          className="px-4 py-3 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center gap-2 transition-colors"
        >
          <ArrowLeft className="h-4 w-4" /> Back to Home
        </button>
        
        {selectedTheme && (
          <motion.button
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.98 }}
            onClick={() => onNext(selectedTheme)}
            className="px-6 py-3 bg-amber-500 hover:bg-amber-600 text-white rounded-md flex items-center gap-2 shadow-md hover:shadow-lg transition-all"
          >
            Next <ChevronRight className="h-4 w-4" />
          </motion.button>
        )}
      </div>
    </div>
  )
}

export default ThemeSelection 