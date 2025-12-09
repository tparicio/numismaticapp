<template>
  <div class="card bg-base-100 shadow-xl max-w-2xl mx-auto">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">Add New Coin</h2>
      
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
            <p class="font-bold text-lg">Tap or Drag & Drop images here</p>
            <p class="text-sm opacity-70">Please select exactly 2 images (Front & Back)</p>
          </div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Front Preview -->
            <div class="relative group">
              <div class="badge badge-primary absolute top-2 left-2 z-10">Front (Obverse)</div>
              <img :src="frontPreview" class="w-full h-48 object-cover rounded-lg shadow-md" />
              <button type="button" @click.stop="removeFile('front')" class="btn btn-circle btn-xs btn-error absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
            </div>

            <!-- Back Preview -->
            <div class="relative group">
              <div class="badge badge-secondary absolute top-2 left-2 z-10">Back (Reverse)</div>
              <img :src="backPreview" class="w-full h-48 object-cover rounded-lg shadow-md" />
              <button type="button" @click.stop="removeFile('back')" class="btn btn-circle btn-xs btn-error absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
            </div>
          </div>

          <!-- Swap Button -->
          <div v-if="frontFile && backFile" class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20">
             <button type="button" @click.stop="swapImages" class="btn btn-circle btn-primary shadow-lg tooltip" data-tip="Swap Front/Back">
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

        <!-- Group Selector -->
        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Group / Book / Tray</span>
            <span class="label-text-alt">Select existing or type to create new</span>
          </label>
          <div class="relative">
            <input 
              type="text" 
              v-model="groupName" 
              list="groups-list"
              placeholder="e.g. 'Spanish Silver', 'Tray 1'" 
              class="input input-bordered w-full" 
              @change="saveLastGroup"
            />
            <datalist id="groups-list">
              <option v-for="g in groups" :key="g.id" :value="g.name" />
            </datalist>
          </div>
        </div>

        <!-- User Notes -->
        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Your Notes</span>
          </label>
          <textarea 
            v-model="userNotes" 
            class="textarea textarea-bordered h-24" 
            placeholder="Purchase date, price, condition notes..."
          ></textarea>
        </div>

        <div class="alert alert-info shadow-lg" v-if="loading">
          <div>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current flex-shrink-0 w-6 h-6 animate-spin"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>Processing images... (AI Analysis is currently disabled)</span>
          </div>
        </div>

        <div class="card-actions justify-end mt-6">
          <button type="submit" class="btn btn-primary w-full sm:w-auto" :disabled="loading || !frontFile || !backFile">
            Save Coin
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const fileInput = ref(null)
const frontFile = ref(null)
const backFile = ref(null)
const frontPreview = ref(null)
const backPreview = ref(null)
const groupName = ref('')
const userNotes = ref('')
const groups = ref([])
const loading = ref(false)
const isDragging = ref(false)
const error = ref(null)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

onMounted(async () => {
  // Load last used group
  const lastGroup = localStorage.getItem('last_group_name')
  if (lastGroup) {
    groupName.value = lastGroup
  }

  // Fetch groups
  try {
    const res = await axios.get(`${API_URL}/groups`)
    groups.value = res.data || []
  } catch (e) {
    console.error("Failed to load groups", e)
  }
})

const saveLastGroup = () => {
  if (groupName.value) {
    localStorage.setItem('last_group_name', groupName.value)
  }
}

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
      if (files.length > 0) error.value = "Please select image files only."
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
      error.value = "Please select exactly 2 images (Front & Back)."
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
      error.value = "Both Front and Back images are required."
      return
  }

  loading.value = true
  error.value = null
  
  // Save last group before submit
  saveLastGroup()

  const formData = new FormData()
  formData.append('front_image', frontFile.value)
  formData.append('back_image', backFile.value)
  if (groupName.value) {
    formData.append('group_name', groupName.value)
  }
  if (userNotes.value) {
    formData.append('user_notes', userNotes.value)
  }

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
    error.value = 'Error uploading coin: ' + (e.response?.data?.error || e.message)
  } finally {
    loading.value = false
  }
}
</script>
