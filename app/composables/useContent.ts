import { computed, onMounted } from 'vue'
import { useContentStore, type Page, type PageSection, type BlogPost, type Faq } from '~/stores/content'

/**
 * Composable untuk mengelola konten global (Pages, Blog Posts, FAQs)
 * dengan menjembatani langsung ke Pinia content store.
 */
export const useContent = () => {
  const contentStore = useContentStore()

  // Ambil data halaman dari backend pada saat inisialisasi di client
  onMounted(() => {
    contentStore.fetchPages()
  })

  const pages = computed(() => contentStore.pages)
  const blogPosts = computed(() => contentStore.blogPosts)
  const faqs = computed(() => contentStore.faqs)

  const addPage = async (title: string, slug: string, status: 'published' | 'draft' = 'draft') => {
    const cleanSlug = slug.startsWith('/') ? slug : `/${slug}`
    return await contentStore.addPage({ title, slug: cleanSlug, status })
  }

  const updatePageSections = async (slug: string, sections: PageSection[]) => {
    await contentStore.updatePageSections(slug, sections)
  }

  const updatePageMeta = async (id: string | number, updates: Partial<Omit<Page, 'id' | 'sections'>>) => {
    await contentStore.updatePage(id, updates)
  }

  const deletePage = async (id: string | number) => {
    await contentStore.deletePage(id)
  }

  return {
    pages,
    blogPosts,
    faqs,
    addPage,
    updatePageSections,
    updatePageMeta,
    deletePage
  }
}
