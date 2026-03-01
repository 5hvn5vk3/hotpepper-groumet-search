package types

// GourmetSearchParams はグルメ検索APIへ渡すリクエストパラメータをまとめた struct（構造体）。
//
// 【struct とは】
// struct は複数の値（フィールド）をひとまとめにするデータ型。
// 関連する情報をひとつの「箱」に入れることで、関数の引数を整理しやすくなる。
//
// このファイルには json タグが付いていないが、これは外部APIへの送信時に
// 別レイヤー（handler/service 層）でクエリパラメータへ変換するためである。
type GourmetSearchParams struct {
	// Address は住所による絞り込みキーワード。空文字の場合は条件なし扱い。
	Address string

	// Genre はジャンルコード（例: "G001" = 居酒屋）。空文字の場合は全ジャンル対象。
	Genre string

	// Keyword はフリーワード検索文字列。
	Keyword string

	// Start は取得開始インデックス（1始まり）。ページネーションに使用する。
	// int は符号付き整数型。Go では int のビット幅は実行環境（32/64bit）に依存する。
	Start int

	// Count は1回のリクエストで取得する件数。ホットペッパーAPIの上限は100件。
	Count int

	// Lat は検索中心地点の緯度（latitude）。float64 は64ビット浮動小数点数型。
	// 緯度・経度のように小数部が必要な値には float64 を使う。
	Lat float64

	// Lng は検索中心地点の経度（longitude）。
	Lng float64

	// Range は中心地点からの検索半径（ホットペッパーAPIの定義値: 1〜5）。
	// 1=300m, 2=500m, 3=1000m, 4=2000m, 5=3000m に対応する。
	Range int
}
