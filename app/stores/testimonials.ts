import { defineStore } from "pinia";

export type TestimonialStatus = "draft" | "pending" | "published" | "archived";

export interface Testimonial {
  id: string;
  userId: string;
  userName: string;
  userImage: string;
  userRole: string;
  content: string;
  rating: number;
  tags: string;
  status: TestimonialStatus;
  isFeatured: boolean;
  addedBy: string;
  addedAt: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

interface TestimonialsState {
  testimonials: Testimonial[];
  isLoading: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

const initialTestimonials: Testimonial[] = [
  {
    id: "11111111-1111-1111-1111-111111111101",
    userId: "user-001",
    userName: "Ahmad Wijaya",
    userImage: "https://i.pravatar.cc/150?u=ahmad",
    userRole: "Student",
    content:
      "Best driving school in town! The instructors are very patient and professional. I passed my SIM A test on the first try.",
    rating: 5,
    tags: "SIM A,Professional,Patient",
    status: "published",
    isFeatured: true,
    addedBy: "admin-001",
    addedAt: "2026-04-01T10:00:00Z",
    sortOrder: 1,
    createdAt: "2026-04-01T10:00:00Z",
    updatedAt: "2026-04-01T10:00:00Z",
  },
  {
    id: "11111111-1111-1111-1111-111111111102",
    userId: "user-002",
    userName: "Siti Rahayu",
    userImage: "https://i.pravatar.cc/150?u=siti",
    userRole: "Student",
    content:
      "Excellent training program with flexible schedules. Night sessions were very helpful for learning driving at night.",
    rating: 5,
    tags: "Flexible,Night Session,Helpful",
    status: "published",
    isFeatured: true,
    addedBy: "admin-001",
    addedAt: "2026-04-03T14:30:00Z",
    sortOrder: 2,
    createdAt: "2026-04-03T14:30:00Z",
    updatedAt: "2026-04-03T14:30:00Z",
  },
  {
    id: "11111111-1111-1111-1111-111111111103",
    userId: "user-003",
    userName: "Budi Santoso",
    userImage: "https://i.pravatar.cc/150?u=budi",
    userRole: "Student",
    content:
      "Very clean cars and modern teaching methods. The 10x package gave me enough practice to be confident on the road.",
    rating: 4,
    tags: "Clean Cars,Modern,Confident",
    status: "published",
    isFeatured: false,
    addedBy: "admin-001",
    addedAt: "2026-04-05T09:15:00Z",
    sortOrder: 3,
    createdAt: "2026-04-05T09:15:00Z",
    updatedAt: "2026-04-05T09:15:00Z",
  },
  {
    id: "11111111-1111-1111-1111-111111111104",
    userId: "user-004",
    userName: "Dewi Lestari",
    userImage: "https://i.pravatar.cc/150?u=dewi",
    userRole: "Student",
    content:
      "Had a great experience with weekend sessions. The instructors accommodated my schedule perfectly.",
    rating: 5,
    tags: "Weekend,Flexible,Great Experience",
    status: "pending",
    isFeatured: false,
    addedBy: "admin-001",
    addedAt: "2026-04-08T16:45:00Z",
    sortOrder: 4,
    createdAt: "2026-04-08T16:45:00Z",
    updatedAt: "2026-04-08T16:45:00Z",
  },
  {
    id: "11111111-1111-1111-1111-111111111105",
    userId: "user-005",
    userName: "Rizky Pratama",
    userImage: "https://i.pravatar.cc/150?u=rizky",
    userRole: "Student",
    content:
      "The process from registration to certification was seamless. Highly recommended!",
    rating: 4,
    tags: "Seamless,Recommended,Certification",
    status: "pending",
    isFeatured: false,
    addedBy: "admin-001",
    addedAt: "2026-04-10T11:20:00Z",
    sortOrder: 5,
    createdAt: "2026-04-10T11:20:00Z",
    updatedAt: "2026-04-10T11:20:00Z",
  },
];

export const useTestimonialsStore = defineStore("testimonials", {
  state: (): TestimonialsState => ({
    testimonials: [],
    isLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
  }),

  getters: {
    publishedTestimonials: (state) =>
      state.testimonials.filter((t) => t.status === "published"),
    pendingTestimonials: (state) =>
      state.testimonials.filter((t) => t.status === "pending"),
    archivedTestimonials: (state) =>
      state.testimonials.filter((t) => t.status === "archived"),
    featuredTestimonials: (state) =>
      state.testimonials.filter(
        (t) => t.isFeatured && t.status === "published",
      ),
    totalTestimonials: (state) => state.testimonials.length,
    totalPublished: (state) =>
      state.testimonials.filter((t) => t.status === "published").length,
    totalPending: (state) =>
      state.testimonials.filter((t) => t.status === "pending").length,
    averageRating: (state) => {
      const published = state.testimonials.filter(
        (t) => t.status === "published",
      );
      if (published.length === 0) return 0;
      const sum = published.reduce((acc, t) => acc + t.rating, 0);
      return (sum / published.length).toFixed(1);
    },
    getTestimonialById: (state) => (id: string) =>
      state.testimonials.find((t) => t.id === id),
    getTestimonialsByStatus: (state) => (status: TestimonialStatus) =>
      state.testimonials.filter((t) => t.status === status),
  },

  actions: {
    async fetchTestimonials() {
      this.isLoading = true;
      this.error = null;

      try {
        const { testimonialService } =
          await import("~/services/testimonialService");
        const result = await testimonialService.fetchAll();

        this.testimonials = result.testimonials;
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch testimonials";
        console.error("Error fetching testimonials:", err);
        // API failed, use dummy data
        this.testimonials = [...initialTestimonials];
      } finally {
        this.isLoading = false;
      }
    },

    addTestimonial(data: Omit<Testimonial, "id" | "createdAt" | "updatedAt">) {
      const newTestimonial: Testimonial = {
        ...data,
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      this.testimonials.unshift(newTestimonial);
      return newTestimonial;
    },

    updateTestimonial(id: string, data: Partial<Testimonial>) {
      const index = this.testimonials.findIndex((t) => t.id === id);
      if (index !== -1) {
        const existing = this.testimonials[index];
        this.testimonials[index] = {
          ...existing,
          ...data,
          updatedAt: new Date().toISOString(),
        } as Testimonial;
        return this.testimonials[index];
      }
      return null;
    },

    deleteTestimonial(id: string) {
      this.testimonials = this.testimonials.filter((t) => t.id !== id);
    },

    toggleFeatured(id: string) {
      const testimonial = this.testimonials.find((t) => t.id === id);
      if (testimonial) {
        testimonial.isFeatured = !testimonial.isFeatured;
        testimonial.updatedAt = new Date().toISOString();
        return testimonial.isFeatured;
      }
      return null;
    },

    changeStatus(id: string, status: TestimonialStatus) {
      const testimonial = this.testimonials.find((t) => t.id === id);
      if (testimonial) {
        testimonial.status = status;
        testimonial.updatedAt = new Date().toISOString();
        return testimonial.status;
      }
      return null;
    },

    publishTestimonial(id: string) {
      return this.changeStatus(id, "published");
    },

    archiveTestimonial(id: string) {
      return this.changeStatus(id, "archived");
    },

    reorderTestimonials(fromIndex: number, toIndex: number) {
      const testimonial = this.testimonials.splice(fromIndex, 1)[0];
      if (!testimonial) return;
      this.testimonials.splice(toIndex, 0, testimonial);
      // Update sort orders
      this.testimonials.forEach((t, i) => {
        t.sortOrder = i + 1;
        t.updatedAt = new Date().toISOString();
      });
    },

    reset() {
      this.testimonials = [...initialTestimonials];
      this.isLoading = false;
      this.error = null;
    },
  },
});
