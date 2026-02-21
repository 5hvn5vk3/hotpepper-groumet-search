// 環境変数からAPIのベースURLを取得
// import.meta.envはViteの環境変数にアクセスするための特殊なオブジェクト
// VITE_で始まる環境変数のみクライアントサイドで利用可能
// 設定されていない場合は空文字列をデフォルト値として使用
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

// apiGet関数：GETリクエストを送信する汎用関数
// <T>はジェネリック型パラメータ：レスポンスの型を柔軟に指定できる
// async: 非同期関数であることを示す（Promiseを返す）
export async function apiGet<T>(endpoint: string): Promise<T> {
  // テンプレートリテラル（`文字列`）で完全なURLを構築
  // ${変数}で変数を埋め込む
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    // fetch APIでHTTPリクエストを送信
    // await: Promiseの結果を待つ（非同期処理を同期的に書ける）
    const response = await fetch(url, {
      method: "GET", // HTTPメソッドを指定
      headers: {
        // リクエストヘッダーを設定
        "Content-Type": "application/json", // JSON形式のデータを扱うことを示す
      },
    });

    // レスポンスのステータスコードをチェック
    // response.okはステータスコードが200〜299の範囲であればtrue
    if (!response.ok) {
      // エラーオブジェクトを投げる（throw）とcatchブロックに処理が移る
      throw new Error(
        `API request failed: ${response.status} ${response.statusText}`
      );
    }

    // レスポンスボディをJSON形式としてパース（解析）
    // 返り値の型はジェネリック型Tとして推論される
    return response.json();
  } catch (error) {
    // errorがError型のインスタンスかどうかをチェック
    // instanceofは型チェック演算子
    if (error instanceof Error) {
      // Error型の場合はそのまま再スロー
      throw error;
    }
    // 予期しないエラーの場合は新しいErrorオブジェクトを作成
    throw new Error("An unexpected error occurred");
  }
}
