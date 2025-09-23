import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { ThemeRegistry } from './ThemeRegistry'
import { AuthProvider } from '@/lib/auth'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'StaffFind - Employee Matching System',
  description: 'Find the perfect employees for your projects',
  keywords: ['employees', 'matching', 'skills', 'hiring', 'recruitment'],
  authors: [{ name: 'StaffFind Team' }],
}

export const viewport = {
  width: 'device-width',
  initialScale: 1,
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <ThemeRegistry>
          <AuthProvider>
            {children}
          </AuthProvider>
        </ThemeRegistry>
      </body>
    </html>
  )
}
