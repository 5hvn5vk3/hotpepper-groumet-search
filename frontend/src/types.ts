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
  photo: {
    pc: {
      l: string;
      m: string;
      s: string;
    };
  };
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
  lat?: number;
  lng?: number;
  range?: number;
}
export type SearchParams = GourmetSearchParams;

// 非2xx時のエラーレスポンス型（成功型とは分離）
export interface APIErrorResponse {
  error: {
    message: string;
  };
}
