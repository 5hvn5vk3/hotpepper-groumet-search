// このファイルが属するパッケージ名を宣言する。
package types

// この行で構造体の型定義を始める。
type GourmetSearchResponse struct {
	// この行は処理の一部として必要な操作を実行する。
	Results GourmetSearchResults `json:"results"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type GourmetSearchResults struct {
	// この行は処理の一部として必要な操作を実行する。
	APIVersion string `json:"api_version"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsAvailable int `json:"results_available"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsReturned string `json:"results_returned"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsStart int `json:"results_start"`

	// この行は処理の一部として必要な操作を実行する。
	Shop []Shop `json:"shop"`

	// この行は処理の一部として必要な操作を実行する。
	Error []APIError `json:"error"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type APIError struct {
	// この行は処理の一部として必要な操作を実行する。
	Message string `json:"message"`
	// この行は処理の一部として必要な操作を実行する。
	Code int `json:"code"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type Shop struct {
	// この行は処理の一部として必要な操作を実行する。
	ID string `json:"id"`

	// この行は処理の一部として必要な操作を実行する。
	Name string `json:"name"`

	// この行は処理の一部として必要な操作を実行する。
	Address string `json:"address"`

	// この行は処理の一部として必要な操作を実行する。
	Lat float64 `json:"lat"`

	// この行は処理の一部として必要な操作を実行する。
	Lng float64 `json:"lng"`

	// この行は処理の一部として必要な操作を実行する。
	Open string `json:"open"`

	// この行は処理の一部として必要な操作を実行する。
	Close string `json:"close"`

	// この行は処理の一部として必要な操作を実行する。
	Genre ShopGenre `json:"genre"`

	// この行は処理の一部として必要な操作を実行する。
	Catch string `json:"catch"`

	// この行は処理の一部として必要な操作を実行する。
	Access string `json:"access"`

	// この行は処理の一部として必要な操作を実行する。
	BudgetMemo string `json:"budget_memo"`

	// この行は処理の一部として必要な操作を実行する。
	Wifi string `json:"wifi"`

	// この行は処理の一部として必要な操作を実行する。
	PrivateRoom string `json:"private_room"`

	// この行は処理の一部として必要な操作を実行する。
	NonSmoking string `json:"non_smoking"`

	// この行は処理の一部として必要な操作を実行する。
	Parking string `json:"parking"`

	// この行は処理の一部として必要な操作を実行する。
	Lunch string `json:"lunch"`

	// この行は処理の一部として必要な操作を実行する。
	Midnight string `json:"midnight"`

	// この行は処理の一部として必要な操作を実行する。
	ShopDetailMemo string `json:"shop_detail_memo"`

	// この行は処理の一部として必要な操作を実行する。
	URLs ShopURLs `json:"urls"`

	// この行は処理の一部として必要な操作を実行する。
	Photo ShopPhoto `json:"photo"`

	// この行は処理の一部として必要な操作を実行する。
	CreditCard []ShopCreditCard `json:"credit_card"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type ShopGenre struct {
	// この行は処理の一部として必要な操作を実行する。
	Name string `json:"name"`

	// この行は処理の一部として必要な操作を実行する。
	Catch string `json:"catch"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type ShopURLs struct {
	// この行は処理の一部として必要な操作を実行する。
	PC string `json:"pc"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type ShopPhoto struct {
	// この行は処理の一部として必要な操作を実行する。
	PC ShopPhotoPC `json:"pc"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type ShopPhotoPC struct {
	// この行は処理の一部として必要な操作を実行する。
	L string `json:"l"`

	// この行は処理の一部として必要な操作を実行する。
	M string `json:"m"`

	// この行は処理の一部として必要な操作を実行する。
	S string `json:"s"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type ShopCreditCard struct {
	// この行は処理の一部として必要な操作を実行する。
	Code string `json:"code"`

	// この行は処理の一部として必要な操作を実行する。
	Name string `json:"name"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type GenreMasterResponse struct {
	// この行は処理の一部として必要な操作を実行する。
	Results GenreMasterResults `json:"results"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type GenreMasterResults struct {
	// この行は処理の一部として必要な操作を実行する。
	APIVersion string `json:"api_version"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsAvailable int `json:"results_available"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsReturned string `json:"results_returned"`

	// この行は処理の一部として必要な操作を実行する。
	ResultsStart int `json:"results_start"`

	// この行は処理の一部として必要な操作を実行する。
	Genre []Genre `json:"genre"`

	// この行は処理の一部として必要な操作を実行する。
	Error []APIError `json:"error"`
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type Genre struct {
	// この行は処理の一部として必要な操作を実行する。
	Code string `json:"code"`

	// この行は処理の一部として必要な操作を実行する。
	Name string `json:"name"`
	// ここで処理ブロックを終了する。
}
