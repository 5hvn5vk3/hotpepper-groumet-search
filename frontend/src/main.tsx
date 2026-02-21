// Reactライブラリから必要なモジュールをインポート
// StrictModeは開発時に潜在的な問題を検出するためのツール
import { StrictMode } from "react";
// createRootはReact 18の新しいルート作成API
import { createRoot } from "react-dom/client";
// グローバルなCSSスタイルをインポート
import "./index.css";
// メインのアプリケーションコンポーネントをインポート
import App from "./App.tsx";

// DOMの"root"要素を取得し、Reactアプリケーションのルートを作成
// document.getElementById('root')でHTML内のid="root"要素を取得
// !は非nullアサーション演算子（要素が必ず存在することをTypeScriptに伝える）
createRoot(document.getElementById("root")!).render(
  // StrictModeでアプリケーションをラップ
  // 開発モードで2回レンダリングすることで副作用を検出する
  <StrictMode>
    {/* Appコンポーネントをレンダリング */}
    <App />
  </StrictMode>
);
