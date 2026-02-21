// StorageServiceクラス：ブラウザのsessionStorageを管理するサービス
// sessionStorageはブラウザにデータをセッション中のみ保存できるWeb Storage API
class StorageService {
  // getメソッド：sessionStorageからデータを取得
  // <T>はジェネリック型：取得するデータの型を柔軟に指定
  // 戻り値: T型のデータ、または存在しない場合はnull
  get<T>(key: string): T | null {
    try {
      // sessionStorage.getItem()で指定したキーのデータを取得
      // 存在しない場合はnullが返される
      const item = sessionStorage.getItem(key);

      // 三項演算子: itemが存在する場合はJSONパース、存在しない場合はnull
      // JSON.parse(): JSON文字列をJavaScriptオブジェクトに変換
      return item ? JSON.parse(item) : null;
    } catch (error) {
      // try-catch: エラーが発生してもアプリがクラッシュしないように処理
      // console.error(): 開発者ツールのコンソールにエラーを出力
      // テンプレートリテラルでエラーメッセージを作成
      console.error(`Failed to get item from storage: ${key}`, error);
      // エラー時はnullを返す
      return null;
    }
  }

  // setメソッド：sessionStorageにデータを保存
  // <T>: 保存するデータの型
  // void: 戻り値なし
  set<T>(key: string, value: T): void {
    try {
      // JSON.stringify(): JavaScriptオブジェクトをJSON文字列に変換
      // sessionStorage.setItem()でキーと値のペアを保存
      // sessionStorageは文字列のみ保存可能なため、JSON.stringify()が必要
      sessionStorage.setItem(key, JSON.stringify(value));
    } catch (error) {
      // エラーが発生した場合（例：ストレージ容量オーバー）
      console.error(`Failed to set item in storage: ${key}`, error);
    }
  }
}

// StorageServiceのインスタンスを作成してエクスポート
// シングルトンパターン：アプリ全体で1つのインスタンスを共有
export const storageService = new StorageService();
