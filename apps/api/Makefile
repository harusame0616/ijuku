.PHONY: dev migrate-up test-all test-coverage test-coverage-show

DATABASE_URL ?= postgresql://postgres:password@localhost:5432/postgres?sslmode=disable

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

test-all:
	go test -coverpkg=./... -coverprofile=coverage.out ./...

test-coverage-show: test-all
	go tool cover -html=coverage.out

test-coverage:
	@go test -coverpkg=./... -coverprofile=coverage.out ./... > /dev/null 2>&1 || (echo "NG: テスト失敗"; exit 1)
	@go tool cover -func=coverage.out | tail -1 | awk '{gsub(/%/,"",$$3); if ($$3+0 < 80) {print "NG: " $$3 "% < 80%"; exit 1} else {print "OK: " $$3 "%"}}'
