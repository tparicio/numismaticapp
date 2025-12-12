import { createI18n } from 'vue-i18n'
import es from './locales/es.json'
import en from './locales/en.json'

// Detect user language
const getUserLocale = () => {
    const saved = localStorage.getItem('locale')
    if (saved) return saved

    const browser = navigator.language.split('-')[0]
    return browser === 'es' ? 'es' : 'es' // Default to ES as requested if strict fallback needed, or EN
}

// User requested: "try to detect from browser, otherwise use ES"
// So if browser is EN -> EN. If browser is FR -> ES.

const detectDefault = () => {
    const browser = navigator.language.split('-')[0]
    return ['es', 'en'].includes(browser) ? browser : 'es'
}

const i18n = createI18n({
    legacy: false, // Compostion API
    locale: localStorage.getItem('locale') || detectDefault(),
    fallbackLocale: 'es',
    messages: {
        es,
        en
    }
})

export default i18n
