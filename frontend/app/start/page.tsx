"use client"

import { useState } from "react"
import UploadFlow from "@/components/upload-flow"
import Footer from "@/components/footer"
import { useRouter } from "next/navigation"

export default function StartPage() {
  const router = useRouter()

  const handleBackToHome = () => {
    router.push('/')
  }

  return (
    <div className="h-screen w-full bg-amber-50 flex flex-col overflow-hidden">
      <div className="flex-1 overflow-auto">
        <UploadFlow onBack={handleBackToHome} />
      </div>
      <Footer />
    </div>
  )
} 