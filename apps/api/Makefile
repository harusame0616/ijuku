.PHONY: dev

help:
	@echo "make dev:開発サーバー起動（ホットリロード付き）"

dev:
	mkdir -p ./tmp && GOTMPDIR=$(CURDIR)/tmp GOCACHE=$(CURDIR)/tmp/go-cache go tool air
