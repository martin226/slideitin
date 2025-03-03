import '@/lib/promise-polyfill';
import type { Metadata } from "next";
import { GoogleAnalytics } from '@next/third-parties/google'
import localFont from "next/font/local";
import "./globals.css";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export const metadata: Metadata = {
  title: "Slide it In",
  description: "Making a presentation? Upload your documents and instantly get beautiful, presentation-ready PowerPoint slides in < 3 minutes.",
  metadataBase: new URL(process.env.NEXT_PUBLIC_URL || "https://justslideitin.com"),
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
      <GoogleAnalytics gaId="G-0TDB0J43MC" />
    </html>
  );
}
