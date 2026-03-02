package types

import "fmt"

// HotpepperAPIError はホットペッパーAPIが返したエラーコードを保持するカスタムエラー型。
// エラーメッセージは情報漏洩防止のためこの型には持たせない。
type HotpepperAPIError struct {
	Code int
}

func (e *HotpepperAPIError) Error() string {
	return fmt.Sprintf("hotpepper API error: code=%d", e.Code)
}
