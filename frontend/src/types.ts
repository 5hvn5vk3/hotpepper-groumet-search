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
  photo: {
    pc: {
      l: string;
      m: string;
      s: string;
    };
  };
}

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
