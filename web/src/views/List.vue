<template>
  <div>
    <h2 class="text-3xl font-bold mb-6">My Collection</h2>
    
    <div v-if="loading" class="flex justify-center p-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="coins.length === 0" class="text-center p-10 bg-base-100 rounded-box">
      <p>No coins found. Start by adding one!</p>
      <router-link to="/add" class="btn btn-primary mt-4">Add Coin</router-link>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="coin in coins" :key="coin.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow cursor-pointer">
        <figure class="px-4 pt-4 relative group flex justify-center gap-2" @click.stop="openViewer(coin)">
             <!-- Overlay for zoom hint -->
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-20 transition-all flex items-center justify-center z-10 rounded-xl">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                </svg>
            </div>
          <img :src="getThumbnailUrl(coin, 'front')" alt="Coin Front" class="rounded-full h-24 w-24 object-cover shadow-md" />
          <img :src="getThumbnailUrl(coin, 'back')" alt="Coin Back" class="rounded-full h-24 w-24 object-cover shadow-md" />
        </figure>
        <div class="card-body" @click="goToDetail(coin.id)">
          <h2 class="card-title">
            {{ coin.country }} {{ coin.face_value }}
            <div class="badge badge-secondary" v-if="coin.year">{{ coin.year }}</div>
          </h2>
          <p class="text-sm text-gray-500">{{ coin.currency }}</p>
          <div class="card-actions justify-end mt-2">
            <div class="badge badge-outline">{{ coin.grade || 'N/A' }}</div>
            <div class="badge badge-outline">{{ coin.material }}</div>
          </div>
        </div>
      </div>
    </div>

    <ImageViewer 
      :is-open="viewerOpen" 
      :image-url="viewerImage" 
      @close="viewerOpen = false" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import ImageViewer from '../components/ImageViewer.vue'

const coins = ref([])
const loading = ref(true)
const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080' // Base URL for static files

const viewerOpen = ref(false)
const viewerImage = ref('')

const getImageUrl = (path) => {
    if (!path) return 'https://via.placeholder.com/150'
    // Hacky fix for absolute paths from docker
    // If path contains "storage", take everything after
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const getThumbnailUrl = (coin, side = 'front') => {
    // Try to find thumbnail in images array
    if (coin.images && coin.images.length > 0) {
        const thumb = coin.images.find(img => img.image_type === 'thumbnail' && img.side === side)
        if (thumb) {
            return getImageUrl(thumb.path)
        }
    }
    // Fallback to sample image
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const getFullResUrl = (coin) => {
    // Try to find processed crop first, then original
    if (coin.images && coin.images.length > 0) {
        const processed = coin.images.find(img => img.image_type === 'crop' && img.side === 'front')
        if (processed) return getImageUrl(processed.path)
        
        const original = coin.images.find(img => img.image_type === 'original' && img.side === 'front')
        if (original) return getImageUrl(original.path)
    }
    return getImageUrl(coin.sample_image_url_front)
}

const openViewer = (coin) => {
    viewerImage.value = getFullResUrl(coin)
    viewerOpen.value = true
}

const goToDetail = (id) => {
  router.push(`/coin/${id}`)
}

onMounted(async () => {
  try {
    const res = await axios.get(`${API_URL}/coins?limit=50`)
    coins.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
