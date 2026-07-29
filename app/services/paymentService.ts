import type { ApiResponse } from "~/composables/useApiClients";

// Payment status enum
export type PaymentStatus = "pending" | "paid" | "failed" | "cancelled" | "expired";

// Payment method types (matching backend)
export type PaymentMethodType = "credit_card" | "bank_transfer" | "e_wallet" | "qris" | "cstore";

// Frontend payment method options (UI-facing)
export interface PaymentMethodOption {
  id: string;
  name: string;
  description: string;
  icon: string;
  color: "blue" | "green" | "purple" | "orange";
  code: PaymentMethodType;
}

// API Response DTOs
export interface PaymentApiResponse {
  id: string;
  enrollmentId: string;
  userId: string;
  orderId: string;
  amount: number;
  paymentMethod: PaymentMethodType;
  status: PaymentStatus;
  paymentUrl?: string;
  paidAt?: string;
  expiresAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PaymentDetailResponse {
  orderId: string;
  paymentType: PaymentMethodType;
  transactionId: string;
  grossAmount: number;
  status: PaymentStatus;
  paymentUrl: string;
  expiryTime?: string;
}

export interface PaymentListResponse {
  data: PaymentApiResponse[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Create payment request data
export interface CreatePaymentData {
  enrollmentId: string;
  userId: string;
  amount: number;
  paymentMethod: PaymentMethodType;
}

// Pagination result
export interface PaginatedPaymentsResult {
  payments: PaymentApiResponse[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Map API response to frontend model
export const mapApiToPayment = (item: PaymentApiResponse): PaymentApiResponse => {
  return {
    id: item.id,
    enrollmentId: item.enrollmentId,
    userId: item.userId,
    orderId: item.orderId,
    amount: item.amount,
    paymentMethod: item.paymentMethod,
    status: item.status,
    paymentUrl: item.paymentUrl,
    paidAt: item.paidAt,
    expiresAt: item.expiresAt,
    createdAt: item.createdAt,
    updatedAt: item.updatedAt,
  };
};

// Frontend payment method options (matches payment-method.vue)
export const paymentMethodOptions: PaymentMethodOption[] = [
  {
    id: "va",
    name: "Virtual Account (VA)",
    description: "Transfer via BCA, Mandiri, BNI, or BRI virtual account",
    icon: "i-lucide-building",
    color: "blue" as const,
    code: "bank_transfer" as PaymentMethodType,
  },
  {
    id: "qris",
    name: "QRIS",
    description: "Scan QR code with any e-wallet",
    icon: "i-lucide-qr-code",
    color: "green" as const,
    code: "qris" as PaymentMethodType,
  },
  {
    id: "bank_transfer",
    name: "Bank Transfer",
    description: "Direct transfer from your bank account",
    icon: "i-lucide-landmark",
    color: "purple" as const,
    code: "bank_transfer" as PaymentMethodType,
  },
  {
    id: "ewallet",
    name: "E-Wallet",
    description: "GoPay, OVO, DANA, or LinkAja",
    icon: "i-lucide-wallet",
    color: "orange" as const,
    code: "e_wallet" as PaymentMethodType,
  },
];

// Map frontend payment method ID to backend code
export const mapMethodIdToCode = (methodId: string): PaymentMethodType => {
  const method = paymentMethodOptions.find((m) => m.id === methodId);
  return method?.code || "bank_transfer";
};

// Map backend payment method code to frontend ID
export const mapMethodCodeToId = (code: PaymentMethodType): string => {
  const method = paymentMethodOptions.find((m) => m.code === code);
  return method?.id || "va";
};

// Payment service
export const paymentService = {
  // Create a new payment
  async create(data: CreatePaymentData): Promise<PaymentApiResponse | null> {
    const { booking, extractData } = useApiClients();
    try {
      const response = await booking<ApiResponse<PaymentApiResponse>>(
        "/payments/transactions",
        {
          method: "POST",
          body: data,
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // Get payment by order ID
  async getByOrderId(orderId: string): Promise<PaymentApiResponse | null> {
    const { booking, extractData } = useApiClients();
    try {
      const response = await booking<ApiResponse<PaymentApiResponse>>(
        `/payments/order/${orderId}`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // Get payment by ID
  async getById(id: string): Promise<PaymentApiResponse | null> {
    const { booking, extractData } = useApiClients();
    try {
      const response = await booking<ApiResponse<PaymentApiResponse>>(
        `/payments/${id}`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // Get payment details by order ID
  async getDetails(orderId: string): Promise<PaymentDetailResponse | null> {
    const { booking, extractData } = useApiClients();
    try {
      const response = await booking<ApiResponse<PaymentDetailResponse>>(
        `/payments/${orderId}/details`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // List user payments
  async listUserPayments(
    userId?: string,
    page: number = 1,
    limit: number = 20,
  ): Promise<PaginatedPaymentsResult> {
    const { booking, extractPaginatedData } = useApiClients();
    const params = new URLSearchParams();
    if (userId) params.append("userId", userId);
    params.append("page", page.toString());
    params.append("limit", limit.toString());

    const response = await booking<{
      success: boolean;
      message: string;
      data: PaymentApiResponse[];
      total: number;
      page: number;
      limit: number;
      totalPages: number;
    }>(`/payments?${params.toString()}`, {
      method: "GET",
    });

    const data = extractPaginatedData(response);
    return {
      payments: Array.isArray(data) ? data.map(mapApiToPayment) : [],
      pagination: {
        page: response.page,
        limit: response.limit,
        total: response.total,
        totalPages: response.totalPages,
      },
    };
  },

  // Cancel a payment
  async cancel(orderId: string): Promise<boolean> {
    const { booking } = useApiClients();
    try {
      await booking(`/payments/${orderId}/cancel`, {
        method: "POST",
      });
      return true;
    } catch {
      return false;
    }
  },

  // Check payment status
  async checkStatus(orderId: string): Promise<PaymentStatus | null> {
    const { booking, extractData } = useApiClients();
    try {
      const response = await booking<ApiResponse<{ status: PaymentStatus }>>(
        `/payments/${orderId}/status`,
        {
          method: "GET",
        },
      );
      return extractData(response)?.status || null;
    } catch {
      return null;
    }
  },

  // Simulate payment (sandbox / dev)
  async simulate(orderId: string): Promise<boolean> {
    const { booking } = useApiClients();
    try {
      await booking(`/payments/order/${orderId}/simulate`, {
        method: "POST",
      });
      return true;
    } catch {
      try {
        await booking(`/payments/${orderId}/simulate`, {
          method: "POST",
        });
        return true;
      } catch {
        return false;
      }
    }
  },
};
