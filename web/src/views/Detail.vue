<template>
  <div v-if="coin" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Images Section -->
    <div class="space-y-4">
      <div class="card bg-base-100 shadow-xl">
        <figure class="p-6 bg-base-200">
          <img :src="getImageUrl(coin.sample_image_url_front)" class="rounded-full shadow-lg max-w-xs" alt="Front" />
        </figure>
        <div class="card-body py-2 text-center">
          <span class="font-bold">Obverse</span>
        </div>
      </div>
      
      <div class="card bg-base-100 shadow-xl">
        <figure class="p-6 bg-base-200">
          <img :src="getImageUrl(coin.sample_image_url_back)" class="rounded-full shadow-lg max-w-xs" alt="Back" />
        </figure>
        <div class="card-body py-2 text-center">
          <span class="font-bold">Reverse</span>
        </div>
      </div>
    </div>

    <!-- Details Section -->
    <div class="card bg-base-100 shadow-xl h-fit">
      <div class="card-body">
        <h1 class="card-title text-4xl mb-2">{{ coin.country }} {{ coin.face_value }}</h1>
        <div class="flex gap-2 mb-6">
            <div class="badge badge-lg badge-primary">{{ coin.year }}</div>
            <div class="badge badge-lg badge-secondary">{{ coin.currency }}</div>
            <div class="badge badge-lg badge-accent">{{ coin.grade || 'Un-graded' }}</div>
        </div>

        <div class="divider">Details</div>

        <div class="grid grid-cols-2 gap-4">
            <div>
                <span class="font-bold block text-sm text-gray-500">Material</span>
                <span>{{ coin.material }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">KM Code</span>
                <span>{{ coin.km_code || 'N/A' }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">Est. Value</span>
                <span>{{ coin.min_value }} - {{ coin.max_value }}</span>
            </div>
             <div>
                <span class="font-bold block text-sm text-gray-500">Added On</span>
                <span>{{ new Date(coin.created_at).toLocaleDateString() }}</span>
            </div>
        </div>

        <div class="divider">Description</div>
        <p class="text-justify">{{ coin.description }}</p>

        <div class="divider">Notes</div>
        <p class="text-sm italic">{{ coin.notes || 'No notes available.' }}</p>
        
        <div class="card-actions justify-end mt-8">
            <router-link to="/list" class="btn btn-ghost">Back to Gallery</router-link>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="flex justify-center p-20">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'

const route = useRoute()
const coin = ref(null)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080'

const getImageUrl = (path) => {
    if (!path) return 'https://via.placeholder.com/150'
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

onMounted(async () => {
  try {
    const res = await axios.get(`${API_URL}/coins/${route.params.id}`)
    coin.value = res.data
  } catch (e) {
    console.error(e)
  }
})
</script>
