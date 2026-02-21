// Reactのフックをインポート
import { useState, useCallback } from "react";

// モーダルフックが返す値の型定義
// <T>はジェネリック型：モーダルに渡すデータの型を柔軟に指定できる
interface UseModalResult<T> {
  isOpen: boolean; // モーダルが開いているかどうか
  data: T | null; // モーダルに表示するデータ（閉じている時はnull）
  open: (data: T) => void; // モーダルを開く関数
  close: () => void; // モーダルを閉じる関数
}

// useModalフック：モーダルの開閉状態とデータを管理する汎用フック
// <T>はジェネリック型パラメータ：使用時に具体的な型を指定（例：Shop型）
export const useModal = <T>(): UseModalResult<T> => {
  // モーダルの開閉状態を管理（初期値：false = 閉じている）
  const [isOpen, setIsOpen] = useState(false);

  // モーダルに表示するデータを管理（初期値：null）
  const [data, setData] = useState<T | null>(null);

  // open関数：モーダルを開き、表示するデータをセット
  // useCallbackでメモ化：不要な関数の再生成を防ぐ
  const open = useCallback((modalData: T) => {
    // 表示するデータをセット
    setData(modalData);
    // モーダルを開く
    setIsOpen(true);
  }, []); // 依存配列が空：マウント時に一度だけ関数を作成

  // close関数：モーダルを閉じ、データをクリア
  const close = useCallback(() => {
    // モーダルを閉じる
    setIsOpen(false);
    // データをクリア
    setData(null);
  }, []); // 依存配列が空：マウント時に一度だけ関数を作成

  // フックの戻り値：状態と関数をまとめたオブジェクト
  return { isOpen, data, open, close };
};
