"use client"

import { motion } from "framer-motion"
import { FileIcon as FilePresentation, ChevronRight, ChevronDown } from "lucide-react"
import FloatingDocument from "./floating-document"
import Logo from "./logo"
import { useEffect, useRef } from "react"

interface HeroSectionProps {
  onUploadClick: () => void;
}

export function HeroSection({ onUploadClick }: HeroSectionProps) {
  const fadeUpVariants = {
    hidden: { opacity: 0, y: 30 },
    visible: (i: number) => ({
      opacity: 1,
      y: 0,
      transition: {
        duration: 1,
        delay: i * 0.2,
        ease: [0.25, 0.4, 0.25, 1],
      },
    }),
  }

  // Reference for the GitHub button container
  const githubButtonRef = useRef<HTMLDivElement>(null);
  
  // Handle GitHub buttons script loading
  useEffect(() => {
    const script = document.createElement('script');
    script.src = "https://buttons.github.io/buttons.js";
    script.async = true;
    script.defer = true;
    document.body.appendChild(script);
    
    return () => {
      document.body.removeChild(script);
    };
  }, []);

  // Scroll down function for the chevron button
  const scrollToNextSection = () => {
    window.scrollTo({
      top: window.innerHeight,
      behavior: 'smooth'
    });
  };

  return (
    <div className="h-full w-full flex items-center justify-center">
      {/* Floating documents */}
      <div className="absolute inset-0 overflow-hidden">
        <FloatingDocument 
          delay={0.3} 
          rotate={-5} 
          className="hidden sm:block left-[5%] top-[10%] lg:top-[15%]" 
        />

        <FloatingDocument 
          delay={0.5} 
          rotate={5} 
          isPresentation={true} 
          className="hidden sm:block right-[8%] top-[15%] lg:top-[20%]" 
        />

        <FloatingDocument 
          delay={0.7} 
          rotate={-8} 
          className="hidden sm:block left-[12%] bottom-[10%] lg:bottom-[15%]" 
        />

        <FloatingDocument 
          delay={0.9} 
          rotate={7} 
          isPresentation={true} 
          className="hidden sm:block right-[10%] bottom-[15%] lg:bottom-[20%]" 
        />
      </div>

      <div className="relative z-10 container mx-auto px-4 md:px-6">
        <div className="max-w-4xl mx-auto text-center">
          <motion.div
            custom={0}
            variants={fadeUpVariants}
            initial="hidden"
            animate="visible"
            className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white border border-gray-200 shadow-sm mb-4 md:mb-6"
          >
            <div className="w-6 h-6 rounded-full bg-amber-500 flex items-center justify-center">
              <FilePresentation className="w-4 h-4 text-white" />
            </div>
            <span className="text-base md:text-lg text-gray-600 tracking-wide font-medium">PDF, MD, TXT to PowerPoint</span>
          </motion.div>

          <motion.div custom={1} variants={fadeUpVariants} initial="hidden" animate="visible" className="mb-2">
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight">
              <span className="bg-clip-text text-transparent bg-gradient-to-b from-gray-800 to-gray-600">
                Making a Presentation?
              </span>
            </h2>
          </motion.div>

          <motion.div custom={1.5} variants={fadeUpVariants} initial="hidden" animate="visible">
            <div className="flex flex-col sm:flex-row items-center justify-center mb-4 md:mb-6">
              <motion.span
                className="font-medium text-xl sm:text-2xl md:text-4xl text-gray-700 mb-2 sm:mb-0 sm:mr-4 px-3 py-1 bg-gradient-to-b from-amber-50 to-yellow-100 rounded-md shadow-sm transform -rotate-2 sm:-translate-y-8 border border-amber-200/70"
                initial={{ rotate: -12 }}
                animate={{ 
                  rotate: -9,
                }}
                transition={{
                  rotate: {
                    type: "spring",
                    stiffness: 70,
                    damping: 15,
                  },
                }}
                style={{
                  boxShadow: "2px 2px 4px rgba(0,0,0,0.06), 0px 1px 2px rgba(0,0,0,0.04)",
                  textShadow: "0.5px 0.5px 0px rgba(0,0,0,0.05)",
                }}
              >
                just
              </motion.span>
              <h1 className="text-4xl sm:text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight">
                <Logo size="xl" withLink={false} className="inline-block" />
              </h1>
            </div>
          </motion.div>

          <motion.div custom={2} variants={fadeUpVariants} initial="hidden" animate="visible">
            <p className="text-sm sm:text-base md:text-lg lg:text-xl text-gray-600 mb-4 sm:mb-6 leading-relaxed font-light tracking-wide max-w-2xl mx-auto px-2 sm:px-4">
              Upload your documents and instantly get beautiful, presentation-ready PowerPoint slides in &lt; 3 minutes.
            </p>
          </motion.div>

          <motion.div custom={3} variants={fadeUpVariants} initial="hidden" animate="visible">
            <button 
              onClick={onUploadClick}
              className="px-6 py-3 bg-amber-500 hover:bg-amber-600 text-white rounded-full text-base md:text-lg font-medium flex items-center gap-2 mx-auto transition-colors"
            >
              Upload Documents <ChevronRight className="w-4 h-4" />
            </button>
            
            {/* GitHub star button */}
            <div className="mt-6" ref={githubButtonRef}>
              <div className="flex justify-center items-center">
                  <a 
                    className="github-button"
                    href="https://github.com/martin226/slideitin" 
                    data-size="large" 
                    aria-label="Star martin226/slideitin on GitHub"
                  >
                    Star us on GitHub
                  </a>
              </div>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Down chevron at the bottom */}
      <motion.div
        className="absolute bottom-8 left-0 right-0 mx-auto w-10 flex justify-center cursor-pointer"
        onClick={scrollToNextSection}
        initial={{ opacity: 0 }}
        animate={{ 
          opacity: 1,
          y: [0, 10, 0],
        }}
        transition={{
          delay: 2,
          y: {
            repeat: Infinity,
            duration: 2,
            ease: "easeInOut",
          }
        }}
        whileHover={{ scale: 1.1 }}
      >
        <div className="w-10 h-10 rounded-full bg-white shadow-md flex items-center justify-center">
          <ChevronDown className="w-6 h-6 text-amber-500" />
        </div>
      </motion.div>
    </div>
  )
}

export default HeroSection