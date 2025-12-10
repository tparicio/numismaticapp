<template>
  <div v-if="coin" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Images Section -->
    <!-- Images Section -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="flex flex-col sm:flex-row justify-center gap-8 items-center">
            <!-- Front -->
            <div class="text-center">
                <figure class="cursor-zoom-in relative group inline-block" @click="openViewer('front')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getThumbnailUrl(coin, 'front')" class="rounded-full shadow-lg max-w-xs" alt="Front" />
                </figure>
                <div class="mt-4 font-bold text-lg">Obverse</div>
            </div>

            <!-- Back -->
            <div class="text-center">
                <figure class="cursor-zoom-in relative group inline-block" @click="openViewer('back')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getThumbnailUrl(coin, 'back')" class="rounded-full shadow-lg max-w-xs" alt="Back" />
                </figure>
                <div class="mt-4 font-bold text-lg">Reverse</div>
            </div>
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
            <div class="tooltip" :data-tip="getGradeDescription(coin.grade)">
                <div class="badge badge-lg badge-accent cursor-help">{{ coin.grade || 'Un-graded' }}</div>
            </div>
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

  <ImageViewer 
      :is-open="viewerOpen" 
      :image-url="viewerImage" 
      @close="viewerOpen = false" 
    />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'
import { GRADES } from '../constants/grades'
import ImageViewer from '../components/ImageViewer.vue'

const route = useRoute()
const coin = ref(null)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080'

const viewerOpen = ref(false)
const viewerImage = ref('')

const getImageUrl = (path) => {
    if (!path) return 'https://via.placeholder.com/150'
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const getThumbnailUrl = (coin, side) => {
    if (!coin) return ''
    // Try to find thumbnail in images array
    if (coin.images && coin.images.length > 0) {
        const thumb = coin.images.find(img => img.image_type === 'thumbnail' && img.side === side)
        if (thumb) {
            return getImageUrl(thumb.path)
        }
    }
    // Fallback
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const getFullResUrl = (coin, side) => {
    if (!coin) return ''
    // Try to find processed crop first, then original
    if (coin.images && coin.images.length > 0) {
        const processed = coin.images.find(img => img.image_type === 'crop' && img.side === side)
        if (processed) return getImageUrl(processed.path)
        
        const original = coin.images.find(img => img.image_type === 'original' && img.side === side)
        if (original) return getImageUrl(original.path)
    }
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const openViewer = (side) => {
    viewerImage.value = getFullResUrl(coin.value, side)
    viewerOpen.value = true
}

const getGradeDescription = (code) => {
    if (!code) return 'No grade assigned'
    const g = GRADES.find(g => g.code === code)
    return g ? g.description : 'Unknown grade'
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
