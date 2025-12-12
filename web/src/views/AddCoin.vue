<template>
  <div class="card bg-base-100 shadow-xl max-w-2xl mx-auto">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">{{ $t('form.add_title') }}</h2>
      
      <form @submit.prevent="uploadCoin" class="space-y-6">
        
        <!-- Drag & Drop Zone -->
        <div 
          class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors relative"
          :class="{ 'border-primary bg-base-200': isDragging, 'border-base-300 hover:border-primary': !isDragging, 'bg-error/10 border-error': error }"
          @dragover.prevent="isDragging = true"
          @dragleave.prevent="isDragging = false"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <input 
            type="file" 
            ref="fileInput" 
            class="hidden" 
            multiple 
            accept="image/*" 
            @change="handleFileSelect" 
          />
          
          <div v-if="!frontFile && !backFile" class="space-y-2">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto text-base-content/50">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
            </svg>
            <p class="font-bold text-lg">{{ $t('form.drag_drop.main') }}</p>
            <p class="text-sm opacity-70">{{ $t('form.drag_drop.sub') }}</p>
          </div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Front Preview -->
            <div class="relative group">
              <div class="badge badge-primary absolute top-2 left-2 z-10">{{ $t('form.drag_drop.front') }}</div>
              <img :src="frontPreview" class="w-full h-48 object-cover rounded-lg shadow-md" />
              <button type="button" @click.stop="removeFile('front')" class="btn btn-circle btn-xs btn-error absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
            </div>

            <!-- Back Preview -->
            <div class="relative group">
              <div class="badge badge-secondary absolute top-2 left-2 z-10">{{ $t('form.drag_drop.back') }}</div>
              <img :src="backPreview" class="w-full h-48 object-cover rounded-lg shadow-md" />
              <button type="button" @click.stop="removeFile('back')" class="btn btn-circle btn-xs btn-error absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
            </div>
          </div>

          <!-- Swap Button -->
          <div v-if="frontFile && backFile" class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20">
             <button type="button" @click.stop="swapImages" class="btn btn-circle btn-primary shadow-lg tooltip flex items-center justify-center" :data-tip="$t('form.drag_drop.swap')">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
                </svg>
             </button>
          </div>
        </div>
        
        <div v-if="error" class="alert alert-error shadow-lg text-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          <span>{{ error }}</span>
        </div>

        <!-- AI Settings -->
        <GeminiConfig 
          v-model:isOpen="aiSettingsOpen"
          v-model:model="selectedModel"
          v-model:temperature="temperature"
          :available-models="availableModels"
        />

        <!-- Group Selector -->
        <GroupSelector v-model="selectedGroup" />

        <!-- User Notes -->
        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">{{ $t('form.fields.personal_notes') }}</span>
          </label>
          <textarea v-model="userNotes" class="textarea textarea-bordered h-24" :placeholder="$t('form.placeholders.notes')"></textarea>
        </div>

        <div class="alert alert-info shadow-lg" v-if="uploading">
          <div>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current flex-shrink-0 w-6 h-6 animate-spin"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>{{ $t('form.ai_analysis') }}</span>
          </div>
        </div>

        <div class="card-actions justify-end mt-6">
          <button type="submit" class="btn btn-primary w-full sm:w-auto" :disabled="uploading || !frontFile || !backFile">
            <span v-if="uploading" class="loading loading-spinner"></span>
            {{ uploading ? $t('form.processing') : $t('nav.add_coin') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import GroupSelector from '../components/GroupSelector.vue'
import GeminiConfig from '../components/GeminiConfig.vue'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()
const fileInput = ref(null)
const frontFile = ref(null)
const backFile = ref(null)
const frontPreview = ref(null)
const backPreview = ref(null)
const userNotes = ref('')
const selectedGroup = ref('')
const uploading = ref(false)
const isDragging = ref(false)
const error = ref(null)
const aiSettingsOpen = ref(false)

// AI Config
const availableModels = ref([])
const selectedModel = ref('gemini-2.5-flash') // Default fallback
const temperature = ref(0.4)

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

const selectedModelDescription = computed(() => {
  const model = availableModels.value.find(m => m.name === selectedModel.value)
  return model ? model.description : ''
})

onMounted(async () => {
  // Preselect last used group
  const lastGroup = localStorage.getItem('lastSelectedGroup')
  if (lastGroup) {
      selectedGroup.value = lastGroup
  }

  // Load AI Preferences
  const lastModel = localStorage.getItem('lastSelectedModel')
  if (lastModel) selectedModel.value = lastModel
  const lastTemp = localStorage.getItem('lastTemperature')
  if (lastTemp) temperature.value = parseFloat(lastTemp)

  // Fetch Models
  try {
      const res = await axios.get(`${API_URL}/gemini/models`)
      if (res.data && res.data.length > 0) {
          availableModels.value = res.data
          // If stored model not in list (and list not empty), default to first or specific default
          const modelExists = availableModels.value.some(m => m.name === selectedModel.value)
          if (!modelExists) {
             // Try to find gemini-2.5-flash or gemini-1.5-flash
             const cleanDefault = availableModels.value.find(m => m.name.includes('gemini-2.5-flash') || m.name.includes('gemini-1.5-flash'))
             if (cleanDefault) {
                 selectedModel.value = cleanDefault.name
             } else {
                 selectedModel.value = availableModels.value[0].name
             }
          }
      }
  } catch (e) {
      console.error("Failed to fetch models", e)
      // Fallback manual list if API fails? Or just rely on default string
      availableModels.value = [
          { name: 'gemini-2.0-flash', description: 'Fast and versatile' },
          { name: 'gemini-1.5-flash', description: 'Previous generic model' }
      ]
  }
})

const triggerFileInput = () => {
  fileInput.value.click()
}

const handleFileSelect = (e) => {
  processFiles(e.target.files)
}

const handleDrop = (e) => {
  isDragging.value = false
  processFiles(e.dataTransfer.files)
}

const processFiles = (files) => {
  error.value = null
  
  const newFiles = Array.from(files).filter(f => f.type.startsWith('image/'))
  
  if (newFiles.length === 0) {
      if (files.length > 0) error.value = t('form.errors.images_only')
      return
  }

  if (newFiles.length === 2) {
      setFile('front', newFiles[0])
      setFile('back', newFiles[1])
  } else if (newFiles.length === 1) {
      if (!frontFile.value) {
          setFile('front', newFiles[0])
      } else if (!backFile.value) {
          setFile('back', newFiles[0])
      } else {
          error.value = "Both slots are full. Remove one to add another, or select 2 files to replace both."
      }
  } else {
      error.value = t('form.errors.both_images')
  }
}

const setFile = (side, file) => {
    const reader = new FileReader()
    reader.onload = (e) => {
        if (side === 'front') {
            frontFile.value = file
            frontPreview.value = e.target.result
        } else {
            backFile.value = file
            backPreview.value = e.target.result
        }
    }
    reader.readAsDataURL(file)
}

const removeFile = (side) => {
    if (side === 'front') {
        frontFile.value = null
        frontPreview.value = null
    } else {
        backFile.value = null
        backPreview.value = null
    }
}

const swapImages = () => {
    const tempFile = frontFile.value
    const tempPreview = frontPreview.value
    frontFile.value = backFile.value
    frontPreview.value = backPreview.value
    backFile.value = tempFile
    backPreview.value = tempPreview
}

const uploadCoin = async () => {
  if (!frontFile.value || !backFile.value) {
      error.value = t('form.errors.both_images')
      return
  }

  uploading.value = true
  error.value = null
  
  // Save preferences
  if (selectedGroup.value) localStorage.setItem('lastSelectedGroup', selectedGroup.value)
  if (selectedModel.value) localStorage.setItem('lastSelectedModel', selectedModel.value)
  localStorage.setItem('lastTemperature', temperature.value)

  const formData = new FormData()
  formData.append('front_image', frontFile.value)
  formData.append('back_image', backFile.value)
  formData.append('group_name', selectedGroup.value)
  formData.append('user_notes', userNotes.value)
  formData.append('model_name', selectedModel.value)
  formData.append('temperature', temperature.value.toString())

  try {
    const res = await axios.post(`${API_URL}/coins`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    // Redirect to detail view
    router.push(`/coin/${res.data.id}`)
  } catch (e) {
    console.error(e)
    error.value = t('form.errors.upload_failed') + ': ' + (e.response?.data?.error || e.message)
  } finally {
    uploading.value = false
  }
}
</script>
