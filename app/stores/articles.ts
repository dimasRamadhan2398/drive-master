import { defineStore } from "pinia";
import { articleService } from "~/services/articleService";
import type {
  Article,
  CreateArticleData,
  UpdateArticleData,
} from "~/services/articleService";

interface ArticlesState {
  articles: Article[];
  currentArticle: Article | null;
  featuredArticles: Article[];
  spotlightArticle: Article | null;
  searchResults: Article[];
  relatedArticles: Article[];
  isLoading: boolean;
  error: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
  searchQuery: string;
  selectedTags: string[];
  statusFilter: "all" | "draft" | "published" | "archived";
}

const initialArticles: Article[] = [
  {
    id: "1",
    title: "Getting Your Driver's License in Indonesia",
    slug: "getting-drivers-license-indonesia",
    excerpt:
      "A comprehensive guide to obtaining your driver's license in Indonesia, including requirements and examination process.",
    content: "",
    featuredImage: "/images/articles/drivers-license.jpg",
    author: {
      id: "auth-001",
      name: "Admin",
      avatar: "/images/avatars/admin.jpg",
    },
    tags: ["driving", "license", "guide"],
    status: "published",
    viewCount: 1250,
    createdAt: "Jan 15, 2026",
    updatedAt: "Feb 20, 2026",
    publishedAt: "Jan 15, 2026",
  },
  {
    id: "2",
    title: "Defensive Driving Techniques",
    slug: "defensive-driving-techniques",
    excerpt:
      "Learn essential defensive driving techniques that can help you avoid accidents and stay safe on the road.",
    content: "",
    featuredImage: "/images/articles/defensive-driving.jpg",
    author: {
      id: "auth-001",
      name: "Admin",
      avatar: "/images/avatars/admin.jpg",
    },
    tags: ["safety", "driving", "tips"],
    status: "published",
    viewCount: 890,
    createdAt: "Feb 1, 2026",
    updatedAt: "Feb 25, 2026",
    publishedAt: "Feb 1, 2026",
  },
  {
    id: "3",
    title: "Understanding Road Signs in Indonesia",
    slug: "road-signs-indonesia",
    excerpt:
      "A detailed guide to understanding and interpreting road signs you'll encounter while driving in Indonesia.",
    content: "",
    featuredImage: "/images/articles/road-signs.jpg",
    author: {
      id: "auth-001",
      name: "Admin",
      avatar: "/images/avatars/admin.jpg",
    },
    tags: ["traffic", "education", "road"],
    status: "draft",
    viewCount: 0,
    createdAt: "Mar 10, 2026",
    updatedAt: "Mar 10, 2026",
  },
];

