package types

// GourmetSearchParams はグルメサーチAPIに渡すパラメータをまとめた構造体
// リクエストパラメータとして使用される
type GourmetSearchParams struct {
	// 必須：都道府県コード（例：SA11は東京、SA23は大阪）
	ServiceArea string
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
}
