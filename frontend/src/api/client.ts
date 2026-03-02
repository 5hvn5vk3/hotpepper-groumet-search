import type { APIErrorResponse } from "@/types";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

// HTTPステータスとユーザー向けメッセージを保持するカスタムエラークラス
export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function apiGet<T>(endpoint: string): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  try {
    const response = await fetch(url, {
      method: "GET",
    });
    if (!response.ok) {
      const body: APIErrorResponse | null = await response
        .json()
        .catch(() => null);
      const message =
        body?.error?.message ?? "サービスが一時的に利用できません";
      throw new ApiError(response.status, message);
    }
    return response.json();
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error("An unexpected error occurred");
  }
}
