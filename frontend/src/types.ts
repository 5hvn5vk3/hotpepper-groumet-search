export interface ShopCreditCard {
  name: string;
}

export interface ShopGenre {
  name: string;
  catch: string;
}

export interface ShopURLs {
  pc: string;
}

export interface ShopPhoto {
  pc: {
    l: string;
  };
}

export interface ShopListItem {
  id: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  genre: ShopGenre;
  catch: string;
  access: string;
  urls: ShopURLs;
  photo: ShopPhoto;
  credit_card?: ShopCreditCard[];
}

export interface ShopDetailSupplement {
  open?: string;
  close?: string;
  budget_memo?: string;
  wifi?: string;
  private_room?: string;
  non_smoking?: string;
  parking?: string;
  lunch?: string;
  midnight?: string;
  shop_detail_memo?: string;
}

export type ShopDetailView = ShopListItem & ShopDetailSupplement;

export type RestaurantDetailStatus = "idle" | "loading" | "ready" | "failed";
export interface GourmetSearchResponse<TShop> {
  results: {
    results_available: number;
    results_start: number;
    shop: TShop[];
  };
}
export interface Genre {
  code: string;
  name: string;
}
export interface GenreMasterResponse {
  results: {
    genre: Genre[];
  };
}
export interface GourmetSearchParams {
  address?: string;
  genre?: string;
  keyword?: string;
  page?: number;
  count?: number;
  lat?: number;
  lng?: number;
  range?: number;
}

// 非2xx時のエラーレスポンス型（成功型とは分離）
export interface APIErrorResponse {
  error: {
    message: string;
  };
}
