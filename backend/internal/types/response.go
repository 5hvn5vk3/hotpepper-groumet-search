package types

// ===== ホットペッパーグルメ検索API レスポンス型群 =====
//
// 【json タグとは】
// struct のフィールドに `json:"キー名"` と書くと、encoding/json パッケージが
// JSON ↔ Go struct を変換（Marshal/Unmarshal）するときのキー名を指定できる。
// 例: `json:"api_version"` → JSON の "api_version" キーを APIVersion フィールドに対応させる。
// タグがない場合はフィールド名がそのまま使われる（大文字小文字は区別しない）。

// GourmetSearchResponse はグルメ検索APIのトップレベルレスポンス。
// ホットペッパーAPIはすべての結果を "results" キーにネストして返すため、
// それを受け取るラッパー struct として定義している。
type GourmetSearchResponse struct {
	// Results は実際の検索結果をまとめた struct。
	// フィールド型に別の struct を指定することで、JSON のネスト構造を表現できる。
	Results GourmetSearchResults `json:"results"`
}

// GourmetSearchResults はグルメ検索結果の本体。
// ページネーション情報・店舗一覧・エラー一覧を含む。
type GourmetSearchResults struct {
	// APIVersion はホットペッパーAPIのバージョン文字列。
	APIVersion string `json:"api_version"`

	// ResultsAvailable は検索条件に一致する総件数（ページングに関わらず全件数）。
	ResultsAvailable int `json:"results_available"`

	// ResultsReturned は今回のレスポンスに含まれる実際の件数。
	// ホットペッパーAPIは int ではなく文字列で返すため string 型にしている。
	ResultsReturned string `json:"results_returned"`

	// ResultsStart は今回の結果が何件目から始まるか（1始まり）。ページネーション用。
	ResultsStart int `json:"results_start"`

	// Shop は店舗情報のスライス（可変長配列）。
	// Go では配列長が固定の「配列」と、長さが可変の「スライス」を区別する。
	// []Shop のように [] を付けるとスライス型になる。
	Shop []Shop `json:"shop"`

	// Error はAPIが返したエラー情報のスライス。
	// 正常時は空スライスになる。
	Error []APIError `json:"error"`
}

// APIError はホットペッパーAPIが返すエラーオブジェクト。
// errors.go の HotpepperAPIError とは別物で、こちらは JSON デシリアライズ用の DTO。
type APIError struct {
	// Message は人間向けのエラーメッセージ文字列。
	Message string `json:"message"`
	// Code は機械向けのエラーコード番号。
	Code    int    `json:"code"`
}

// Shop は1店舗分の情報を表す struct。
// ホットペッパーAPIの shop オブジェクトに対応する。
type Shop struct {
	// ID は店舗を一意に識別する文字列（例: "J001234567"）。
	ID string `json:"id"`

	// Name は店舗名。
	Name string `json:"name"`

	// Address は店舗の住所文字列。
	Address string `json:"address"`

	// Lat は店舗の緯度。float64 は小数を扱える64ビット浮動小数点数型。
	Lat float64 `json:"lat"`

	// Lng は店舗の経度。
	Lng float64 `json:"lng"`

	// Open は営業時間の文字列表現（例: "月～土: 17:00～翌3:00"）。
	Open string `json:"open"`

	// Close は定休日の文字列表現（例: "日曜日"）。
	Close string `json:"close"`

	// Genre は店舗のジャンル情報。別 struct で名前とキャッチコピーを持つ。
	Genre ShopGenre `json:"genre"`

	// Catch は店舗のキャッチコピー文字列。
	Catch string `json:"catch"`

	// Access は最寄り駅などのアクセス情報。
	Access string `json:"access"`

	// BudgetMemo は予算の補足メモ（例: "飲み放題あり"）。
	BudgetMemo string `json:"budget_memo"`

	// Wifi はWifi提供状況（例: "あり" / "なし"）。
	Wifi string `json:"wifi"`

	// PrivateRoom は個室の有無（例: "あり（2名～4名可）"）。
	PrivateRoom string `json:"private_room"`

	// NonSmoking は禁煙・喫煙の区分（例: "全席禁煙"）。
	NonSmoking string `json:"non_smoking"`

	// Parking は駐車場の有無。
	Parking string `json:"parking"`

	// Lunch はランチ営業の有無。
	Lunch string `json:"lunch"`

	// Midnight は深夜営業の有無。
	Midnight string `json:"midnight"`

	// ShopDetailMemo は店舗の追加メモ。
	ShopDetailMemo string `json:"shop_detail_memo"`

	// URLs は店舗ページの URL 情報。
	URLs ShopURLs `json:"urls"`

	// Photo は店舗写真の URL 情報。
	Photo ShopPhoto `json:"photo"`

	// CreditCard は利用可能なクレジットカードのスライス。
	CreditCard []ShopCreditCard `json:"credit_card"`
}

