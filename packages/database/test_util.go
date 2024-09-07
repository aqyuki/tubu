package database

import (
	"os"
	"testing"
)

// LoadConnectionStringは環境変数からテスト用DBへの接続文字列を取得する
// テスト用DBへの接続文字列はTEST_DATABASE_DSN環境変数に設定されている必要がある
func LoadConnectionString(t *testing.T) string {
	t.Helper()
	return os.Getenv("TEST_DATABASE_DSN")
}
