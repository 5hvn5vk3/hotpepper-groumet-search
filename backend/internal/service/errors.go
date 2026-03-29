package service

import "fmt"

// HotpepperAPIError はホットペッパーAPIが返したエラーコードを保持するカスタムエラー型。
// HotPepperという特定外部サービス固有のエラーであるため service パッケージに定義する。
// エラーメッセージは情報漏洩防止のためこの型には持たせない。
type HotpepperAPIError struct {
	Code int
}

func (e *HotpepperAPIError) Error() string {
	return fmt.Sprintf("hotpepper API error: code=%d", e.Code)
}
