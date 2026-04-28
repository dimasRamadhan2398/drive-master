import { ref } from 'vue'

// Tipe data untuk konsistensi
export interface PageSection {
  id: string
  type: 'hero' | 'feature_grid' | 'cta' | 'testimonial' | 'logocloud'
  props: Record<string, any>
}

export interface Page {
  id: number
  title: string
  slug: string
  lastUpdated: string
  status: 'published' | 'draft'
  sections: PageSection[]
}

export interface BlogPost {
  id: number
  title: string
  author: string
  date: string
  status: 'published' | 'draft'
  views: number
  content: string
  media: { name: string; type: string; size: string; url: string; fileType: 'image' | 'video' }[]
}

export interface Faq {
  id: number
  question: string
  answer: string
  order: number
}

// Data awal (mock data)
const initialPages = ref<Page[]>([
  {
    id: 1,
    title: 'Home',
    slug: '/',
    lastUpdated: 'Apr 10, 2026',
    status: 'published',
    sections: [
      // Anda bisa menambahkan section default di sini jika perlu
      // Contoh:
      // {
      //   id: 'hero-1',
      //   type: 'hero',
      //   props: {
      //     title: 'Ini Judul dari CMS',
      //     description: 'Deskripsi ini juga bisa diubah dari admin.'
      //   }
      // }
    ]
  },
  {
    id: 2,
    title: 'About Us',
    slug: '/about',
    lastUpdated: 'Apr 9, 2026',
    status: 'draft',
    sections: []
  }
])

const initialBlogPosts = ref<BlogPost[]>([
  {
    id: 1,
    title: 'Welcome to Our New Blog!',
    author: 'Admin',
    date: 'Apr 10, 2026',
    status: 'published',
    views: 120,
    content: 'This is the first post on our brand new blog. Stay tuned for more updates!',
    media: []
  }
])

const initialFaqs = ref<Faq[]>([
  { id: 1, question: 'What is an EV?', answer: 'An Electric Vehicle (EV) is a vehicle that uses one or more electric motors for propulsion.', order: 1 },
  { id: 2, question: 'How long does it take to charge?', answer: 'Charging time depends on the charger and the vehicle. A typical home charger can fully charge an EV overnight.', order: 2 }
])

/**
 * Composable untuk mengelola konten global (Pages, Blog Posts, FAQs).
 * Menggunakan useState untuk menciptakan state singleton yang dibagikan di seluruh aplikasi.
 */
export const useContent = () => {
  const pages = useState<Page[]>('content-pages', () => initialPages.value)
  const blogPosts = useState<BlogPost[]>('content-blog-posts', () => initialBlogPosts.value)
  const faqs = useState<Faq[]>('content-faqs', () => initialFaqs.value)

  return {
    pages,
    blogPosts,
    faqs
  }
}
