// このファイルが属するパッケージ名を宣言する。
package types

// この行で必要なパッケージを1つ読み込む。
import "fmt"

// この行は補足のためのコメントを記述する。
// HotpepperAPIError はホットペッパーAPIが返したエラーコードを保持するカスタムエラー型。
// この行は補足のためのコメントを記述する。
// エラーメッセージは情報漏洩防止のためこの型には持たせない。
// この行で構造体の型定義を始める。
type HotpepperAPIError struct {
	// この行は処理の一部として必要な操作を実行する。
	Code int
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (e *HotpepperAPIError) Error() string {
	// この時点で関数の結果を返して処理を終了する。
	return fmt.Sprintf("hotpepper API error: code=%d", e.Code)
	// ここで処理ブロックを終了する。
}
