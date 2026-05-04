import { createClient } from "@supabase/supabase-js";
import * as v from "valibot";
import { mockUser } from "./mock-user.fixture.ts";
import { seedCourses } from "./seed-course.ts";

const env = v.parse(
  v.object({
    url: v.optional(v.pipe(v.string(), v.url()), "http://127.0.0.1:54321"),
    serviceRoleKey: v.optional(v.pipe(v.string(), v.minLength(1)), 'sb_secret_N7UND0UgjKTVK-Uodkm0Hg_xSvEMPvz'),
    databaseUrl: v.optional(
      v.pipe(v.string(), v.minLength(1)),
      "postgresql://postgres:postgres@127.0.0.1:54322/postgres",
    ),
  }),
  {
    url: process.env.NEXT_PUBLIC_SUPABASE_URL ?? process.env.SUPABASE_URL,
    serviceRoleKey: process.env.SUPABASE_SERVICE_ROLE_KEY,
    databaseUrl: process.env.DATABASE_URL,
  },
);

const admin = createClient(env.url, env.serviceRoleKey, {
  auth: { persistSession: false, autoRefreshToken: false },
});

const { data, error } = await admin.auth.admin.createUser({
  email: mockUser.email,
  password: mockUser.password,
  email_confirm: true,
});

if (!error && data.user) {
  console.log(`Seeded user: ${mockUser.email} (${data.user.id})`);
} else if (error && /already|exists|registered/i.test(error.message)) {
  const { data: list, error: listError } = await admin.auth.admin.listUsers({
    perPage: 1000,
  });
  if (listError) throw listError;
  const existing = list.users.find((u) => u.email === mockUser.email);
  console.log(
    `User already exists: ${mockUser.email} (${existing?.id ?? "unknown"})`,
  );
} else {
  throw error ?? new Error("createUser が user を返しませんでした");
}

await seedCourses(env.databaseUrl);
