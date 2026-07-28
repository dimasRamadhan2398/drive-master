import { defineStore } from "pinia";
import {
  contentService,
  dataUriToFileUpload,
  type Publishing,
  type Attractiveness,
  type FileUpload,
} from "../services/contentService";

// ============================================================
// Local types
// ============================================================

export interface PageSection {
  id: string;
  type: string;
  props?: Record<string, unknown>;
  data: unknown;
}

export interface Page {
  id: string | number;
  title: string;
  slug: string;
  lastUpdated: string;
  status: "published" | "draft";
  sections: PageSection[];
}

/** Internal BlogPost representation used by the store / UI */
export interface BlogPost {
  id: string;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  featuredImage: string;
  excerpt: string;
  date: string;
  status: "published" | "draft" | "archived";
  views: number;
  content: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
  viewCount: number;
  likeCount: number;
  readingTime: number;
  createdAt: string;
  updatedAt: string;
}

// -------------------------------------------------------
// FAQ Types
// -------------------------------------------------------
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

// ============================================================
// Store state
// ============================================================

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

export const useContentStore = defineStore("content", {
  state: (): ContentState => ({
    pages: initialPages,
    blogPosts: [],
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
    draftPosts: (state) =>
      state.blogPosts.filter((p) => p.status === "draft"),
    totalPages: (state) => state.pages.length,
    totalPosts: (state) => state.blogPosts.length,
    getPageById: (state) => (id: string | number) => state.pages.find((p) => String(p.id) === String(id)),
    getPageBySlug: (state) => (slug: string) => state.pages.find((p) => p.slug === slug),
    getPostById: (state) => (id: string) => state.blogPosts.find((p) => p.id === id),
    // FAQ getters
    activeFaqs: (state) => state.faqs.filter((f) => f.isActive),
    sortedFaqs: (state) => [...state.faqs].sort((a, b) => a.sortOrder - b.sortOrder),
    totalFaqs: (state) => state.faqs.length,
    getFaqById: (state) => (id: string) => state.faqs.find((f) => f.id === id),
  },

  actions: {
    // ----------------------------------------------------------
    _mapResponseToPage(pageApi: any): Page {
      let parsedSections: PageSection[] = [];
      if (pageApi.sections) {
        try {
          if (typeof pageApi.sections === "string") {
            parsedSections = JSON.parse(pageApi.sections);
          } else if (Array.isArray(pageApi.sections)) {
            parsedSections = pageApi.sections;
          }
        } catch (e) {
          console.error("Failed to parse page sections:", e);
        }
      }
      return {
        id: pageApi.id,
        title: pageApi.title,
        slug: pageApi.slug,
        status: pageApi.status,
        lastUpdated: pageApi.lastUpdated || this.formatDate(new Date(pageApi.updatedAt || Date.now())),
        sections: parsedSections,
      };
    },

    async fetchPages() {
      this.isLoading = true;
      try {
        const rawPages = await contentService.fetchPages();
        this.pages = rawPages.map((p) => this._mapResponseToPage(p));
      } catch (e) {
        console.error("Failed to fetch pages:", e);
        this.error = "Failed to fetch pages";
      } finally {
        this.isLoading = false;
      }
    },

    async addPage(data: { title: string; slug: string; status: "published" | "draft" }) {
      this.isLoading = true;
      try {
        const created = await contentService.createPage(data);
        if (created) {
          const mapped = this._mapResponseToPage(created);
          this.pages.push(mapped);
          return mapped;
        }
      } catch (e) {
        console.error("Failed to create page:", e);
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    async updatePage(id: string | number, data: Partial<Page>) {
      this.isLoading = true;
      try {
        const payload: any = {};
        if (data.title !== undefined) payload.title = data.title;
        if (data.slug !== undefined) payload.slug = data.slug;
        if (data.status !== undefined) payload.status = data.status;
        if (data.sections !== undefined) payload.sections = JSON.stringify(data.sections);

        const updated = await contentService.updatePage(String(id), payload);
        if (updated) {
          const mapped = this._mapResponseToPage(updated);
          const index = this.pages.findIndex((p) => String(p.id) === String(id));
          if (index !== -1) {
            this.pages[index] = mapped;
          }
          return mapped;
        }
      } catch (e) {
        console.error("Failed to update page:", e);
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    async deletePage(id: string | number) {
      this.isLoading = true;
      try {
        const ok = await contentService.deletePage(String(id));
        if (ok) {
          this.pages = this.pages.filter((p) => String(p.id) !== String(id));
        }
      } catch (e) {
        console.error("Failed to delete page:", e);
      } finally {
        this.isLoading = false;
      }
    },

    async updatePageSections(slug: string, sections: PageSection[]) {
      const page = this.pages.find((p) => p.slug === slug);
      if (page) {
        await this.updatePage(page.id, { sections });
      }
    },

    async togglePageStatus(id: string | number) {
      const page = this.pages.find((p) => String(p.id) === String(id));
      if (page) {
        const nextStatus = page.status === "published" ? "draft" : "published";
        const updated = await this.updatePage(id, { status: nextStatus });
        return updated ? updated.status : null;
      }
      return null;
    },

    // ----------------------------------------------------------
    // Blog Post actions
    // ----------------------------------------------------------

    /**
     * Create a new blog post. Sends a JSON payload matching Go's
     * CreateBlogPostRequest. If featuredImage is a data URI string,
     * it is converted to a FileUpload object for upload.
     */
    async addPost(data: {
      title: string;
      slug?: string;
      author?: string;
      authorId: string;
      leadParagraph?: string;
      /** Either a data URI (newly selected file) or an existing URL */
      featuredImage?: File | string;
      /** Raw filename of the selected file (needed for format validation) */
      featuredImageFileName?: string;
      content?: string;
      status?: "draft" | "published" | "archived";
      publishing?: Partial<Publishing>;
      attractiveness?: Partial<Attractiveness>;
    }) {
      this.isLoading = true;
      try {
        // Build featuredImage FileUpload if a data URI was provided
        let featuredImagePayload: FileUpload | undefined;
        if (data.featuredImage && typeof data.featuredImage === "string") {
          const parsed = dataUriToFileUpload(
            data.featuredImage,
            data.featuredImageFileName || "image.jpg",
          );
          if (parsed) featuredImagePayload = parsed;
          // If not a data URI it's an existing URL – skip (backend ignores it)
        }

        const status = data.status || data.publishing?.status || "draft";

        const created = await contentService.createBlogPost({
          title: data.title,
          slug: data.slug,
          author: data.author,
          authorId: data.authorId,
          leadParagraph: data.leadParagraph,
          content: data.content,
          status,
          featuredImage: featuredImagePayload,
          publishing: {
            status: status as "draft" | "published" | "archived",
            publishedAt: data.publishing?.publishedAt ?? null,
            scheduledAt: data.publishing?.scheduledAt ?? null,
          },
          attractiveness: {
            isFeatured: data.attractiveness?.isFeatured ?? false,
            isSpotlight: data.attractiveness?.isSpotlight ?? false,
            priority: data.attractiveness?.priority ?? 0,
            highlight: data.attractiveness?.highlight ?? false,
          },
        });

        if (created) {
          const newPost = this._mapResponseToPost(created);
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

    /**
     * Update an existing blog post. Same JSON payload approach.
     */
    async updatePost(
      id: string,
      data: {
        title?: string;
        slug?: string;
        author?: string;
        authorId: string;
        leadParagraph?: string;
        /** Either a data URI (new upload) or an existing URL (no-op) */
        featuredImage?: File | string;
        featuredImageFileName?: string;
        content?: string;
        status?: "draft" | "published" | "archived";
        publishing?: Partial<Publishing>;
        attractiveness?: Partial<Attractiveness>;
      },
    ) {
      this.isLoading = true;
      try {
        let featuredImagePayload: FileUpload | undefined;
        if (data.featuredImage && typeof data.featuredImage === "string") {
          const parsed = dataUriToFileUpload(
            data.featuredImage,
            data.featuredImageFileName || "image.jpg",
          );
          if (parsed) featuredImagePayload = parsed;
        }

        const status = data.status || data.publishing?.status;

        const updated = await contentService.updateBlogPost(id, {
          title: data.title,
          slug: data.slug,
          author: data.author,
          authorId: data.authorId,
          leadParagraph: data.leadParagraph,
          content: data.content,
          status,
          featuredImage: featuredImagePayload,
          publishing: data.publishing
            ? {
                status: (data.publishing.status || status || "draft") as
                  | "draft"
                  | "published"
                  | "archived",
                publishedAt: data.publishing.publishedAt ?? null,
                scheduledAt: data.publishing.scheduledAt ?? null,
              }
            : undefined,
          attractiveness: data.attractiveness
            ? {
                isFeatured: data.attractiveness.isFeatured ?? false,
                isSpotlight: data.attractiveness.isSpotlight ?? false,
                priority: data.attractiveness.priority ?? 0,
                highlight: data.attractiveness.highlight ?? false,
              }
            : undefined,
        });

        if (updated) {
          const index = this.blogPosts.findIndex((p) => p.id === id);
          const updatedPost = this._mapResponseToPost(updated);
          if (index !== -1) {
            this.blogPosts[index] = updatedPost;
          }
          return updatedPost;
        }
      } catch (e) {
        console.error("Failed to update blog post:", e);
        this.error = "Failed to update blog post";
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    async deletePost(id: string) {
      this.isLoading = true;
      try {
        const success = await contentService.deleteBlogPost(id);
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
        this.blogPosts = (result?.posts ?? []).map((p: any) => ({
          id: p.id,
          title: p.title,
          slug: p.slug || "",
          author: p.author || "Admin",
          leadParagraph: p.leadParagraph || p.excerpt || "",
          featuredImage: p.featuredImage || "",
          excerpt: p.excerpt || p.leadParagraph || "",
          date: this.formatDate(new Date(p.createdAt)),
          status: (p.publishing?.status || p.status || "draft") as
            | "published"
            | "draft"
            | "archived",
          views: p.viewCount || 0,
          content: p.content || "",
          publishing: {
            status: (p.publishing?.status || p.status || "draft") as
              | "draft"
              | "published"
              | "archived",
            publishedAt: p.publishing?.publishedAt ?? null,
            scheduledAt: p.publishing?.scheduledAt ?? null,
          },
          attractiveness: {
            isFeatured: p.attractiveness?.isFeatured ?? false,
            isSpotlight: p.attractiveness?.isSpotlight ?? false,
            priority: p.attractiveness?.priority ?? 0,
            highlight: p.attractiveness?.highlight ?? false,
          },
          viewCount: p.viewCount || 0,
          likeCount: p.likeCount || 0,
          readingTime: p.readingTime || 1,
          createdAt: p.createdAt || new Date().toISOString(),
          updatedAt: p.updatedAt || new Date().toISOString(),
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

    async incrementBlogPostViewCount(id: string) {
      try {
        await contentService.incrementBlogPostViewCount(id);
        const post = this.blogPosts.find((p) => p.id === id);
        if (post) {
          post.viewCount++;
          post.views++;
        }
      } catch {
        // Silent fail – don't disrupt UX
      }
    },

    // ----------------------------------------------------------
    // Helper: map a BlogPostResponse to the internal BlogPost type
    // ----------------------------------------------------------
    _mapResponseToPost(r: any): BlogPost {
      return {
        id: r.id,
        title: r.title,
        slug: r.slug || "",
        author: r.author || "Admin",
        leadParagraph: r.leadParagraph || r.excerpt || "",
        featuredImage: r.featuredImage || "",
        excerpt: r.excerpt || r.leadParagraph || "",
        date: this.formatDate(new Date(r.createdAt || Date.now())),
        status: (r.publishing?.status || r.status || "draft") as
          | "published"
          | "draft"
          | "archived",
        views: r.viewCount || 0,
        content: r.content || "",
        publishing: {
          status: (r.publishing?.status || "draft") as
            | "draft"
            | "published"
            | "archived",
          publishedAt: r.publishing?.publishedAt ?? null,
          scheduledAt: r.publishing?.scheduledAt ?? null,
        },
        attractiveness: {
          isFeatured: r.attractiveness?.isFeatured ?? false,
          isSpotlight: r.attractiveness?.isSpotlight ?? false,
          priority: r.attractiveness?.priority ?? 0,
          highlight: r.attractiveness?.highlight ?? false,
        },
        viewCount: r.viewCount || 0,
        likeCount: r.likeCount || 0,
        readingTime: r.readingTime || 1,
        createdAt: r.createdAt || new Date().toISOString(),
        updatedAt: r.updatedAt || new Date().toISOString(),
      };
    },

    // ----------------------------------------------------------
    // Helpers
    // ----------------------------------------------------------
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

    // ----------------------------------------------------------
    // FAQ Actions
    // ----------------------------------------------------------
    async fetchFaqs() {
      this.isLoading = true;
      this.error = null;
      try {
        const faqs = await contentService.fetchFaqs();
        this.faqs = faqs;
        return faqs;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to fetch FAQs";
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
        this.error = err instanceof Error ? err.message : "Failed to create FAQ";
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
          if (index !== -1) this.faqs[index] = faq;
          if (this.currentFaq?.id === id) this.currentFaq = faq;
          return faq;
        }
        return null;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to update FAQ";
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
          if (this.currentFaq?.id === id) this.currentFaq = null;
        }
        return success;
      } catch (err) {
        this.error = err instanceof Error ? err.message : "Failed to delete FAQ";
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
        this.faqs.forEach((f, i) => (f.sortOrder = i + 1));
        try {
          await Promise.all(
            this.faqs.map((f, i) => contentService.reorderFaq(f.id, i + 1)),
          );
        } catch (error) {
          console.error("Failed to sync FAQ order with server:", error);
        }
      }
    },

    reset() {
      this.pages = [...initialPages];
      this.blogPosts = [];
      this.faqs = [];
      this.currentFaq = null;
      this.isLoading = false;
      this.error = null;
    },
  },
});
