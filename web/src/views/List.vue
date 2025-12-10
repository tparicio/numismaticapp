<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-3xl font-bold">My Collection</h2>
      <div class="join">
        <button class="join-item btn" :class="{ 'btn-active': viewMode === 'grid' }" @click="viewMode = 'grid'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
        </button>
        <button class="join-item btn" :class="{ 'btn-active': viewMode === 'table' }" @click="viewMode = 'table'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" /></svg>
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="flex justify-center p-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="coins.length === 0" class="text-center p-10 bg-base-100 rounded-box">
      <p>No coins found. Start by adding one!</p>
      <router-link to="/add" class="btn btn-primary mt-4">Add Coin</router-link>
    </div>

    <!-- Grid View -->
    <div v-else-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="coin in coins" :key="coin.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow cursor-pointer">
        <figure class="px-4 pt-4 relative group flex justify-center gap-2">
          <div class="relative group/img cursor-zoom-in" @click.stop="openViewer(coin, 'front')">
             <!-- Overlay for zoom hint -->
            <div class="relative group/img cursor-zoom-in" @click.stop="openViewer(coin, 'front')">
             <!-- Overlay for zoom hint -->
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover/img:bg-opacity-20 transition-all flex items-center justify-center z-10 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white opacity-0 group-hover/img:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                </svg>
            </div>
            <img :src="getThumbnailUrl(coin, 'front')" alt="Coin Front" class="rounded-full h-24 w-24 object-cover shadow-md hover:scale-110 transition-transform duration-300" />
          </div>
          <div class="relative group/img cursor-zoom-in" @click.stop="openViewer(coin, 'back')">
             <!-- Overlay for zoom hint -->
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover/img:bg-opacity-20 transition-all flex items-center justify-center z-10 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white opacity-0 group-hover/img:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                </svg>
            </div>
            <img :src="getThumbnailUrl(coin, 'back')" alt="Coin Back" class="rounded-full h-24 w-24 object-cover shadow-md hover:scale-110 transition-transform duration-300" />
          </div>
        </figure>
        <div class="card-body" @click="goToDetail(coin.id)">
          <h2 class="card-title flex-col items-start gap-1">
            <span v-if="coin.name" class="text-lg font-bold text-primary">{{ coin.name }}</span>
            <span class="text-base font-normal">
                {{ coin.country }} {{ coin.face_value }}
                <div class="badge badge-secondary ml-2" v-if="coin.year && coin.year !== 0">{{ coin.year }}</div>
            </span>
          </h2>
          <p class="text-sm text-gray-500">{{ coin.currency }}</p>
          <div class="card-actions justify-end mt-2">
            <div class="badge badge-outline" v-if="coin.grade">{{ coin.grade }}</div>
            <div class="badge badge-outline">{{ coin.material }}</div>
            <div class="flex gap-2 ml-auto">
                <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-square btn-sm btn-ghost text-info">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                </button>
                <button @click.stop="confirmDelete(coin)" class="btn btn-square btn-sm btn-ghost text-error">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                    </svg>
                </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Table View -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-box shadow-xl">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>Images</th>
            <th>Name</th>
            <th>Mint</th>
            <th>Mintage</th>
            <th>Country</th>
            <th>Value</th>
            <th>Year</th>
            <th>Currency</th>
            <th>Grade</th>
            <th>Material</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="coin in coins" :key="coin.id" class="hover cursor-pointer" @click="goToDetail(coin.id)">
            <td>
                <div class="flex gap-2">
                    <div class="avatar cursor-zoom-in" @click.stop="openViewer(coin, 'front')">
                        <div class="w-12 h-12 rounded-full overflow-hidden">
                            <img :src="getThumbnailUrl(coin, 'front')" alt="Front" class="hover:scale-110 transition-transform duration-300" />
                        </div>
                    </div>
                    <div class="avatar cursor-zoom-in" @click.stop="openViewer(coin, 'back')">
                        <div class="w-12 h-12 rounded-full overflow-hidden">
                            <img :src="getThumbnailUrl(coin, 'back')" alt="Back" class="hover:scale-110 transition-transform duration-300" />
                        </div>
                    </div>
                </div>
            </td>
            <td class="font-bold text-primary">{{ coin.name || '-' }}</td>
            <td>{{ coin.mint || '-' }}</td>
            <td>{{ formatMintage(coin.mintage) }}</td>
            <td class="font-semibold">{{ coin.country }}</td>
            <td>{{ coin.face_value }}</td>
            <td>{{ (coin.year && coin.year !== 0) ? coin.year : '-' }}</td>
            <td>{{ coin.currency }}</td>
            <td><div class="badge badge-ghost" v-if="coin.grade">{{ coin.grade }}</div><span v-else>-</span></td>
            <td>{{ coin.material }}</td>
            <td>
                <div class="flex gap-1">
                    <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-square btn-sm btn-ghost text-info">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                        </svg>
                    </button>
                    <button @click.stop="confirmDelete(coin)" class="btn btn-square btn-sm btn-ghost text-error">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                        </svg>
                    </button>
                </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Delete Modal -->
    <dialog id="delete_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-error">Delete Coin</h3>
        <p class="py-4">Are you sure you want to delete <span class="font-bold">{{ coinToDelete?.name || 'this coin' }}</span>? This action cannot be undone.</p>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import ImageViewer from '../components/ImageViewer.vue'
import { formatMintage } from '../utils/formatters'

const coins = ref([])
const loading = ref(true)
const viewMode = ref('grid')
const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const STORAGE_URL = 'http://localhost:8080' // Base URL for static files

const viewerOpen = ref(false)
const viewerImage = ref('')

// Delete Modal State
const deleteModalOpen = ref(false)
const coinToDelete = ref(null)
const deleting = ref(false)

const confirmDelete = (coin) => {
    coinToDelete.value = coin
    deleteModalOpen.value = true
}

const deleteCoin = async () => {
    if (!coinToDelete.value) return
    deleting.value = true
    try {
        await axios.delete(`${API_URL}/coins/${coinToDelete.value.id}`)
        coins.value = coins.value.filter(c => c.id !== coinToDelete.value.id)
        deleteModalOpen.value = false
        coinToDelete.value = null
    } catch (e) {
        console.error("Failed to delete coin", e)
        alert("Failed to delete coin")
    } finally {
        deleting.value = false
    }
}

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

const getFullResUrl = (coin, side = 'front') => {
    // Try to find processed crop first, then original
    if (coin.images && coin.images.length > 0) {
        const processed = coin.images.find(img => img.image_type === 'crop' && img.side === side)
        if (processed) return getImageUrl(processed.path)
        
        const original = coin.images.find(img => img.image_type === 'original' && img.side === side)
        if (original) return getImageUrl(original.path)
    }
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const openViewer = (coin, side) => {
    viewerImage.value = getFullResUrl(coin, side)
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
