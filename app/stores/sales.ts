import { defineStore } from "pinia";

export interface Transaction {
  id: number;
  studentName: string;
  packageId: number;
  purchaseDate: string;
  amount: number;
  status: "Completed" | "Pending" | "Refunded";
}

export interface PackageSummary {
  id: number;
  name: string;
  price: number;
  color: string;
  totalSales: number;
  revenue: number;
}

interface SalesState {
  transactions: Transaction[];
  startDate: string;
  endDate: string;
  isLoading: boolean;
  error: string | null;
}

const initialTransactions: Transaction[] = [
  {
    id: 101,
    studentName: "John Doe",
    packageId: 1,
    purchaseDate: "2026-03-10",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 102,
    studentName: "Jane Smith",
    packageId: 1,
    purchaseDate: "2026-03-12",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 103,
    studentName: "Budi Santoso",
    packageId: 2,
    purchaseDate: "2026-03-15",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 104,
    studentName: "Amanda Chen",
    packageId: 1,
    purchaseDate: "2026-03-20",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 105,
    studentName: "David Lee",
    packageId: 2,
    purchaseDate: "2026-03-22",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 106,
    studentName: "Sarah Putri",
    packageId: 3,
    purchaseDate: "2026-04-01",
    amount: 2250000,
    status: "Completed",
  },
  {
    id: 107,
    studentName: "Michael Brown",
    packageId: 1,
    purchaseDate: "2026-04-02",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 108,
    studentName: "Emily Davis",
    packageId: 4,
    purchaseDate: "2026-04-05",
    amount: 2650000,
    status: "Completed",
  },
  {
    id: 109,
    studentName: "Ricky Wijaya",
    packageId: 2,
    purchaseDate: "2026-04-10",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 110,
    studentName: "Anita Sari",
    packageId: 3,
    purchaseDate: "2026-04-12",
    amount: 2250000,
    status: "Completed",
  },
];

export const useSalesStore = defineStore("sales", {
  state: (): SalesState => ({
    transactions: initialTransactions,
    startDate: "",
    endDate: "",
    isLoading: false,
    error: null,
  }),

  getters: {
    filteredTransactions: (state) => {
      return state.transactions.filter((t) => {
        const pDate = new Date(t.purchaseDate);
        if (state.startDate && pDate < new Date(state.startDate)) return false;
        if (state.endDate && pDate > new Date(state.endDate)) return false;
        return true;
      });
    },

    transactionsByPackage: (state) => (packageId: string) => {
      return state.transactions.filter((t) => t.packageId === packageId);
    },

    totalRevenue: (state) => {
      return state.transactions.reduce((sum, t) => sum + t.amount, 0);
    },

    filteredTotalRevenue(): number {
      return this.filteredTransactions.reduce((sum, t) => sum + t.amount, 0);
    },

    totalSales: (state) => state.transactions.length,

    filteredTotalSales(): number {
      return this.filteredTransactions.length;
    },

    packageSummary(): PackageSummary[] {
      // This assumes packages are imported or defined elsewhere
      // For now, we'll create a dynamic summary from transactions
      const packageMap = new Map<
        number,
        { totalSales: number; revenue: number }
      >();

      this.filteredTransactions.forEach((t) => {
        const existing = packageMap.get(t.packageId) || {
          totalSales: 0,
          revenue: 0,
        };
        packageMap.set(t.packageId, {
          totalSales: existing.totalSales + 1,
          revenue: existing.revenue + t.amount,
        });
      });

      return Array.from(packageMap.entries()).map(([id, data]) => ({
        id,
        name: `Package ${id}`,
        price: 0,
        color: "neutral",
        totalSales: data.totalSales,
        revenue: data.revenue,
      }));
    },

    revenueByPackage(): {
      packageId: number;
      revenue: number;
      sales: number;
    }[] {
      const map = new Map<number, { revenue: number; sales: number }>();

      this.filteredTransactions.forEach((t) => {
        const existing = map.get(t.packageId) || { revenue: 0, sales: 0 };
        map.set(t.packageId, {
          revenue: existing.revenue + t.amount,
          sales: existing.sales + 1,
        });
      });

      return Array.from(map.entries())
        .map(([packageId, data]) => ({ packageId, ...data }))
        .sort((a, b) => b.revenue - a.revenue);
    },
  },

  actions: {
    setDateRange(start: string, end: string) {
      this.startDate = start;
      this.endDate = end;
    },

    clearDateFilter() {
      this.startDate = "";
      this.endDate = "";
    },

    addTransaction(data: Omit<Transaction, "id">) {
      const newId = Math.max(...this.transactions.map((t) => t.id), 0) + 1;
      const transaction: Transaction = {
        ...data,
        id: newId,
      };
      this.transactions.unshift(transaction);
      return transaction;
    },

    getTransactionById(id: number) {
      return this.transactions.find((t) => t.id === id);
    },

    getTransactionsByPackage(packageId: number) {
      return this.transactions.filter((t) => t.packageId === packageId);
    },

    getPackageRevenue(packageId: number) {
      return this.transactions
        .filter((t) => t.packageId === packageId)
        .reduce((sum, t) => sum + t.amount, 0);
    },

    getPackageSales(packageId: number) {
      return this.transactions.filter((t) => t.packageId === packageId).length;
    },

    formatPrice(price: number) {
      return new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        maximumFractionDigits: 0,
      }).format(price);
    },

    formatDate(dateString: string) {
      return new Date(dateString).toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
    },

    reset() {
      this.transactions = [...initialTransactions];
      this.startDate = "";
      this.endDate = "";
      this.isLoading = false;
      this.error = null;
    },
  },
});
