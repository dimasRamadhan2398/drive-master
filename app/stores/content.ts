import { defineStore } from "pinia";
import { contentService } from "~/services/contentService";
import type {
  BlogPostMedia,
  CreateBlogPostData,
  UpdateBlogPostData,
  CreateBlogArticleData,
} from "~/services/contentService";

// Re-export types from contentService for convenience
export type {
  BlogPostMedia,
  CreateBlogPostData,
  UpdateBlogPostData,
  CreateBlogArticleData,
};

export interface PageSection {
  id: string;
  type: string;
  props?: Record<string, unknown>;
  data: unknown;
}

export interface Page {
  id: number;
  title: string;
  slug: string;
  lastUpdated: string;
  status: "published" | "draft";
  sections: PageSection[];
}

export interface BlogPost {
  id: number;
  title: string;
  slug?: string;
  author: string;
  excerpt?: string;
  date: string;
  status: "published" | "draft" | "archived";
  views: number;
  content: string;
  media: BlogPostMedia[];
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
  viewCount?: number;
  likeCount?: number;
  shareCount?: number;
  readingTime?: number;
  createdAt?: string;
  updatedAt?: string;
}

// FAQ Types
export interface Faq {
  id: string;
  question: string;
  answer: string;
  sortOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateFaqData {
  question: string;
  answer: string;
  sortOrder?: number;
  isActive?: boolean;
}

export interface UpdateFaqData {
  question?: string;
  answer?: string;
  sortOrder?: number;
  isActive?: boolean;
}

interface ContentState {
  pages: Page[];
  blogPosts: BlogPost[];
  faqs: Faq[];
  currentFaq: Faq | null;
  isLoading: boolean;
  error: string | null;
}

const initialPages: Page[] = [
  {
    id: 1,
    title: "Home",
    slug: "/",
    lastUpdated: "Apr 10, 2026",
    status: "published",
    sections: [],
  },
  {
    id: 2,
    title: "About Us",
    slug: "/about",
    lastUpdated: "Apr 9, 2026",
    status: "draft",
    sections: [],
  },
];

const initialBlogPosts: BlogPost[] = [
  {
    id: 1,
    title: "Welcome to Our New Blog!",
    author: "Admin",
    date: "Apr 10, 2026",
    status: "published",
    views: 120,
    content:
      "This is the first post on our brand new blog. Stay tuned for more updates!",
    media: [],
    publishing: {
      status: "published",
    },
    attractiveness: {
      isFeatured: false,
      isSpotlight: false,
      priority: 0,
      highlight: false,
    },
  },
];

export const useContentStore = defineStore("content", {
  state: (): ContentState => ({
    pages: initialPages,
    blogPosts: initialBlogPosts,
    faqs: [],
    currentFaq: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    publishedPages: (state) =>
      state.pages.filter((p) => p.status === "published"),
    draftPages: (state) => state.pages.filter((p) => p.status === "draft"),
    publishedPosts: (state) =>
      state.blogPosts.filter((p) => p.status === "published"),
    draftPosts: (state) => state.blogPosts.filter((p) => p.status === "draft"),
    totalPages: (state) => state.pages.length,
    totalPosts: (state) => state.blogPosts.length,
    getPageById: (state) => (id: number) =>
      state.pages.find((p) => p.id === id),
    getPageBySlug: (state) => (slug: string) =>
      state.pages.find((p) => p.slug === slug),
    getPostById: (state) => (id: number) =>
      state.blogPosts.find((p) => p.id === id),
    // FAQ getters
    activeFaqs: (state) => state.faqs.filter((f) => f.isActive),
    sortedFaqs: (state) =>
      [...state.faqs].sort((a, b) => a.sortOrder - b.sortOrder),
    totalFaqs: (state) => state.faqs.length,
    getFaqById: (state) => (id: string) => state.faqs.find((f) => f.id === id),
  },

  actions: {
    // Page actions
    addPage(data: {
      title: string;
      slug: string;
      status: "published" | "draft";
    }) {
      const newId = Math.max(...this.pages.map((p) => p.id), 0) + 1;
      const newPage: Page = {
        ...data,
        id: newId,
        lastUpdated: this.formatDate(new Date()),
        sections: [],
      };
      this.pages.push(newPage);
      return newPage;
    },

    updatePage(id: number, data: Partial<Page>) {
      const index = this.pages.findIndex((p) => p.id === id);
      if (index !== -1) {
        const existing = this.pages[index];
        if (!existing) return null;
        this.pages[index] = {
          id: existing.id,
          title: data.title ?? existing.title,
          slug: data.slug ?? existing.slug,
          lastUpdated: this.formatDate(new Date()),
          status: data.status ?? existing.status,
          sections: data.sections ?? existing.sections,
        };
        return this.pages[index];
      }
      return null;
    },

    deletePage(id: number) {
      this.pages = this.pages.filter((p) => p.id !== id);
    },

    togglePageStatus(id: number) {
      const page = this.pages.find((p) => p.id === id);
      if (page) {
        page.status = page.status === "published" ? "draft" : "published";
        page.lastUpdated = this.formatDate(new Date());
        return page.status;
      }
      return null;
    },

    // Blog Post actions
    async addPost(data: CreateBlogArticleData) {
      this.isLoading = true;
      try {
        const created = await contentService.createBlogPost(data);

        if (created) {
          const newPost: BlogPost = {
            id: created.id as any,
            title: created.title,
            slug: created.slug,
            author: created.author,
            content: created.content,
            excerpt: created.excerpt,
            date: this.formatDate(new Date(created.createdAt)),
            status: created.publishing?.status || "draft",
            views: created.viewCount || 0,
            media: created.media || [],
            publishing: created.publishing,
            attractiveness: created.attractiveness,
          };
          this.blogPosts.unshift(newPost);
          return newPost;
        }
      } catch (e) {
        console.error("Failed to create blog post:", e);
        this.error = "Failed to create blog post";
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    async updatePost(id: string | number, data: UpdateBlogPostData) {
      this.isLoading = true;
      try {
        const updated = await contentService.updateBlogPost(String(id), data);

        if (updated) {
          const index = this.blogPosts.findIndex((p) => p.id === id);
          if (index !== -1) {
            const existing = this.blogPosts[index];
            if (!existing) return null;
            this.blogPosts[index] = {
              ...existing,
              title: updated.title,
              slug: updated.slug,
              author: updated.author,
              content: updated.content,
              excerpt: updated.excerpt,
              date: this.formatDate(new Date(updated.updatedAt)),
              status: (updated.publishing?.status as any) || existing.status,
              views: updated.viewCount || existing.views,
              media: updated.media || existing.media,
              publishing: (updated.publishing as any) || existing.publishing,
              attractiveness: updated.attractiveness || existing.attractiveness,
            };
            return this.blogPosts[index];
          }
        }
      } catch (e) {
        console.error("Failed to update blog post:", e);
        this.error = "Failed to update blog post";
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    async deletePost(id: string | number) {
      this.isLoading = true;
      try {
        const success = await contentService.deleteBlogPost(String(id));
        if (success) {
          this.blogPosts = this.blogPosts.filter((p) => p.id !== id);
          return true;
        }
      } catch (e) {
        console.error("Failed to delete blog post from server:", e);
        this.error = "Failed to delete blog post";
      } finally {
        this.isLoading = false;
      }
      return false;
    },

    togglePostStatus(id: number) {
      const post = this.blogPosts.find((p) => p.id === id);
      if (post) {
        post.status = post.status === "published" ? "draft" : "published";
        post.date = this.formatDate(new Date());
        return post.status;
      }
      return null;
    },

    incrementPostViews(id: number) {
      const post = this.blogPosts.find((p) => p.id === id);
      if (post) {
        post.views++;
      }
    },

    async fetchBlogPosts(
      params: {
        page?: number;
        limit?: number;
        search?: string;
        status?: string;
      } = {},
    ) {
      this.isLoading = true;
      this.error = null;
      try {
        const result = await contentService.fetchBlogPosts(params);
        // Use the posts array from the result
        const posts = result?.posts || [];
        this.blogPosts = posts.map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          content: "",
          excerpt: p.excerpt || p.leadParagraph || "",
          date: this.formatDate(new Date(p.createdAt)),
          status: (p.status as "published" | "draft" | "archived") || "draft",
          views: p.viewCount || 0,
          media: [],
          publishing: {
            status: (p.status as "published" | "draft" | "archived") || "draft",
          },
          attractiveness: {
            isFeatured: false,
            isSpotlight: false,
            priority: 0,
            highlight: false,
          },
        }));
        return result;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch blog posts";
        console.error("Error fetching blog posts:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    // Helper functions
    formatDate(date: Date): string {
      return date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
      });
    },

    generateSlug(title: string): string {
      return (
        "/" +
        title
          .toLowerCase()
          .replace(/[^a-z0-9]+/g, "-")
          .replace(/^-|-$/g, "")
      );
    },

    // Media helpers
    addMediaToPost(postId: number, media: BlogPostMedia) {
      const post = this.blogPosts.find((p) => p.id === postId);
      if (post) {
        post.media.push(media);
        return true;
      }
      return false;
    },

    removeMediaFromPost(postId: number, mediaIndex: number) {
      const post = this.blogPosts.find((p) => p.id === postId);
      if (post && mediaIndex >= 0 && mediaIndex < post.media.length) {
        post.media.splice(mediaIndex, 1);
        return true;
      }
      return false;
    },

    // FAQ Actions
    async fetchFaqs() {
      this.isLoading = true;
      this.error = null;

      try {
        const faqs = await contentService.fetchFaqs();
        this.faqs = faqs;
        return faqs;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch FAQs";
        console.error("Error fetching FAQs:", err);
        return [];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchFaqById(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const faq = await contentService.fetchFaqById(id);
        if (faq) {
          this.currentFaq = faq;
          return faq;
        }
        return null;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to fetch FAQ";
        console.error("Error fetching FAQ:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async createFaq(data: CreateFaqData) {
      this.isLoading = true;
      this.error = null;

      try {
        const faq = await contentService.createFaq(data);
        if (faq) {
          this.faqs.unshift(faq);
          return faq;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create FAQ";
        console.error("Error creating FAQ:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async updateFaq(id: string, data: UpdateFaqData) {
      this.isLoading = true;
      this.error = null;

      try {
        const faq = await contentService.updateFaq(id, data);
        if (faq) {
          const index = this.faqs.findIndex((f) => f.id === id);
          if (index !== -1) {
            this.faqs[index] = faq;
          }
          if (this.currentFaq?.id === id) {
            this.currentFaq = faq;
          }
          return faq;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to update FAQ";
        console.error("Error updating FAQ:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async deleteFaq(id: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const success = await contentService.deleteFaq(id);
        if (success) {
          this.faqs = this.faqs.filter((f) => f.id !== id);
          if (this.currentFaq?.id === id) {
            this.currentFaq = null;
          }
        }
        return success;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to delete FAQ";
        console.error("Error deleting FAQ:", err);
        return false;
      } finally {
        this.isLoading = false;
      }
    },

    setCurrentFaq(faq: Faq | null) {
      this.currentFaq = faq;
    },

    async reorderFaqs(fromIndex: number, toIndex: number) {
      const faq = this.faqs.splice(fromIndex, 1)[0];
      if (faq) {
        this.faqs.splice(toIndex, 0, faq);
        // Update sort order locally first
        this.faqs.forEach((f, i) => (f.sortOrder = i + 1));

        // Then sync with API
        try {
          await Promise.all(
            this.faqs.map((f, i) => contentService.reorderFaq(f.id, i + 1)),
          );
        } catch (error) {
          console.error("Failed to sync FAQ order with server:", error);
        }
      }
    },

    // Reset state
    reset() {
      this.pages = [...initialPages];
      this.blogPosts = [...initialBlogPosts];
      this.faqs = [];
      this.currentFaq = null;
      this.isLoading = false;
      this.error = null;
    },
  },
});
