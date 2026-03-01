package types

// fmt パッケージは文字列のフォーマット（整形）に使う標準ライブラリ。
// fmt.Sprintf は C 言語の sprintf に似ており、書式文字列と値を組み合わせて
// 新しい文字列を返す（標準出力には出力しない）。
import "fmt"

// HotpepperAPIError はホットペッパーAPIが返したエラーコードを保持するカスタムエラー型。
//
// Go では標準の error インターフェースを満たせばどんな struct もエラーとして扱える。
// error インターフェースは Error() string という 1 つのメソッドだけを要求する。
//
// 【なぜ struct でエラーを定義するか】
// errors.New や fmt.Errorf で作るシンプルなエラーと違い、struct にすると
// エラー固有のデータ（ここでは Code）を型安全に保持・取り出しできる。
//
// 【Code だけ持つ設計の理由】
// APIが返す人間向けメッセージをそのままログに流すと、内部情報が外部へ漏れる
// 可能性があるため、コードのみを保持してメッセージの解釈はアプリ側に委ねる。
type HotpepperAPIError struct {
	// Code はホットペッパーAPIが返したエラーコード（整数値）。
	// int は Go の基本型で、64ビット環境では 64 ビット整数として扱われる。
	Code int
}

// Error は error インターフェースを満たすためのメソッド。
// Go のインターフェースは「このメソッドを持っていれば自動的にそのインターフェース型」
// という暗黙的な実装（implicit implementation）を採用している。
//
// 【ポインタレシーバ *HotpepperAPIError について】
// `(e *HotpepperAPIError)` のように * を付けるとポインタレシーバになる。
// ポインタレシーバはレシーバのアドレスを受け取るため、フィールドの変更が
// 呼び出し元に反映される。また、大きな struct のコピーコストを避けられる。
// error インターフェースを変数に代入するときも *HotpepperAPIError 型として
// 扱われるため、型アサーション（errors.As など）で元の型を取り出せる。
//
// 【fmt.Sprintf の書式】
// %d は整数を 10 進数文字列に変換するフォーマット動詞。
// 例: Code=503 の場合 → "hotpepper API error: code=503"
func (e *HotpepperAPIError) Error() string {
	return fmt.Sprintf("hotpepper API error: code=%d", e.Code)
}
