"use client"

import { cn } from "@/lib/utils"

export function NotebookLine({ className }: { className?: string }) {
  return <div className={cn("w-full h-px bg-blue-200/30", className)} />
}

export default NotebookLine 