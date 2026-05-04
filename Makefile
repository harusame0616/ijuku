.DEFAULT_GOAL := help

.PHONY: help dev-api migrate-new migrate-up supabase-start supabase-stop db-reset db-seed query-generate api-test-all api-test-coverage api-test-coverage-show web-dev web-lint web-test-all web-test-coverage web-test-coverage-show spell-check

help: ## このヘルプを表示
	@awk 'BEGIN {FS = ":[^#]*?## "} /^[a-zA-Z][a-zA-Z_-]*:.*?## / {printf "make %s: %s\n", $$1, $$2}' $(firstword $(MAKEFILE_LIST))

dev-api: ## API 開発サーバー起動（ホットリロード付き）
	$(MAKE) -C apps/api dev

migrate-new: ## マイグレーションファイルを新規作成（file=<name> 必須）
	$(MAKE) -C apps/api migrate-new file=$(file)

migrate-up: ## マイグレーション実行
	$(MAKE) -C apps/api migrate-up

supabase-start: ## Supabase ローカル環境を起動
	$(MAKE) -C packages/supabase supabase-start

supabase-stop: ## Supabase ローカル環境を停止
	$(MAKE) -C packages/supabase supabase-stop

db-reset: ## DB をリセットしてマイグレーションを再適用
	$(MAKE) -C packages/supabase supabase-reset
	$(MAKE) -C apps/api migrate-up
	$(MAKE) -C packages/supabase supabase-seed

db-seed: ## Mock 認証用テストユーザーとサンプルコースを投入（Service Role キー必須）
	$(MAKE) -C packages/supabase supabase-seed

query-generate: ## sqlc クエリコードを生成
	$(MAKE) -C apps/api query-generate

api-test-all: ## API の全テスト実行（カバレッジ付き）
	$(MAKE) -C apps/api test-all

api-test-coverage: ## API のカバレッジが 80% 以上かチェック
	$(MAKE) -C apps/api test-coverage

api-test-coverage-show: ## API のカバレッジをブラウザで表示
	$(MAKE) -C apps/api test-coverage-show

web-dev: ## web 開発サーバー起動
	$(MAKE) -C apps/web dev

web-lint: ## web の ESLint 実行
	$(MAKE) -C apps/web lint

web-test-all: ## web の全テスト実行（カバレッジ付き）
	$(MAKE) -C apps/web test-all

web-test-coverage: ## web のカバレッジが 80% 以上かチェック
	$(MAKE) -C apps/web test-coverage

web-test-coverage-show: ## web のカバレッジをブラウザで表示
	$(MAKE) -C apps/web test-coverage-show

spell-check: ## スペルチェックの実行
	pnpm cspell "**"
