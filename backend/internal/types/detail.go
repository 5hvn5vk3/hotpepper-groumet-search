package types

// ShopDetailSupplement は詳細モーダルで追加表示する補助情報 DTO（Data Transfer Object）。
//
// 【DTO とは】
// DTO はデータを層間で受け渡すための「入れ物」専用の struct。
// ビジネスロジックは持たず、データの形だけを定義する。
//
// 【omitempty オプションとは】
// json タグの値に ",omitempty" を付けると、フィールドがゼロ値（string なら空文字 ""、
// int なら 0、bool なら false など）のとき、JSON へのシリアライズ時にそのキー自体を
// 出力しない。これにより「値がない項目はレスポンス JSON に含めない」という設計を実現する。
// 例: Wifi が "" → JSON に "wifi" キーが現れない。
//
// 【credit_card を含めない理由】
// credit_card は一覧取得時（GourmetSearchResults.Shop）ですでにフロントが保持しており、
// 詳細APIの呼び出しコストを下げるため、このDTOには重複して含めない設計にしている。
type ShopDetailSupplement struct {
	// Open は営業時間の文字列（例: "月～土: 17:00～翌3:00"）。
	// omitempty により空文字のときはレスポンス JSON に含まれない。
	Open string `json:"open,omitempty"`

	// Close は定休日の文字列（例: "日曜日"）。
	Close string `json:"close,omitempty"`

	// BudgetMemo は予算の補足メモ（例: "飲み放題あり"）。
	// json タグのキー名をスネークケース（budget_memo）にすることで、
	// フロントエンドの JavaScript 慣習に合わせている。
	BudgetMemo string `json:"budget_memo,omitempty"`

	// Wifi はWifi提供状況（例: "あり" / "なし"）。
	Wifi string `json:"wifi,omitempty"`

	// PrivateRoom は個室の有無（例: "あり（2名～4名可）"）。
	PrivateRoom string `json:"private_room,omitempty"`

	// NonSmoking は禁煙・喫煙の区分（例: "全席禁煙"）。
	NonSmoking string `json:"non_smoking,omitempty"`

	// Parking は駐車場の有無。
	Parking string `json:"parking,omitempty"`

	// Lunch はランチ営業の有無。
	Lunch string `json:"lunch,omitempty"`

	// Midnight は深夜営業の有無。
	Midnight string `json:"midnight,omitempty"`

	// ShopDetailMemo は店舗の追加メモ。
	ShopDetailMemo string `json:"shop_detail_memo,omitempty"`
}
