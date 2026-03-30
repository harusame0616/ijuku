.PHONY: dev-api migrate-up api-test-all api-test-coverage api-test-coverage-show

help:
	@echo "make dev-api: API 開発サーバー起動（ホットリロード付き）"
	@echo "make migrate-up: マイグレーション実行"
	@echo "make api-test-all: 全テスト実行（カバレッジ付き）"
	@echo "make api-test-coverage: カバレッジが80%以上かチェック"
	@echo "make api-test-coverage-show: カバレッジをブラウザで表示"

dev-api:
	$(MAKE) -C apps/api dev

migrate-up:
	$(MAKE) -C apps/api migrate-up

api-test-all:
	$(MAKE) -C apps/api test-all

api-test-coverage:
	$(MAKE) -C apps/api test-coverage

api-test-coverage-show:
	$(MAKE) -C apps/api test-coverage-show
