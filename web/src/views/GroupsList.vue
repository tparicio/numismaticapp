<template>
  <div>
    <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
      <h2 class="text-3xl font-bold">{{ $t('groups.title') }}</h2>
      <button class="btn btn-primary" @click="openModal()">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 mr-2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        {{ $t('groups.add_button') }}
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center p-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Empty State -->
    <div v-else-if="groups.length === 0" class="text-center p-10 bg-base-100 rounded-box shadow-xl">
      <p class="text-lg opacity-70 mb-4">{{ $t('groups.empty_state') }}</p>
      <button class="btn btn-primary btn-outline" @click="openModal()">{{ $t('groups.add_button') }}</button>
    </div>

    <!-- Groups Table -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-box shadow-xl">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th class="w-16">ID</th>
            <th>{{ $t('groups.name') }}</th>
            <th>{{ $t('groups.description') }}</th>
            <th class="w-24">{{ $t('groups.count') }}</th>
            <th class="text-right">{{ $t('list.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="group in groups" :key="group.id" class="hover">
            <td class="opacity-50 font-mono text-xs">{{ group.id }}</td>
            <td class="font-bold text-lg">{{ group.name }}</td>
            <td class="opacity-70">{{ group.description || '-' }}</td>
            <td>
                <span class="badge badge-secondary badge-outline" v-if="group.coin_count">{{ group.coin_count }}</span>
                <span class="opacity-30" v-else>-</span>
            </td>
            <td class="text-right">
              <div class="flex justify-end gap-2">
                <router-link :to="`/list?group_id=${group.id}`" class="btn btn-square btn-sm btn-ghost text-primary" :title="$t('groups.view_collection')">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </router-link>
                <button class="btn btn-square btn-sm btn-ghost text-info" @click="openModal(group)" :title="$t('common.edit')">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                  </svg>
                </button>
                <button class="btn btn-square btn-sm btn-ghost text-secondary" @click="openImagesModal(group)" :title="$t('groups.images')">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                  </svg>
                </button>
                <button class="btn btn-square btn-sm btn-ghost text-error" @click="confirmDelete(group)" :title="$t('common.delete')">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <dialog id="group_modal" class="modal" :class="{ 'modal-open': showModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ isEditing ? $t('groups.edit_title') : $t('groups.add_title') }}</h3>
        
        <div class="form-control w-full mb-4">
          <label class="label">
            <span class="label-text font-semibold">{{ $t('groups.name') }}</span>
          </label>
          <input type="text" v-model="form.name" class="input input-bordered w-full" :class="{'input-error': errors.name}" placeholder="My Collection" />
          <label class="label" v-if="errors.name">
            <span class="label-text-alt text-error">{{ errors.name }}</span>
          </label>
        </div>

        <div class="form-control w-full mb-6">
          <label class="label">
            <span class="label-text font-semibold">{{ $t('groups.description') }}</span>
          </label>
          <textarea v-model="form.description" class="textarea textarea-bordered h-24" placeholder="Optional description..."></textarea>
        </div>

        <div class="modal-action">
          <button class="btn" @click="closeModal">{{ $t('common.cancel') }}</button>
          <button class="btn btn-primary" @click="saveGroup" :disabled="saving">
            <span v-if="saving" class="loading loading-spinner"></span>
            {{ isEditing ? $t('common.update') : $t('common.save') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeModal">close</button>
      </form>
    </dialog>

    <!-- Delete Confirmation Modal -->
    <dialog id="delete_group_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-error">{{ $t('groups.delete_modal.title') }}</h3>
        <p class="py-4">
          {{ $t('groups.delete_modal.confirm') }} <span class="font-bold">{{ groupToDelete?.name }}</span>?
          <br/>
          <span class="text-sm opacity-70">{{ $t('groups.delete_modal.warning') }}</span>
        </p>
        <div class="modal-action">
          <button class="btn" @click="deleteModalOpen = false">{{ $t('common.cancel') }}</button>
          <button class="btn btn-error" @click="deleteGroup" :disabled="deleting">
            <span v-if="deleting" class="loading loading-spinner"></span>
            {{ $t('common.delete') }}
          </button>
        </div>
      </div>
    </dialog>

    <!-- Images Modal Component -->
    <GroupImagesModal 
        :is-open="imagesModalOpen" 
        :group="currentGroupImages" 
        @close="closeImagesModal" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import ImageViewer from '../components/ImageViewer.vue'
import GroupImagesModal from '../components/GroupImagesModal.vue'

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

const groups = ref([])
const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)

// Wrapper for error handling
const handleError = (e) => {
    console.error(e)
    alert(e.response?.data?.error || "An error occurred")
}

const fetchGroups = async () => {
    loading.value = true
    try {
        const res = await axios.get(`${API_URL}/groups`)
        groups.value = res.data
    } catch (e) {
        handleError(e)
    } finally {
        loading.value = false
    }
}

// Modal State
const showModal = ref(false)
const editingGroup = ref(null)
const isEditing = computed(() => !!editingGroup.value)
const form = ref({
    name: '',
    description: ''
})
const errors = ref({})

const openModal = (group = null) => {
    errors.value = {}
    if (group) {
        editingGroup.value = group
        form.value = { 
            name: group.name, 
            description: group.description || '' 
        }
    } else {
        editingGroup.value = null
        form.value = { name: '', description: '' }
    }
    showModal.value = true
}

const closeModal = () => {
    showModal.value = false
    editingGroup.value = null
}

const validate = () => {
    errors.value = {}
    if (!form.value.name || form.value.name.trim().length < 3) {
        errors.value.name = "Name must be at least 3 characters"
        return false
    }
    return true
}

const saveGroup = async () => {
    if (!validate()) return

    saving.value = true
    try {
        if (isEditing.value) {
            const res = await axios.put(`${API_URL}/groups/${editingGroup.value.id}`, form.value)
            // Update in list
            const index = groups.value.findIndex(g => g.id === editingGroup.value.id)
            if (index !== -1) {
                groups.value[index] = res.data
            }
        } else {
            const res = await axios.post(`${API_URL}/groups`, form.value)
            groups.value.push(res.data)
        }
        closeModal()
    } catch (e) {
        handleError(e)
    } finally {
        saving.value = false
    }
}

// Delete State
const deleteModalOpen = ref(false)
const groupToDelete = ref(null)

const confirmDelete = (group) => {
    groupToDelete.value = group
    deleteModalOpen.value = true
}

const deleteGroup = async () => {
    if (!groupToDelete.value) return
    deleting.value = true
    try {
        await axios.delete(`${API_URL}/groups/${groupToDelete.value.id}`)
        groups.value = groups.value.filter(g => g.id !== groupToDelete.value.id)
        deleteModalOpen.value = false
        groupToDelete.value = null
    } catch (e) {
        handleError(e)
    } finally {
        deleting.value = false
    }
}

// Images Logic
const imagesModalOpen = ref(false)
const currentGroupImages = ref(null) // The group object being viewed

const getImageUrl = (path) => {
    if (!path) return '/broken_coin.png'
    if (path.includes('storage/')) {
        return `${API_URL.replace('/api/v1', '')}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const openImagesModal = async (group) => {
    currentGroupImages.value = group
    imagesModalOpen.value = true
}

const closeImagesModal = () => {
    imagesModalOpen.value = false
    currentGroupImages.value = null
}

onMounted(() => {
    fetchGroups()
})
</script>
