package types

// GourmetSearchParams はグルメサーチAPIに渡すパラメータをまとめた構造体
// リクエストパラメータとして使用される
type GourmetSearchParams struct {
	// オプション：住所のキーワード（部分一致検索）
	Address string
	// オプション：ジャンルコード（例：G001は居酒屋）
	Genre string
	// オプション：フリーワード検索（店名、キャッチ等）
	Keyword string
	// ページング：検索開始位置（1始まり、1〜1000の範囲）
	Start int
	// ページング：取得件数（1〜100の範囲、デフォルト20）
	Count int
	// 緯度（lat/lng検索を行う場合に使用）
	Lat float64
	// 経度（lat/lng検索を行う場合に使用）
	Lng float64
	// 検索範囲（Hotpepper APIのrangeパラメータ、1-5で指定）
	Range int
}
