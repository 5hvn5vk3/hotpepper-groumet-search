package types

// GourmetSearchResponse はグルメサーチAPIのレスポンス全体を表す構造体
type GourmetSearchResponse struct {
	Results GourmetSearchResults `json:"results"`
}

// GourmetSearchResults は検索結果の詳細を表す構造体
type GourmetSearchResults struct {
	// APIバージョン
	APIVersion string `json:"api_version"`
	// 検索条件に該当する総件数
	ResultsAvailable int `json:"results_available"`
	// 今回のレスポンスで返却された件数（APIから文字列で返される）
	ResultsReturned string `json:"results_returned"`
	// 開始位置
	ResultsStart int `json:"results_start"`
	// 店舗情報の配列
	Shop []Shop `json:"shop"`
}

// Shop は店舗情報を表す構造体
type Shop struct {
	// 店舗ID
	ID string `json:"id"`
	// 店舗名
	Name string `json:"name"`
	// 住所
	Address string `json:"address"`
	// 緯度
	Lat float64 `json:"lat"`
	// 経度
	Lng float64 `json:"lng"`
	// 営業時間（APIからのフィールド名に合わせる）
	Open string `json:"open"`
	// ジャンル情報
	Genre ShopGenre `json:"genre"`
	// キャッチコピー
	Catch string `json:"catch"`
	// アクセス情報
	Access string `json:"access"`
	// URL情報
	URLs ShopURLs `json:"urls"`
	// 写真情報
	Photo ShopPhoto `json:"photo"`
	// クレジットカード情報（type=credit_card 追加時）
	CreditCard []ShopCreditCard `json:"credit_card"`
}

// ShopGenre は店舗のジャンル情報を表す構造体
type ShopGenre struct {
	// ジャンル名
	Name string `json:"name"`
	// ジャンルキャッチ
	Catch string `json:"catch"`
}

// ShopURLs は店舗のURL情報を表す構造体
type ShopURLs struct {
	// PC向けURL
	PC string `json:"pc"`
}

// ShopPhoto は店舗の写真情報を表す構造体
type ShopPhoto struct {
	// PC向け写真情報
	PC ShopPhotoPC `json:"pc"`
}

// ShopPhotoPC はPC向け写真の各サイズのURLを表す構造体
type ShopPhotoPC struct {
	// 大サイズ画像URL
	L string `json:"l"`
	// 中サイズ画像URL
	M string `json:"m"`
	// 小サイズ画像URL
	S string `json:"s"`
}

// ShopCreditCard は利用可能なクレジットカード情報を表す構造体
type ShopCreditCard struct {
	// カードコード
	Code string `json:"code"`
	// カード名
	Name string `json:"name"`
}

// GenreMasterResponse はジャンルマスタAPIのレスポンス全体を表す構造体
type GenreMasterResponse struct {
	Results GenreMasterResults `json:"results"`
}

// GenreMasterResults はジャンルマスタの詳細を表す構造体
type GenreMasterResults struct {
	// APIバージョン
	APIVersion string `json:"api_version"`
	// 総件数
	ResultsAvailable int `json:"results_available"`
	// 返却された件数（APIから文字列で返される）
	ResultsReturned string `json:"results_returned"`
	// 開始位置
	ResultsStart int `json:"results_start"`
	// ジャンル情報の配列
	Genre []Genre `json:"genre"`
}

// Genre はジャンル情報を表す構造体
type Genre struct {
	// ジャンルコード
	Code string `json:"code"`
	// ジャンル名
	Name string `json:"name"`
}
