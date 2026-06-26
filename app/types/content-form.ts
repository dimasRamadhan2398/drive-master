import type { BlogPostMedia } from "~/services/contentService";

export interface PostFormData {
  id: number;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  content: string;
  status: "draft" | "published" | "archived";
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
}
