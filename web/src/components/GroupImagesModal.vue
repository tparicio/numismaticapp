<template>
  <div v-if="isOpen">
    <!-- Modal Dialog -->
    <dialog class="modal modal-open">
      <div class="modal-box w-11/12 max-w-4xl">
        <h3 class="font-bold text-lg mb-4 flex items-center justify-between">
            <span>{{ $t('groups.images_title') }}: {{ group?.name }}</span>
            <button class="btn btn-circle btn-sm btn-ghost" @click="close">✕</button>
        </h3>
        
        <div class="flex justify-end mb-4">
             <input type="file" ref="groupImageInput" class="hidden" @change="uploadGroupImage" accept="image/*">
             <button class="btn btn-primary btn-sm gap-2" @click="$refs.groupImageInput.click()" :disabled="uploadingImage">
                <span v-if="uploadingImage" class="loading loading-spinner loading-xs"></span>
                {{ $t('common.upload') }}
            </button>
        </div>

        <div v-if="loading" class="flex justify-center p-10">
             <span class="loading loading-spinner loading-lg"></span>
        </div>
        <div v-else-if="!images || images.length === 0" class="text-center py-10 bg-base-200 rounded-lg border border-dashed border-base-300">
             <p class="opacity-50 font-bold">{{ $t('groups.images_empty') }}</p>
        </div>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
             <div v-for="img in images" :key="img.id" class="relative group aspect-square bg-base-100 rounded-lg overflow-hidden shadow border border-base-200 cursor-zoom-in">
                <img :src="getImageUrl(img.path)" class="w-full h-full object-cover transition-transform group-hover:scale-105" @click="openViewer(img.path)">
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-start justify-end p-2 pointer-events-none">
                    <button @click.stop="confirmDelete(img.id)" class="btn btn-xs btn-circle btn-error pointer-events-auto" :class="{'loading': deletingImageId === img.id}">
                        <svg v-if="deletingImageId !== img.id" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
             </div>
        </div>

        <div class="modal-action">
          <button class="btn" @click="close">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="close">close</button>
      </form>
    </dialog>
    
    <!-- Delete Confirmation Modal -->
    <dialog class="modal" :class="{'modal-open': deleteConfirmationId}">
        <div class="modal-box">
            <h3 class="font-bold text-lg">{{ $t('common.delete') }}</h3>
            <p class="py-4">{{ $t('common.delete_modal.confirm') }}?</p>
            <div class="modal-action">
                <button class="btn" @click="deleteConfirmationId = null">{{ $t('common.cancel') }}</button>
                <button class="btn btn-error" @click="performDelete">{{ $t('common.delete') }}</button>
            </div>
        </div>
    </dialog>

    <ImageViewer :is-open="viewerOpen" :image-url="viewerImage" @close="viewerOpen = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import axios from 'axios'
import ImageViewer from './ImageViewer.vue'

const props = defineProps({
  isOpen: Boolean,
  group: Object
})

const emit = defineEmits(['close'])

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

const images = ref([])
const loading = ref(false)
const uploadingImage = ref(false)
const deletingImageId = ref(null)
const groupImageInput = ref(null)

// Viewer Logic
const viewerOpen = ref(false)
const viewerImage = ref('')

const getImageUrl = (path) => {
    if (!path) return '/broken_coin.png'
    if (path.includes('storage/')) {
        return `${API_URL.replace('/api/v1', '')}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const fetchGroupImages = async () => {
    if (!props.group) return
    loading.value = true
    try {
        const res = await axios.get(`${API_URL}/groups/${props.group.id}/images`)
        images.value = res.data || []
    } catch (e) {
        console.error("Failed to fetch group images", e)
    } finally {
        loading.value = false
    }
}

const uploadGroupImage = async (event) => {
    const file = event.target.files[0]
    if (!file || !props.group) return

    uploadingImage.value = true
    const formData = new FormData()
    formData.append('image', file)

    try {
        await axios.post(`${API_URL}/groups/${props.group.id}/images`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
        await fetchGroupImages()
    } catch (e) {
        console.error("Upload failed", e)
        alert('Failed to upload image')
    } finally {
        uploadingImage.value = false
        event.target.value = ''
    }
}

const deleteConfirmationId = ref(null)

const confirmDelete = (id) => {
    deleteConfirmationId.value = id
}

const performDelete = async () => {
    if (!deleteConfirmationId.value) return
    const imageId = deleteConfirmationId.value
    deleteConfirmationId.value = null // Close modal
    
    deletingImageId.value = imageId
    try {
        await axios.delete(`${API_URL}/groups/${props.group.id}/images/${imageId}`)
        images.value = images.value.filter(img => img.id !== imageId)
    } catch (e) {
        console.error("Delete failed", e)
        alert('Failed to delete image')
    } finally {
        deletingImageId.value = null
    }
}

const openViewer = (path) => {
    viewerImage.value = getImageUrl(path)
    viewerOpen.value = true
}

const close = () => {
    emit('close')
}

watch(() => props.isOpen, (newVal) => {
    if (newVal && props.group) {
        fetchGroupImages()
    }
})
</script>
