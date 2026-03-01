// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"math"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"strconv"
	// この行は処理の一部として必要な操作を実行する。
	"strings"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/service"
	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で構造体の型定義を始める。
type ValidationError struct {
	// この行は処理の一部として必要な操作を実行する。
	Message string
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (e *ValidationError) Error() string {

	// この時点で関数の結果を返して処理を終了する。
	return e.Message
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type GourmetHandler struct {
	// この行は処理の一部として必要な操作を実行する。
	service *service.HotpepperService
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func NewGourmetHandler(service *service.HotpepperService) *GourmetHandler {

	// この時点で関数の結果を返して処理を終了する。
	return &GourmetHandler{service: service}
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (h *GourmetHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// 条件を評価し、真のときだけ次の処理を実行する。
	if r.Method != http.MethodGet {

		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	params, parseErr := h.parseParams(r)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if parseErr != nil {

		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusBadRequest, parseErr.Error())
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	response, err := h.service.SearchGourmet(params)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この行は処理の一部として必要な操作を実行する。
		handleServiceError(w, err)
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := json.Marshal(response)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.Header().Set("Content-Type", "application/json")

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.WriteHeader(http.StatusOK)

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.Write(body)
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (h *GourmetHandler) parseParams(r *http.Request) (types.GourmetSearchParams, error) {

	// この行で新しい変数を宣言しつつ値を代入する。
	query := r.URL.Query()

	// この行で新しい変数を宣言しつつ値を代入する。
	address := strings.TrimSpace(query.Get("address"))
	// この行で新しい変数を宣言しつつ値を代入する。
	keyword := strings.TrimSpace(query.Get("keyword"))

	// この行で新しい変数を宣言しつつ値を代入する。
	lat, lng, parseErr := parseLatLngStrict(query.Get("lat"), query.Get("lng"))
	// 条件を評価し、真のときだけ次の処理を実行する。
	if parseErr != nil {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, parseErr
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	rng, parseErr := parseOptionalIntStrict("range", query.Get("range"), 0)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if parseErr != nil {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, parseErr
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if rng != 0 && (rng < 1 || rng > 5) {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, &ValidationError{Message: "range must be 1, 2, 3, 4, or 5"}
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	start, parseErr := parseOptionalIntStrict("start", query.Get("start"), 1)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if parseErr != nil {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, parseErr
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	count, parseErr := parseOptionalIntStrict("count", query.Get("count"), 20)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if parseErr != nil {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, parseErr
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	hasLocation := lat != 0 && lng != 0
	// 条件を評価し、真のときだけ次の処理を実行する。
	if rng != 0 && !hasLocation {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, &ValidationError{Message: "range requires lat and lng"}
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	hasAddress := address != ""
	// この行で新しい変数を宣言しつつ値を代入する。
	hasKeyword := keyword != ""
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !hasLocation && !hasAddress && !hasKeyword {
		// この時点で関数の結果を返して処理を終了する。
		return types.GourmetSearchParams{}, &ValidationError{Message: "either lat/lng, address, or keyword must be provided"}
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return types.GourmetSearchParams{
		// この行は処理の一部として必要な操作を実行する。
		Address: address,
		// この行は処理の一部として必要な操作を実行する。
		Genre: query.Get("genre"),
		// この行は処理の一部として必要な操作を実行する。
		Keyword: keyword,
		// この行は処理の一部として必要な操作を実行する。
		Start: start,
		// この行は処理の一部として必要な操作を実行する。
		Count: count,
		// この行は処理の一部として必要な操作を実行する。
		Lat: lat,
		// この行は処理の一部として必要な操作を実行する。
		Lng: lng,
		// この行は処理の一部として必要な操作を実行する。
		Range: rng,
		// この行は処理の一部として必要な操作を実行する。
	}, nil
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func parseOptionalIntStrict(name, value string, fallback int) (int, error) {
	// この行で新しい変数を宣言しつつ値を代入する。
	trimmed := strings.TrimSpace(value)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if trimmed == "" {
		// この時点で関数の結果を返して処理を終了する。
		return fallback, nil
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	parsed, err := strconv.Atoi(trimmed)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return 0, &ValidationError{Message: name + " must be an integer"}
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return parsed, nil
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func parseLatLngStrict(latStr, lngStr string) (float64, float64, error) {
	// この行で新しい変数を宣言しつつ値を代入する。
	latTrimmed := strings.TrimSpace(latStr)
	// この行で新しい変数を宣言しつつ値を代入する。
	lngTrimmed := strings.TrimSpace(lngStr)

	// この行で新しい変数を宣言しつつ値を代入する。
	hasLat := latTrimmed != ""
	// この行で新しい変数を宣言しつつ値を代入する。
	hasLng := lngTrimmed != ""

	// 条件を評価し、真のときだけ次の処理を実行する。
	if hasLat != hasLng {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, &ValidationError{Message: "lat and lng must be provided together"}
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if !hasLat {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, nil
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	lat, err := strconv.ParseFloat(latTrimmed, 64)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, &ValidationError{Message: "lat must be a number"}
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, &ValidationError{Message: "lat must be a finite number"}
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	lng, err := strconv.ParseFloat(lngTrimmed, 64)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, &ValidationError{Message: "lng must be a number"}
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if math.IsNaN(lng) || math.IsInf(lng, 0) {
		// この時点で関数の結果を返して処理を終了する。
		return 0, 0, &ValidationError{Message: "lng must be a finite number"}
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return lat, lng, nil
	// ここで処理ブロックを終了する。
}
