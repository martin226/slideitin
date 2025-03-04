"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import Footer from "@/components/footer";
import Logo from "@/components/logo";
import { ArrowLeft } from "lucide-react";

const container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.3,
    },
  },
};

const item = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 },
};

export default function Privacy() {
  return (
    <div className="min-h-screen flex flex-col bg-amber-50">
        <div className="flex justify-center pt-10">
          <Logo size="lg" withAnimation={true} />
        </div>

      <main className="container mx-auto px-6 pt-8 pb-16">
        <motion.div className="max-w-4xl mx-auto" initial="hidden" animate="show" variants={item}>
          <Link
          href="/"
          className="w-fit px-4 py-3 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center gap-2 transition-colors"
        >
          <ArrowLeft className="h-4 w-4" /> Back to Home
        </Link>
        </motion.div>
        
        <motion.div 
          className="max-w-4xl mx-auto py-8 md:py-12"
          initial="hidden"
          animate="show"
          variants={container}
        >
          <motion.div variants={item}>
            <h1 className="font-playfair text-5xl md:text-6xl lg:text-7xl font-bold text-center mb-4">
              Privacy Policy
            </h1>
            <p className="text-center text-muted-foreground text-lg md:text-xl mb-12 max-w-2xl mx-auto">
              We believe privacy policies should be clear, honest, and actually readable. Here&apos;s ours.
            </p>
          </motion.div>

          <motion.div variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 mb-16 hover:shadow-md transition-all duration-300 border-2 border-amber-400/30">
            <h2 className="text-2xl md:text-3xl font-bold mb-6">Key Points</h2>
            <ul className="space-y-4">
              <li className="flex items-start gap-4">
                <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <span className="text-lg">We process your uploaded content through Gemini&apos;s API to transform your documents into beautiful slide decks.</span>
              </li>
              <li className="flex items-start gap-4">
                <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="text-lg">Your uploaded documents are temporarily cached and deleted permanently within one day.</span>
              </li>
              <li className="flex items-start gap-4">
                <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="text-lg">We prioritize your privacy and never sell, rent, or store your data beyond what is necessary.</span>
              </li>
            </ul>
          </motion.div>

          <motion.div 
            className="space-y-12 md:space-y-16"
            variants={container}
          >
            {/* Introduction */}
            <motion.section variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 hover:shadow-md transition-all duration-300 border border-slate-100">
              <div className="flex items-start gap-4 mb-6">
                <div className="bg-amber-500/10 rounded-lg p-3">
                  <svg className="w-6 h-6 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h2 className="text-2xl md:text-3xl font-bold">Introduction</h2>
              </div>
              <div className="prose prose-slate lg:prose-lg">
                <p className="text-muted-foreground leading-relaxed">
                  Welcome to Slide it In (&quot;we,&quot; &quot;our,&quot; or &quot;us&quot;). We value your privacy and are committed to protecting your personal information. 
                  This Privacy Policy explains how we collect, use, and protect the information you provide when using our service at{" "}
                  <a href="https://justslideitin.com" className="text-amber-500 hover:text-amber-600 no-underline border-b-2 border-amber-500/20 hover:border-amber-600/90 transition-colors">
                    justslideitin.com
                  </a>.
                </p>
              </div>
            </motion.section>

            {/* Information We Collect */}
            <motion.section variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 hover:shadow-md transition-all duration-300 border border-slate-100">
              <div className="flex items-start gap-4 mb-6">
                <div className="bg-amber-500/10 rounded-lg p-3">
                  <svg className="w-6 h-6 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <h2 className="text-2xl md:text-3xl font-bold">Information We Collect</h2>
              </div>
              <div className="prose prose-slate lg:prose-lg">
                <p className="text-muted-foreground leading-relaxed">When you use our service to transform documents into slide decks, we collect:</p>
                <ul className="space-y-4 mt-4">
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="text-muted-foreground">Document Content: The files you upload to our service for transformation into presentations</span>
                  </li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mt-4">
                  We do not collect personal information like email addresses or require account creation to use our service.
                </p>
              </div>
            </motion.section>

            {/* How We Use Your Information */}
            <motion.section variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 hover:shadow-md transition-all duration-300 border border-slate-100">
              <div className="flex items-start gap-4 mb-6">
                <div className="bg-amber-500/10 rounded-lg p-3">
                  <svg className="w-6 h-6 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                </div>
                <h2 className="text-2xl md:text-3xl font-bold">How We Use Your Information</h2>
              </div>
              <div className="prose prose-slate lg:prose-lg">
                <p className="text-muted-foreground leading-relaxed">The primary use of your uploaded document content is to transform it into slide decks. Specifically:</p>
                <ul className="space-y-4 mt-4">
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                    <span className="text-muted-foreground">Your document is processed through Google&apos;s Gemini 1.5 Flash API to analyze and extract presentation content</span>
                  </li>
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                    <span className="text-muted-foreground">The transformed presentation is made available for download by you</span>
                  </li>
                </ul>
              </div>
            </motion.section>

            {/* Data Retention and Deletion */}
            <motion.section variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 hover:shadow-md transition-all duration-300 border border-slate-100">
              <div className="flex items-start gap-4 mb-6">
                <div className="bg-amber-500/10 rounded-lg p-3">
                  <svg className="w-6 h-6 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h2 className="text-2xl md:text-3xl font-bold">Data Retention and Deletion</h2>
              </div>
              <div className="prose prose-slate lg:prose-lg">
                <ul className="space-y-4">
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span className="text-muted-foreground">Uploaded documents are cached temporarily and never stored for more than one day</span>
                  </li>
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span className="text-muted-foreground">After one day, all uploaded documents are permanently deleted from our servers</span>
                  </li>
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span className="text-muted-foreground">No backups or copies of your documents are stored beyond this period</span>
                  </li>
                </ul>
              </div>
            </motion.section>

            {/* Security Measures */}
            <motion.section variants={item} className="bg-amber-50 rounded-2xl shadow-sm p-8 md:p-10 hover:shadow-md transition-all duration-300 border border-slate-100">
              <div className="flex items-start gap-4 mb-6">
                <div className="bg-amber-500/10 rounded-lg p-3">
                  <svg className="w-6 h-6 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                </div>
                <h2 className="text-2xl md:text-3xl font-bold">Security Measures</h2>
              </div>
              <div className="prose prose-slate lg:prose-lg">
                <p className="text-muted-foreground leading-relaxed">We implement industry-standard security measures to protect your data, including:</p>
                <ul className="space-y-4 mt-4">
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    <span className="text-muted-foreground">Encrypted transmission of uploaded documents</span>
                  </li>
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    <span className="text-muted-foreground">Secure deletion processes to ensure documents are erased after one day</span>
                  </li>
                  <li className="flex items-start gap-4">
                    <svg className="w-6 h-6 text-amber-500 flex-shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    <span className="text-muted-foreground">Transport Layer Security (TLS) to protect data in transit</span>
                  </li>
                </ul>
              </div>
            </motion.section>
          </motion.div>
        </motion.div>
      </main>
      <Footer />
    </div>
  );
} 