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
        
        <!-- Lang Selector -->
        <select v-model="$i18n.locale" class="select select-ghost select-sm w-16" @change="saveLocale">
            <option value="es">ES</option>
            <option value="en">EN</option>
        </select>

        <!-- Theme Toggle -->
        <button class="btn btn-ghost btn-circle" @click="toggleTheme">
            <svg v-if="theme === 'light'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" /></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
        </button>
      </div>
    </div>
    <div class="container mx-auto p-2 sm:p-4">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const theme = ref('light')
const mobileMenuOpen = ref(false)
const dropdownRef = ref(null)

const toggleMobileMenu = () => {
    mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
    mobileMenuOpen.value = false
}

const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', theme.value)
    localStorage.setItem('theme', theme.value)
}

const saveLocale = (e) => {
    localStorage.setItem('locale', e.target.value)
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

onMounted(() => {
    const savedTheme = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
    theme.value = savedTheme
    document.documentElement.setAttribute('data-theme', savedTheme)
    
    // Add click outside listener
    document.addEventListener('click', handleClickOutside)
})

// Cleanup listener on unmount
import { onUnmounted } from 'vue'
onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>
