import { createBrowserClient } from '@supabase/ssr'
import { getSupabaseConfig } from './config'

export function createClient() {
  const supabsePublicConfig = getSupabaseConfig()

  return createBrowserClient(
    supabsePublicConfig.url,
    supabsePublicConfig.anonKey
  )
}
