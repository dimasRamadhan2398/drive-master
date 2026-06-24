import type {
  ApiResponse,
} from "~/composables/useApiClients";
import type { Faq, CreateFaqData, UpdateFaqData } from "~/stores/content";

// Blog Post Types (mirroring backend DTOs - blog_post.go)
export interface BlogPostMedia {
  id?: string;
  name: string;
  type: string;
  size: string;
  url: string;
  fileType: "image" | "video";
  order?: number;
}

export interface Publishing {
  status: "draft" | "published" | "archived";
  publishedAt?: string;
  scheduledAt?: string;
}

export interface Attractiveness {
  isFeatured: boolean;
  isSpotlight: boolean;
  priority: number;
  highlight: boolean;
}

// Request DTOs (matching backend CreateBlogPostRequest, UpdateBlogPostRequest)
export interface CreateBlogPostData {
  title: string;
  slug?: string;
  author: string;
  content: string;
  media?: BlogPostMedia[];
  publishing: {
    status: "draft" | "published" | "archived";
    publishedAt?: string;
    scheduledAt?: string;
  };
  attractiveness: {
    isFeatured: boolean;
    isSpotlight: boolean;
    priority: number;
    highlight: boolean;
  };
}

export interface UpdateBlogPostData {
  title?: string;
  slug?: string;
  author?: string;
  content?: string;
  media?: BlogPostMedia[];
  publishing?: {
    status?: "draft" | "published" | "archived";
    publishedAt?: string;
    scheduledAt?: string;
  };
  attractiveness?: {
    isFeatured?: boolean;
    isSpotlight?: boolean;
    priority?: number;
    highlight?: boolean;
  };
}

// Response DTOs (matching backend BlogPostResponse)
export interface BlogPostResponse {
  id: string;
  title: string;
  slug: string;
  author: string;
  content: string;
  excerpt: string;
  media: BlogPostMedia[];
  publishing: Publishing;
  attractiveness: Attractiveness;
  viewCount: number;
  likeCount: number;
  shareCount: number;
  readingTime: number;
  createdAt: string;
  updatedAt: string;
}

// Brief response for lists (matching backend BlogPostBrief)
export interface BlogPostBrief {
  id: string;
  title: string;
  slug: string;
  author: string;
  excerpt: string;
  featuredImage?: string;
  readingTime: number;
  status: string;
  publishedAt?: string;
  viewCount: number;
  createdAt: string;
}

