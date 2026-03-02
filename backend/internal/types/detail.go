package types

// ShopDetailSupplement は詳細モーダルで追加表示する補助情報DTO。
// credit_card は一覧由来で保持する方針のため、このDTOには含めない。
type ShopDetailSupplement struct {
	Open           string `json:"open,omitempty"`
	Close          string `json:"close,omitempty"`
	BudgetMemo     string `json:"budget_memo,omitempty"`
	Wifi           string `json:"wifi,omitempty"`
	PrivateRoom    string `json:"private_room,omitempty"`
	NonSmoking     string `json:"non_smoking,omitempty"`
	Parking        string `json:"parking,omitempty"`
	Lunch          string `json:"lunch,omitempty"`
	Midnight       string `json:"midnight,omitempty"`
	ShopDetailMemo string `json:"shop_detail_memo,omitempty"`
}
