import type {
    GourmetSearchParams,
    GourmetSearchResponse,
    ShopDetailSupplement,
    ShopDetailView,
    ShopListItem,
} from "@/types";
import { apiGet } from "@/api/client";

const pickShopDetailSupplement = (
    detail: ShopDetailView,
): ShopDetailSupplement => ({
    open: detail.open,
    close: detail.close,
    budget_memo: detail.budget_memo,
    wifi: detail.wifi,
    private_room: detail.private_room,
    non_smoking: detail.non_smoking,
    parking: detail.parking,
    lunch: detail.lunch,
    midnight: detail.midnight,
    shop_detail_memo: detail.shop_detail_memo,
});

export const searchRestaurants = async (
    params: GourmetSearchParams,
): Promise<GourmetSearchResponse<ShopListItem>> => {
    const hasLatLng = typeof params.lat === "number" && typeof params.lng === "number";
    const address = params.address?.trim() ?? "";
    const keyword = params.keyword?.trim() ?? "";
    const hasAddress = address !== "";
    const hasKeyword = keyword !== "";
    if (!hasLatLng && !hasAddress && !hasKeyword) {
        return {
            results: {
                results_available: 0,
                results_start: 1,
                shop: [],
            },
        };
    }
    const page = Math.max(1, params.page ?? 1);
    const count = Math.min(Math.max(params.count ?? 20, 1), 100);
    const start = (page - 1) * count + 1;
    const queryParams = new URLSearchParams({
        start: start.toString(),
        count: count.toString(),
    });
    if (hasAddress) {
        queryParams.set("address", address);
    }
    if (params.genre) {
        queryParams.set("genre", params.genre);
    }
    if (hasKeyword) {
        queryParams.set("keyword", keyword);
    }
    if (hasLatLng) {
        queryParams.set("lat", params.lat!.toString());
        queryParams.set("lng", params.lng!.toString());
        if (params.range) {
            queryParams.set("range", params.range.toString());
        }
    }
    return apiGet<GourmetSearchResponse<ShopListItem>>(
        `/api/gourmet?${queryParams}`,
    );
};

export const fetchRestaurantDetail = async (
    id: string,
): Promise<ShopDetailSupplement | null> => {
    const shopID = id.trim();
    if (shopID === "") {
        return null;
    }
    const queryParams = new URLSearchParams({ id: shopID });
    const response = await apiGet<GourmetSearchResponse<ShopDetailView>>(
        `/api/gourmet?${queryParams}`,
    );
    const detail = response.results.shop[0];
    if (!detail) {
        return null;
    }
    return pickShopDetailSupplement(detail);
};
