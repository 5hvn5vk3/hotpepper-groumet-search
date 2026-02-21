// ===================================================================
// RestaurantListコンポーネントのテスト
// ===================================================================
// 【目的】RestaurantListコンポーネントの動作を検証
// 【Triangulation】最小の3点でリスト表示とイベントをテスト
// ===================================================================

// Vitestのテスト関数をインポート
import { describe, it, expect, vi, beforeEach } from "vitest";

// React Testing Library
import { render, screen, fireEvent } from "@testing-library/react";

// テスト対象のコンポーネント
import { RestaurantList } from "./RestaurantList";

// Shop型
import type { Shop } from "@/types";

// ===================================================================
// Triangulation Point の概要
// ===================================================================
// Point 1: 空配列 - 検索結果がない場合のメッセージ表示
// Point 2: データあり - レストランカードの表示
// Point 3: クリックイベント - onSelectRestaurantコールバックの呼び出し
// ===================================================================

describe("RestaurantList", () => {
  // モック関数の準備
  let mockOnSelectRestaurant: ReturnType<typeof vi.fn>;

  // テスト用のレストランデータ
  const mockRestaurants: Shop[] = [
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
    {
      id: "J001234568",
      name: "イタリアンレストラン",
      address: "東京都新宿区",
      lat: 35.6938,
      lng: 139.7036,
      genre: { name: "イタリアン", catch: "本格派" },
      catch: "本場の味を堪能",
      access: "新宿駅徒歩3分",
      urls: { pc: "https://example.com/restaurant2" },
      photo: {
        pc: {
          l: "https://example.com/photo2_l.jpg",
          m: "https://example.com/photo2_m.jpg",
          s: "https://example.com/photo2_s.jpg",
        },
      },
    },
  ];

  // 各テストの前に実行される準備処理
  beforeEach(() => {
    // モック関数を作成
    mockOnSelectRestaurant = vi.fn() as any;
  });

  // ===================================================================
  // Triangulation Point 1: 空配列の場合
  // ===================================================================
  // 【目的】レストランが0件の場合、適切な空メッセージが表示されることを確認
  // 【重要性】境界値（空データ）での動作保証、ユーザーへの適切なフィードバック
  // 【差分】restaurants=[] → 空メッセージ表示
  it("レストランが0件の場合、検索結果なしメッセージを表示", () => {
    // Arrange（準備）: 空配列でレンダリング
    render(
      <RestaurantList
        restaurants={[]}
        onSelectRestaurant={mockOnSelectRestaurant as any}
      />
    );

    // Assert（検証）: 空メッセージが表示されているか
    expect(screen.getByText("検索結果がありません")).toBeInTheDocument();

    // レストランカードは表示されていないことを確認
    const images = screen.queryAllByRole("img");
    expect(images.length).toBe(0);
  });

  // ===================================================================
  // Triangulation Point 2: データがある場合
  // ===================================================================
  // 【目的】レストランデータがある場合、カードが正しく表示されることを確認
  // 【重要性】基本的なリスト表示機能の動作保証
  // 【差分】restaurants=[{...}, {...}] → カード2枚表示
  it("レストランが1件以上ある場合、カード形式で表示", () => {
    // Arrange（準備）: 2件のレストランでレンダリング
    render(
      <RestaurantList
        restaurants={mockRestaurants}
        onSelectRestaurant={mockOnSelectRestaurant as any}
      />
    );

    // Assert（検証）: 各レストランの情報が表示されているか

    // 1件目のレストラン
    expect(screen.getByText("テスト居酒屋")).toBeInTheDocument();
    expect(screen.getByText("居酒屋")).toBeInTheDocument();
    expect(screen.getByText("美味しい料理とお酒")).toBeInTheDocument();
    expect(screen.getByText("東京都渋谷区")).toBeInTheDocument();

    // 2件目のレストラン
    expect(screen.getByText("イタリアンレストラン")).toBeInTheDocument();
    expect(screen.getByText("イタリアン")).toBeInTheDocument();
    expect(screen.getByText("本場の味を堪能")).toBeInTheDocument();
    expect(screen.getByText("東京都新宿区")).toBeInTheDocument();

    // 画像が2枚表示されていることを確認
    const images = screen.getAllByRole("img");
    expect(images.length).toBe(2);
    expect(images[0]).toHaveAttribute("alt", "テスト居酒屋");
    expect(images[1]).toHaveAttribute("alt", "イタリアンレストラン");
  });

  // ===================================================================
  // Triangulation Point 3: クリックイベント
  // ===================================================================
  // 【目的】カードをクリックした時、適切なコールバックが呼ばれることを確認
  // 【重要性】ユーザーインタラクションの動作保証、詳細画面への遷移ロジック
  // 【差分】カードクリック → onSelectRestaurant(restaurant) 呼び出し
  it("レストランカードをクリックするとonSelectRestaurantが呼ばれる", () => {
    // Arrange（準備）: レストランリストをレンダリング
    render(
      <RestaurantList
        restaurants={mockRestaurants}
        onSelectRestaurant={mockOnSelectRestaurant as any}
      />
    );

    // Act（実行）: 1件目のカードをクリック
    // getByTextで要素を見つけてクリック（カード全体がクリッカブル）
    const firstCard = screen
      .getByText("テスト居酒屋")
      .closest("div")?.parentElement;
    if (firstCard) {
      fireEvent.click(firstCard);
    }

    // Assert（検証）: onSelectRestaurantが正しい引数で呼ばれたか
    expect(mockOnSelectRestaurant).toHaveBeenCalledTimes(1);
    expect(mockOnSelectRestaurant).toHaveBeenCalledWith(mockRestaurants[0]);

    // モックをリセット
    mockOnSelectRestaurant.mockClear();

    // Act（実行）: 2件目のカードをクリック
    const secondCard = screen
      .getByText("イタリアンレストラン")
      .closest("div")?.parentElement;
    if (secondCard) {
      fireEvent.click(secondCard);
    }

    // Assert（検証）: 2件目のレストランで呼ばれたか
    expect(mockOnSelectRestaurant).toHaveBeenCalledTimes(1);
    expect(mockOnSelectRestaurant).toHaveBeenCalledWith(mockRestaurants[1]);
  });
});

// ===================================================================
// Triangulationのまとめ
// ===================================================================
// Point 1: 空配列 → 空メッセージ表示（境界値テスト）
// Point 2: データあり → カード表示（基本機能）
// Point 3: クリック → コールバック呼び出し（インタラクション）
// ===================================================================
// この3つのテストで、RestaurantListの主要な動作が網羅的にテストされる
// （条件分岐、表示ロジック、イベントハンドラー）
// ===================================================================
