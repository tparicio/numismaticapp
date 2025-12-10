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
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full z-10">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getThumbnailUrl(coin, 'front')" class="rounded-full shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" alt="Front" />
                </figure>
                <div class="mt-4 font-bold text-lg">Obverse</div>
            </div>

            <!-- Back -->
            <div class="text-center">
                <figure class="cursor-zoom-in relative group inline-block" @click="openViewer('back')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full z-10">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getThumbnailUrl(coin, 'back')" class="rounded-full shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" alt="Back" />
                </figure>
                <div class="mt-4 font-bold text-lg">Reverse</div>
            </div>
        </div>
      </div>
    </div>

    <!-- Details Section -->
    <div class="card bg-base-100 shadow-xl h-fit">
      <div class="card-body">
        <h2 v-if="coin.name" class="text-2xl font-bold text-primary mb-1">{{ coin.name }}</h2>
        <h1 class="card-title text-4xl mb-2">{{ coin.country }} {{ coin.face_value }}</h1>
        <div class="flex gap-2 mb-6">
            <div class="badge badge-lg badge-primary" v-if="coin.year && coin.year !== 0">{{ coin.year }}</div>
            <div class="badge badge-lg badge-secondary">{{ coin.currency }}</div>
            <div class="tooltip" :data-tip="getGradeDescription(coin.grade)" v-if="coin.grade">
                <div class="badge badge-lg badge-accent cursor-help">{{ coin.grade }}</div>
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
                <span class="font-bold block text-sm text-gray-500">Mint (Ceca)</span>
                <span>{{ coin.mint || 'N/A' }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">Mintage (Tirada)</span>
                <span>{{ formatMintage(coin.mintage) }}</span>
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
        
        <div class="card-actions justify-end mt-8 gap-2">
            <router-link to="/list" class="btn btn-ghost">Back to Gallery</router-link>
            <router-link :to="`/edit/${coin.id}`" class="btn btn-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                </svg>
                Edit
            </router-link>
            <button @click="deleteModalOpen = true" class="btn btn-error">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                Delete
            </button>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="flex justify-center p-20">
    <span class="loading loading-spinner loading-lg"></span>
  </div>

  <!-- Delete Modal -->
  <dialog id="delete_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-error">Delete Coin</h3>
      <p class="py-4">Are you sure you want to delete <span class="font-bold">{{ coin?.name || 'this coin' }}</span>? This action cannot be undone.</p>
      <div class="modal-action">
        <button class="btn" @click="deleteModalOpen = false">Cancel</button>
        <button class="btn btn-error" @click="deleteCoin" :disabled="deleting">
          <span v-if="deleting" class="loading loading-spinner"></span>
          Delete
        </button>
      </div>
    </div>
  </dialog>

  <ImageViewer 
      :is-open="viewerOpen" 
      :image-url="viewerImage" 
      @close="viewerOpen = false" 
    />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import { GRADES } from '../constants/grades'
import ImageViewer from '../components/ImageViewer.vue'
import { formatMintage } from '../utils/formatters'

const route = useRoute()
const router = useRouter()
const coin = ref(null)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080'

const viewerOpen = ref(false)
const viewerImage = ref('')

// Delete Modal State
const deleteModalOpen = ref(false)
const deleting = ref(false)

const deleteCoin = async () => {
    if (!coin.value) return
    deleting.value = true
    try {
        await axios.delete(`${API_URL}/coins/${coin.value.id}`)
        router.push('/list')
    } catch (e) {
        console.error("Failed to delete coin", e)
        alert("Failed to delete coin")
    } finally {
        deleting.value = false
        deleteModalOpen.value = false
    }
}

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
