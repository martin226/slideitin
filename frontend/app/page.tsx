"use client"

import { useRouter } from "next/navigation"
import { cn } from "@/lib/utils"
import HeroSection from "@/components/hero-section"
import Footer from "@/components/footer"
import NotebookLine from "@/components/notebook-line"

export default function Page() {
  const router = useRouter()

  const handleUploadClick = () => {
    router.push('/start')
  }

  // Create an array of notebook lines for the background
  const notebookLines = Array.from({ length: 30 }, (_, i) => (
    <NotebookLine key={i} className={cn(i % 5 === 0 ? "bg-blue-200/50" : "")} />
  ))

  return (
    <div className="h-screen w-full bg-amber-50 flex flex-col overflow-hidden">
      {/* Notebook background only for hero section */}
      <div className="absolute inset-0 flex flex-col justify-around py-8 overflow-hidden z-0">
        {notebookLines}
      </div>
      <div className="absolute hidden sm:block left-16 md:left-24 top-0 bottom-0 w-0.5 bg-rose-400/30 z-0" />
      
      {/* Hero content */}
      <div className="flex-1 z-10 relative overflow-hidden">
        <HeroSection onUploadClick={handleUploadClick} />
      </div>
      
      {/* Footer */}
      <div className="z-10 relative">
        <Footer />
      </div>
    </div>
  )
}