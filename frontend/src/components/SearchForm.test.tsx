import { beforeEach, describe, expect, it, vi } from "vitest";
import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";

import { SearchForm } from "./SearchForm";
const LOCATION_STORAGE_KEY = "searchCurrentLocation";
const LOCATION_REFRESH_INTERVAL_MS = 3 * 60 * 1000;

describe("SearchForm", () => {
  let mockOnSearch: ReturnType<typeof vi.fn>;
  let getCurrentPosition: ReturnType<typeof vi.fn>;
  let permissionsQuery: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    sessionStorage.clear();
    vi.restoreAllMocks();
    mockOnSearch = vi.fn() as any;
    getCurrentPosition = vi.fn();
    permissionsQuery = vi
      .fn()
      .mockResolvedValue({ state: "prompt" } as PermissionStatus);

    Object.defineProperty(navigator, "geolocation", {
      value: { getCurrentPosition },
      configurable: true,
    });
    Object.defineProperty(navigator, "permissions", {
      value: { query: permissionsQuery },
      configurable: true,
    });

    vi.spyOn(window.performance, "getEntriesByType").mockImplementation(
      (entryType: string) => {
        if (entryType === "navigation") {
          return [{ type: "navigate" } as PerformanceNavigationTiming];
        }
        return [];
      },
    );
  });

  it("位置情報・住所・キーワードが全て未入力なら画面内エラーを表示して送信しない", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    fireEvent.click(screen.getByRole("button", { name: "検索（おススメ順）" }));

    expect(
      screen.getByText(
        "位置情報が取得できないため検索できません。住所またはキーワード検索をお試しください",
      ),
    ).toBeInTheDocument();
    expect(mockOnSearch).not.toHaveBeenCalled();
  });

  it("住所だけでも検索できる", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    fireEvent.change(screen.getByPlaceholderText("例: 新宿"), {
      target: { value: "梅田" },
    });
    fireEvent.click(screen.getByRole("button", { name: "検索（おススメ順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: "梅田",
      genre: undefined,
      keyword: undefined,
      lat: undefined,
      lng: undefined,
      range: undefined,
    });
  });

  it("キーワードだけでも検索できる", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    fireEvent.change(screen.getByPlaceholderText("例: 個室"), {
      target: { value: "飲み放題" },
    });
    fireEvent.click(screen.getByRole("button", { name: "検索（おススメ順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: "飲み放題",
      lat: undefined,
      lng: undefined,
      range: undefined,
    });
  });

  it("ページアクセス時の位置情報取得だけでも検索できる", () => {
    getCurrentPosition = vi.fn((success: PositionCallback) => {
      success({
        coords: {
          latitude: 35.6895,
          longitude: 139.6917,
          accuracy: 1,
          altitude: null,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
          toJSON: () => ({}),
        },
        timestamp: Date.now(),
        toJSON: () => ({}),
      } as GeolocationPosition);
    });

    Object.defineProperty(navigator, "geolocation", {
      value: { getCurrentPosition },
      configurable: true,
    });

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    expect(getCurrentPosition).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: undefined,
      lat: 35.6895,
      lng: 139.6917,
      range: 3,
    });
  });

  it("位置情報取得時にタイムスタンプ付きでセッションストレージへ保存する", () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    getCurrentPosition = vi.fn((success: PositionCallback) => {
      success({
        coords: {
          latitude: 35.6895,
          longitude: 139.6917,
          accuracy: 1,
          altitude: null,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
          toJSON: () => ({}),
        },
        timestamp: Date.now(),
        toJSON: () => ({}),
      } as GeolocationPosition);
    });

    Object.defineProperty(navigator, "geolocation", {
      value: { getCurrentPosition },
      configurable: true,
    });

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    expect(sessionStorage.getItem(LOCATION_STORAGE_KEY)).toBe(
      JSON.stringify({
        lat: 35.6895,
        lng: 139.6917,
        timestamp: 1700000000000,
      }),
    );
  });

  it("セッションストレージに保存済みの位置情報を初期表示に反映する", async () => {
    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 34.6937,
        lng: 135.5023,
        timestamp: 1700000000000,
      }),
    );

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);
    await act(async () => {});

    expect(getCurrentPosition).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    await waitFor(() =>
      expect(mockOnSearch).toHaveBeenCalledWith({
        address: undefined,
        genre: undefined,
        keyword: undefined,
        lat: 34.6937,
        lng: 135.5023,
        range: 3,
      }),
    );
  });

  it("リロード時かつ前回取得から3分未満なら位置情報を再取得しない", () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    vi.spyOn(window.performance, "getEntriesByType").mockImplementation(
      (entryType: string) => {
        if (entryType === "navigation") {
          return [{ type: "reload" } as PerformanceNavigationTiming];
        }
        return [];
      },
    );

    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 35.1,
        lng: 139.1,
        timestamp: 1700000000000 - (LOCATION_REFRESH_INTERVAL_MS - 1),
      }),
    );

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    expect(getCurrentPosition).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: undefined,
      lat: 35.1,
      lng: 139.1,
      range: 3,
    });
  });

  it("リロード時かつ前回取得から3分以上なら位置情報を再取得する", () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    vi.spyOn(window.performance, "getEntriesByType").mockImplementation(
      (entryType: string) => {
        if (entryType === "navigation") {
          return [{ type: "reload" } as PerformanceNavigationTiming];
        }
        return [];
      },
    );

    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 35.1,
        lng: 139.1,
        timestamp: 1700000000000 - LOCATION_REFRESH_INTERVAL_MS,
      }),
    );

    getCurrentPosition = vi.fn((success: PositionCallback) => {
      success({
        coords: {
          latitude: 35.6895,
          longitude: 139.6917,
          accuracy: 1,
          altitude: null,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
          toJSON: () => ({}),
        },
        timestamp: Date.now(),
        toJSON: () => ({}),
      } as GeolocationPosition);
    });

    Object.defineProperty(navigator, "geolocation", {
      value: { getCurrentPosition },
      configurable: true,
    });

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    expect(getCurrentPosition).toHaveBeenCalledTimes(1);
    expect(sessionStorage.getItem(LOCATION_STORAGE_KEY)).toBe(
      JSON.stringify({
        lat: 35.6895,
        lng: 139.6917,
        timestamp: 1700000000000,
      }),
    );
  });

  it("ローディング中は検索ボタンが無効化される", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={true} />);

    const searchButton = screen.getByRole("button", { name: "検索中..." });
    expect(searchButton).toBeDisabled();
  });

  it("selectedGenreが指定されている場合は検索パラメータに反映する", () => {
    render(
      <SearchForm
        onSearch={mockOnSearch as any}
        isLoading={false}
        selectedGenre="G001"
      />,
    );

    fireEvent.change(screen.getByPlaceholderText("例: 新宿"), {
      target: { value: "梅田" },
    });
    fireEvent.click(screen.getByRole("button", { name: "検索（おススメ順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: "梅田",
      genre: "G001",
      keyword: undefined,
      lat: undefined,
      lng: undefined,
      range: undefined,
    });
  });

  it("検索ボタン押下時に権限がgrantedかつ保存位置が3分以上前なら位置情報を再取得してから検索する", async () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    permissionsQuery.mockResolvedValue({
      state: "granted",
    } as PermissionStatus);

    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 35.1,
        lng: 139.1,
        timestamp: 1700000000000 - LOCATION_REFRESH_INTERVAL_MS,
      }),
    );

    getCurrentPosition = vi.fn((success: PositionCallback) => {
      success({
        coords: {
          latitude: 35.6895,
          longitude: 139.6917,
          accuracy: 1,
          altitude: null,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
          toJSON: () => ({}),
        },
        timestamp: Date.now(),
        toJSON: () => ({}),
      } as GeolocationPosition);
    });
    Object.defineProperty(navigator, "geolocation", {
      value: { getCurrentPosition },
      configurable: true,
    });

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    await waitFor(() =>
      expect(permissionsQuery).toHaveBeenCalledWith({ name: "geolocation" }),
    );
    await waitFor(() =>
      expect(mockOnSearch).toHaveBeenCalledWith({
        address: undefined,
        genre: undefined,
        keyword: undefined,
        lat: 35.6895,
        lng: 139.6917,
        range: 3,
      }),
    );
    expect(getCurrentPosition).toHaveBeenCalledTimes(1);
  });

  it("検索ボタン押下時に権限がgranted以外なら位置情報を再取得せず検索する", async () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    permissionsQuery.mockResolvedValue({ state: "prompt" } as PermissionStatus);

    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 35.1,
        lng: 139.1,
        timestamp: 1700000000000 - LOCATION_REFRESH_INTERVAL_MS,
      }),
    );

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    await waitFor(() =>
      expect(mockOnSearch).toHaveBeenCalledWith({
        address: undefined,
        genre: undefined,
        keyword: undefined,
        lat: 35.1,
        lng: 139.1,
        range: 3,
      }),
    );
    expect(getCurrentPosition).not.toHaveBeenCalled();
  });

  it("検索ボタン押下時にPermissions APIが未対応なら位置情報を再取得せず検索する", () => {
    vi.spyOn(Date, "now").mockReturnValue(1700000000000);
    Object.defineProperty(navigator, "permissions", {
      value: undefined,
      configurable: true,
    });

    sessionStorage.setItem(
      LOCATION_STORAGE_KEY,
      JSON.stringify({
        lat: 35.1,
        lng: 139.1,
        timestamp: 1700000000000 - LOCATION_REFRESH_INTERVAL_MS,
      }),
    );

    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);
    fireEvent.click(screen.getByRole("button", { name: "検索（距離順）" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: undefined,
      lat: 35.1,
      lng: 139.1,
      range: 3,
    });
    expect(getCurrentPosition).not.toHaveBeenCalled();
  });
});
