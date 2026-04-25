.PHONY: dev test-all test-coverage test-coverage-show lint

help:
	@echo "make dev: 開発サーバー起動"
	@echo "make test-all: 全テスト実行（カバレッジ付き）"
	@echo "make test-coverage: カバレッジが 80% 以上かチェック"
	@echo "make test-coverage-show: カバレッジをブラウザで表示"
	@echo "make lint: ESLint の実行"

dev:
	pnpm dev

lint:
	pnpm lint

test-all:
	pnpm vitest run --coverage

test-coverage:
	@pnpm vitest run --coverage --silent > /dev/null 2>&1 || (echo "NG: テスト失敗"; exit 1)
	@node -e "const s=require('./coverage/coverage-summary.json').total.lines.pct; if(s<80){console.log('NG: '+s+'% < 80%');process.exit(1)}else{console.log('OK: '+s+'%')}"

test-coverage-show: test-all
	open coverage/index.html
