package types

// HotpepperGenre はジャンルマスタAPIの1件。
// "Genre" という汎用名を避け、HotPepperとの文脈を名前に含める。
type HotpepperGenre struct {
	Code string `json:"code"`

	Name string `json:"name"`
}

// GenreMasterResult は service 層が handler に返すジャンルマスタ結果の公開型。
// JSON シリアライズもこの型が担い、フロントエンドへのレスポンス形式を定義する。
type GenreMasterResult struct {
	Results GenreMasterResultData `json:"results"`
}

type GenreMasterResultData struct {
	Genre []HotpepperGenre `json:"genre"`
}
