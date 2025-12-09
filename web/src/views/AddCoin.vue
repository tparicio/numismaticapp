<template>
  <div class="card bg-base-100 shadow-xl max-w-2xl mx-auto">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">Add New Coin</h2>
      
      <form @submit.prevent="uploadCoin" class="space-y-6">
        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Front Image (Obverse)</span>
          </label>
          <input type="file" @change="handleFrontChange" class="file-input file-input-bordered w-full" accept="image/*" required />
        </div>

        <div class="form-control w-full">
          <label class="label">
            <span class="label-text">Back Image (Reverse)</span>
          </label>
          <input type="file" @change="handleBackChange" class="file-input file-input-bordered w-full" accept="image/*" required />
        </div>

        <div class="alert alert-info shadow-lg" v-if="loading">
          <div>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current flex-shrink-0 w-6 h-6 animate-spin"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>Processing images and analyzing with AI... This may take a few seconds.</span>
          </div>
        </div>

        <div class="card-actions justify-end mt-6">
          <button type="submit" class="btn btn-primary" :disabled="loading || !frontFile || !backFile">
            Analyze & Save
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const frontFile = ref(null)
const backFile = ref(null)
const loading = ref(false)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const handleFrontChange = (e) => {
  frontFile.value = e.target.files[0]
}

const handleBackChange = (e) => {
  backFile.value = e.target.files[0]
}

const uploadCoin = async () => {
  if (!frontFile.value || !backFile.value) return

  loading.value = true
  const formData = new FormData()
  formData.append('front_image', frontFile.value)
  formData.append('back_image', backFile.value)

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
    alert('Error uploading coin: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}
</script>
