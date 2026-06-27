import type { ApiResponse } from "~/composables/useApiClients";
import type { Faq, CreateFaqData, UpdateFaqData } from "~/stores/content";

// ============================================================
// Types matching Go backend DTOs exactly (blog_post.go)
// ============================================================

/** Mirrors dto.FileUpload in Go – base64-encoded file sent in JSON body */
export interface FileUpload {
  data: string;     // base64 encoded file data (without data: URI prefix)
  fileName: string; // e.g. "cover.jpg"
  mimeType: string; // e.g. "image/jpeg"
}

/** Mirrors dto.Publishing */
export interface Publishing {
  status: "draft" | "published" | "archived";
  publishedAt?: string | null;
  scheduledAt?: string | null;
}

/** Mirrors dto.Attractiveness */
export interface Attractiveness {
  isFeatured: boolean;
  isSpotlight: boolean;
  priority: number;
  highlight: boolean;
}

// -------------------------------------------------------
// Request DTOs (matching Go CreateBlogPostRequest)
// -------------------------------------------------------
export interface CreateBlogPostData {
  title: string;
  slug?: string;
  author?: string;
  leadParagraph?: string;
  content?: string;
  authorId: string;
  /** Top-level status mirror (backend reads this for form-data compat) */
  status?: string;
  featuredImage?: FileUpload;
  publishing?: Publishing;
  attractiveness?: Attractiveness;
}

// -------------------------------------------------------
// Request DTOs (matching Go UpdateBlogPostRequest)
// -------------------------------------------------------
export interface UpdateBlogPostData {
  title?: string;
  slug?: string;
  author?: string;
  leadParagraph?: string;
  content?: string;
  authorId: string;
  status?: string;
  featuredImage?: FileUpload;
  publishing?: Publishing;
  attractiveness?: Attractiveness;
}

// -------------------------------------------------------
// Response DTOs (matching Go BlogPostResponse)
// -------------------------------------------------------
export interface BlogPostResponse {
  id: string;
  title: string;
  slug: string;
  author: string;
  content: string;
  excerpt: string;
  featuredImage: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
  viewCount: number;
  likeCount: number;
  shareCount: number;
  readingTime: number;
  createdAt: string;
  updatedAt: string;
}

// -------------------------------------------------------
// List response (matching Go BlogPostBrief)
// -------------------------------------------------------
export interface BlogPostBrief {
  id: string;
  title: string;
  slug: string;
  author: string;
  excerpt: string;
  content?: string;
  featuredImage?: string;
  readingTime: number;
  status: string;
  publishedAt?: string | null;
  viewCount: number;
  createdAt: string;
  /** Nested publishing object – present when the list returns BlogPostResponse */
  publishing?: Publishing;
  /** Nested attractiveness object – present when the list returns BlogPostResponse */
  attractiveness?: Attractiveness;
}

