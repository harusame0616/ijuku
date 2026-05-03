.PHONY: supabase-start supabase-stop supabase-reset supabase-seed


supabase-start:
	pnpm dlx supabase start

supabase-stop:
	pnpm dlx supabase stop

supabase-reset:
	pnpm dlx supabase stop
	docker volume ls -q --filter label=com.supabase.cli.project=jukubox | xargs -r docker volume rm
	pnpm dlx supabase start

supabase-seed:
	node --experimental-strip-types --no-warnings=ExperimentalWarning seed.ts
