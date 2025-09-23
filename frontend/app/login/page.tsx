'use client'

import { AuthCard } from '@/components'
import { useRouter } from 'next/navigation'

export default function LoginPage() {
  const router = useRouter()

  const handleSuccess = () => {
    router.push('/')
  }

  return <AuthCard onSuccess={handleSuccess} />
}
