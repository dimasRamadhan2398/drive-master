import type { ApiResponse } from "~/composables/useApiClients";

export interface TransactionApiResponse {
  id: string;
  paymentId: string;
  type: "charge" | "refund" | "reversal";
  status: "pending" | "success" | "failed" | "reversed";
  amount: number;
  currency: string;
  gateway: string;
  gatewayTxnId: string;
  processedAt?: string;
  createdAt: string;
  updatedAt: string;
  payment?: {
    id: string;
    orderId: string;
    bookingId?: string; // This corresponds to enrollmentId / package reference
    userId: string;
    description?: string;
    status: string;
  };
}

export interface PaginatedTransactionsResult {
  transactions: TransactionApiResponse[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export const salesService = {
  // GET /transactions - List transactions with pagination
  async listTransactions(
    page: number = 1,
    limit: number = 20,
  ): Promise<PaginatedTransactionsResult> {
    const { core, extractData } = useApiClients();
    const params = new URLSearchParams();
    params.append("page", page.toString());
    params.append("limit", limit.toString());

    try {
      const response = await core<
        ApiResponse<{
          data: TransactionApiResponse[];
          total: number;
          page: number;
          limit: number;
        }>
      >(`/transactions?${params.toString()}`, {
        method: "GET",
      });

      const result = extractData(response);
      const totalPages = Math.ceil((result?.total || 0) / (result?.limit || limit));

      return {
        transactions: Array.isArray(result?.data) ? result.data : [],
        pagination: {
          page: result?.page || page,
          limit: result?.limit || limit,
          total: result?.total || 0,
          totalPages: totalPages || 1,
        },
      };
    } catch (error) {
      console.error("Failed to fetch transactions:", error);
      return {
        transactions: [],
        pagination: {
          page,
          limit,
          total: 0,
          totalPages: 1,
        },
      };
    }
  },
};
