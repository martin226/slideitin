"use client"

import UploadFlow from "@/components/upload-flow"
import Footer from "@/components/footer"
import { useRouter } from "next/navigation"
import Logo from "@/components/logo"

export default function StartPage() {
  const router = useRouter()

  const handleBackToHome = () => {
    router.push('/')
  }

  return (
    <div className="h-screen w-full bg-amber-50 flex flex-col overflow-hidden">
      {/* Put the logo and content in the same scrollable container */}
      <div className="flex-1 overflow-auto flex flex-col">
        <div className="flex justify-center pt-10">
          <Logo size="lg" withAnimation={true} />
        </div>
        
        <div className="flex-1">
          <UploadFlow onBack={handleBackToHome} />
        </div>
        <Footer />
      </div>
    </div>
  )
} 