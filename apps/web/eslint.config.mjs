import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";
import neostandard from "neostandard";
import reactHooks from "eslint-plugin-react-hooks";
import vitest from "@vitest/eslint-plugin";
import tailwindcss from "eslint-plugin-tailwindcss";
import prettierConfig from "eslint-config-prettier";
import * as importX from "eslint-plugin-import-x";
import unusedImports from "eslint-plugin-unused-imports";
import unicorn from "eslint-plugin-unicorn";
import noRelativeImportPaths from "eslint-plugin-no-relative-import-paths";

const eslintConfig = defineConfig([
  // Next.js ルール
  ...nextVitals,
  ...nextTs,

  // ベース: neostandard（フォーマットは Prettier に委ねるため noStyle: true）
  ...neostandard({ noStyle: true, ts: true }),

  // React hooks
  reactHooks.configs.flat.recommended,

  // Tailwind CSS 品質ルール
  tailwindcss.configs["flat/recommended"],

  // インポート順序・循環依存検出
  importX.flatConfigs.recommended,
  importX.flatConfigs.typescript,

  // カスタムルール
  {
    plugins: {
      "unused-imports": unusedImports,
      "no-relative-import-paths": noRelativeImportPaths,
    },
    rules: {
      // [unused-imports] 未使用インポートを自動削除可能な形で検出する
      // @typescript-eslint/no-unused-vars と競合するため無効化して代替させる
      "@typescript-eslint/no-unused-vars": "off",
      "no-unused-vars": "off",
      "unused-imports/no-unused-imports": "error",
      "unused-imports/no-unused-vars": [
        "warn",
        {
          vars: "all",
          varsIgnorePattern: "^_",
          args: "after-used",
          argsIgnorePattern: "^_",
        },
      ],

      // [no-relative-import-paths] @/ エイリアスによる絶対インポートを強制（同一フォルダは許可）
      "no-relative-import-paths/no-relative-import-paths": [
        "error",
        { allowSameFolder: true, rootDir: ".", prefix: "@" },
      ],

      // [import-x] TypeScript コンパイラが解決するため未解決チェックを無効化
      "import-x/no-unresolved": "off",

      // [tailwindcss] クラス順序は prettier-plugin-tailwindcss に委ねるため無効化
      // no-custom-classname は Tailwind v4 の CSS ベーストークンを解決できないため無効化
      "tailwindcss/classnames-order": "off",
      "tailwindcss/no-custom-classname": "off",

      // [no-restricted-syntax] enum の代わりに union 型か const オブジェクトを使う
      "no-restricted-syntax": [
        "error",
        {
          selector: "TSEnumDeclaration",
          message:
            "enum は禁止です。union 型か const オブジェクトを使ってください。",
        },
      ],

      // [@typescript-eslint] import 文には必ず type を付ける
      "@typescript-eslint/consistent-type-imports": [
        "error",
        { prefer: "type-imports", fixStyle: "inline-type-imports" },
      ],
      // [@typescript-eslint] 型定義は interface に統一する
      "@typescript-eslint/consistent-type-definitions": ["error", "interface"],
      // [@typescript-eslint] 関数の戻り値型を明示する
      "@typescript-eslint/explicit-function-return-type": "error",
      // [@typescript-eslint] 命名規則: camelCase 基本、関数は PascalCase も許可（React コンポーネント用）
      // objectLiteralProperty は外部 API レスポンスなどの snake_case キーに対応するため無効化
      "@typescript-eslint/naming-convention": [
        "error",
        { selector: "default", format: ["camelCase"] },
        {
          selector: "variable",
          format: ["camelCase", "UPPER_CASE"],
        },
        {
          selector: "parameter",
          format: ["camelCase"],
          leadingUnderscore: "allow",
        },
        { selector: "function", format: ["camelCase", "PascalCase"] },
        { selector: "typeLike", format: ["PascalCase"] },
        { selector: "objectLiteralProperty", format: null },
        {
          selector: "import",
          format: ["camelCase", "PascalCase", "UPPER_CASE"],
        },
      ],
    },
  },

  // Unicorn: モダン JS のベストプラクティスとファイル命名
  {
    plugins: { unicorn },
    rules: {
      ...unicorn.configs["flat/recommended"].rules,
      // ファイル名は kebab-case に統一する
      "unicorn/filename-case": ["error", { cases: { kebabCase: true } }],
      // 略語の展開（props/Props は React の規約として許可）
      "unicorn/prevent-abbreviations": [
        "error",
        {
          allowList: {
            props: true,
            Props: true,
          },
        },
      ],
      // null を使う外部 API・ライブラリとの互換性のため無効化
      "unicorn/no-null": "off",
    },
  },

  // Prettier: 競合するスタイルルールを無効化するため必ず最後に配置
  prettierConfig,

  // テストファイル: Vitest プラグイン導入 + テスト向けのルール調整
  {
    files: ["**/*.test.ts", "**/*.test.tsx"],
    plugins: { vitest },
    rules: {
      ...vitest.configs.recommended.rules,
      // vitest の test.extend fixture は空のオブジェクト分割パターン `({}, provide)` を要求する
      "no-empty-pattern": "off",
      // テストでは戻り値型の明示を不要とする
      "@typescript-eslint/explicit-function-return-type": "off",
    },
  },

  // shadcn/ui コンポーネントは shadcn CLI が生成するためアップストリームの規約に従う
  {
    files: ["components/ui/**"],
    rules: {
      "@typescript-eslint/explicit-function-return-type": "off",
      "@typescript-eslint/naming-convention": "off",
    },
  },

  globalIgnores([
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "coverage/**",
    "node_modules/**",
  ]),
]);

export default eslintConfig;
