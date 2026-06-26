import { defineStore } from "pinia";
import { salesService } from "~/services/salesService";
import { useStudentsStore } from "~/stores/students";
import { usePackagesStore } from "~/stores/packages";

export interface Transaction {
  id: string | number;
  studentName: string;
  packageId: string | number;
  purchaseDate: string;
  amount: number;
  status: "Completed" | "Pending" | "Refunded";
}

export interface PackageSummary {
  id: string | number;
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
    packageId: "11111111-1111-1111-1111-111111111101",
    purchaseDate: "2026-03-10",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 102,
    studentName: "Jane Smith",
    packageId: "11111111-1111-1111-1111-111111111101",
    purchaseDate: "2026-03-12",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 103,
    studentName: "Budi Santoso",
    packageId: "11111111-1111-1111-1111-111111111102",
    purchaseDate: "2026-03-15",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 104,
    studentName: "Amanda Chen",
    packageId: "11111111-1111-1111-1111-111111111101",
    purchaseDate: "2026-03-20",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 105,
    studentName: "David Lee",
    packageId: "11111111-1111-1111-1111-111111111102",
    purchaseDate: "2026-03-22",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 106,
    studentName: "Sarah Putri",
    packageId: "11111111-1111-1111-1111-111111111103",
    purchaseDate: "2026-04-01",
    amount: 2250000,
    status: "Completed",
  },
  {
    id: 107,
    studentName: "Michael Brown",
    packageId: "11111111-1111-1111-1111-111111111101",
    purchaseDate: "2026-04-02",
    amount: 1750000,
    status: "Completed",
  },
  {
    id: 108,
    studentName: "Emily Davis",
    packageId: "11111111-1111-1111-1111-111111111104",
    purchaseDate: "2026-04-05",
    amount: 2650000,
    status: "Completed",
  },
  {
    id: 109,
    studentName: "Ricky Wijaya",
    packageId: "11111111-1111-1111-1111-111111111102",
    purchaseDate: "2026-04-10",
    amount: 1950000,
    status: "Completed",
  },
  {
    id: 110,
    studentName: "Anita Sari",
    packageId: "11111111-1111-1111-1111-111111111103",
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

    transactionsByPackage: (state) => (packageId: string | number) => {
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
      const packageMap = new Map<
        string | number,
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
      packageId: string | number;
      revenue: number;
      sales: number;
    }[] {
      const map = new Map<string | number, { revenue: number; sales: number }>();

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

    async fetchTransactions(page: number = 1, limit: number = 100) {
      this.isLoading = true;
      this.error = null;
      try {
        const result = await salesService.listTransactions(page, limit);
        
        // Resolve student and packages stores
        const studentsStore = useStudentsStore();
        const packagesStore = usePackagesStore();
        
        // Map API response to Transaction model
        this.transactions = result.transactions.map((item) => {
          // Resolve student name
          const student = studentsStore.students.find((s) => s.id === item.payment?.userId) || 
                          studentsStore.allStudents.find((s) => s.id === item.payment?.userId);
          const studentName = student ? student.name : (item.payment?.description || "Student");
          
          // Resolve package ID (matching amount/price)
          const pkg = packagesStore.packages.find(
            (p) => p.price === item.amount || p.discountPrice === item.amount
          );
          const packageId = pkg ? pkg.id : (item.payment?.bookingId || "");
          
          // Map status
          let status: "Completed" | "Pending" | "Refunded" = "Pending";
          if (item.status === "success") {
            status = "Completed";
          } else if (item.status === "failed" || item.status === "reversed") {
            status = "Refunded";
          }

          return {
            id: item.id,
            studentName,
            packageId,
            purchaseDate: item.processedAt || item.createdAt,
            amount: item.amount,
            status,
          };
        });
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to fetch transactions";
        console.error("Error fetching transactions:", err);
      } finally {
        this.isLoading = false;
      }
    },

    addTransaction(data: Omit<Transaction, "id">) {
      const newId = this.transactions.length > 0 ? Math.max(...this.transactions.map((t) => typeof t.id === 'number' ? t.id : 0), 0) + 1 : 1;
      const transaction: Transaction = {
        ...data,
        id: newId,
      };
      this.transactions.unshift(transaction);
      return transaction;
    },

    getTransactionById(id: string | number) {
      return this.transactions.find((t) => t.id === id);
    },

    getTransactionsByPackage(packageId: string | number) {
      return this.transactions.filter((t) => t.packageId === packageId);
    },

    getPackageRevenue(packageId: string | number) {
      return this.transactions
        .filter((t) => t.packageId === packageId)
        .reduce((sum, t) => sum + t.amount, 0);
    },

    getPackageSales(packageId: string | number) {
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
