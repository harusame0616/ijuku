import * as v from 'valibot'

const supabasePublicConfigSchema = v.object({
    url: v.string(),
    anonKey: v.string(),
})

type SupabasePublicConfig = v.InferOutput<typeof supabasePublicConfigSchema>

const supabasePublicConfig = v.parse(supabasePublicConfigSchema, {
    url: process.env.NEXT_PUBLIC_SUPABASE_URL,
    anonKey: process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY,
})

export function getSupabaseConfig(): SupabasePublicConfig {
    return { ...supabasePublicConfig }
}
