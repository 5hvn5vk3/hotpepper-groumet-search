package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/service"
	"backend/internal/types"
)

// ValidationError はバリデーション（入力値検証）失敗を表す独自エラー型。
// Go の error インターフェースは Error() string メソッドを持つだけでよい。
// 独自エラー型を定義することで、errors.As で型ごとに分岐処理できる。
type ValidationError struct {
	Message string
}

// Error() は error インターフェースを満たすためのメソッド。
// Go では interface の実装は明示的な "implements" 宣言が不要で、
// 必要なメソッドを持つだけで自動的に interface を実装したとみなされる（duck typing）。
func (e *ValidationError) Error() string {

	return e.Message
}

// GourmetHandler はグルメ検索エンドポイントのハンドラー構造体。
// service フィールドに HotpepperService を持つことで、
// ハンドラーとビジネスロジックを分離している（依存性の注入 / Dependency Injection）。
type GourmetHandler struct {
	service *service.HotpepperService
}

// NewGourmetHandler はコンストラクタ関数。
// Go には class/constructor がないため、慣習として New〇〇 という関数でインスタンスを生成する。
// 呼び出し側は HotpepperService を外から渡すため、テスト時にモックに差し替えやすい。
func NewGourmetHandler(service *service.HotpepperService) *GourmetHandler {

	return &GourmetHandler{service: service}
}

// Handle は HTTP リクエストを受け取り、グルメ検索を行い、結果を JSON で返すメソッド。
// シグネチャ func(http.ResponseWriter, *http.Request) は標準ライブラリの http.HandlerFunc 型と
// 一致するため、http.HandleFunc に直接渡せる。
// w（ResponseWriter）はレスポンスを書き込むための出力先、r（Request）はリクエストの全情報を持つ。
func (h *GourmetHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// GETメソッド以外のリクエストを拒否する。
	// 405 Method Not Allowed はリクエストのHTTPメソッドが許可されていないことを示す。
	if r.Method != http.MethodGet {

		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// リクエストのクエリパラメータをパース・バリデーションする。
	// バリデーション失敗時は 400 Bad Request を返す（ユーザー側の入力ミスのため）。
	params, parseErr := h.parseParams(r)

	if parseErr != nil {

		writeErrorJSON(w, http.StatusBadRequest, parseErr.Error())
		return
	}

	// service層に検索処理を委譲する。
	// service層はHTTPを知らず、純粋なビジネスロジックだけを担当する分離設計になっている。
	response, err := h.service.SearchGourmet(params)

	if err != nil {
		handleServiceError(w, err)
		return
	}

	// json.Marshal はGoの構造体・マップをJSONバイト列に変換する標準関数。
	// 変換失敗はほぼ起こらないが、チャネルや関数など変換不可能な型が含まれる場合に失敗する。
	body, err := json.Marshal(response)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		return
	}

	// レスポンスヘッダーの Content-Type を application/json にセットする。
	// ブラウザや受信側がレスポンスをJSONとして正しく解釈するために必要。
	w.Header().Set("Content-Type", "application/json")

	// 200 OK はリクエストが正常に処理されたことを示す最も基本的なステータスコード。
	w.WriteHeader(http.StatusOK)

	w.Write(body)
}

