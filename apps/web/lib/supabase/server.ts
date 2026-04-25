import { createServerClient } from '@supabase/ssr'
import { cookies } from 'next/headers'
import { getSupabaseConfig } from './config'

export async function createClient(): Promise<ReturnType<typeof createServerClient>> {
  const cookieStore = await cookies()
  const supabsePublicConfig = getSupabaseConfig()

  return createServerClient(
    supabsePublicConfig.url,
    supabsePublicConfig.anonKey,
    {
      cookies: {
        getAll() {
          return cookieStore.getAll()
        },
        setAll(cookiesToSet) {
          try {
            for (const { name, value, options } of cookiesToSet) {
              cookieStore.set(name, value, options)
            }
          } catch {
            // Server Component からの呼び出し時は無視
          }
        },
      },
    },
  )
}
