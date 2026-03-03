package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// ParsePositiveIntEnv は環境変数 name を正の整数として読み込む。
// 未設定・空文字・不正な値の場合は fallback を返す。
func ParsePositiveIntEnv(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		log.Printf("invalid %s=%q; using default %d", name, value, fallback)
		return fallback
	}

	return parsed
}
