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
      <div v-for="coin in coins" :key="coin.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow cursor-pointer" @click="goToDetail(coin.id)">
        <figure class="px-4 pt-4 bg-base-200">
            <!-- Show front image. Assuming API serves static files at /storage or we use full URL if stored as such.
                 In our backend we stored full path. We need to fix this to be a URL.
                 For now, let's assume we can map it or the backend returns a relative path.
                 Ideally backend should return a proper URL.
                 Let's assume the backend returns absolute path, we need to strip it or backend should have returned relative.
                 Let's fix backend later, for now let's try to construct URL.
                 If path is /app/storage/..., we need to map it to http://localhost:8080/storage/...
            -->
          <img :src="getImageUrl(coin.sample_image_url_front)" alt="Coin Front" class="rounded-xl h-48 w-48 object-cover" />
        </figure>
        <div class="card-body">
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const coins = ref([])
const loading = ref(true)
const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080' // Base URL for static files

const getImageUrl = (path) => {
    if (!path) return 'https://via.placeholder.com/150'
    // Hacky fix for absolute paths from docker
    // If path contains "storage", take everything after
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
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
