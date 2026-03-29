package types

// GourmetSearchResult は service 層が handler に返すグルメ検索結果の公開型。
// JSON シリアライズもこの型が担い、フロントエンドへのレスポンス形式を定義する。
type GourmetSearchResult struct {
	Results GourmetSearchResultData `json:"results"`
}

type GourmetSearchResultData struct {
	ResultsAvailable int    `json:"results_available"`
	ResultsStart     int    `json:"results_start"`
	Shop             []Shop `json:"shop"`
}

type Shop struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Address string `json:"address"`

	Lat float64 `json:"lat"`

	Lng float64 `json:"lng"`

	Open string `json:"open"`

	Close string `json:"close"`

	Genre ShopGenre `json:"genre"`

	Catch string `json:"catch"`

	Access string `json:"access"`

	BudgetMemo string `json:"budget_memo"`

	Wifi string `json:"wifi"`

	PrivateRoom string `json:"private_room"`

	NonSmoking string `json:"non_smoking"`

	Parking string `json:"parking"`

	Lunch string `json:"lunch"`

	Midnight string `json:"midnight"`

	ShopDetailMemo string `json:"shop_detail_memo"`

	URLs ShopURLs `json:"urls"`

	Photo ShopPhotos `json:"photo"`

	CreditCard []ShopCreditCard `json:"credit_card"`
}

type ShopGenre struct {
	Name string `json:"name"`

	Catch string `json:"catch"`
}

type ShopURLs struct {
	PC string `json:"pc"`
}

// ShopPhotos は店舗写真URLを束ねる型。
type ShopPhotos struct {
	PC ShopPhotoPC `json:"pc"`
}

type ShopPhotoPC struct {
	L string `json:"l"`

	M string `json:"m"`

	S string `json:"s"`
}

type ShopCreditCard struct {
	Code string `json:"code"`

	Name string `json:"name"`
}
