import type { Publishing, Attractiveness } from "~/services/contentService";

/** PostFormData is the shape passed into PostBlogModal for edit mode */
export interface PostFormData {
  id: string;
  title: string;
  slug: string;
  author: string;
  leadParagraph: string;
  content: string;
  featuredImage: string;
  publishing: Publishing;
  attractiveness: Attractiveness;
}
