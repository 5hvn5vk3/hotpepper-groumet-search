// 型定義をインポート
import type { GenreMasterResponse, Genre } from "@/types";
// API関数をインポート
import { apiGet } from "@/api/client";
// ストレージサービスをインポート（キャッシュ管理用）
import { storageService } from "@/services/storageService";

// ローカルストレージのキャッシュキー
// constで定数として定義し、タイポを防ぐ
const GENRES_CACHE_KEY = "genres";

// getGenres関数：ジャンル一覧を取得（キャッシュ機能付き）
// 戻り値: Genre型の配列
export const getGenres = async (): Promise<Genre[]> => {
  // まずキャッシュから取得を試みる
  // storageService.getでローカルストレージからデータを取得
  // <Genre[]>でジェネリック型を指定し、正しい型でデータを取得
  const cached = storageService.get<Genre[]>(GENRES_CACHE_KEY);

  // キャッシュにデータがある場合、それを返す（API呼び出しをスキップ）
  // これによりネットワーク通信を削減し、パフォーマンスを向上
  if (cached) {
    return cached;
  }

  // キャッシュがない場合、APIから取得
  // apiGet関数でGETリクエストを送信
  const response = await apiGet<GenreMasterResponse>("/api/genre");

  // レスポンスからジャンル配列を抽出
  // response.results.genreの階層構造でジャンルデータにアクセス
  const genres = response.results.genre;

  // 取得したデータをキャッシュに保存
  // 次回以降はAPIを呼ばずにキャッシュから取得できる
  storageService.set(GENRES_CACHE_KEY, genres);

  // ジャンル配列を返す
  return genres;
};
