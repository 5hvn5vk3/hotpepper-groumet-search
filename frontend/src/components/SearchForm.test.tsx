import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";

import { SearchForm } from "./SearchForm";
import * as useGenresHook from "@/hooks/useGenres";

vi.mock("@/hooks/useGenres");

const mockGenres = [
  { code: "G001", name: "居酒屋" },
  { code: "G002", name: "イタリアン" },
];

describe("SearchForm", () => {
  let mockOnSearch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnSearch = vi.fn() as any;
    vi.spyOn(useGenresHook, "useGenres").mockReturnValue({
      genres: mockGenres,
      isLoading: false,
      error: null,
    });
  });

  it("位置情報・住所・キーワードが全て未入力なら画面内エラーを表示して送信しない", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    fireEvent.click(screen.getByRole("button", { name: "検索" }));

    expect(
      screen.getByText(
        "位置情報が取得できないため検索できません。住所またはキーワード検索をお試しください"
      )
    ).toBeInTheDocument();
    expect(mockOnSearch).not.toHaveBeenCalled();
  });

  it("住所だけでも検索できる", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={false} />);

    fireEvent.change(screen.getByPlaceholderText("例: 新宿"), {
      target: { value: "梅田" },
    });
    fireEvent.click(screen.getByRole("button", { name: "検索" }));

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
    fireEvent.click(screen.getByRole("button", { name: "検索" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: "飲み放題",
      lat: undefined,
      lng: undefined,
      range: undefined,
    });
  });

  it("位置情報だけでも検索できる", () => {
    const getCurrentPosition = vi.fn((success: PositionCallback) => {
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

    fireEvent.click(screen.getByRole("button", { name: "現在地を取得" }));
    fireEvent.click(screen.getByRole("button", { name: "検索" }));

    expect(mockOnSearch).toHaveBeenCalledWith({
      address: undefined,
      genre: undefined,
      keyword: undefined,
      lat: 35.6895,
      lng: 139.6917,
      range: 3,
    });
  });

  it("ローディング中は検索ボタンが無効化される", () => {
    render(<SearchForm onSearch={mockOnSearch as any} isLoading={true} />);

    const searchButton = screen.getByRole("button", { name: "検索中..." });
    expect(searchButton).toBeDisabled();
  });
});
