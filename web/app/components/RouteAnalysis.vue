<template>
  <div v-if="isOpen" @click.self="closeModal" class="fixed inset-0 bg-slate-950/40 backdrop-blur-sm z-[2000] flex items-center justify-center p-4">
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
            <!-- Distance -->
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.total_distance') }}</div>
              <div class="text-lg font-bold text-brand-violet-600 dark:text-brand-violet-400 mt-1 font-mono">
                {{ (totalDistance / 1000).toFixed(1) }} km
              </div>
            </div>
            
            <!-- Duration -->
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('planner.duration') }}</div>
              <div class="text-lg font-bold text-slate-900 dark:text-white mt-1 font-mono">
                {{ formatDuration(totalDuration) }}
              </div>
            </div>

            <!-- Cost -->
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.total_cost') }}</div>
              <div class="text-lg font-bold text-slate-900 dark:text-white mt-1 font-mono">
                {{ (totalFuelCost + totalTollCost).toFixed(2) }} €
              </div>
            </div>

            <!-- CO2 Footprint -->
            <div class="p-4 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl group relative">
              <div class="text-xs text-slate-400 dark:text-slate-500 font-medium flex items-center gap-1">
                <span>{{ t('analysis.co2_emissions') }}</span>
                <!-- Info Trigger -->
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <path d="M12 16v-4"></path>
                  <path d="M12 8h.01"></path>
                </svg>
              </div>
              <div class="text-lg font-bold text-emerald-600 dark:text-emerald-400 mt-1 font-mono">
                {{ estimatedCO2.toFixed(1) }} kg
              </div>
              <!-- Tooltip -->
              <div class="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 w-52 p-2 bg-slate-950/95 backdrop-blur text-white text-[9px] rounded-lg shadow-lg border border-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none z-50 leading-relaxed">
                {{ t('analysis.co2_calculation_info') }}
              </div>
            </div>
          </div>
        </div>

        <!-- Budget Breakdown and Diagnostics -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <!-- Budget Breakdown -->
          <div class="p-5 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl space-y-4">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
              {{ t('analysis.cost_breakdown') }}
            </h4>
            
            <div class="space-y-3">
              <!-- Fuel Share -->
              <div class="space-y-1">
                <div class="flex justify-between text-xs font-semibold">
                  <span class="text-slate-600 dark:text-slate-400">{{ t('analysis.fuel_share') }}</span>
                  <span class="font-mono text-slate-900 dark:text-white">{{ totalFuelCost.toFixed(2) }} € ({{ fuelSharePercent }}%)</span>
                </div>
                <div class="h-2 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden">
                  <div class="h-full bg-brand-violet-600 rounded-full" :style="{ width: fuelSharePercent + '%' }"></div>
                </div>
              </div>

              <!-- Toll Share -->
              <div class="space-y-1">
                <div class="flex justify-between text-xs font-semibold">
                  <span class="text-slate-600 dark:text-slate-400">{{ t('analysis.toll_share') }}</span>
                  <span class="font-mono text-slate-900 dark:text-white">{{ totalTollCost.toFixed(2) }} € ({{ tollSharePercent }}%)</span>
                </div>
                <div class="h-2 w-full bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden">
                  <div class="h-full bg-brand-pink-500 rounded-full" :style="{ width: tollSharePercent + '%' }"></div>
                </div>
              </div>
            </div>
            
            <!-- Pricing profiles details -->
            <div class="text-[10px] text-slate-400 dark:text-slate-500 leading-snug border-t border-slate-100 dark:border-slate-800/60 pt-3">
              {{ t('analysis.pricing_details', { consumption: fuelConsumption.toFixed(1), fuelType: fuelType, fuelPrice: fuelPrice.toFixed(2), tollPrice: tollPrice.toFixed(2) }) }}
            </div>
          </div>

          <!-- Optimization Diagnosis -->
          <div class="p-5 bg-brand-violet-50/20 dark:bg-brand-violet-950/10 border border-brand-violet-100/50 dark:border-brand-violet-950/30 rounded-xl space-y-3">
            <h4 class="text-sm font-semibold text-brand-violet-800 dark:text-brand-violet-400 flex items-center gap-1.5">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              {{ t('analysis.insights_title') }}
            </h4>
            <p class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
              {{ diagnosticText }}
            </p>
            <div class="text-[10px] text-slate-400 dark:text-slate-500 pt-2 border-t border-brand-violet-100/20">
              {{ t('analysis.active_criterion') }} <span class="px-2 py-0.5 rounded bg-brand-violet-100 dark:bg-brand-violet-950/60 text-brand-violet-700 dark:text-brand-violet-300 font-bold uppercase tracking-wider text-[9px]">{{ t(`settings.criterion_${calculatedCriterion}`) }}</span>
            </div>
          </div>
        </div>

        <!-- Step by Step Segment list -->
        <div class="p-5 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
          <h4 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-3">
            {{ t('analysis.leg_title') }}
          </h4>
          <div class="max-h-[32rem] overflow-y-auto space-y-3 pr-1">
            <div
              v-for="(segment, idx) in routeSegments"
              :key="idx"
              class="p-4 bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 rounded-xl space-y-3"
            >
              <!-- Leg Header -->
              <div class="flex items-center justify-between gap-2.5 flex-wrap">
                <div class="flex items-center gap-1.5">
                  <div class="w-5 h-5 rounded-full flex items-center justify-center bg-brand-violet-600 text-white text-[10px] font-extrabold shrink-0">
                    {{ idx + 1 }}
                  </div>
                  <span class="font-bold text-slate-900 dark:text-white text-sm truncate max-w-[120px] sm:max-w-[180px]">
                    {{ getShortName(segment.from.name) }}
                  </span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                    <polyline points="12 5 19 12 12 19"></polyline>
                  </svg>
                  <span class="font-bold text-slate-900 dark:text-white text-sm truncate max-w-[120px] sm:max-w-[180px]">
                    {{ getShortName(segment.to.name) }}
                  </span>
                </div>
                
                <!-- Badges -->
                <div class="flex items-center gap-1">
                  <span
                    v-if="idx === mostExpensiveLegIndex"
                    class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-rose-100 dark:bg-rose-950/50 text-rose-600 dark:text-rose-400 border border-rose-200/20"
                  >
                    {{ t('analysis.leg_expensive') }}
                  </span>
                  <span
                    v-if="idx === longestLegIndex"
                    class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-amber-100 dark:bg-amber-950/50 text-amber-600 dark:text-amber-400 border border-amber-200/20"
                  >
                    {{ t('analysis.leg_longest') }}
                  </span>
                </div>
              </div>

              <!-- Leg Details grid -->
              <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs border-t border-slate-100 dark:border-slate-800/60 pt-3">
                <!-- Dist / Dur -->
                <div class="space-y-0.5">
                  <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.leg_dist_dur') }}</div>
                  <div class="font-bold text-slate-800 dark:text-slate-200">
                    <span class="font-mono">{{ (segment.distance / 1000).toFixed(1) }} km</span> • <span class="font-mono text-slate-500">{{ formatDuration(segment.duration) }}</span>
                  </div>
                </div>

                <!-- Fuel cost -->
                <div class="space-y-0.5">
                  <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.leg_fuel') }}</div>
                  <div class="font-bold text-slate-800 dark:text-slate-200 font-mono">
                    {{ segment.fuel.toFixed(2) }} €
                  </div>
                </div>

                <!-- Toll cost -->
                <div class="space-y-0.5">
                  <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.leg_toll') }}</div>
                  <div class="font-bold text-slate-800 dark:text-slate-200 font-mono">
                    {{ segment.toll.toFixed(2) }} €
                  </div>
                </div>

                <!-- Total cost -->
                <div class="space-y-0.5">
                  <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium">{{ t('analysis.leg_total_cost') }}</div>
                  <div class="font-bold text-brand-violet-600 dark:text-brand-violet-400 font-mono">
                    {{ segment.totalCost.toFixed(2) }} €
                  </div>
                </div>
              </div>
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
  totalDuration: number
  totalFuelCost: number
  totalTollCost: number
  calculatedCriterion: string
  fuelType: string
  fuelConsumption: number
  fuelPrice: number
  tollPrice: number
  timeValue: number
  distanceMatrix?: Map<string, Map<string, number>>
  durationMatrix?: Map<string, Map<string, number>>
  fuelCostMatrix?: Map<string, Map<string, number>>
  tollCostMatrix?: Map<string, Map<string, number>>
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

