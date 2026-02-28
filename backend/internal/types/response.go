package types

type GourmetSearchResponse struct {
	Results GourmetSearchResults `json:"results"`
}

type GourmetSearchResults struct {
	APIVersion string `json:"api_version"`

	ResultsAvailable int `json:"results_available"`

	ResultsReturned string `json:"results_returned"`

	ResultsStart int `json:"results_start"`

	Shop []Shop `json:"shop"`

	Error *APIError `json:"error"`
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
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

	Photo ShopPhoto `json:"photo"`

	CreditCard []ShopCreditCard `json:"credit_card"`
}

type ShopGenre struct {
	Name string `json:"name"`

	Catch string `json:"catch"`
}

type ShopURLs struct {
	PC string `json:"pc"`
}

type ShopPhoto struct {
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

type GenreMasterResponse struct {
	Results GenreMasterResults `json:"results"`
}

type GenreMasterResults struct {
	APIVersion string `json:"api_version"`

	ResultsAvailable int `json:"results_available"`

	ResultsReturned string `json:"results_returned"`

	ResultsStart int `json:"results_start"`

	Genre []Genre `json:"genre"`

	Error *APIError `json:"error"`
}

type Genre struct {
	Code string `json:"code"`

	Name string `json:"name"`
}