export interface BlogPostListResponse {
  posts: BlogPostBrief[];
  articles?: BlogPostBrief[]; // Backend may return data in 'articles' field
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Matching backend CreateBlogPostRequest
export interface CreateBlogArticleData {
  title: string;
  slug?: string;
  author?: string;
  leadParagraph?: string;
  content?: string;
  authorId: string;
  media?: BlogPostMedia[];
  publishing?: {
    status: "draft" | "published" | "archived";
  };
  attractiveness?: {
    isFeatured: boolean;
    isSpotlight: boolean;
    priority: number;
    highlight: boolean;
  };
  tags?: string[];
  featuredImage?: string;
  categoryId?: string;
  status: string;
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

export const contentService = {
  // ==================== BLOG POST METHODS ====================
  // Endpoint: /articles/blog

  // GET /articles/blog - Get all blog posts
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
      const response: any = await core(url, {
        method: "GET",
      });

      // Handle different response structures
      let posts: BlogPostBrief[] = [];

      // Case 1: Backend response { success, message, data: { data: [...], pagination: {...} } }
      // The array is at response.data.data
      if (response?.data?.data && Array.isArray(response.data.data)) {
        posts = response.data.data.map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          excerpt: p.leadParagraph || "",
          featuredImage: p.featuredImage || "",
          readingTime: p.readingTime || 5,
          status: p.status || "draft",
          publishedAt: p.publishedAt,
          viewCount: p.viewCount || 0,
          createdAt: p.createdAt || new Date().toISOString(),
        }));
        return {
          posts,
          total: response.data.pagination?.total || posts.length,
          page: response.data.pagination?.page || 1,
          limit: response.data.pagination?.limit || 10,
          totalPages: response.data.pagination?.totalPages || 1,
        };
      }

      // Case 2: Paginated response with data array at response.data
      if (response?.data && Array.isArray(response.data) && !Array.isArray(response.data.data)) {
        posts = response.data.map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          excerpt: p.leadParagraph || p.excerpt || "",
          featuredImage: p.featuredImage || "",
          readingTime: p.readingTime || 5,
          status: p.status || p.publishing?.status || "draft",
          publishedAt: p.publishedAt || p.publishing?.publishedAt,
          viewCount: p.viewCount || 0,
          createdAt: p.createdAt || new Date().toISOString(),
        }));
        return {
          posts,
          total: response.pagination?.total || posts.length,
          page: response.pagination?.page || 1,
          limit: response.pagination?.limit || 10,
          totalPages: response.pagination?.totalPages || 1,
        };
      }

      // Case 3: Direct array response
      if (Array.isArray(response)) {
        posts = response.map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          excerpt: p.leadParagraph || p.excerpt || "",
          featuredImage: p.featuredImage || "",
          readingTime: p.readingTime || 5,
          status: p.status || "draft",
          publishedAt: p.publishedAt,
          viewCount: p.viewCount || 0,
          createdAt: p.createdAt || new Date().toISOString(),
        }));
        return {
          posts,
          total: posts.length,
          page: 1,
          limit: posts.length,
          totalPages: 1,
        };
      }

      // Case 4: Response with articles array (alternative backend format)
      if (response?.articles && Array.isArray(response.articles)) {
        posts = response.articles.map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          excerpt: p.leadParagraph || p.excerpt || "",
          featuredImage: p.featuredImage || "",
          readingTime: p.readingTime || 5,
          status: p.status || "draft",
          publishedAt: p.publishedAt,
          viewCount: p.viewCount || 0,
          createdAt: p.createdAt || new Date().toISOString(),
        }));
        return {
          posts,
          total: response.total || posts.length,
          page: response.page || 1,
          limit: response.limit || 10,
          totalPages: response.totalPages || 1,
        };
      }

      console.warn("Unexpected response structure:", response);
      return {
        posts: [],
        total: 0,
        page: 1,
        limit: 10,
        totalPages: 0,
      };
    } catch (error) {
      console.error("Error fetching blog posts:", error);
      return {
        posts: [],
        total: 0,
        page: 1,
        limit: 10,
        totalPages: 0,
      };
    }
  },

  // GET /articles/blog/:id - Get blog post by ID
  async fetchBlogPostById(id: string): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse> | BlogPostResponse | any>(
        `/articles/blog/${id}`,
        {
          method: "GET",
        },
      );
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // POST /articles/blog - Create blog post
  async createBlogPost(
    data: CreateBlogArticleData,
  ): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse> | BlogPostResponse | any>(
        "/articles/blog",
        {
          method: "POST",
          body: data,
        },
      );
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // PUT /articles/blog/:id - Update blog post
  async updateBlogPost(
    id: string,
    data: UpdateBlogPostData,
  ): Promise<BlogPostResponse | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse> | BlogPostResponse | any>(
        `/articles/blog/${id}`,
        {
          method: "PUT",
          body: data,
        },
      );
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // DELETE /articles/blog/:id - Delete blog post
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
  // Endpoint: /articles/faq

  // GET /articles/faq - Get all FAQs
  async fetchFaqs(): Promise<Faq[]> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq[]> | Faq[]>("/articles/faq", {
        method: "GET",
      });
      // Handle both wrapped and unwrapped responses
      if (Array.isArray(response)) {
        return response;
      }
      if (response && Array.isArray((response as any).data)) {
        return (response as ApiResponse<Faq[]>).data;
      }
      return [];
    } catch {
      return [];
    }
  },

  // GET /articles/faq/:id - Get FAQ by ID
  async fetchFaqById(id: string): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq> | Faq | any>(`/articles/faq/${id}`, {
        method: "GET",
      });
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // POST /articles/faq - Create FAQ
  async createFaq(data: CreateFaqData): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq> | Faq | any>("/articles/faq", {
        method: "POST",
        body: data,
      });
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // PUT /articles/faq/:id - Update FAQ
  async updateFaq(id: string, data: UpdateFaqData): Promise<Faq | null> {
    const { core } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq> | Faq | any>(`/articles/faq/${id}`, {
        method: "PUT",
        body: data,
      });
      if (response && typeof response === "object") {
        return (response as any).data || response;
      }
      return null;
    } catch {
      return null;
    }
  },

  // DELETE /articles/faq/:id - Delete FAQ
  async deleteFaq(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/faq/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // POST /articles/faq/:id/reorder - Reorder FAQ
  async reorderFaq(id: string, orderNum: number): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/faq/${id}/reorder`, {
        method: "POST",
        body: {
          newOrder: orderNum,
        },
      });
      return true;
    } catch {
      return false;
    }
  },

  // ==================== BLOG POST VIEW COUNT ====================

  // POST /articles/:id/view - Increment view count for blog post
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
