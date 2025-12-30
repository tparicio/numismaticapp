<template>
  <div v-if="isOpen" class="fixed inset-0 z-[9999] flex items-center justify-center bg-black bg-opacity-90" @click.self="close">
    <button class="absolute top-4 right-4 btn btn-circle btn-ghost text-white z-[9999]" @click="close">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
    </button>

    <div class="relative w-full h-full overflow-hidden flex items-center justify-center" 
         @wheel.prevent="handleWheel"
         @mousedown="startDrag"
         @mousemove="onDrag"
         @mouseup="stopDrag"
         @mouseleave="stopDrag"
         @touchstart="startTouch"
         @touchmove="onTouch"
         @touchend="stopDrag">
      
      <img 
        :src="imageUrl" 
        class="max-h-[90vh] max-w-[90vw] object-contain transition-transform duration-75 ease-linear cursor-move"
        :style="{ transform: `translate(${position.x}px, ${position.y}px) scale(${scale})` }"
        alt="Full size coin"
        draggable="false"
      />

      <div class="absolute bottom-8 left-1/2 transform -translate-x-1/2 flex gap-2 bg-black bg-opacity-50 p-2 rounded-full">
        <button class="btn btn-sm btn-circle btn-ghost text-white" @click="zoomOut">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 10a1 1 0 011-1h8a1 1 0 110 2H6a1 1 0 01-1-1z" clip-rule="evenodd" />
            </svg>
        </button>
        <span class="text-white self-center text-sm w-12 text-center">{{ Math.round(scale * 100) }}%</span>
        <button class="btn btn-sm btn-circle btn-ghost text-white" @click="zoomIn">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
        </button>
        <button class="btn btn-sm btn-circle btn-ghost text-white ml-2" @click="reset">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  isOpen: Boolean,
  imageUrl: String
})

const emit = defineEmits(['close'])

const scale = ref(1)
const position = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const startPos = ref({ x: 0, y: 0 })

const close = () => {
  emit('close')
  reset()
}

const reset = () => {
  scale.value = 1
  position.value = { x: 0, y: 0 }
}

const handleWheel = (e) => {
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  const newScale = Math.max(0.1, Math.min(5, scale.value + delta))
  scale.value = newScale
}

const zoomIn = () => {
  scale.value = Math.min(5, scale.value + 0.25)
}

const zoomOut = () => {
  scale.value = Math.max(0.1, scale.value - 0.25)
}

const startDrag = (e) => {
  isDragging.value = true
  startPos.value = { x: e.clientX - position.value.x, y: e.clientY - position.value.y }
}

const onDrag = (e) => {
  if (!isDragging.value) return
  position.value = {
    x: e.clientX - startPos.value.x,
    y: e.clientY - startPos.value.y
  }
}

const stopDrag = () => {
  isDragging.value = false
}

// Touch support
const startTouch = (e) => {
    if (e.touches.length === 1) {
        isDragging.value = true
        startPos.value = { x: e.touches[0].clientX - position.value.x, y: e.touches[0].clientY - position.value.y }
    }
}

const onTouch = (e) => {
    if (!isDragging.value || e.touches.length !== 1) return
    position.value = {
        x: e.touches[0].clientX - startPos.value.x,
        y: e.touches[0].clientY - startPos.value.y
    }
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>
