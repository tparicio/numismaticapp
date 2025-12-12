<template>
  <div class="collapse collapse-arrow border border-base-300 bg-base-100 rounded-box">
    <input type="checkbox" :checked="isOpen" @change="$emit('update:isOpen', $event.target.checked)" /> 
    <div class="collapse-title text-xl font-medium flex items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-primary">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
      </svg>
      {{ title || $t('form.ai_settings') || 'Configuración IA' }}
    </div>
    <div class="collapse-content space-y-4">
       <!-- Model Selector -->
       <div class="form-control w-full">
         <label class="label">
           <span class="label-text">Gemini Model</span>
         </label>
         <select :value="model" @input="$emit('update:model', $event.target.value)" class="select select-bordered w-full">
           <option v-for="m in availableModels" :key="m.name" :value="m.name">
             {{ m.name }}
           </option>
         </select>
         <label class="label" v-if="selectedModelDescription">
           <span class="label-text-alt text-base-content/70">{{ selectedModelDescription }}</span>
         </label>
       </div>

       <!-- Temperature Slider -->
       <div class="form-control w-full">
         <label class="label">
           <span class="label-text">Creatividad (Temperatura)</span>
           <span class="label-text-alt">{{ temperature }}</span>
         </label>
         <input type="range" min="0" max="1" step="0.1" :value="temperature" @input="$emit('update:temperature', parseFloat($event.target.value))" class="range range-primary range-xs" />
         <div class="w-full flex justify-between text-xs px-2 mt-1">
           <span>Preciso</span>
           <span>Creativo</span>
         </div>
       </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  model: String,
  temperature: Number,
  availableModels: {
    type: Array,
    default: () => []
  },
  isOpen: {
    type: Boolean,
    default: false
  },
  title: String
})

defineEmits(['update:model', 'update:temperature', 'update:isOpen'])

const selectedModelDescription = computed(() => {
  const m = props.availableModels.find(m => m.name === props.model)
  return m ? m.description : ''
})
</script>
