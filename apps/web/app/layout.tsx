import type { Metadata } from "next";
import { Inter_Tight } from "next/font/google";
import "./globals.css";

const interTight = Inter_Tight({
  variable: "--font-inter-tight",
  subsets: ["latin"]
})

export const metadata: Metadata = {
  title: "Weiss",
  description: "",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html className={`${interTight.variable} antialised`} lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
