<template>
  <div class="container mx-auto p-4 max-w-2xl">
    <div class="flex items-center gap-4 mb-6">
      <router-link to="/list" class="btn btn-circle btn-ghost">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
        </svg>
      </router-link>
      <h1 class="text-3xl font-bold">Edit Coin</h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <form v-else @submit.prevent="submit" class="space-y-6">
      <!-- Images (Read-only for MVP) -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title text-sm uppercase text-base-content/70 mb-4">Images (Cannot be changed)</h2>
          <div class="flex gap-4 justify-center">
            <div class="avatar">
              <div class="w-24 rounded-xl">
                <img :src="getThumbnail('front')" />
              </div>
            </div>
            <div class="avatar">
              <div class="w-24 rounded-xl">
                <img :src="getThumbnail('back')" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Details -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body space-y-4">
          <h2 class="card-title text-sm uppercase text-base-content/70">Details</h2>
          
          <div class="form-control w-full">
            <label class="label"><span class="label-text">Name</span></label>
            <input v-model="form.name" type="text" class="input input-bordered w-full" required />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-control w-full">
              <label class="label"><span class="label-text">Mint</span></label>
              <input v-model="form.mint" type="text" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">Mintage</span></label>
              <input v-model.number="form.mintage" type="number" class="input input-bordered w-full" />
            </div>
          </div>

          <div class="form-control w-full">
            <label class="label"><span class="label-text">Group</span></label>
            <select v-model="form.group_name" class="select select-bordered w-full">
              <option value="">None</option>
              <option v-for="g in groups" :key="g.id" :value="g.name">{{ g.name }}</option>
            </select>
          </div>

          <div class="form-control w-full">
            <label class="label"><span class="label-text">Personal Notes</span></label>
            <textarea v-model="form.user_notes" class="textarea textarea-bordered h-24"></textarea>
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <router-link to="/list" class="btn btn-ghost">Cancel</router-link>
        <button type="submit" class="btn btn-primary" :disabled="saving">
          <span v-if="saving" class="loading loading-spinner"></span>
          Save Changes
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const loading = ref(true)
const saving = ref(false)
const coin = ref(null)
const groups = ref([])

const form = ref({
  name: '',
  mint: '',
  mintage: 0,
  group_name: '',
  user_notes: ''
})

onMounted(async () => {
  try {
    const [coinRes, groupsRes] = await Promise.all([
      axios.get(`${API_URL}/coins/${route.params.id}`),
      axios.get(`${API_URL}/groups`)
    ])
    
    coin.value = coinRes.data
    groups.value = groupsRes.data

    // Populate form
    form.value.name = coin.value.name || ''
    form.value.mint = coin.value.mint || ''
    form.value.mintage = coin.value.mintage || 0
    form.value.user_notes = coin.value.personal_notes || ''
    
    if (coin.value.group_id) {
      const g = groups.value.find(g => g.id === coin.value.group_id)
      if (g) form.value.group_name = g.name
    }

  } catch (e) {
    console.error("Failed to load data", e)
    alert("Failed to load coin details")
    router.push('/list')
  } finally {
    loading.value = false
  }
})

const getThumbnail = (side) => {
  if (!coin.value || !coin.value.images) return 'https://placehold.co/100x100?text=No+Image'
  const img = coin.value.images.find(i => i.image_type === 'thumbnail' && i.side === side)
  if (img) return `${API_URL.replace('/api/v1', '')}/${img.path}`
  return 'https://placehold.co/100x100?text=No+Image'
}

const submit = async () => {
  saving.value = true
  try {
    await axios.put(`${API_URL}/coins/${route.params.id}`, form.value)
    router.push(`/coin/${route.params.id}`)
  } catch (e) {
    console.error("Failed to update coin", e)
    alert("Failed to update coin")
  } finally {
    saving.value = false
  }
}
</script>
