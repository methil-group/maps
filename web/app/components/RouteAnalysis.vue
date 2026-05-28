<template>
  <div v-if="isOpen" class="fixed inset-0 bg-slate-950/40 backdrop-blur-sm z-[2000] flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 rounded-2xl w-full max-w-4xl max-h-[85vh] overflow-hidden flex flex-col shadow-2xl animate-pop-in">
      
      <!-- Modal Header -->
      <div class="p-6 pb-4 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-lg bg-brand-violet-100 dark:bg-brand-violet-950 text-brand-violet-600 dark:text-brand-violet-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="20" x2="18" y2="10"></line>
              <line x1="12" y1="20" x2="12" y2="4"></line>
              <line x1="6" y1="20" x2="6" y2="14"></line>
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">
              {{ t('analysis.modal_title') }}
            </h3>
            <p class="text-xs text-slate-400 dark:text-slate-500 font-medium">
              {{ t('analysis.modal_desc') }}
            </p>
          </div>
        </div>
        <button @click="closeModal" class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto flex-1 space-y-6">
        
        <!-- Summary Cards Grid -->
        <div>
          <h4 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-3">
            {{ t('analysis.summary') }}
          </h4>
          <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.total_distance') }}</div>
              <div class="text-lg font-bold text-brand-violet-600 dark:text-brand-violet-400 mt-1 font-mono">
                {{ (totalDistance / 1000).toFixed(2) }} km
              </div>
            </div>
            
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.num_stops') }}</div>
              <div class="text-lg font-bold text-slate-900 dark:text-white mt-1">
                {{ locations.length }} {{ locations.length !== 1 ? t('common.location_plural') : t('common.location_singular') }}
              </div>
            </div>

            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.avg_distance') }}</div>
              <div class="text-lg font-bold text-slate-900 dark:text-white mt-1 font-mono">
                {{ (totalDistance / (locations.length > 1 ? locations.length : 1) / 1000).toFixed(2) }} km
              </div>
            </div>

            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.route_type') }}</div>
              <div class="text-lg font-bold text-slate-900 dark:text-white mt-1">
                {{ t('analysis.one_way') }}
              </div>
            </div>
          </div>
        </div>

        <!-- Step by Step Segment list -->
        <div class="p-5 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
          <h4 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-3 flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-brand-violet-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 18l6-6-6-6"></path>
            </svg>
            {{ t('analysis.step_by_step') }}
          </h4>
          <div class="max-h-60 overflow-y-auto space-y-2 pr-1">
            <div
              v-for="(segment, idx) in routeSegments"
              :key="idx"
              class="flex items-center gap-3 p-3 bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 rounded-lg"
            >
              <div class="w-6 h-6 rounded-full flex items-center justify-center bg-brand-violet-600 text-white text-[10px] font-extrabold shrink-0">
                {{ idx + 1 }}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-1.5 flex-wrap">
                  <span class="font-semibold text-slate-900 dark:text-white text-sm truncate max-w-[140px] md:max-w-[200px]">
                    {{ getShortName(segment.from.name) }}
                  </span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                    <polyline points="12 5 19 12 12 19"></polyline>
                  </svg>
                  <span class="font-semibold text-slate-900 dark:text-white text-sm truncate max-w-[140px] md:max-w-[200px]">
                    {{ getShortName(segment.to.name) }}
                  </span>
                </div>
                <div class="text-xs text-slate-400 dark:text-slate-500 mt-0.5 flex items-center gap-2">
                  <span>{{ t('analysis.distance') }}: <span class="font-bold text-slate-700 dark:text-slate-300 font-mono">{{ (segment.distance / 1000).toFixed(2) }} km</span></span>
                  <span>•</span>
                  <span>{{ t('analysis.running_total') }}: <span class="font-bold text-slate-700 dark:text-slate-300 font-mono">{{ (segment.cumulativeDistance / 1000).toFixed(2) }} km</span></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Adjacency Matrix -->
        <div class="p-5 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
          <h4 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-1.5">
            {{ t('analysis.matrix') }}
          </h4>
          <p class="text-xs text-slate-400 dark:text-slate-500 mb-4 leading-relaxed">
            {{ t('analysis.matrix_desc') }}
          </p>
          <div class="overflow-x-auto border border-slate-100 dark:border-slate-800 rounded-lg">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="bg-slate-100 dark:bg-slate-900 text-slate-500 dark:text-slate-400 border-b border-slate-100 dark:border-slate-800">
                  <th class="p-3 font-semibold w-28 shrink-0 bg-slate-100/50 dark:bg-slate-900/50 border-r border-slate-100 dark:border-slate-800">
                    {{ t('analysis.from_to') }}
                  </th>
                  <th v-for="loc in locations" :key="loc.id" class="p-3 font-semibold text-center whitespace-nowrap">
                    {{ getShortName(loc.name) }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="fromLoc in locations" :key="fromLoc.id" class="border-b border-slate-100 dark:border-slate-800 last:border-0 hover:bg-white/50 dark:hover:bg-slate-900/20">
                  <td class="p-3 font-bold bg-slate-100/20 dark:bg-slate-900/10 border-r border-slate-100 dark:border-slate-800 text-slate-900 dark:text-white max-w-28 truncate">
                    {{ getShortName(fromLoc.name) }}
                  </td>
                  <td
                    v-for="toLoc in locations"
                    :key="toLoc.id"
                    class="p-3 text-center font-mono text-slate-600 dark:text-slate-400"
                    :class="fromLoc.id === toLoc.id ? 'bg-slate-100/40 dark:bg-slate-900/40 text-slate-300 dark:text-slate-600' : ''"
                  >
                    <span v-if="fromLoc.id === toLoc.id">—</span>
                    <span v-else>{{ ((distanceMatrix?.get(fromLoc.id)?.get(toLoc.id) || 0) / 1000).toFixed(2) }} km</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Logic description -->
        <div class="p-5 bg-brand-violet-50/50 dark:bg-brand-violet-950/20 border border-brand-violet-100 dark:border-brand-violet-950/40 rounded-xl">
          <h4 class="text-sm font-semibold text-brand-violet-800 dark:text-brand-violet-400 mb-3 flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            {{ t('analysis.why_optimal') }}
          </h4>
          <div class="text-xs space-y-3 leading-relaxed text-slate-600 dark:text-slate-400">
            <div class="flex items-start gap-2">
              <div class="w-4 h-4 rounded-full bg-brand-violet-100 dark:bg-brand-violet-950 text-brand-violet-700 dark:text-brand-violet-400 flex items-center justify-center font-bold text-[9px] mt-0.5 shrink-0">
                1
              </div>
              <p>
                <span class="font-bold text-slate-800 dark:text-slate-300">{{ t('analysis.reason_1_title') }}:</span>
                {{ t('analysis.reason_1_desc') }}
              </p>
            </div>
            
            <div class="flex items-start gap-2">
              <div class="w-4 h-4 rounded-full bg-brand-violet-100 dark:bg-brand-violet-950 text-brand-violet-700 dark:text-brand-violet-400 flex items-center justify-center font-bold text-[9px] mt-0.5 shrink-0">
                2
              </div>
              <p>
                <span class="font-bold text-slate-800 dark:text-slate-300">{{ t('analysis.reason_2_title') }}:</span>
                {{ t('analysis.reason_2_desc', { count: locations.length > 0 ? getFactorialVal().toLocaleString() : 0 }) }}
              </p>
            </div>

            <div class="flex items-start gap-2">
              <div class="w-4 h-4 rounded-full bg-brand-violet-100 dark:bg-brand-violet-950 text-brand-violet-700 dark:text-brand-violet-400 flex items-center justify-center font-bold text-[9px] mt-0.5 shrink-0">
                3
              </div>
              <p>
                <span class="font-bold text-slate-800 dark:text-slate-300">{{ t('analysis.reason_3_title') }}:</span>
                {{ t('analysis.reason_3_desc') }}
              </p>
            </div>
          </div>
        </div>

      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTrans } from '../composables/useTrans'
import type { Location } from '../utils/routeOptimizer'

const props = defineProps<{
  isOpen: boolean
  locations: Location[]
  totalDistance: number
  distanceMatrix?: Map<string, Map<string, number>>
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { t } = useTrans()

const closeModal = () => {
  emit('close')
}

const getShortName = (name: string) => {
  return name.split(',')[0].trim()
}

// Compute step by step details
const routeSegments = computed(() => {
  const segments = []
  let cumulativeDistance = 0
  
  for (let i = 0; i < props.locations.length - 1; i++) {
    const from = props.locations[i]
    const to = props.locations[i + 1]
    const distance = props.distanceMatrix?.get(from.id)?.get(to.id) || 0
    cumulativeDistance += distance
    
    segments.push({
      from,
      to,
      distance,
      cumulativeDistance
    })
  }
  return segments
})

const getFactorialVal = () => {
  const n = props.locations.length
  if (n <= 1) return 1
  return factorial(n - 1)
}

function factorial(num: number): number {
  if (num <= 1) return 1
  return num * factorial(num - 1)
}
</script>