export const useArticlesStore = defineStore("articles", {
  state: (): ArticlesState => ({
    articles: [],
    currentArticle: null,
    featuredArticles: [],
    spotlightArticle: null,
    searchResults: [],
    relatedArticles: [],
    isLoading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
    searchQuery: "",
    selectedTags: [],
    statusFilter: "all",
  }),

  getters: {
    publishedArticles: (state) =>
      state.articles.filter((a) => a.status === "published"),
    draftArticles: (state) =>
      state.articles.filter((a) => a.status === "draft"),
    archivedArticles: (state) =>
      state.articles.filter((a) => a.status === "archived"),
    totalArticles: (state) => state.articles.length,
    getArticleById: (state) => (id: string) =>
      state.articles.find((a) => a.id === id),
    getArticleBySlug: (state) => (slug: string) =>
      state.articles.find((a) => a.slug === slug),
    allTags: (state) => {
      const tags = new Set<string>();
      state.articles.forEach((article) => {
        article.tags.forEach((tag) => tags.add(tag));
      });
      return Array.from(tags).sort();
    },
    topViewedArticles: (state) =>
      [...state.articles].sort((a, b) => b.viewCount - a.viewCount).slice(0, 10),
    recentArticles: (state) =>
      [...state.articles]
        .sort(
          (a, b) =>
            new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
        )
        .slice(0, 10),
  },

  actions: {
    // Admin CRUD actions
    async fetchArticles(page = 1, resetPage = true) {
      this.isLoading = true;
      this.error = null;

      try {
        if (resetPage) {
          this.pagination.page = page;
        }

        const params: {
          page: number;
          limit: number;
          search?: string;
          status?: string;
        } = {
          page: this.pagination.page,
          limit: this.pagination.limit,
        };

        if (this.searchQuery) {
          params.search = this.searchQuery;
        }

        if (this.statusFilter !== "all") {
          params.status = this.statusFilter;
        }

        const result = await articleService.fetchAll(params);

        this.articles = result.articles;
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch articles";
        console.error("Error fetching articles:", err);
        this.articles = [...initialArticles];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchArticleById(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.fetchById(id);
        if (article) {
          this.currentArticle = article;
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch article";
        console.error("Error fetching article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async fetchArticleBySlug(slug: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.fetchBySlug(slug);
        if (article) {
          this.currentArticle = article;
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch article";
        console.error("Error fetching article by slug:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async createArticle(data: CreateArticleData) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.create(data);
        if (article) {
          this.articles.unshift(article);
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create article";
        console.error("Error creating article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async updateArticle(id: string, data: UpdateArticleData) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.update(id, data);
        if (article) {
          const index = this.articles.findIndex((a) => a.id === id);
          if (index !== -1) {
            this.articles[index] = article;
          }
          if (this.currentArticle?.id === id) {
            this.currentArticle = article;
          }
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to update article";
        console.error("Error updating article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async deleteArticle(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const success = await articleService.delete(id);
        if (success) {
          this.articles = this.articles.filter((a) => a.id !== id);
          if (this.currentArticle?.id === id) {
            this.currentArticle = null;
          }
        }
        return success;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to delete article";
        console.error("Error deleting article:", err);
        return false;
      } finally {
        this.isLoading = false;
      }
    },

    // Public actions
    async searchArticles(query: string, page = 1) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await articleService.search({ query, page });
        this.searchResults = result.articles;
        this.pagination = result.pagination;
        this.searchQuery = query;
        return result;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to search articles";
        console.error("Error searching articles:", err);
        return { articles: [], pagination: this.pagination };
      } finally {
        this.isLoading = false;
      }
    },

    async fetchFeaturedArticles(limit = 5) {
      this.isLoading = true;
      this.error = null;

      try {
        const articles = await articleService.getFeatured(limit);
        this.featuredArticles = articles;
        return articles;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch featured articles";
        console.error("Error fetching featured articles:", err);
        return [];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchSpotlightArticle() {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.getSpotlight();
        this.spotlightArticle = article;
        return article;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch spotlight article";
        console.error("Error fetching spotlight article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async fetchArticlesByTag(tag: string, page = 1) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await articleService.getByTag(tag, { page });
        this.searchResults = result.articles;
        this.pagination = result.pagination;
        return result;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch articles by tag";
        console.error("Error fetching articles by tag:", err);
        return { articles: [], pagination: this.pagination };
      } finally {
        this.isLoading = false;
      }
    },

    async fetchRelatedArticles(id: string, limit = 5) {
      this.isLoading = true;
      this.error = null;

      try {
        const articles = await articleService.getRelated(id, limit);
        this.relatedArticles = articles;
        return articles;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch related articles";
        console.error("Error fetching related articles:", err);
        return [];
      } finally {
        this.isLoading = false;
      }
    },

    async incrementViewCount(id: string) {
      try {
        await articleService.incrementViewCount(id);
        const article = this.articles.find((a) => a.id === id);
        if (article) {
          article.viewCount++;
        }
        if (this.currentArticle?.id === id) {
          this.currentArticle!.viewCount++;
        }
      } catch {
        // Silent fail for view count
      }
    },

    async publishArticle(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.publish(id);
        if (article) {
          const index = this.articles.findIndex((a) => a.id === id);
          if (index !== -1) {
            this.articles[index] = article;
          }
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to publish article";
        console.error("Error publishing article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async archiveArticle(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const article = await articleService.archive(id);
        if (article) {
          const index = this.articles.findIndex((a) => a.id === id);
          if (index !== -1) {
            this.articles[index] = article;
          }
          return article;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to archive article";
        console.error("Error archiving article:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // Filter and search helpers
    setSearchQuery(query: string) {
      this.searchQuery = query;
    },

    setStatusFilter(status: "all" | "draft" | "published" | "archived") {
      this.statusFilter = status;
    },

    setSelectedTags(tags: string[]) {
      this.selectedTags = tags;
    },

    setPage(page: number) {
      this.pagination.page = page;
      this.fetchArticles(page, false);
    },

    setCurrentArticle(article: Article | null) {
      this.currentArticle = article;
    },

    clearSearchResults() {
      this.searchResults = [];
      this.searchQuery = "";
    },

    reset() {
      this.articles = [...initialArticles];
      this.currentArticle = null;
      this.featuredArticles = [];
      this.spotlightArticle = null;
      this.searchResults = [];
      this.relatedArticles = [];
      this.isLoading = false;
      this.error = null;
      this.pagination = {
        page: 1,
        limit: 10,
        total: 0,
        totalPages: 0,
      };
      this.searchQuery = "";
      this.selectedTags = [];
      this.statusFilter = "all";
    },
  },
});