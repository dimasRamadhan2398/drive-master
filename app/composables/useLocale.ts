export const useLocale = () => {
  const { locale, locales, setLocale } = useI18n()

  const currentLocale = computed(() => locale.value)
  const availableLocales = computed(() => locales.value)

  const switchLocale = (code: string) => {
    setLocale(code)
  }

  const isCurrentLocale = (code: string) => locale.value === code

  return {
    locale: currentLocale,
    locales: availableLocales,
    switchLocale,
    isCurrentLocale,
  }
}