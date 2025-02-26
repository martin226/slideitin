"use client"

import { Wand2 } from "lucide-react"
import { motion } from "framer-motion"

const Success = () => {
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
        We're working our magic to transform your documents into beautiful slides.
      </motion.p>
      <motion.div 
        initial={{ width: 0 }}
        animate={{ width: "100%" }}
        transition={{ delay: 0.4 }}
        className="max-w-md h-2 mx-auto bg-gray-200 rounded-full overflow-hidden"
      >
        <motion.div 
          initial={{ width: "0%" }}
          animate={{ width: "75%" }}
          transition={{ 
            delay: 0.5,
            duration: 1.5, 
            ease: "easeInOut",
            repeat: Infinity,
            repeatType: "reverse" 
          }}
          className="h-full bg-gradient-to-r from-amber-400 to-amber-600 rounded-full"
        ></motion.div>
      </motion.div>
      
      <motion.p
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.6 }}
        className="text-sm text-gray-500 mt-6"
      >
        This usually takes about 1-2 minutes depending on document size
      </motion.p>
    </div>
  )
}

export default Success 