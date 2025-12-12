<template>
  <div v-if="coin" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Images Section -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        
        <!-- Image Source Toggles -->
        <div class="flex justify-center mb-6">
            <div class="join">
                <button 
                    class="btn join-item" 
                    :class="{ 'btn-primary': activeImageSource === 'processed' }"
                    @click="activeImageSource = 'processed'"
                >
                    {{ $t('details.toggles.processed') }}
                </button>
                <button 
                    class="btn join-item" 
                    :class="{ 'btn-primary': activeImageSource === 'original' }"
                    @click="activeImageSource = 'original'"
                >
                    {{ $t('details.toggles.original') }}
                </button>
            </div>
        </div>

        <div class="flex flex-col sm:flex-row justify-center gap-8 items-center">
            <!-- Front -->
            <div class="text-center">
                <figure class="cursor-zoom-in relative group inline-block" @click="openViewer('front')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full z-10">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getCurrentImageUrl('front')" class="rounded-full shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" alt="Front" @error="handleImageError" />
                </figure>
                <div class="mt-4 font-bold text-lg">{{ $t('details.obverse') }}</div>
            </div>

            <!-- Back -->
            <div class="text-center">
                <figure class="cursor-zoom-in relative group inline-block" @click="openViewer('back')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center rounded-full z-10">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getCurrentImageUrl('back')" class="rounded-full shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" alt="Back" @error="handleImageError" />
                </figure>
                <div class="mt-4 font-bold text-lg">{{ $t('details.reverse') }}</div>
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

        <div class="divider">{{ $t('details.sections.details') }}</div>

        <div class="grid grid-cols-2 gap-4">
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.material') }}</span>
                <span>{{ coin.material }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.km') }}</span>
                <span class="flex items-center gap-2">
                    {{ coin.km_code || 'N/A' }}
                    <a v-if="getNumistaUrl()" 
                       :href="getNumistaUrl()" 
                       target="_blank" 
                       class="btn btn-xs btn-outline btn-info gap-1"
                       title="Ver en Numista"
                    >
                        N
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                        </svg>
                    </a>
                </span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.mint') }}</span>
                <span>{{ coin.mint || 'N/A' }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.mintage') }}</span>
                <span>{{ formatMintage(coin.mintage) }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.est_value') }}</span>
                <span>{{ coin.min_value }} - {{ coin.max_value }}</span>
            </div>
             <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.added_on') }}</span>
                <span>{{ new Date(coin.created_at).toLocaleDateString() }}</span>
            </div>
        </div>

        <div class="divider">{{ $t('details.sections.description') }}</div>
        <p class="text-justify">{{ coin.description }}</p>

        <div v-if="coin.gemini_model" class="mt-4 flex flex-col text-xs text-gray-400">
           <div class="flex items-center justify-between">
               <div class="flex gap-2 items-center">
                   <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
                   </svg>
                   <span>AI Generated by {{ coin.gemini_model }} (Temp: {{ coin.gemini_temperature }})</span>
               </div>
               <div class="tooltip" :data-tip="$t('details.reprocess_tooltip') || 'Reprocesar'">
                   <button @click="openReprocessModal" class="btn btn-xs btn-outline btn-primary flex flex-row items-center gap-1 flex-nowrap">
                       <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 flex-shrink-0">
                           <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                       </svg>
                       <span class="whitespace-nowrap">{{ $t('common.reprocess') || 'Reprocesar' }}</span>
                   </button>
               </div>
           </div>
           <div v-if="coin.gemini_details && coin.gemini_details.error" class="mt-1 text-error">
               Warning: {{ coin.gemini_details.error }}
           </div>
        </div>

        <div class="divider">{{ $t('details.sections.notes') }}</div>
        <p class="text-sm italic">{{ coin.notes || $t('details.no_notes') }}</p>
        
        <div class="card-actions justify-end mt-8 gap-2">
            <router-link to="/list" class="btn btn-ghost">{{ $t('details.back_gallery') }}</router-link>
            <router-link :to="`/edit/${coin.id}`" class="btn btn-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                </svg>
                {{ $t('common.edit') }}
            </router-link>
            <button @click="deleteModalOpen = true" class="btn btn-error">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                {{ $t('common.delete') }}
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
      <h3 class="font-bold text-lg text-error">{{ $t('list.delete_modal.title') }}</h3>
      <p class="py-4">{{ $t('list.delete_modal.confirm') }} <span class="font-bold">{{ coin?.name || $t('common.unknown') }}</span>? {{ $t('list.delete_modal.warning') }}</p>
      <div class="modal-action">
        <button class="btn" @click="deleteModalOpen = false">{{ $t('common.cancel') }}</button>
        <button class="btn btn-error" @click="deleteCoin" :disabled="deleting">
          <span v-if="deleting" class="loading loading-spinner"></span>
          {{ $t('common.delete') }}
        </button>
      </div>
    </div>
  </dialog>

  <!-- Reprocess Modal -->
  <dialog id="reprocess_modal" class="modal" :class="{ 'modal-open': reprocessModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-primary flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
          </svg>
          {{ $t('details.reprocess_modal.title') || 'Reprocesar con IA' }}
      </h3>
      
      <div class="py-4 space-y-4">
          <div class="alert alert-warning text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
            <span>{{ $t('details.reprocess_modal.warning') || '¡Advertencia! Esto sobreescribirá todos los detalles de la moneda con información nueva generada por la IA.' }}</span>
          </div>

          <GeminiConfig 
            v-model:isOpen="aiSettingsOpen"
            v-model:model="selectedModel"
            v-model:temperature="temperature"
            :available-models="availableModels"
            :title="$t('form.ai_settings') || 'Configuración del Modelo'"
          />
      </div>

      <div class="modal-action">
        <button class="btn" @click="reprocessModalOpen = false">{{ $t('common.cancel') }}</button>
        <button class="btn btn-primary" @click="reprocessCoin" :disabled="reprocessing">
          <span v-if="reprocessing" class="loading loading-spinner"></span>
          {{ $t('common.process') || 'Procesar' }}
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
import { normalizeGrade } from '../utils/grades'
import ImageViewer from '../components/ImageViewer.vue'
import GeminiConfig from '../components/GeminiConfig.vue'
import { formatMintage } from '../utils/formatters'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const route = useRoute()
const router = useRouter()
const coin = ref(null)
const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const STORAGE_URL = ''

const viewerOpen = ref(false)
const viewerImage = ref('')
const activeImageSource = ref('processed') // processed, original

// Delete Modal State
const deleteModalOpen = ref(false)
const deleting = ref(false)

// Reprocess Modal State
const reprocessModalOpen = ref(false)
const reprocessing = ref(false)
const selectedModel = ref('gemini-2.5-flash')
const temperature = ref(0.4)
const availableModels = ref([])
const aiSettingsOpen = ref(true)

const openReprocessModal = async () => {
    reprocessModalOpen.value = true
    // Initialize with current values if available
    if (coin.value.gemini_model) selectedModel.value = coin.value.gemini_model
    if (coin.value.gemini_temperature) temperature.value = coin.value.gemini_temperature
    
    // Fetch models if empty
    if (availableModels.value.length === 0) {
        try {
            const res = await axios.get(`${API_URL}/gemini/models`)
            availableModels.value = res.data
        } catch (e) {
            console.error("Failed to fetch models", e)
        }
    }
}

const reprocessCoin = async () => {
    if (!coin.value) return
    reprocessing.value = true
    try {
        const res = await axios.post(`${API_URL}/coins/${coin.value.id}/analyze`, {
            model_name: selectedModel.value,
            temperature: temperature.value
        })
        coin.value = res.data
        reprocessModalOpen.value = false
        // Optional: Show success toast
    } catch (e) {
        console.error(e)
        alert('Failed to reprocess coin: ' + (e.response?.data?.error || e.message))
    } finally {
        reprocessing.value = false
    }
}

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
    if (!path) return '/broken_coin.png'
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const getNumistaUrl = () => {
    if (coin.value && coin.value.numista_number) {
        return `https://es.numista.com/${coin.value.numista_number}`
    }
    return null
}

const getCurrentImageUrl = (side) => {
    if (!coin.value) return '/broken_coin.png'
    
    // Try to find specific image type in images array
    if (coin.value.images && coin.value.images.length > 0) {
        let typeToFind = 'crop'
        if (activeImageSource.value === 'original') {
            typeToFind = 'original'
        }
        
        const img = coin.value.images.find(img => img.image_type === typeToFind && img.side === side)
        if (img) {
            return getImageUrl(img.path)
        }
    }

    // Fallback logic
    return '/broken_coin.png'
}

const handleImageError = (e) => {
    e.target.src = '/broken_coin.png'
}

const openViewer = (side) => {
    viewerImage.value = getCurrentImageUrl(side)
    viewerOpen.value = true
}

const getGradeDescription = (code) => {
    if (!code) return t('common.unknown')
    const base = normalizeGrade(code)
    return t(`grades.${base}.desc`)
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
