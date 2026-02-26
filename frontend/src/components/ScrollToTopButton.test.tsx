import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ScrollToTopButton } from "./ScrollToTopButton";

describe("ScrollToTopButton", () => {
  beforeEach(() => {
    Object.defineProperty(window, "scrollY", {
      configurable: true,
      writable: true,
      value: 0,
    });

    Object.defineProperty(window, "scrollTo", {
      configurable: true,
      writable: true,
      value: vi.fn(),
    });
  });

  it("規定のスクロール量を超えると表示される", async () => {
    render(<ScrollToTopButton />);

    expect(
      screen.queryByRole("button", { name: "ページ上部へ戻る" }),
    ).toBeNull();

    window.scrollY = 320;
    fireEvent.scroll(window);

    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "ページ上部へ戻る" }),
      ).toBeInTheDocument();
    });
  });

  it("クリックするとページ最上部へスムーズスクロールする", async () => {
    window.scrollY = 500;
    render(<ScrollToTopButton />);

    const button = await screen.findByRole("button", {
      name: "ページ上部へ戻る",
    });

    fireEvent.click(button);

    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 0,
      behavior: "smooth",
    });
  });
});
