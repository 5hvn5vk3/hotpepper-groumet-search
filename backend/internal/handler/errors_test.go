package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/types"
)

// TestMaskAPIKeyInLog は 2-a のテスト。
// maskAPIKeyInLog がAPIキーを [REDACTED] に置換することをテーブル駆動で検証する。
//
// 【テーブル駆動テスト（Table-Driven Tests）とは？】
// Go コミュニティで標準的に使われるテストパターン。
// 入力と期待値のペアを「テーブル（構造体スライス）」にまとめ、
// 同じ検証ロジックを複数ケースに適用する。
// メリット：
//   - 新しいケースはスライスへの追加だけで済む（コード重複なし）
//   - テストケース一覧が視認しやすく、カバレッジの抜け漏れを発見しやすい
//   - 名前付きフィールドでケースの意図を明示できる
func TestMaskAPIKeyInLog(t *testing.T) {
	// 各要素は name（ケース名）・err（入力）・want（期待値）を持つ無名構造体のスライス。
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "? 区切り",
			err:  fmt.Errorf("url?key=SECRET&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
		{
			name: "& 区切り",
			err:  fmt.Errorf("url?lat=35&key=SECRET"),
			want: "url?lat=35&key=[REDACTED]",
		},
		{
			name: "キーなし",
			err:  fmt.Errorf("network timeout"),
			want: "network timeout",
		},
		{
			name: "nil",
			err:  nil,
			want: "",
		},
		{
			name: "記号入りキー",
			err:  fmt.Errorf("url?key=A+B/C&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
		{
			name: "key のみのパラメータ",
			err:  fmt.Errorf("url?key=SECRET"),
			want: "url?key=[REDACTED]",
		},
	}

	// 【t.Run によるサブテスト】
	// t.Run("名前", func(t *testing.T) { ... }) でテーブルの各ケースを独立したサブテストとして実行する。
	// メリット：
	//   - 失敗時に「TestMaskAPIKeyInLog/? 区切り」のようにどのケースが失敗したか名前付きで分かる
	//   - go test -run TestMaskAPIKeyInLog/キーなし のように特定ケースだけ実行できる
	//   - 内部で t.Parallel() を呼べば並列実行も可能
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := maskAPIKeyInLog(tc.err)
			if got != tc.want {
				// 【t.Error vs t.Fatal の使い分け】
				// t.Error/t.Errorf：テストを「失敗」マークするが実行を続行する。
				//   → 他のケースの結果も確認したい場合に使う（ここはケースの最後なので継続可）。
				// t.Fatal/t.Fatalf：即座にテスト関数を終了する（runtime.Goexit を呼ぶ）。
				//   → nilポインタ参照など後続処理が危険になる場合、またはそれ以降の検証が
				//     意味をなさない場合に使う。
				t.Errorf("maskAPIKeyInLog(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}

// TestHandleServiceError_MaskAPIKey は 2-b のテスト。
// handleServiceError がログに APIキーを出力しないことを検証する。
// log がグローバル状態を持つため t.Parallel() は付けない。
func TestHandleServiceError_MaskAPIKey(t *testing.T) {
	// ① 元の出力先を退避し、テスト終了時に確実に戻す
	//
	// 【t.Cleanup とは？】
	// t.Cleanup(func) に渡した関数は、テスト終了時（Pass/Fail どちらでも）に必ず呼ばれる。
	// defer と似ているが、テストヘルパー関数の中で登録しても
	// 呼び出し元テストの終了時に実行される点が異なる。
	// グローバルな状態（ここでは log の出力先）を変更する場合は必ず Cleanup で元に戻すこと。
	// そうしないと他のテストへ影響する「テスト汚染」が起きる。
	original := log.Writer()
	t.Cleanup(func() { log.SetOutput(original) })

	// ② ログをメモリバッファに横取りする
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// ③ レスポンスを httptest.Recorder で捕捉する
	//
	// 【httptest.NewRecorder とは？】
	// 実際の net.Conn を使わず、HTTP レスポンスをメモリ上に記録する ResponseWriter の実装。
	// テスト後に rec.Code（ステータスコード）・rec.Body（ボディ）・rec.Header()（ヘッダー）を
	// 検査できる。本番では http.ResponseWriter が渡されるが、
	// テストでは httptest.ResponseRecorder で代替する。
	rec := httptest.NewRecorder()

	// ④ APIキーを含むエラーで handleServiceError を呼ぶ
	handleServiceError(rec, fmt.Errorf("url?key=SECRETKEY&lat=35.0"))

	// ⑤ 検証
	// 【strings.Contains によるエラーメッセージの部分一致検証】
	// ログ出力は実装の詳細によって前後に付加情報が入ることがある。
	// そのため完全一致（==）ではなく strings.Contains で部分一致を確認するのが堅牢。
	logOutput := buf.String()
	if strings.Contains(logOutput, "SECRETKEY") {
		t.Errorf("ログに生のAPIキー 'SECRETKEY' が含まれている: %s", logOutput)
	}
	if !strings.Contains(logOutput, "[REDACTED]") {
		t.Errorf("ログに '[REDACTED]' が含まれていない: %s", logOutput)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("HTTPステータス %d を期待したが %d だった", http.StatusInternalServerError, rec.Code)
	}
}

// TestHandleServiceError_ErrorCodeMapping は 2-c のテスト。
// HotpepperAPIError のコードが正しい HTTP ステータスとメッセージにマッピングされることをテーブル駆動で検証する。
func TestHandleServiceError_ErrorCodeMapping(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		wantStatus     int
		wantMessage    string
	}{
		{
			name:        "code=3000 → 400",
			code:        3000,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "検索条件が正しくありません",
		},
		{
			name:        "code=1000 → 502",
			code:        1000,
			wantStatus:  http.StatusBadGateway,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "code=2000 → 500",
			code:        2000,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "code=9999（未定義） → 500",
			code:        9999,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handleServiceError(rec, &types.HotpepperAPIError{Code: tc.code})

			if rec.Code != tc.wantStatus {
				t.Errorf("HTTPステータス %d を期待したが %d だった", tc.wantStatus, rec.Code)
			}

			var body map[string]map[string]string
			// 【t.Fatal の使いどころ】
			// JSON デコードに失敗した場合、後続の body["error"] アクセスが
			// 正しく行えないため、ここで即座にテストを終了するのが適切。
			// 「後続処理が意味をなさない場合」は t.Fatal/t.Fatalf を使う。
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("レスポンスボディの JSON デコードに失敗: %v", err)
			}
			gotMessage := body["error"]["message"]
			if gotMessage != tc.wantMessage {
				t.Errorf("error.message = %q, want %q", gotMessage, tc.wantMessage)
			}
		})
	}
}
