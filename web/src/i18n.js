import { createI18n } from 'vue-i18n'
import es from './locales/es.json'
import en from './locales/en.json'

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
