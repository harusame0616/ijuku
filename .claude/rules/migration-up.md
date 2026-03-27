---
paths:
  - apps/api/migrations/*.up.sql
---

# MIGRATION UP

## テーブル

- \_created_at カラムと \_updated_at カラムを用意
  - 初期値は now()
- update 時に update_meta_updated_at() を呼び出すトリガーを設定
- CHECK 制約と文字数上限の制限は行わない
  - アプリケーション側で担保
- 初期値は不要
  - 常にアプリケーション側で生成
  - カラム追加時に既存のデータのマイグレーションが必要な場合のみ一時的な初期値を用意
  - マイグレーション完了後初期値は削除する
  - \_created_at と \_updated_at はアプリケーションから参照しないメタカラムのため例外
- 制約の名前を明示

## FUNCTION

- CREATE 時の OR REPLACE は不要

## EXTENSION

- CREATE 時の IF NOT EXISTS は不要

## INDEX

- CREATE INDEX は USING で種類を明示