// parseParams はリクエストのクエリパラメータを取り出し、型変換・バリデーションを行って
// types.GourmetSearchParams 構造体を返す。
// バリデーション失敗時は *ValidationError を返し、ハンドラー側で 400 として扱う。
func (h *GourmetHandler) parseParams(r *http.Request) (types.GourmetSearchParams, error) {

	// r.URL.Query() はクエリ文字列をパースして url.Values（map[string][]string）を返す。
	// 存在しないキーに対して Get("key") を呼ぶと空文字 "" が返る（nil ではない）。
	query := r.URL.Query()

	// strings.TrimSpace は文字列の先頭・末尾にある空白文字（スペース・タブ・改行）を除去する。
	// ユーザー入力には意図せず空白が混入する場合があるため、安全のために正規化する。
	address := strings.TrimSpace(query.Get("address"))
	keyword := strings.TrimSpace(query.Get("keyword"))

	// lat/lng は浮動小数点数（float64）なので専用のパース関数で変換する。
	// 片方だけ指定された場合などの複合バリデーションもここで行う。
	lat, lng, parseErr := parseLatLngStrict(query.Get("lat"), query.Get("lng"))
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	// range・start・count はオプションの整数パラメータ。
	// parseOptionalIntStrict は空文字ならデフォルト値を、変換不能なら ValidationError を返す。
	rng, parseErr := parseOptionalIntStrict("range", query.Get("range"), 0)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}
	// ホットペッパーAPIの range パラメータは 1〜5 の整数のみ有効。
	// 0 は「未指定」扱いのため除外した上で範囲チェックを行う。
	if rng != 0 && (rng < 1 || rng > 5) {
		return types.GourmetSearchParams{}, &ValidationError{Message: "range must be 1, 2, 3, 4, or 5"}
	}

	start, parseErr := parseOptionalIntStrict("start", query.Get("start"), 1)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	count, parseErr := parseOptionalIntStrict("count", query.Get("count"), 20)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	// range が指定されているのに lat/lng が指定されていない場合は意味不明なリクエストなので弾く。
	hasLocation := lat != 0 && lng != 0
	if rng != 0 && !hasLocation {
		return types.GourmetSearchParams{}, &ValidationError{Message: "range requires lat and lng"}
	}

	// 検索条件が何も指定されていない場合はエラー。
	// lat/lng・address・keyword のうち少なくとも1つは必要。
	hasAddress := address != ""
	hasKeyword := keyword != ""
	if !hasLocation && !hasAddress && !hasKeyword {
		return types.GourmetSearchParams{}, &ValidationError{Message: "either lat/lng, address, or keyword must be provided"}
	}

	return types.GourmetSearchParams{
		Address: address,
		Genre:   query.Get("genre"),
		Keyword: keyword,
		Start:   start,
		Count:   count,
		Lat:     lat,
		Lng:     lng,
		Range:   rng,
	}, nil
}

// parseOptionalIntStrict はオプションの整数クエリパラメータを安全に取り出すヘルパー関数。
//   - value が空文字の場合は fallback（デフォルト値）を返す
//   - 空でない場合は strconv.Atoi で整数に変換し、失敗したら ValidationError を返す
//
// strconv.Atoi（Ascii TO Integer）は文字列を int に変換する標準関数。
// 変換に失敗すると error を返すため、必ずエラーハンドリングが必要。
func parseOptionalIntStrict(name, value string, fallback int) (int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, &ValidationError{Message: name + " must be an integer"}
	}

	return parsed, nil
}

// parseLatLngStrict は緯度（lat）・経度（lng）の文字列を float64 に変換するヘルパー関数。
// 以下のバリデーションを行う:
//   1. lat と lng は必ずセットで指定（片方だけはNG）
//   2. 両方空文字の場合は 0, 0, nil を返す（未指定として扱う）
//   3. 数値に変換できない場合は ValidationError
//   4. NaN（非数）・Inf（無限大）は不正な座標値として拒否する
func parseLatLngStrict(latStr, lngStr string) (float64, float64, error) {
	latTrimmed := strings.TrimSpace(latStr)
	lngTrimmed := strings.TrimSpace(lngStr)

	hasLat := latTrimmed != ""
	hasLng := lngTrimmed != ""

	// XOR的なチェック: どちらか一方だけ指定されている場合はエラー
	if hasLat != hasLng {
		return 0, 0, &ValidationError{Message: "lat and lng must be provided together"}
	}

	// 両方未指定の場合は「位置情報なし」として正常終了
	if !hasLat {
		return 0, 0, nil
	}

	// strconv.ParseFloat は文字列を浮動小数点数に変換する。第2引数 64 は float64 精度を指定。
	lat, err := strconv.ParseFloat(latTrimmed, 64)
	if err != nil {
		return 0, 0, &ValidationError{Message: "lat must be a number"}
	}
	// math.IsNaN は「非数（Not a Number）」かどうかを検査する。
	// "NaN" という文字列を ParseFloat すると nan が返るケースに対応。
	// math.IsInf は正または負の無限大かどうかを検査する（第2引数0は正負両方を確認）。
	// 緯度は -90〜90、経度は -180〜180 の有限値でなければならない。
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return 0, 0, &ValidationError{Message: "lat must be a finite number"}
	}

	lng, err := strconv.ParseFloat(lngTrimmed, 64)
	if err != nil {
		return 0, 0, &ValidationError{Message: "lng must be a number"}
	}
	if math.IsNaN(lng) || math.IsInf(lng, 0) {
		return 0, 0, &ValidationError{Message: "lng must be a finite number"}
	}

	return lat, lng, nil
}
