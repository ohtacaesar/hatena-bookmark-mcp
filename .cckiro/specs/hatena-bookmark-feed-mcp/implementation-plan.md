# はてなブックマークフィードMCP - 実装計画書

## 実装の全体フロー

### Phase 1: プロジェクト基盤構築
1. **プロジェクト初期化**
   - `go.mod` ファイルの作成
   - ディレクトリ構造の構築
   - 依存関係の追加

2. **基本設定ファイル**
   - `Makefile` の作成
   - `.gitignore` の作成

### Phase 2: 型定義とインターフェース
1. **型定義ファイル作成**
   - `internal/types/bookmark.go` - データ構造体の定義
   - リクエスト/レスポンス型の実装

2. **インターフェース定義**
   - サービス層のインターフェース
   - パーサー層のインターフェース

### Phase 3: コアロジック実装
1. **RSSパーサー実装**
   - `internal/parser/rss.go`
   - XML解析ロジック
   - タグ抽出ロジック
   - 日付パースロジック

2. **ブックマークサービス実装**
   - `internal/service/bookmark.go`
   - HTTP クライアント
   - URL構築ロジック
   - パラメータ検証

3. **エラーハンドリング**
   - `internal/errors/handler.go`
   - エラー分類と処理
   - ログ出力

### Phase 4: ユーティリティ機能
1. **バリデーション**
   - `internal/utils/validator.go`
   - 入力パラメータ検証

2. **キャッシュ機能（オプション）**
   - `internal/utils/cache.go`
   - メモリキャッシュ実装

### Phase 5: MCPサーバー実装
1. **メインサーバー**
   - `cmd/main.go`
   - MCP SDK統合
   - ツール登録
   - リクエストハンドリング

2. **ログ設定**
   - slog設定
   - ログレベル制御

### Phase 6: テスト実装
1. **ユニットテスト**
   - パーサーテスト
   - サービステスト
   - バリデーションテスト

2. **統合テスト**
   - エンドツーエンドテスト
   - モックデータを使用した動作確認

## 実装順序と詳細

### 1. プロジェクト基盤構築

#### 1.1 ディレクトリ作成
```bash
mkdir -p hatena-bookmark-mcp/{cmd,internal/{types,service,parser,errors,utils},test/fixtures}
```

#### 1.2 go.mod 初期化
```bash
cd hatena-bookmark-mcp
go mod init hatena-bookmark-mcp
```

#### 1.3 依存関係追加
```bash
go get github.com/modelcontextprotocol/go-sdk@latest
go get golang.org/x/net@latest
```

### 2. 型定義実装

#### 2.1 `internal/types/bookmark.go`
```go
// 主要な構造体の実装
- GetHatenaBookmarksParams
- GetHatenaBookmarksResponse  
- BookmarkItem
- FilterParams
- RSS構造体 (Channel, Item等)
```

### 3. RSSパーサー実装

#### 3.1 `internal/parser/rss.go`
```go
// 実装内容:
1. XML → Go構造体への変換
2. dc:subject からタグ抽出
3. 日付文字列のISO8601変換
4. エラーハンドリング付きパース処理
```

**実装手順:**
1. RSS XMLの構造体定義
2. `encoding/xml` を使ったパース処理
3. タグ抽出ロジック（正規表現使用）
4. 日付変換ロジック
5. エラーハンドリング

### 4. ブックマークサービス実装

#### 4.1 `internal/service/bookmark.go`
```go
// 実装内容:
1. HTTP クライアント設定
2. URL構築（クエリパラメータ含む）
3. リクエスト送信と レスポンス処理
4. パラメータ検証
5. タイムアウト設定
```

**実装手順:**
1. サービス構造体の初期化
2. URL構築メソッド
3. HTTP リクエスト処理
4. レスポンス処理とエラーハンドリング
5. パラメータバリデーション

### 5. エラーハンドリング実装

#### 5.1 `internal/errors/handler.go`
```go
// 実装内容:
1. エラータイプの分類
2. MCPエラー形式への変換
3. ログ出力
4. エラーメッセージの統一
```

### 6. MCPサーバー実装

#### 6.1 `cmd/main.go`
```go
// 実装内容:
1. slogロガー初期化
2. MCP サーバー初期化
3. ツール登録 (get_hatena_bookmarks)
4. リクエストハンドラー実装
5. サーバー起動
```

**実装手順:**
1. 依存関係注入の設定
2. MCPツール定義
3. ハンドラー関数実装
4. サーバー起動処理

### 7. ユーティリティ実装

#### 7.1 `internal/utils/validator.go`
```go
// 実装内容:
1. ユーザー名検証
2. 日付形式検証
3. URL形式検証
4. ページ番号検証
```

#### 7.2 `internal/utils/cache.go` (オプション)
```go
// 実装内容:
1. インメモリキャッシュ
2. TTL設定
3. キー生成
```

### 8. テスト実装

#### 8.1 テストファイル作成
```
test/
├── service_test.go      # サービス層テスト
├── parser_test.go       # パーサーテスト
├── validator_test.go    # バリデーションテスト
└── fixtures/
    └── sample-rss.xml   # テスト用RSSデータ
```

#### 8.2 テスト項目
1. **正常系テスト**
   - 基本的なブックマーク取得
   - フィルタリング機能
   - ページネーション

2. **異常系テスト**
   - 不正なパラメータ
   - ネットワークエラー
   - 不正なXML形式

## 開発環境セットアップ

### 必要なツール
```bash
# Go 1.21+のインストール確認
go version

# 開発用ツール
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
```

### 実行環境テスト
```bash
# ビルド確認
make build

# テスト実行
make test

# 動作確認
make run
```

## 実装時の注意事項

### 1. エラーハンドリング
- すべてのHTTPリクエストでタイムアウト設定
- XML解析エラーの適切な処理
- ログ出力での情報漏洩防止

### 2. パフォーマンス
- HTTP クライアントの再利用
- 不要なメモリ使用の回避
- 適切なバッファサイズの設定

### 3. セキュリティ
- 入力値の厳格な検証
- ログでのセンシティブ情報マスク
- レート制限の実装

### 4. 保守性
- 適切なコメント記述
- 一貫したエラーメッセージ
- 設定値の外部化

## 完成目標

### 最小機能（MVP）
1. 基本的なブックマーク取得
2. MCPツールとしての動作
3. エラーハンドリング

### 拡張機能
1. 全フィルタリング機能
2. キャッシュ機能
3. 詳細なログ出力
4. 包括的なテストカバレッジ

## 想定実装期間
- **Phase 1-2**: 基盤構築・型定義 (1-2時間)
- **Phase 3-4**: コアロジック実装 (3-4時間)  
- **Phase 5**: MCPサーバー実装 (2-3時間)
- **Phase 6**: テスト実装 (2-3時間)
- **総計**: 8-12時間

各フェーズの完了後に動作確認を行い、次のフェーズに進む段階的な開発を実施します。