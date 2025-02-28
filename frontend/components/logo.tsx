"use client";

import { Pacifico } from "next/font/google";
import { cn } from "@/lib/utils";
import Link from "next/link";
import { motion } from "framer-motion";

const pacifico = Pacifico({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-pacifico",
});

interface LogoProps {
  size?: "sm" | "md" | "lg" | "xl";
  withLink?: boolean;
  withAnimation?: boolean;
  className?: string;
}

export default function Logo({
  size = "md",
  withLink = true,
  withAnimation = false,
  className,
}: LogoProps) {
  const sizeClasses = {
    sm: "text-2xl md:text-3xl",
    md: "text-3xl md:text-4xl",
    lg: "text-4xl md:text-5xl lg:text-6xl",
    xl: "text-4xl sm:text-5xl md:text-7xl lg:text-8xl",
  };

  const logoContent = (
    <h2
      className={cn(
        sizeClasses[size],
        "font-bold",
        pacifico.className,
        "bg-clip-text text-transparent bg-gradient-to-r from-amber-500 via-orange-500 to-rose-500",
        className
      )}
    >
      Slide it In
    </h2>
  );

  const content = withAnimation ? (
    <motion.div
      initial={{ opacity: 0, x: -40 }}
      animate={{ opacity: 1, x: 0 }}
      className="inline-block"
    >
      {logoContent}
    </motion.div>
  ) : (
    <div className="inline-block">{logoContent}</div>
  );

  if (withLink) {
    return <Link href="/">{content}</Link>;
  }

  return content;
} 