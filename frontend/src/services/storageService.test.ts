import { describe, it, expect, beforeEach, vi } from "vitest";
import { storageService } from "./storageService";
describe("StorageService", () => {
    beforeEach(() => {
        sessionStorage.clear();
        vi.clearAllMocks();
    });
    describe("get", () => {
        it("存在しないキーの場合nullを返す", () => {
            const result = storageService.get("nonExistent");
            expect(result).toBeNull();
        });
        it("オブジェクトデータを正しく取得できる", () => {
            const testData = { name: "test", value: 123 };
            sessionStorage.setItem("testKey", JSON.stringify(testData));
            const result = storageService.get<typeof testData>("testKey");
            expect(result).toEqual(testData);
        });
        it("無効なJSONの場合nullを返す", () => {
            sessionStorage.setItem("testKey", "invalid json");
            const result = storageService.get("testKey");
            expect(result).toBeNull();
        });
    });
    describe("set", () => {
        it("オブジェクトデータを保存できる", () => {
            const testData = { name: "test", value: 123 };
            storageService.set("testKey", testData);
            const stored = sessionStorage.getItem("testKey");
            expect(stored).toBe(JSON.stringify(testData));
        });
    });
    describe("remove", () => {
        it("指定したキーのデータを削除できる", () => {
            sessionStorage.setItem("testKey", JSON.stringify({ value: 1 }));
            storageService.remove("testKey");
            expect(sessionStorage.getItem("testKey")).toBeNull();
        });
    });
});
