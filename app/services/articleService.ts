import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

export interface ArticleAuthor {
  id: string;
  name: string;
  avatar?: string;
}

export interface Article {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  featuredImage?: string;
  author: ArticleAuthor | "Admin Drive Master";
  tags: string[];
  status: "draft" | "published" | "archived";
  viewCount: number;
  createdAt: string;
  updatedAt: string;
  publishedAt?: string;
}

export interface CreateArticleData {
  title: string;
  slug?: string;
  excerpt?: string;
  content: string;
  featuredImage?: string;
  tags?: string[];
}

export interface UpdateArticleData {
  title?: string;
  slug?: string;
  excerpt?: string;
  content?: string;
  featuredImage?: string;
  tags?: string[];
}

export interface PaginatedArticlesResult {
  articles: Article[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export const articleService = {
  // GET /articles - Get all articles (admin)
  async fetchAll(
    params: {
      page?: number;
      limit?: number;
      search?: string;
      status?: string;
    } = {},
  ): Promise<PaginatedArticlesResult> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.search) queryParams.set("search", params.search);
    if (params.status) queryParams.set("status", params.status);

    const queryString = queryParams.toString();
    const url = `/articles${queryString ? `?${queryString}` : ""}`;

    const response = await core<PaginatedResponse<Article>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      articles: Array.isArray(data) ? data : [],
      pagination,
    };
  },

  // GET /articles/:id - Get article by ID (admin)
  async fetchById(id: string): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>(`/articles/${id}`, {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // GET /articles/slug/:slug - Get article by slug (public)
  async fetchBySlug(slug: string): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>(
        `/articles/slug/${slug}`,
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /articles - Create article (admin)
  async create(data: CreateArticleData): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>("/articles", {
        method: "POST",
        body: data,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // PUT /articles/:id - Update article (admin)
  async update(id: string, data: UpdateArticleData): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>(`/articles/${id}`, {
        method: "PUT",
        body: data,
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // DELETE /articles/:id - Delete article (admin)
  async delete(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/${id}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // GET /articles/search - Search articles (public)
  async search(
    params: { query: string; page?: number; limit?: number } = {
      query: "",
    },
  ): Promise<PaginatedArticlesResult> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    queryParams.set("query", params.query);
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));

    const response = await core<PaginatedResponse<Article>>(
      `/articles/search?${queryParams.toString()}`,
      { method: "GET" },
    );

    const { data, pagination } = extractPaginatedData(response);
    return {
      articles: Array.isArray(data) ? data : [],
      pagination,
    };
  },

  // GET /articles/featured - Get featured articles (public)
  async getFeatured(limit = 5): Promise<Article[]> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article[]>>(
        `/articles/featured?limit=${limit}`,
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  // GET /articles/spotlight - Get spotlight article (public)
  async getSpotlight(): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>("/articles/spotlight", {
        method: "GET",
      });
      return extractData(response);
    } catch {
      return null;
    }
  },

  // GET /articles/tag/:tag - Get articles by tag (public)
  async getByTag(
    tag: string,
    params: { page?: number; limit?: number } = {},
  ): Promise<PaginatedArticlesResult> {
    const { core, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));

    const queryString = queryParams.toString();
    const url = `/articles/tag/${encodeURIComponent(tag)}${queryString ? `?${queryString}` : ""}`;

    const response = await core<PaginatedResponse<Article>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      articles: Array.isArray(data) ? data : [],
      pagination,
    };
  },

  // GET /articles/:id/related - Get related articles (public)
  async getRelated(id: string, limit = 5): Promise<Article[]> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article[]>>(
        `/articles/${id}/related?limit=${limit}`,
        { method: "GET" },
      );
      return extractData(response);
    } catch {
      return [];
    }
  },

  // POST /articles/:id/view - Increment view count (public)
  async incrementViewCount(id: string): Promise<boolean> {
    const { core } = useApiClients();
    try {
      await core(`/articles/${id}/view`, { method: "POST" });
      return true;
    } catch {
      return false;
    }
  },

  // POST /articles/:id/publish - Publish article (admin)
  async publish(id: string): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>(
        `/articles/${id}/publish`,
        { method: "POST" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /articles/:id/archive - Archive article (admin)
  async archive(id: string): Promise<Article | null> {
    const { core, extractData } = useApiClients();
    try {
      const response = await core<ApiResponse<Article>>(
        `/articles/${id}/archive`,
        { method: "POST" },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },
};
