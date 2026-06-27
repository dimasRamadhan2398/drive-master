import type { Testimonial, TestimonialStatus } from "~/stores/testimonials";
import type { ApiResponse } from "~/composables/useApiClients";

export interface TestimonialApiResponse {
  id: string;
  user_id: string;
  user_name: string;
  user_image: string;
  user_role: string;
  content: string;
  rating: number;
  tags: string;
  status: TestimonialStatus;
  is_featured: boolean;
  added_by: string;
  added_at: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface PaginatedTestimonialsResult {
  testimonials: Testimonial[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface CreateTestimonialData {
  userId: string;
  userName: string;
  userImage?: string;
  userRole?: string;
  content: string;
  rating: number;
  tags?: string;
  status?: TestimonialStatus;
  isFeatured?: boolean;
  addedBy: string;
}

export interface UpdateTestimonialData {
  userName?: string;
  userImage?: string;
  userRole?: string;
  content?: string;
  rating?: number;
  tags?: string;
  status?: TestimonialStatus;
  isFeatured?: boolean;
  sortOrder?: number;
}

export const mapApiToTestimonial = (
  item: TestimonialApiResponse,
): Testimonial => {
  return {
    id: String(item.id),
    userId: item.user_id,
    userName: item.user_name,
    userImage: item.user_image || "",
    userRole: item.user_role || "Student",
    content: item.content,
    rating: item.rating,
    tags: item.tags || "",
    status: item.status,
    isFeatured: item.is_featured || false,
    addedBy: item.added_by,
    addedAt: item.added_at,
    sortOrder: item.sort_order || 0,
    createdAt: item.created_at,
    updatedAt: item.updated_at,
  };
};

export const testimonialService = {
  // Fetch all testimonials
  async fetchAll(): Promise<PaginatedTestimonialsResult> {
    const { core, extractData } = useApiClients();

    const response = await core<{
      success: boolean;
      message: string;
      data: TestimonialApiResponse[];
    }>("/testimonials", {
      method: "GET",
    });

    const data = extractData(response);
    return {
      testimonials: Array.isArray(data) ? data.map(mapApiToTestimonial) : [],
      pagination: {
        page: 1,
        limit: data.length,
        total: data.length,
        totalPages: 1,
      },
    };
  },

  async fetchById(id: string): Promise<Testimonial | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<TestimonialApiResponse>>(
        `/testimonials/${id}`,
        {
          method: "GET",
        },
      );
      return mapApiToTestimonial(extractData(response));
    } catch {
      return null;
    }
  },

  async fetchPublished(): Promise<Testimonial[]> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<{
        success: boolean;
        message: string;
        data: TestimonialApiResponse[];
      }>("/testimonials/published", { method: "GET" });
      return (extractData(response) as TestimonialApiResponse[]).map(
        mapApiToTestimonial,
      );
    } catch {
      return [];
    }
  },

  async fetchFeatured(): Promise<Testimonial[]> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<{
        success: boolean;
        message: string;
        data: TestimonialApiResponse[];
      }>("/testimonials/featured", { method: "GET" });
      return (extractData(response) as TestimonialApiResponse[]).map(
        mapApiToTestimonial,
      );
    } catch {
      return [];
    }
  },

  async create(data: CreateTestimonialData): Promise<Testimonial | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<TestimonialApiResponse>>(
        "/testimonials",
        {
          method: "POST",
          body: data,
        },
      );
      return mapApiToTestimonial(extractData(response));
    } catch {
      return null;
    }
  },

  async update(
    id: string,
    data: UpdateTestimonialData,
  ): Promise<Testimonial | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<TestimonialApiResponse>>(
        `/testimonials/${id}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return mapApiToTestimonial(extractData(response));
    } catch {
      return null;
    }
  },

  async delete(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/testimonials/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  async publish(id: string): Promise<boolean> {
    return this.updateStatus(id, "published");
  },

  async archive(id: string): Promise<boolean> {
    return this.updateStatus(id, "archived");
  },

  async updateStatus(id: string, status: TestimonialStatus): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/testimonials/${id}/status`, {
        method: "PATCH",
        body: { status },
      });
      return true;
    } catch {
      return false;
    }
  },

  async toggleFeatured(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/testimonials/${id}/toggle-featured`, { method: "PATCH" });
      return true;
    } catch {
      return false;
    }
  },
};