const formatDuration = (sec: number) => {
  const hours = Math.floor(sec / 3600)
  const minutes = Math.floor((sec % 3600) / 60)
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes} min`
}

// CO2 footprint calculations
const estimatedCO2 = computed(() => {
  const liters = (props.totalDistance / 100000) * props.fuelConsumption
  let factor = 2.31 // default petrol
  if (props.fuelType === 'Gazole') factor = 2.68
  else if (props.fuelType === 'E85') factor = 1.51
  else if (props.fuelType === 'GPLc') factor = 1.66
  return liters * factor
})

// Cost share percentages
const fuelSharePercent = computed(() => {
  const total = props.totalFuelCost + props.totalTollCost
  if (total === 0) return 50
  return Math.round((props.totalFuelCost / total) * 100)
})

const tollSharePercent = computed(() => {
  const total = props.totalFuelCost + props.totalTollCost
  if (total === 0) return 50
  return Math.round((props.totalTollCost / total) * 100)
})

// Compute step by step details with cost, duration and accumulators
const routeSegments = computed(() => {
  const segments = []
  let cumulativeDistance = 0
  let cumulativeDuration = 0
  let cumulativeFuel = 0
  let cumulativeToll = 0
  
  for (let i = 0; i < props.locations.length - 1; i++) {
    const from = props.locations[i]
    const to = props.locations[i + 1]
    
    const distance = props.distanceMatrix?.get(from.id)?.get(to.id) || 0
    const duration = props.durationMatrix?.get(from.id)?.get(to.id) || 0
    const fuel = props.fuelCostMatrix?.get(from.id)?.get(to.id) || 0
    const toll = props.tollCostMatrix?.get(from.id)?.get(to.id) || 0
    
    cumulativeDistance += distance
    cumulativeDuration += duration
    cumulativeFuel += fuel
    cumulativeToll += toll
    
    segments.push({
      from,
      to,
      distance,
      duration,
      fuel,
      toll,
      totalCost: fuel + toll,
      cumulativeDistance,
      cumulativeDuration,
      cumulativeFuel,
      cumulativeToll
    })
  }
  return segments
})

// Find longest and most expensive legs
const mostExpensiveLegIndex = computed(() => {
  if (routeSegments.value.length <= 1) return -1
  let maxCost = -1
  let maxIdx = -1
  routeSegments.value.forEach((seg, idx) => {
    if (seg.totalCost > maxCost) {
      maxCost = seg.totalCost
      maxIdx = idx
    }
  })
  return maxIdx
})

const longestLegIndex = computed(() => {
  if (routeSegments.value.length <= 1) return -1
  let maxDist = -1
  let maxIdx = -1
  routeSegments.value.forEach((seg, idx) => {
    if (seg.distance > maxDist) {
      maxDist = seg.distance
      maxIdx = idx
    }
  })
  return maxIdx
})

const diagnosticText = computed(() => {
  return t(`analysis.diagnostic_${props.calculatedCriterion}`, { timeValue: props.timeValue })
})
</script>
