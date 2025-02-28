"use client"

import UploadFlow from "@/components/upload-flow"
import Footer from "@/components/footer"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Pacifico } from "next/font/google"
import { cn } from "@/lib/utils"
import { motion } from "framer-motion"

const pacifico = Pacifico({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-pacifico",
})

export default function StartPage() {
  const router = useRouter()

  const handleBackToHome = () => {
    router.push('/')
  }

  return (
    <div className="h-screen w-full bg-amber-50 flex flex-col overflow-hidden">
      {/* Put the logo and content in the same scrollable container */}
      <div className="flex-1 overflow-auto flex flex-col">
        <motion.div 
          initial={{ opacity: 0, x: -40 }}
          animate={{ opacity: 1, x: 0 }}
          className="flex justify-center pt-10"
        >
          <Link href="/" className="inline-block">
            <h2 className={cn(
              "text-4xl md:text-5xl lg:text-6xl font-bold",
              pacifico.className,
              "bg-clip-text text-transparent bg-gradient-to-r from-amber-500 via-orange-500 to-rose-500"
            )}>
              Slide it In
            </h2>
          </Link>
        </motion.div>
        
        <UploadFlow onBack={handleBackToHome} />
      </div>
      <Footer />
    </div>
  )
} 