import { defineStore } from "pinia";
import { paymentService, paymentMethodOptions } from "~/services/paymentService";
import type {
  PaymentApiResponse,
  PaymentStatus,
  PaymentMethodType,
  CreatePaymentData,
  PaymentMethodOption,
} from "~/services/paymentService";

interface PaymentsState {
  payments: PaymentApiResponse[];
  currentPayment: PaymentApiResponse | null;
  paymentDetails: {
    orderId: string;
    paymentType: PaymentMethodType;
    transactionId: string;
    grossAmount: number;
    status: PaymentStatus;
    paymentUrl: string;
    expiryTime?: string;
  } | null;
  isLoading: boolean;
  isCreating: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Initial fallback data for when API is unavailable
const initialPayments: PaymentApiResponse[] = [];

export const usePaymentsStore = defineStore("payments", {
  state: (): PaymentsState => ({
    payments: [],
    currentPayment: null,
    paymentDetails: null,
    isLoading: false,
    isCreating: false,
    error: null,
    pagination: {
      page: 1,
      limit: 20,
      total: 0,
      totalPages: 0,
    },
  }),

  getters: {
    pendingPayments: (state) =>
      state.payments.filter((p) => p.status === "pending"),

    completedPayments: (state) =>
      state.payments.filter((p) => p.status === "paid"),

    failedPayments: (state) =>
      state.payments.filter((p) => p.status === "failed"),

    cancelledPayments: (state) =>
      state.payments.filter((p) => p.status === "cancelled"),

    expiredPayments: (state) =>
      state.payments.filter((p) => p.status === "expired"),

    totalRevenue: (state) =>
      state.payments
        .filter((p) => p.status === "paid")
        .reduce((sum, p) => sum + p.amount, 0),

    getPaymentById: (state) => (id: string) =>
      state.payments.find((p) => p.id === id),

    getPaymentByOrderId: (state) => (orderId: string) =>
      state.payments.find((p) => p.orderId === orderId),

    paymentMethodOptions: () => paymentMethodOptions,

    isPaymentPending: (state) =>
      state.currentPayment?.status === "pending",

    isPaymentCompleted: (state) =>
      state.currentPayment?.status === "paid",

    isPaymentExpired: (state) =>
      state.currentPayment?.status === "expired",
  },

  actions: {
    // Create a new payment
    async createPayment(data: CreatePaymentData) {
      this.isCreating = true;
      this.error = null;

      try {
        const payment = await paymentService.create(data);
        if (payment) {
          this.payments.unshift(payment);
          this.currentPayment = payment;
          return payment;
        }
        this.error = "Failed to create payment";
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create payment";
        return null;
      } finally {
        this.isCreating = false;
      }
    },

    // Fetch payment by order ID
    async fetchPaymentByOrderId(orderId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const payment = await paymentService.getByOrderId(orderId);
        if (payment) {
          this.currentPayment = payment;
          // Update in list if exists
          const index = this.payments.findIndex((p) => p.id === payment.id);
          if (index !== -1) {
            this.payments[index] = payment;
          } else {
            this.payments.unshift(payment);
          }
          return payment;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch payment";
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch payment details
    async fetchPaymentDetails(orderId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const details = await paymentService.getDetails(orderId);
        if (details) {
          this.paymentDetails = details;
          return details;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch payment details";
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch user payments
    async fetchUserPayments(userId?: string, page: number = 1) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await paymentService.listUserPayments(userId, page);
        if (page === 1) {
          this.payments = result.payments;
        } else {
          this.payments = [...this.payments, ...result.payments];
        }
        this.pagination = result.pagination;
        return result;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch payments";
        // Fallback to empty array
        this.payments = [...initialPayments];
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // Cancel a payment
    async cancelPayment(orderId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const success = await paymentService.cancel(orderId);
        if (success) {
          // Update local state
          if (this.currentPayment?.orderId === orderId) {
            this.currentPayment.status = "cancelled";
          }
          const payment = this.payments.find((p) => p.orderId === orderId);
          if (payment) {
            payment.status = "cancelled";
          }
        }
        return success;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to cancel payment";
        return false;
      } finally {
        this.isLoading = false;
      }
    },

    // Check payment status
    async checkPaymentStatus(orderId: string) {
      try {
        const status = await paymentService.checkStatus(orderId);
        if (status && this.currentPayment?.orderId === orderId) {
          this.currentPayment.status = status;
        }
        return status;
      } catch {
        return null;
      }
    },

    // Set current payment manually
    setCurrentPayment(payment: PaymentApiResponse | null) {
      this.currentPayment = payment;
    },

    // Clear payment details
    clearPaymentDetails() {
      this.paymentDetails = null;
    },

    // Reset store
    reset() {
      this.payments = [];
      this.currentPayment = null;
      this.paymentDetails = null;
      this.isLoading = false;
      this.isCreating = false;
      this.error = null;
      this.pagination = {
        page: 1,
        limit: 20,
        total: 0,
        totalPages: 0,
      };
    },
  },
});

// Export payment method options type for convenience
export type { PaymentMethodOption };
