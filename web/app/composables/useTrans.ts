import { ref } from 'vue'
import en from '../locales/en.json'
import fr from '../locales/fr.json'

const currentLocale = ref<'en' | 'fr'>('fr') // Default to French
const translations = { en, fr }

export function useTrans() {
  const t = (key: string, params?: Record<string, any>) => {
    const keys = key.split('.')
    let val: any = translations[currentLocale.value]
    
    for (const k of keys) {
      if (val && typeof val === 'object') {
        val = val[k]
      } else {
        return key
      }
    }
    
    if (typeof val !== 'string') return key
    
    if (params) {
      let res = val
      for (const [pk, pv] of Object.entries(params)) {
        res = res.replace(new RegExp(`{{\\s*${pk}\\s*}}`, 'g'), String(pv))
      }
      return res
    }
    
    return val
  }
  
  const toggleLocale = () => {
    currentLocale.value = currentLocale.value === 'en' ? 'fr' : 'en'
  }
  
  return {
    t,
    locale: currentLocale,
    toggleLocale,
    setLocale: (l: 'en' | 'fr') => {
      currentLocale.value = l
    }
  }
}
