# Repository Guidelines

## Project Structure & Module Organization
- Root
  - `main_goskymoderator.go`: メインエントリ。
  - `config/`: 実行時設定（例: `config/*.yaml`）。秘密情報は置かない。
  - `go.mod` / `go.sum`: 依存関係管理。
  - `bskymoderator`: ビルド生成物（配布用）。コミット対象外を推奨。
- Package: Go標準のモジュール構成。新規コードは用途別ディレクトリへ。

## Build, Test, and Development Commands
- Build: `go build -o bskymoderator .`（実行バイナリを生成）
- Run (source): `go run ./main_goskymoderator.go`（ローカル実行）
- Test: `go test ./...`（全パッケージのテスト実行）
- Vet: `go vet ./...`（静的診断）
- Format: `gofmt -s -w .`（整形）、`goimports -w .`（import 整列）

## Coding Style & Naming Conventions
- 整形: `gofmt`必須、PR前に差分なしを保証。
- 命名: パッケージ=短く小文字、公開APIは `PascalCase`、内部は `camelCase`。
- エラー: `errors.Join/Is/As` を活用、メッセージは小文字開始・文末句点なし。
- 依存: 標準を優先、不要依存は追加しない。

## Testing Guidelines
- Framework: 標準 `testing`。テストファイルは `*_test.go`。
- 命名: `TestXxx/Subtest` を機能単位で分割。
- カバレッジ: 目安 80%（クリティカルパスを優先）。`go test -cover ./...`
- Table-driven tests を推奨、I/Oや失敗系も含める。

## Commit & Pull Request Guidelines
- Commits: 英語・Conventional風 `type(scope): subject`、50文字以内・命令形・終止記号なし。
  - Types: `feat|fix|docs|refactor|perf|test|build|ci|chore|style|revert`
  - Body: 箇条書きで要点・理由・影響範囲、72桁目安。課題連携は `Refs: #123` / `Closes: #123`
- PR: 目的・変更点・確認手順・リスクを簡潔に。必要ならスクリーンショット/ログを添付。小さく頻繁に出す。

## Security & Configuration Tips
- 秘密情報は環境変数で注入（例: `BSKY_APP_PASSWORD`）。`.gitignore`徹底。
- `config/` はテンプレート/サンプルのみをコミット。実値はローカルに保持。
- ログに個人情報を出力しない。失敗時は原因と再試行方針を記録。

