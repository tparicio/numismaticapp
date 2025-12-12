<template>
  <div class="container mx-auto p-4 max-w-2xl">
    <div class="flex items-center gap-4 mb-6">
      <router-link to="/list" class="btn btn-circle btn-ghost">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
        </svg>
      </router-link>
      <h1 class="text-3xl font-bold">{{ $t('form.edit_title') }}</h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <form v-else @submit.prevent="submit" class="space-y-6">
      <!-- Images (Read-only for MVP) -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title text-sm uppercase text-base-content/70 mb-4">{{ $t('form.sections.images') }}</h2>
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

      <!-- Basic Details -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body space-y-4">
          <h2 class="card-title text-sm uppercase text-base-content/70">{{ $t('form.sections.basic') }}</h2>
          
          <div class="form-control w-full">
            <label class="label"><span class="label-text">{{ $t('form.fields.name') }}</span></label>
            <input v-model="form.name" type="text" class="input input-bordered w-full" required />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.country') }}</span></label>
              <input v-model="form.country" type="text" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.year') }}</span></label>
              <input v-model.number="form.year" type="number" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.face_value') }}</span></label>
              <input v-model="form.face_value" type="text" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.currency') }}</span></label>
              <input v-model="form.currency" type="text" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.material') }}</span></label>
              <input v-model="form.material" type="text" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.grade') }}</span></label>
              <input v-model="form.grade" type="text" class="input input-bordered w-full" />
            </div>
          </div>

           <div class="form-control w-full">
            <label class="label"><span class="label-text">{{ $t('form.fields.description') }}</span></label>
            <textarea v-model="form.description" class="textarea textarea-bordered h-24"></textarea>
          </div>
        </div>
      </div>

      <!-- Technical Details -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body space-y-4">
          <h2 class="card-title text-sm uppercase text-base-content/70">{{ $t('form.sections.technical') }}</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.mint') }}</span></label>
              <input v-model="form.mint" type="text" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.mintage') }}</span></label>
              <input v-model.number="form.mintage" type="number" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.km_code') }}</span></label>
              <input v-model="form.km_code" type="text" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.weight') }}</span></label>
              <input v-model.number="form.weight_g" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.diameter') }}</span></label>
              <input v-model.number="form.diameter_mm" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.thickness') }}</span></label>
              <input v-model.number="form.thickness_mm" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.edge') }}</span></label>
              <input v-model="form.edge" type="text" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.shape') }}</span></label>
              <input v-model="form.shape" type="text" class="input input-bordered w-full" />
            </div>
          </div>

           <div class="form-control w-full">
            <label class="label"><span class="label-text">{{ $t('form.fields.technical_notes') }}</span></label>
            <textarea v-model="form.technical_notes" class="textarea textarea-bordered h-24"></textarea>
          </div>
        </div>
      </div>

       <!-- Valuation & Acquisition -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body space-y-4">
          <h2 class="card-title text-sm uppercase text-base-content/70">{{ $t('form.sections.valuation') }}</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.min_value') }}</span></label>
              <input v-model.number="form.min_value" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
            <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.max_value') }}</span></label>
              <input v-model.number="form.max_value" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.price_paid') }}</span></label>
              <input v-model.number="form.price_paid" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.sold_price') }}</span></label>
              <input v-model.number="form.sold_price" type="number" step="0.01" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.acquired_at') }}</span></label>
              <input v-model="form.acquired_at" type="date" class="input input-bordered w-full" />
            </div>
             <div class="form-control w-full">
              <label class="label"><span class="label-text">{{ $t('form.fields.sold_at') }}</span></label>
              <input v-model="form.sold_at" type="date" class="input input-bordered w-full" />
            </div>
          </div>
        </div>
      </div>

      <!-- Organization -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body space-y-4">
          <h2 class="card-title text-sm uppercase text-base-content/70">{{ $t('form.sections.organization') }}</h2>
          
          <div class="form-control w-full">
            <GroupSelector v-model="form.group_name" />
          </div>

          <div class="form-control w-full">
            <label class="label"><span class="label-text">{{ $t('form.fields.personal_notes') }}</span></label>
            <textarea v-model="form.personal_notes" class="textarea textarea-bordered h-24"></textarea>
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <router-link to="/list" class="btn btn-ghost">{{ $t('common.cancel') }}</router-link>
        <button type="submit" class="btn btn-primary" :disabled="saving">
          <span v-if="saving" class="loading loading-spinner"></span>
          {{ $t('common.save') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import GroupSelector from '../components/GroupSelector.vue'

const route = useRoute()
const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

const loading = ref(true)
const saving = ref(false)
const coin = ref(null)

const form = ref({
  name: '',
  country: '',
  year: 0,
  face_value: '',
  currency: '',
  material: '',
  grade: '',
  description: '',
  mint: '',
  mintage: 0,
  km_code: '',
  weight_g: 0,
  diameter_mm: 0,
  thickness_mm: 0,
  edge: '',
  shape: '',
  technical_notes: '',
  min_value: 0,
  max_value: 0,
  price_paid: 0,
  sold_price: 0,
  acquired_at: '',
  sold_at: '',
  group_name: '',
  personal_notes: ''
})

const formatDateForInput = (dateString) => {
    if (!dateString) return ''
    return new Date(dateString).toISOString().split('T')[0]
}

const parseDateForSubmit = (dateString) => {
    if (!dateString) return null
    return new Date(dateString).toISOString()
}

onMounted(async () => {
  try {
    const coinRes = await axios.get(`${API_URL}/coins/${route.params.id}`)
    
    coin.value = coinRes.data

    // Populate form
    const c = coin.value
    form.value = {
        name: c.name || '',
        country: c.country || '',
        year: c.year || 0,
        face_value: c.face_value || '',
        currency: c.currency || '',
        material: c.material || '',
        grade: c.grade || '',
        description: c.description || '',
        mint: c.mint || '',
        mintage: c.mintage || 0,
        km_code: c.km_code || '',
        weight_g: c.weight_g || 0,
        diameter_mm: c.diameter_mm || 0,
        thickness_mm: c.thickness_mm || 0,
        edge: c.edge || '',
        shape: c.shape || '',
        technical_notes: c.technical_notes || '',
        min_value: c.min_value || 0,
        max_value: c.max_value || 0,
        price_paid: c.price_paid || 0,
        sold_price: c.sold_price || 0,
        acquired_at: formatDateForInput(c.acquired_at),
        sold_at: formatDateForInput(c.sold_at),
        personal_notes: c.personal_notes || '',
        group_name: '' // Will be populated below if exists
    }
    
    // Fetch group details if assigned
     if (c.group_id) {
         try {
             // We need to fetch groups to map ID to name, or just use the separate Groups endpoint
             // To simplify, GroupSelector handles its own fetching. We just need to set the name.
             // But we only have ID here.
             // Let's fetch the specific group or all groups
             const groupsRes = await axios.get(`${API_URL}/groups`)
             const groups = groupsRes.data || []
             const g = groups.find(g => g.id === c.group_id)
             if (g) form.value.group_name = g.name
         } catch (e) {
             console.error("Failed to fetch groups for mapping", e)
         }
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
    const payload = {
        ...form.value,
        acquired_at: parseDateForSubmit(form.value.acquired_at),
        sold_at: parseDateForSubmit(form.value.sold_at)
    }
    await axios.put(`${API_URL}/coins/${route.params.id}`, payload)
    router.push(`/coin/${route.params.id}`)
  } catch (e) {
    console.error("Failed to update coin", e)
    alert("Failed to update coin")
  } finally {
    saving.value = false
  }
}
</script>
