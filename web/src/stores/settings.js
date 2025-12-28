import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useSettingsStore = defineStore('settings', () => {
    // State
    const theme = ref(localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
    const locale = ref(localStorage.getItem('locale') || 'es')
    const privacyMode = ref(localStorage.getItem('privacy_mode') === 'true')

    // Actions
    const toggleTheme = () => {
        theme.value = theme.value === 'light' ? 'dark' : 'light'
    }

    const setLocale = (newLocale) => {
        locale.value = newLocale
    }

    const togglePrivacyMode = () => {
        privacyMode.value = !privacyMode.value
    }

    // Persistence & Side Effects
    watch(theme, (newTheme) => {
        localStorage.setItem('theme', newTheme)
        document.documentElement.setAttribute('data-theme', newTheme)
    }, { immediate: true })

    watch(locale, (newLocale) => {
        localStorage.setItem('locale', newLocale)
    })

    watch(privacyMode, (newMode) => {
        localStorage.setItem('privacy_mode', newMode)
    })

    return {
        theme,
        locale,
        privacyMode,
        toggleTheme,
        setLocale,
        togglePrivacyMode
    }
})
