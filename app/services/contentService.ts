import type {
  ApiResponse,
  PaginatedResponse,
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

export const contentService = {
  // ==================== BLOG POST METHODS ====================
  // Endpoint: /articles/blog

  // GET /articles/blog - Get all blog posts
  async fetchBlogPosts(
    params: BlogPostFilterParams = {},
  ): Promise<BlogPostListResponse> {
    const { core, extractPaginatedData } = useApiClients();

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

    const response = await core<PaginatedResponse<BlogPostResponse>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      posts: Array.isArray(data)
        ? data.map<BlogPostBrief>((p: BlogPostResponse) => ({
            id: p.id,
            title: p.title,
            slug: p.slug,
            author: p.author,
            excerpt: p.excerpt,
            featuredImage: p.media.find((m) => m.fileType === "image")?.url,
            readingTime: p.readingTime,
            status: p.publishing.status,
            publishedAt: p.publishing.publishedAt,
            viewCount: p.viewCount,
            createdAt: p.createdAt,
          }))
        : [],
      ...pagination,
    };
  },

  // GET /articles/blog/:id - Get blog post by ID
  async fetchBlogPostById(id: string): Promise<BlogPostResponse | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse>>(
        `/articles/blog/${id}`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /articles/blog - Create blog post
  async createBlogPost(
    data: CreateBlogPostData,
  ): Promise<BlogPostResponse | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse>>(
        "/articles/blog",
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

  // PUT /articles/blog/:id - Update blog post
  async updateBlogPost(
    id: string,
    data: UpdateBlogPostData,
  ): Promise<BlogPostResponse | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<BlogPostResponse>>(
        `/articles/blog/${id}`,
        {
          method: "PUT",
          body: data,
        },
      );
      return extractData(response);
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
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq[]>>("/articles/faq", {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /articles/faq/:id - Get FAQ by ID
  async fetchFaqById(id: string): Promise<Faq | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq>>(`/articles/faq/${id}`, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /articles/faq - Create FAQ
  async createFaq(data: CreateFaqData): Promise<Faq | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq>>("/articles/faq", {
        method: "POST",
        body: data,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PUT /articles/faq/:id - Update FAQ
  async updateFaq(id: string, data: UpdateFaqData): Promise<Faq | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Faq>>(`/articles/faq/${id}`, {
        method: "PUT",
        body: data,
      });
      return extractData(response);
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
};
