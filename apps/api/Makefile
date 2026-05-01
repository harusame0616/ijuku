.PHONY: dev migrate-up test-all test-coverage test-coverage-show

DATABASE_URL ?= postgresql://postgres:postgres@127.0.0.1:54322/postgres?sslmode=disable

SUPABASE_URL ?= http://127.0.0.1:54321
SUPABASE_ANON_KEY ?= sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH
SUPABASE_JWT_SECRET ?= super-secret-jwt-token-with-at-least-32-characters-long
WEB_ORIGIN ?= http://localhost:60000
COOKIE_SECURE ?= false

export DATABASE_URL
export SUPABASE_URL
export SUPABASE_ANON_KEY
export SUPABASE_JWT_SECRET
export WEB_ORIGIN
export COOKIE_SECURE

help:
	@echo "make dev: 開発サーバー起動（ホットリロード付き）"
	@echo "make migrate-up: マイグレーション実行"
	@echo "make test-all: 全テスト実行（カバレッジ付き）"
	@echo "make test-coverage: カバレッジが80%以上かチェック"
	@echo "make test-coverage-show: カバレッジをブラウザで表示"

dev:
	mkdir -p ./tmp && GOTMPDIR=$(CURDIR)/tmp GOCACHE=$(CURDIR)/tmp/go-cache go tool air

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

db-reset:
	pnpm dlx supabase --workdir ../../packages/supabase db reset
	$(MAKE) migrate-up

query-generate:
	sqlc generate

test-all:
	go test -coverpkg=./... -coverprofile=coverage.out ./...

test-coverage-show: test-all
	go tool cover -html=coverage.out

test-coverage:
	@go test -count=1 -coverpkg=./... -coverprofile=coverage.out ./... > /dev/null 2>&1 || (echo "NG: テスト失敗"; exit 1)
	@go tool cover -func=coverage.out | tail -1 | awk '{gsub(/%/,"",$$3); if ($$3+0 < 80) {print "NG: " $$3 "% < 80%"; exit 1} else {print "OK: " $$3 "%"}}'
