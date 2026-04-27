# API仕様書: API キー一覧取得 (List API Keys)

## 1. 概要・エンドポイント定義

ログイン中ユーザーが発行済みの API キー一覧を取得します。鍵の平文は返さず、識別用に末尾 4 文字（plain_suffix）と作成日・有効期限のみを返します。

| 項目       | 内容                                          |
| :--------- | :-------------------------------------------- |
| **Method** | `GET`                                         |
| **Path**   | `/v1/users/{userID}/settings/apikeys`         |
| **認証**   | 要 (Supabase JWT / Bearer Token)              |

---

## 2. リクエスト

### 2.1. パスパラメータ

| パラメータ名 | 型   | 必須 | 説明                                                |
| :----------- | :--- | :--- | :-------------------------------------------------- |
| `userID`     | UUID | Yes  | 取得対象のユーザー ID。 JWT の `sub` と一致が必須。 |

### 2.2. ヘッダー

| ヘッダー名        | 値の形式            | 必須 | 説明                                       |
| :---------------- | :------------------ | :--- | :----------------------------------------- |
| **Authorization** | `Bearer {JWT}`      | Yes  | Supabase 発行の JWT (アクセストークン)。   |

リクエストボディはありません。

---

## 3. レスポンス

### 3.1. 成功時 (`200 OK`)

```json
{
  "apiKeys": [
    {
      "apiKeyId": "11111111-1111-1111-1111-111111111111",
      "suffix": "a3f9",
      "createdAt": "2026-01-10T00:00:00Z",
      "expiredAt": "2027-01-10T00:00:00Z"
    },
    {
      "apiKeyId": "22222222-2222-2222-2222-222222222222",
      "suffix": "c5d1",
      "createdAt": "2026-04-01T00:00:00Z",
      "expiredAt": null
    }
  ]
}
```

| フィールド  | 型              | 説明                                                              |
| :---------- | :-------------- | :---------------------------------------------------------------- |
| `apiKeyId`  | UUID            | API キーの識別子。                                                |
| `suffix`    | String          | API キー平文の末尾 4 文字。                                       |
| `createdAt` | String          | RFC 3339 形式の作成日時 (UTC)。                                   |
| `expiredAt` | String \| null  | RFC 3339 形式の有効期限 (UTC)。 無期限の場合は `null`。           |

API キーが 0 件の場合は `apiKeys` が空配列で返ります。

### 3.2. エラー時

| HTTP Status              | エラーコード (errorCode)  | 発生条件 / 理由                                       |
| :----------------------- | :------------------------ | :---------------------------------------------------- |
| **400 Bad Request**      | `INPUT_VALIDATION_ERROR`  | `userID` が UUID 形式でない場合。                     |
| **401 Unauthorized**     | `UNAUTHORIZED`            | Authorization ヘッダー欠落、または JWT が無効な場合。 |
| **403 Forbidden**        | `FORBIDDEN`               | JWT の `sub` がパスの `userID` と一致しない場合。     |
| **500 Internal Server Error** | `SERVER_INTERNAL_ERROR` | DB エラーなど予期しないサーバーエラー。              |

エラーレスポンス例:

```json
{
  "errorCode": "FORBIDDEN",
  "message": "forbidden"
}
```

---

## 4. 認証 (Authentication)

本エンドポイントは Supabase の JWT を用いた認証で保護されており、 [認証.md](./認証.md) で説明している API キー (Bearer) 認証とは別系統です。 JWT の `sub` クレームと URL パスの `userID` が一致する場合のみ閲覧が許可されます。

---

## 5. リクエスト例 (cURL)

```bash
curl -X GET \
  "https://api.example.com/v1/users/00000000-0000-0000-0000-000000000001/settings/apikeys" \
  -H "Authorization: Bearer eyJhbGciOi..."
```
