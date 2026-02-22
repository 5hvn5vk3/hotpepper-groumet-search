// 型定義をインポート
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";
// API関数をインポート
import { apiGet } from "@/api/client";

// searchRestaurants関数：検索パラメータを元にレストランを検索
// params: 検索条件（都道府県、住所、ジャンル、キーワード、ページ番号など）
// 戻り値: 検索結果のレスポンスデータ
export const searchRestaurants = async (
  params: GourmetSearchParams,
): Promise<GourmetSearchResponse> => {
  // serviceAreaがない場合でも、lat/lngがあれば検索を行う
  const hasLatLng =
    typeof params.lat === "number" && typeof params.lng === "number";
  if ((!params.serviceArea || params.serviceArea.trim() === "") && !hasLatLng) {
    // どちらも指定がなければ空の結果を返す
    return {
      results: {
        api_version: "1.0",
        results_available: 0,
        results_returned: "0",
        results_start: 1,
        shop: [],
      },
    };
  }

  // ページ番号のバリデーション：最小値を1に制限
  // Math.max(1, x)は1とxのうち大きい方を返す
  // params.page ?? 1は、params.pageがnullまたはundefinedの場合に1を使用
  const page = Math.max(1, params.page ?? 1);

  // 1ページあたりの件数のバリデーション：1〜100の範囲に制限
  // Math.min(max, x)はmaxとxのうち小さい方を返す
  const count = Math.min(Math.max(params.count ?? 20, 1), 100);

  // API用の開始位置を計算
  // 例：2ページ目で20件/ページの場合 → (2-1) * 20 + 1 = 21
  const start = (page - 1) * count + 1;

  // URLSearchParamsオブジェクトでクエリパラメータを構築
  // URLSearchParamsはクエリ文字列を簡単に作成できるWeb標準API
  const queryParams = new URLSearchParams({
    service_area: params.serviceArea ?? "", // 必須：都道府県コード（未指定なら空）
    start: start.toString(), // 開始位置を文字列に変換
    count: count.toString(), // 取得件数を文字列に変換
  });

  // オプショナルパラメータの追加（値がある場合のみ）
  if (params.address) {
    queryParams.set("address", params.address); // 住所キーワード
  }
  if (params.genre) {
    queryParams.set("genre", params.genre); // ジャンルコード
  }
  if (params.keyword) {
    queryParams.set("keyword", params.keyword); // フリーワード
  }
  // 緯度経度が指定されていればクエリに追加
  if (hasLatLng) {
    queryParams.set("lat", params.lat!.toString());
    queryParams.set("lng", params.lng!.toString());
    if (params.range) {
      queryParams.set("range", params.range.toString());
    }
  }

  // apiGet関数を使用してAPIリクエストを送信
  // テンプレートリテラルでURLを構築：/api/gourmet?パラメータ
  // GourmetSearchResponse型でレスポンスの型を指定
  return apiGet<GourmetSearchResponse>(`/api/gourmet?${queryParams}`);
};
