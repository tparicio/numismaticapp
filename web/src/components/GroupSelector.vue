<template>
  <div class="form-control w-full">
    <label class="label">
      <span class="label-text">Group / Collection</span>
    </label>
    <div class="join w-full">
      <select 
        class="select select-bordered join-item w-full" 
        :value="modelValue" 
        @change="$emit('update:modelValue', $event.target.value)"
      >
        <option value="">None</option>
        <option value="__NEW__" class="font-bold text-primary">+ Create New Group...</option>
        <option disabled>──────────────</option>
        <option v-for="group in groups" :key="group.id" :value="group.name">
          {{ group.name }}
        </option>
      </select>
    </div>

    <!-- New Group Modal -->
    <dialog id="new_group_modal" class="modal" :class="{ 'modal-open': showModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Create New Group</h3>
        <p class="text-sm opacity-70 mt-2">Organize your coins into collections.</p>
        
        <div class="form-control w-full mt-4">
          <label class="label">
            <span class="label-text">Group Name</span>
          </label>
          <input 
            type="text" 
            v-model="newGroupName" 
            class="input input-bordered w-full" 
            ref="inputRef"
            @keyup.enter="createGroup"
            placeholder="e.g. Roman Empire, My 2024 Findings"
            :class="{ 'input-error': error }"
          />
          <label class="label" v-if="error">
            <span class="label-text-alt text-error">{{ error }}</span>
          </label>
        </div>

        <div class="form-control w-full mt-2">
            <label class="label"><span class="label-text">Description (Optional)</span></label>
            <textarea v-model="newGroupDesc" class="textarea textarea-bordered h-20" placeholder="Optional description..."></textarea>
        </div>

        <div class="modal-action">
          <button class="btn" @click="closeModal">Cancel</button>
          <button class="btn btn-primary" @click="createGroup" :disabled="creating">
            <span v-if="creating" class="loading loading-spinner"></span>
            Create Group
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeModal">close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import axios from 'axios'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const groups = ref([])
const showModal = ref(false)
const newGroupName = ref('')
const newGroupDesc = ref('')
const creating = ref(false)
const error = ref(null)
const inputRef = ref(null)

const fetchGroups = async () => {
  try {
    const res = await axios.get(`${API_URL}/groups`)
    console.log("GroupSelector: raw response", res.data)
    
    const rawGroups = res.data || []
    
    // Normalize and filter
    groups.value = rawGroups.map(g => ({
        id: g.id || g.ID,
        name: g.name || g.Name || 'Unnamed Group',
        description: g.description || g.Description
    })).filter(g => g.name && g.name.trim() !== '')
    
    console.log("GroupSelector: processed groups", groups.value)
  } catch (e) {
    console.error("Failed to fetch groups", e)
  }
}

onMounted(() => {
  fetchGroups()
})

watch(() => props.modelValue, (newVal) => {
  if (newVal === '__NEW__') {
    showModal.value = true
    emit('update:modelValue', '') // Reset selection while modal is open
    nextTick(() => {
        inputRef.value?.focus()
    })
  }
})

const closeModal = () => {
  showModal.value = false
  newGroupName.value = ''
  newGroupDesc.value = ''
  error.value = null
}

const createGroup = async () => {
  if (!newGroupName.value || !newGroupName.value.trim()) {
    error.value = "Group name is required"
    return
  }
  
  const name = newGroupName.value.trim()
  
  // Check local duplicate
  if (groups.value.some(g => g.name.toLowerCase() === name.toLowerCase())) {
      error.value = "Group already exists"
      return
  }

  creating.value = true
  error.value = null

  try {
    const res = await axios.post(`${API_URL}/groups`, {
        name: name,
        description: newGroupDesc.value
    })
    
    // Refresh list and select new group
    await fetchGroups()
    emit('update:modelValue', res.data.name)
    closeModal()
  } catch (e) {
    console.error("Failed to create group", e)
    error.value = e.response?.data?.error || "Failed to create group"
  } finally {
    creating.value = false
  }
}
</script>