// ShopGenre は店舗ジャンルの名称とキャッチコピーを持つ struct。
type ShopGenre struct {
	// Name はジャンル名（例: "居酒屋"）。
	Name string `json:"name"`

	// Catch はジャンルのキャッチコピー（例: "もつ鍋・水炊き"）。
	Catch string `json:"catch"`
}

// ShopURLs は店舗ページの URL をまとめた struct。
type ShopURLs struct {
	// PC はPC向け店舗ページの URL。
	PC string `json:"pc"`
}

// ShopPhoto は店舗写真の URL 情報をまとめた struct。
// ホットペッパーAPIは PC 向けと mobile 向けを分けて返すが、ここでは PC のみ使用。
type ShopPhoto struct {
	// PC はPC向け写真 URL のセット。
	PC ShopPhotoPC `json:"pc"`
}

// ShopPhotoPC はPC向け写真を大・中・小の3サイズ分保持する struct。
type ShopPhotoPC struct {
	// L は大サイズ（large）の写真 URL。
	L string `json:"l"`

	// M は中サイズ（medium）の写真 URL。
	M string `json:"m"`

	// S は小サイズ（small）の写真 URL。一覧サムネイルに使用する。
	S string `json:"s"`
}

// ShopCreditCard は1種類のクレジットカード情報を表す struct。
type ShopCreditCard struct {
	// Code はカードの識別コード（例: "VISA"）。
	Code string `json:"code"`

	// Name はカードの表示名（例: "VISA"）。
	Name string `json:"name"`
}

// ===== ジャンルマスターAPI レスポンス型群 =====

// GenreMasterResponse はジャンルマスターAPIのトップレベルレスポンス。
type GenreMasterResponse struct {
	// Results は実際のジャンル一覧データをまとめた struct。
	Results GenreMasterResults `json:"results"`
}

// GenreMasterResults はジャンルマスターAPIの結果本体。
type GenreMasterResults struct {
	// APIVersion はAPIバージョン文字列。
	APIVersion string `json:"api_version"`

	// ResultsAvailable は利用可能なジャンルの総数。
	ResultsAvailable int `json:"results_available"`

	// ResultsReturned は今回返却されたジャンル数（文字列型）。
	ResultsReturned string `json:"results_returned"`

	// ResultsStart は返却開始インデックス（1始まり）。
	ResultsStart int `json:"results_start"`

	// Genre はジャンル情報のスライス。
	Genre []Genre `json:"genre"`

	// Error はAPIエラー情報のスライス。正常時は空。
	Error []APIError `json:"error"`
}

// Genre は1ジャンル分のコードと名称を持つ struct。
type Genre struct {
	// Code はジャンル識別コード（例: "G001"）。検索パラメータに使用する。
	Code string `json:"code"`

	// Name はジャンルの表示名（例: "居酒屋"）。UI に表示する。
	Name string `json:"name"`
}
