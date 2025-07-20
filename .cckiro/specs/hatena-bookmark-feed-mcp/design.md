# はてなブックマークフィードMCP - 設計書

## システムアーキテクチャ

### 全体構成
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────┐
│   Claude Code   │────│  MCP Server      │────│ Hatena Bookmark API │
│   (Client)      │    │  (This System)   │    │   (RSS Feed)        │
└─────────────────┘    └──────────────────┘    └─────────────────────┘
```

### コンポーネント設計

#### 1. MCPサーバー（メインモジュール）
- **ファイル**: `cmd/main.go`
- **責務**: MCP接続の初期化、ツール登録、リクエスト処理の調整

#### 2. ブックマークサービス（ビジネスロジック）
- **ファイル**: `internal/service/bookmark.go`
- **責務**: はてなブックマークAPIとの通信、データ変換

#### 3. RSSパーサー（データ処理）
- **ファイル**: `internal/parser/rss.go`
- **責務**: RSS XMLの解析、構造化データへの変換

#### 4. 型定義（データモデル）
- **ファイル**: `internal/types/bookmark.go`
- **責務**: Go構造体定義、インターフェース定義

## データフロー

```
1. Claude Code → MCP Server
   - ツール呼び出し（get_hatena_bookmarks）
   - パラメータ（username, tag, date, url, page）

2. MCP Server → Bookmark Service
   - パラメータ検証
   - URLの構築

3. Bookmark Service → Hatena API
   - HTTP GETリクエスト
   - RSS XMLの取得

4. RSS Parser
   - XML → JavaScript Object
   - データの正規化

5. MCP Server → Claude Code
   - 構造化JSONレスポンス
```

## API設計

### MCPツール定義

#### `get_hatena_bookmarks`
```go
type GetHatenaBookmarksParams struct {
    Username string `json:"username"`           // 必須: はてなブックマークユーザー名
    Tag      string `json:"tag,omitempty"`      // オプション: フィルタリングタグ
    Date     string `json:"date,omitempty"`     // オプション: 日付フィルタ（YYYYMMDD）
    URL      string `json:"url,omitempty"`      // オプション: URLフィルタ
    Page     int    `json:"page,omitempty"`     // オプション: ページ番号（デフォルト: 1）
}

type GetHatenaBookmarksResponse struct {
    User       string          `json:"user"`
    Page       int             `json:"page"`
    TotalCount int             `json:"total_count"`
    Filters    *FilterParams   `json:"filters,omitempty"`
    Bookmarks  []BookmarkItem  `json:"bookmarks"`
}

type FilterParams struct {
    Tag  string `json:"tag,omitempty"`
    Date string `json:"date,omitempty"`
    URL  string `json:"url,omitempty"`
}

type BookmarkItem struct {
    Title        string   `json:"title"`
    URL          string   `json:"url"`
    BookmarkedAt string   `json:"bookmarked_at"` // ISO 8601形式
    Tags         []string `json:"tags"`
    Comment      string   `json:"comment,omitempty"`
}
```

## パッケージ設計

### BookmarkService
```go
package service

import (
    "context"
    "log/slog"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "hatena-bookmark-mcp/internal/types"
)

type BookmarkService struct {
    baseURL string
    logger  *slog.Logger
    client  *http.Client
}

func NewBookmarkService(logger *slog.Logger) *BookmarkService
func (s *BookmarkService) GetBookmarks(ctx context.Context, params types.GetHatenaBookmarksParams) (*types.GetHatenaBookmarksResponse, error)
func (s *BookmarkService) buildRequestURL(params types.GetHatenaBookmarksParams) string
func (s *BookmarkService) validateParams(params types.GetHatenaBookmarksParams) error
```

### RSSParser
```go
package parser

import (
    "context"
    "log/slog"
    "hatena-bookmark-mcp/internal/types"
)

type RSSParser struct {
    logger *slog.Logger
}

func NewRSSParser(logger *slog.Logger) *RSSParser
func (p *RSSParser) ParseRSSFeed(ctx context.Context, xmlContent []byte) (*types.ParsedRSSData, error)
func (p *RSSParser) extractBookmarkItems(channel *Channel) ([]types.BookmarkItem, error)
func (p *RSSParser) extractTags(description string) []string
func (p *RSSParser) parseDate(dateString string) (string, error)
```

### ErrorHandler
```go
package errors

import (
    "log/slog"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ErrorHandler struct {
    logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandler
func (h *ErrorHandler) HandleNetworkError(err error) *mcp.Error
func (h *ErrorHandler) HandleParsingError(err error) *mcp.Error
func (h *ErrorHandler) HandleValidationError(message string) *mcp.Error
```

## エラーハンドリング設計

### エラータイプ
1. **バリデーションエラー**: 不正なパラメータ
2. **ネットワークエラー**: API接続失敗
3. **パースエラー**: RSS解析失敗
4. **APIエラー**: はてなブックマーク側のエラー

### エラーレスポンス形式
```go
type McpError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

## セキュリティ設計

### 入力検証
- ユーザー名: 英数字とハイフンのみ許可
- タグ: 不正文字の除去
- URL: 適切なURL形式の検証
- 日付: YYYYMMDD形式の検証

### レート制限
- 同一ユーザーからの連続リクエストを制限
- 1秒間に1リクエストまで

## パフォーマンス設計

### キャッシュ戦略
- RSS取得結果を5分間キャッシュ
- 同一パラメータのリクエストはキャッシュから返却

### タイムアウト設定
- HTTP リクエスト: 10秒
- RSS解析: 5秒

## 技術スタック

### 開発環境
- **言語**: Go 1.21+
- **ビルドツール**: go build
- **パッケージマネージャー**: go mod

### 主要ライブラリ
```go
require (
    github.com/modelcontextprotocol/go-sdk v0.1.0
    golang.org/x/net v0.17.0  // HTML/XMLパーサー
)
```

### 標準ライブラリ
```go
import (
    "context"
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log/slog"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)
```

## ファイル構成

```
hatena-bookmark-mcp/
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── main.go                    # MCPサーバーメイン
├── internal/
│   ├── service/
│   │   └── bookmark.go           # ブックマーク取得サービス
│   ├── parser/
│   │   └── rss.go               # RSS解析
│   ├── types/
│   │   └── bookmark.go          # 型定義
│   ├── errors/
│   │   └── handler.go           # エラーハンドリング
│   └── utils/
│       ├── cache.go             # キャッシュ機能
│       └── validator.go         # バリデーション
├── pkg/
├── test/
│   ├── service_test.go
│   ├── parser_test.go
│   └── fixtures/
│       └── sample-rss.xml
└── README.md
```

## 設定ファイル設計

### go.mod
```go
module hatena-bookmark-mcp

go 1.21

require (
    github.com/modelcontextprotocol/go-sdk v0.1.0
    golang.org/x/net v0.17.0
)
```

### Makefile
```makefile
.PHONY: build test clean run

build:
	go build -o bin/hatena-bookmark-mcp cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin/

run:
	go run cmd/main.go

install:
	go install ./cmd/...
```

### MCP設定例
```json
{
  "mcpServers": {
    "hatena-bookmark": {
      "command": "./bin/hatena-bookmark-mcp",
      "args": [],
      "env": {
        "LOG_LEVEL": "info"
      }
    }
  }
}
```