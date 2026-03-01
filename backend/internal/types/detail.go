// このファイルが属するパッケージ名を宣言する。
package types

// この行は補足のためのコメントを記述する。
// ShopDetailSupplement は詳細モーダルで追加表示する補助情報DTO。
// この行は補足のためのコメントを記述する。
// credit_card は一覧由来で保持する方針のため、このDTOには含めない。
// この行で構造体の型定義を始める。
type ShopDetailSupplement struct {
	// この行は処理の一部として必要な操作を実行する。
	Open string `json:"open,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	Close string `json:"close,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	BudgetMemo string `json:"budget_memo,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	Wifi string `json:"wifi,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	PrivateRoom string `json:"private_room,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	NonSmoking string `json:"non_smoking,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	Parking string `json:"parking,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	Lunch string `json:"lunch,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	Midnight string `json:"midnight,omitempty"`
	// この行は処理の一部として必要な操作を実行する。
	ShopDetailMemo string `json:"shop_detail_memo,omitempty"`
	// ここで処理ブロックを終了する。
}
