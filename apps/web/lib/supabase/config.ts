import * as v from 'valibot'

const supabsePublicConfig = v.parse(v.object({
    url: v.string(),
    anonKey: v.string(),
}), {
    url: process.env.NEXT_PUBLIC_SUPABASE_URL,
    anonKey: process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY,
})

export function getSupabaseConfig() {
    return { ...supabsePublicConfig }
}
