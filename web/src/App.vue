<template>
  <div class="min-h-screen bg-base-200">
    <div class="navbar bg-base-100 shadow-lg">
      <div class="navbar-start">
        <div class="dropdown" ref="dropdownRef">
          <label tabindex="0" class="btn btn-ghost lg:hidden" @click="toggleMobileMenu">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h8m-8 6h16" /></svg>
          </label>
          <ul v-show="mobileMenuOpen" tabindex="0" class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52 z-[100]">
            <li><router-link to="/" @click="closeMobileMenu">{{ $t('nav.home') }}</router-link></li>
            <li><router-link to="/list" @click="closeMobileMenu">{{ $t('nav.collection') }}</router-link></li>
            <li><router-link to="/groups" @click="closeMobileMenu">{{ $t('nav.groups') }}</router-link></li>
            <li><router-link to="/add" @click="closeMobileMenu">{{ $t('nav.add_coin') }}</router-link></li>
          </ul>
        </div>
        <router-link to="/" class="btn btn-ghost normal-case text-xl flex items-center gap-2 font-['Cinzel']">
            <img src="/icon.png" alt="Logo" class="w-8 h-8" />
            <span class="hidden sm:inline"><span class="font-bold">Numismatic</span><span class="font-normal opacity-80">App</span></span>
        </router-link>
      </div>
      <div class="navbar-center hidden lg:flex">
        <ul class="menu menu-horizontal px-1">
          <li><router-link to="/">{{ $t('nav.home') }}</router-link></li>
          <li><router-link to="/list">{{ $t('nav.collection') }}</router-link></li>
          <li><router-link to="/groups">{{ $t('nav.groups') }}</router-link></li>
        </ul>
      </div>
      <div class="navbar-end gap-2">
        <router-link to="/add" class="hidden sm:inline-flex btn btn-primary btn-sm text-white mr-2">{{ $t('nav.add_coin') }}</router-link>
        
        <!-- Settings Dropdown -->
        <div class="dropdown dropdown-end">
            <label tabindex="0" class="btn btn-ghost btn-circle">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </label>
            <ul tabindex="0" class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-64 z-[100]">
                <!-- Theme -->
                <li>
                    <div class="form-control">
                        <label class="label cursor-pointer gap-2 w-full justify-between">
                            <span class="label-text">{{ $t('settings.dark_mode') }}</span>
                            <input type="checkbox" class="toggle toggle-sm" :checked="theme === 'dark'" @change="settingsStore.toggleTheme" />
                        </label>
                    </div>
                </li>
                <!-- Language -->
                <li>
                    <div class="form-control">
                        <label class="label cursor-pointer gap-2 w-full justify-between">
                            <span class="label-text">{{ $t('settings.language') }}</span>
                            <select v-model="settingsStore.locale" class="select select-bordered select-xs">
                                <option value="es">ES</option>
                                <option value="en">EN</option>
                            </select>
                        </label>
                    </div>
                </li>
                <!-- Privacy -->
                <li>
                    <div class="form-control">
                        <label class="label cursor-pointer gap-2 w-full justify-between">
                            <span class="label-text">{{ $t('settings.hide_values') }}</span>
                            <input type="checkbox" class="toggle toggle-sm" :checked="privacyMode" @change="settingsStore.togglePrivacyMode" />
                        </label>
                    </div>
                </li>
                <div class="divider my-1"></div>
                <!-- Export -->
                <li class="menu-title">
                    <span>{{ $t('settings.export') }}</span>
                </li>
                <li><a @click="exportCSV">CSV</a></li>
                <li><a @click="exportSQL">SQL</a></li>
            </ul>
        </div>
      </div>
    </div>
    <div class="container mx-auto p-2 sm:p-4">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useSettingsStore } from './stores/settings'

const route = useRoute()
const { locale } = useI18n()
const settingsStore = useSettingsStore()
const { theme, privacyMode } = storeToRefs(settingsStore) // We use store refs for reactivity

const mobileMenuOpen = ref(false)
const dropdownRef = ref(null)

// Sync locale from store to i18n
watch(() => settingsStore.locale, (newLocale) => {
    locale.value = newLocale
}, { immediate: true })

const toggleMobileMenu = () => {
    mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
    mobileMenuOpen.value = false
}

// Close mobile menu when route changes
watch(() => route.path, () => {
    closeMobileMenu()
})

// Close mobile menu when clicking outside
const handleClickOutside = (event) => {
    if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
        closeMobileMenu()
    }
}

// Export functions (placeholders for now)
const exportCSV = () => {
    window.location.href = '/api/v1/export/csv'
}

const exportSQL = () => {
    window.location.href = '/api/v1/export/sql'
}

onMounted(() => {
    // Add click outside listener
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>