export interface BlogPostListResponse {
  posts: BlogPostBrief[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface BlogPostFilterParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  author?: string;
  isFeatured?: boolean;
  isSpotlight?: boolean;
  sortBy?: "createdAt" | "updatedAt" | "viewCount" | "title";
  sortOrder?: "asc" | "desc";
}

// ============================================================
// Helper: parse a data: URI into a FileUpload object
// Returns null if the value is not a data URI (e.g. it's already a URL)
// ============================================================
export function dataUriToFileUpload(
  dataUri: string,
  fileName: string,
): FileUpload | null {
  // e.g. "data:image/jpeg;base64,/9j/4AAQ..."
  const match = dataUri.match(/^data:([^;]+);base64,(.+)$/);
  if (!match) return null;
  return {
    mimeType: match[1]!,
    data: match[2]!,
    fileName,
  };
}

// ============================================================
// Service
// ============================================================
export const contentService = {
  // ==================== BLOG POST METHODS ====================

  /** GET /articles/blog – returns paginated list */
  async fetchBlogPosts(
    params: BlogPostFilterParams = {},
  ): Promise<BlogPostListResponse> {
    const { core } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.search) queryParams.set("search", params.search);
    if (params.status) queryParams.set("status", params.status);
    if (params.author) queryParams.set("author", params.author);
    if (params.isFeatured !== undefined)
      queryParams.set("isFeatured", String(params.isFeatured));
    if (params.isSpotlight !== undefined)
      queryParams.set("isSpotlight", String(params.isSpotlight));
    if (params.sortBy) queryParams.set("sortBy", params.sortBy);
    if (params.sortOrder) queryParams.set("sortOrder", params.sortOrder);

    const queryString = queryParams.toString();
    const url = `/articles/blog${queryString ? `?${queryString}` : ""}`;

    try {
      const response: any = await core(url, { method: "GET" });

      // Backend wraps the response as:
      // { success, message, data: { data: BlogPostResponse[], pagination: { total, page, limit, totalPages } } }
      const inner = response?.data ?? response;
      const rawPosts: any[] = Array.isArray(inner?.data)
        ? inner.data
        : Array.isArray(inner)
          ? inner
          : [];

      const pagination = inner?.pagination ?? {};

      const posts: BlogPostBrief[] = rawPosts.map((p: any) => ({
        id: p.id,
        title: p.title,
        slug: p.slug || "",
        author: p.author || "Admin",
        excerpt: p.excerpt || p.leadParagraph || "",
        content: p.content || "",
        featuredImage: p.featuredImage || "",
        readingTime: p.readingTime || 1,
        status: p.publishing?.status || p.status || "draft",
        publishedAt: p.publishing?.publishedAt ?? p.publishedAt ?? null,
        viewCount: p.viewCount || 0,
        createdAt: p.createdAt || new Date().toISOString(),
        publishing: p.publishing,
        attractiveness: p.attractiveness,
      }));

      return {
        posts,
        total: pagination.total ?? posts.length,
        page: pagination.page ?? 1,
        limit: pagination.limit ?? 10,
        totalPages: pagination.totalPages ?? 1,
      };
    } catch (error) {
      console.error("Error fetching blog posts:", error);
      return { posts: [], total: 0, page: 1, limit: 10, totalPages: 0 };
    }
  },

  /** GET /articles/blog/:id */
  async fetchBlogPostById(id: string): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core(`/articles/blog/${id}`, {
        method: "GET",
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  /** POST /articles/blog – JSON body with optional FileUpload for featuredImage */
  async createBlogPost(
    data: CreateBlogPostData,
  ): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core("/articles/blog", {
        method: "POST",
        body: JSON.stringify(data),
        headers: { "Content-Type": "application/json" },
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  /** PUT /articles/blog/:id – JSON body with optional FileUpload for featuredImage */
  async updateBlogPost(
    id: string,
    data: UpdateBlogPostData,
  ): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core(`/articles/blog/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
        headers: { "Content-Type": "application/json" },
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  /** DELETE /articles/blog/:id */
  async deleteBlogPost(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/blog/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // ==================== FAQ METHODS ====================

  async fetchFaqs(): Promise<Faq[]> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq[]> | Faq[]>("/articles/faq", {
        method: "GET",
      });
      if (Array.isArray(response)) return response;
      if (response && Array.isArray((response as any).data)) {
        return (response as ApiResponse<Faq[]>).data;
      }
      return [];
    } catch {
      return [];
    }
  },

  async fetchFaqById(id: string): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core(`/articles/faq/${id}`, {
        method: "GET",
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  async createFaq(data: CreateFaqData): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core("/articles/faq", {
        method: "POST",
        body: data,
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  async updateFaq(id: string, data: UpdateFaqData): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response: any = await core(`/articles/faq/${id}`, {
        method: "PUT",
        body: data,
      });
      if (response && typeof response === "object") {
        return response.data ?? response;
      }
      return null;
    } catch {
      return null;
    }
  },

  async deleteFaq(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/faq/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  async reorderFaq(id: string, orderNum: number): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/faq/${id}/reorder`, {
        method: "POST",
        body: { newOrder: orderNum },
      });
      return true;
    } catch {
      return false;
    }
  },

  // ==================== VIEW COUNT ====================

  async incrementBlogPostViewCount(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/${id}/view`, { method: "POST" });
      return true;
    } catch {
      return false;
    }
  },
};
