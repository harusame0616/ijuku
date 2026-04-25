import { createBrowserClient } from '@supabase/ssr'
import { getSupabaseConfig } from './config'

export function createClient(): ReturnType<typeof createBrowserClient> {
  const supabsePublicConfig = getSupabaseConfig()

  return createBrowserClient(
    supabsePublicConfig.url,
    supabsePublicConfig.anonKey
  )
}
