"use client"

import { motion } from "framer-motion"
import { FileText, FileIcon as FilePresentation } from "lucide-react"
import { cn } from "@/lib/utils"

export function FloatingDocument({
  className,
  delay = 0,
  rotate = 0,
  isPresentation = false,
}: {
  className?: string
  delay?: number
  rotate?: number
  isPresentation?: boolean
}) {
  return (
    <motion.div
      initial={{
        opacity: 0,
        y: 20,
        rotate: rotate - 5,
      }}
      animate={{
        opacity: 1,
        y: 0,
        rotate,
      }}
      transition={{
        duration: 1.8,
        delay,
        ease: [0.23, 0.86, 0.39, 0.96],
        opacity: { duration: 1 },
      }}
      className={cn(
        "absolute shadow-lg rounded-md bg-white border border-gray-200 p-3 flex flex-col gap-2",
        isPresentation 
          ? "w-32 h-24 xs:w-36 xs:h-28 sm:w-40 sm:h-32 md:w-48 md:h-36 lg:w-56 lg:h-48" 
          : "w-28 h-40 xs:w-32 xs:h-44 sm:w-36 sm:h-48 md:w-40 md:h-56 lg:w-48 lg:h-64",
        className,
      )}
    >
      <motion.div
        animate={{
          y: [0, 5, 0],
        }}
        transition={{
          duration: 8,
          repeat: Number.POSITIVE_INFINITY,
          ease: "easeInOut",
        }}
        className="relative w-full h-full"
      >
        {isPresentation ? (
          <>
            <div className="w-full h-2/3 bg-amber-100 rounded-sm mb-2.5" />
            <div className="w-3/4 h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-1/2 h-3 bg-gray-200 rounded-full" />
          </>
        ) : (
          <>
            <div className="w-full h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-3/4 h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-full h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-2/3 h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-3/4 h-3 bg-gray-200 rounded-full mb-2" />
            <div className="w-1/2 h-3 bg-gray-200 rounded-full" />
          </>
        )}
      </motion.div>
      <div
        className={cn(
          "absolute -bottom-3 -right-3 rounded-full p-2",
          isPresentation ? "bg-rose-500" : "bg-amber-500",
        )}
      >
        {isPresentation ? (
          <FilePresentation className="w-4 h-4 text-white" />
        ) : (
          <FileText className="w-4 h-4 text-white" />
        )}
      </div>
    </motion.div>
  )
}

export default FloatingDocument 