import { defineStore } from "pinia";

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

export interface BlogPostMedia {
  name: string;
  type: string;
  size: string;
  url: string;
  fileType: "image" | "video";
}

export interface BlogPost {
  id: number;
  title: string;
  author: string;
  date: string;
  status: "published" | "draft";
  views: number;
  content: string;
  media: BlogPostMedia[];
}

export interface Faq {
  id: number;
  question: string;
  answer: string;
  order: number;
}

interface ContentState {
  pages: Page[];
  blogPosts: BlogPost[];
  faqs: Faq[];
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
  },
];

const initialFaqs: Faq[] = [
  {
    id: 1,
    question: "What is an EV?",
    answer:
      "An Electric Vehicle (EV) is a vehicle that uses one or more electric motors for propulsion.",
    order: 1,
  },
  {
    id: 2,
    question: "How long does it take to charge?",
    answer:
      "Charging time depends on the charger and the vehicle. A typical home charger can fully charge an EV overnight.",
    order: 2,
  },
];

export const useContentStore = defineStore("content", {
  state: (): ContentState => ({
    pages: initialPages,
    blogPosts: initialBlogPosts,
    faqs: initialFaqs,
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
    sortedFaqs: (state) => [...state.faqs].sort((a, b) => a.order - b.order),
    totalPages: (state) => state.pages.length,
    totalPosts: (state) => state.blogPosts.length,
    totalFaqs: (state) => state.faqs.length,
    getPageById: (state) => (id: number) =>
      state.pages.find((p) => p.id === id),
    getPageBySlug: (state) => (slug: string) =>
      state.pages.find((p) => p.slug === slug),
    getPostById: (state) => (id: number) =>
      state.blogPosts.find((p) => p.id === id),
    getFaqById: (state) => (id: number) => state.faqs.find((f) => f.id === id),
  },

  actions: {
    // Page actions
    addPage(data: { title: string; slug: string; status: "published" | "draft" }) {
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
        this.pages[index] = {
          ...this.pages[index],
          ...data,
          lastUpdated: this.formatDate(new Date()),
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
    addPost(data: {
      title: string;
      author: string;
      content: string;
      status: "published" | "draft";
      media?: BlogPostMedia[];
    }) {
      const newId = Math.max(...this.blogPosts.map((p) => p.id), 0) + 1;
      const newPost: BlogPost = {
        ...data,
        id: newId,
        date: this.formatDate(new Date()),
        views: 0,
        media: data.media || [],
      };
      this.blogPosts.push(newPost);
      return newPost;
    },

    updatePost(id: number, data: Partial<BlogPost>) {
      const index = this.blogPosts.findIndex((p) => p.id === id);
      if (index !== -1) {
        this.blogPosts[index] = {
          ...this.blogPosts[index],
          ...data,
          date: this.formatDate(new Date()),
        };
        return this.blogPosts[index];
      }
      return null;
    },

    deletePost(id: number) {
      this.blogPosts = this.blogPosts.filter((p) => p.id !== id);
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

    // FAQ actions
    addFaq(data: { question: string; answer: string }) {
      const newId = Math.max(...this.faqs.map((f) => f.id), 0) + 1;
      const newFaq: Faq = {
        ...data,
        id: newId,
        order: this.faqs.length + 1,
      };
      this.faqs.push(newFaq);
      return newFaq;
    },

    updateFaq(id: number, data: Partial<Faq>) {
      const index = this.faqs.findIndex((f) => f.id === id);
      if (index !== -1) {
        this.faqs[index] = { ...this.faqs[index], ...data };
        return this.faqs[index];
      }
      return null;
    },

    deleteFaq(id: number) {
      this.faqs = this.faqs.filter((f) => f.id !== id);
      // Reorder remaining FAQs
      this.faqs.forEach((f, i) => (f.order = i + 1));
    },

    reorderFaqs(fromIndex: number, toIndex: number) {
      const faq = this.faqs.splice(fromIndex, 1)[0];
      this.faqs.splice(toIndex, 0, faq);
      // Update order numbers
      this.faqs.forEach((f, i) => (f.order = i + 1));
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
        "/" + title.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "")
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

    // Reset state
    reset() {
      this.pages = [...initialPages];
      this.blogPosts = [...initialBlogPosts];
      this.faqs = [...initialFaqs];
      this.isLoading = false;
      this.error = null;
    },
  },
});