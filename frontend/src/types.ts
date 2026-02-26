export interface ShopCreditCard {
  code: string;
  name: string;
}

export interface Shop {
  id: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  genre: {
    name: string;
    catch: string;
  };
  catch: string;
  access: string;
  urls: {
    pc: string;
  };
  // 営業時間
  open?: string;
  // 定休日・営業時間補足
  close?: string;
  // 予算メモ
  budget_memo?: string;
  // 設備・条件
  wifi?: string;
  private_room?: string;
  non_smoking?: string;
  parking?: string;
  lunch?: string;
  midnight?: string;
  // 店舗説明
  shop_detail_memo?: string;
  photo: {
    pc: {
      l: string;
      m: string;
      s: string;
    };
  };
  // type=lite+credit_card で追加される利用可能カード情報
  credit_card?: ShopCreditCard[];
}

export type RestaurantDetailStatus = "idle" | "loading" | "ready" | "failed";

export interface GourmetSearchResponse {
  results: {
    api_version: string;
    results_available: number;
    results_returned: string;
    results_start: number;
    shop: Shop[];
  };
}

export interface Genre {
  code: string;
  name: string;
}

export interface GenreMasterResponse {
  results: {
    api_version: string;
    results_available: number;
    results_returned: string;
    results_start: number;
    genre: Genre[];
  };
}

export interface GourmetSearchParams {
  address?: string;
  genre?: string;
  keyword?: string;
  page?: number;
  count?: number;
  // 緯度経度で検索する場合
  lat?: number;
  lng?: number;
  // Hotpepper APIのrange（1〜5）
  range?: number;
}

// 互換性のためのエイリアス：既存コードがSearchParamsを参照している場合に対応
export type SearchParams = GourmetSearchParams;
