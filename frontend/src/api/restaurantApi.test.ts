import { beforeEach, describe, expect, it, vi } from "vitest";
import { fetchRestaurantDetail, searchRestaurants } from "./restaurantApi";
import { apiGet } from "./client";
import type { GourmetSearchResponse, GourmetSearchParams } from "@/types";
vi.mock("./client");
describe("searchRestaurants", () => {
    const mockResponse: GourmetSearchResponse = {
        results: {
            api_version: "1.0",
            results_available: 100,
            results_returned: "20",
            results_start: 1,
            shop: [],
        },
    };
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(apiGet).mockResolvedValue(mockResponse);
    });
    it("位置情報・住所・キーワードが全て未指定ならAPI呼び出しをスキップする", async () => {
        const params: GourmetSearchParams = {
            genre: "G001",
        };
        const result = await searchRestaurants(params);
        expect(apiGet).not.toHaveBeenCalled();
        expect(result.results.results_available).toBe(0);
        expect(result.results.shop).toEqual([]);
    });
    it("住所のみで検索できる", async () => {
        await searchRestaurants({ address: "新宿" });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("address=%E6%96%B0%E5%AE%BF"));
        expect(apiGet).toHaveBeenCalledWith(expect.not.stringContaining("service_area"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=1"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=20"));
    });
    it("キーワードのみで検索できる", async () => {
        await searchRestaurants({ keyword: "居酒屋" });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("keyword=%E5%B1%85%E9%85%92%E5%B1%8B"));
    });
    it("位置情報とオプションを含めて検索できる", async () => {
        const params: GourmetSearchParams = {
            lat: 35.6895,
            lng: 139.6917,
            range: 3,
            address: "新宿",
            genre: "G001",
            keyword: "居酒屋",
        };
        await searchRestaurants(params);
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("lat=35.6895"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("lng=139.6917"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("range=3"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("genre=G001"));
    });
    it("ページ2・count10ならstart=11になる", async () => {
        await searchRestaurants({ keyword: "寿司", page: 2, count: 10 });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=11"));
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=10"));
    });
    it("pageが0以下なら1として扱う", async () => {
        await searchRestaurants({ keyword: "寿司", page: 0 });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("start=1"));
    });
    it("countが0以下なら1、100超なら100に補正する", async () => {
        await searchRestaurants({ keyword: "寿司", count: 0 });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=1"));
        vi.clearAllMocks();
        vi.mocked(apiGet).mockResolvedValue(mockResponse);
        await searchRestaurants({ keyword: "寿司", count: 150 });
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("count=100"));
    });
});
describe("fetchRestaurantDetail", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });
    it("idが空ならAPI呼び出しをスキップしてnullを返す", async () => {
        const result = await fetchRestaurantDetail(" ");
        expect(apiGet).not.toHaveBeenCalled();
        expect(result).toBeNull();
    });
    it("id指定で詳細を取得し、先頭のshopを返す", async () => {
        vi.mocked(apiGet).mockResolvedValue({
            results: {
                api_version: "1.0",
                results_available: 1,
                results_returned: "1",
                results_start: 1,
                shop: [
                    {
                        id: "J001234567",
                        name: "テスト居酒屋",
                        address: "東京都渋谷区",
                        lat: 35.6595,
                        lng: 139.7004,
                        genre: { name: "居酒屋", catch: "気軽に楽しめる" },
                        catch: "美味しい料理とお酒",
                        access: "渋谷駅徒歩5分",
                        urls: { pc: "https://example.com/restaurant1" },
                        photo: {
                            pc: {
                                l: "https://example.com/photo_l.jpg",
                                m: "https://example.com/photo_m.jpg",
                                s: "https://example.com/photo_s.jpg",
                            },
                        },
                    },
                ],
            },
        } as GourmetSearchResponse);
        const result = await fetchRestaurantDetail("J001234567");
        expect(apiGet).toHaveBeenCalledWith(expect.stringContaining("/api/gourmet?id=J001234567"));
        expect(result?.id).toBe("J001234567");
    });
});
