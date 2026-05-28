<template>
  <div class="relative w-full h-full">
    <!-- Map Container -->
    <div ref="mapContainer" class="w-full h-full rounded-2xl overflow-hidden shadow-inner border border-slate-100 dark:border-slate-800"></div>
    
    <!-- Long press hint overlay -->
    <div class="absolute bottom-4 left-4 bg-slate-900/80 backdrop-blur-md text-white text-xs px-3 py-1.5 rounded-full pointer-events-none flex items-center gap-1.5 shadow-md border border-white/10 z-[1000] opacity-80">
      <span class="w-2 h-2 rounded-full bg-brand-pink-500 animate-pulse"></span>
      <span>Appui long sur la carte pour ajouter une étape</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch, toRaw } from 'vue'
import type { Location } from '../utils/routeOptimizer'

const props = defineProps<{
  locations: Location[]
  optimalRoute: [number, number][]
  theme: 'light' | 'dark'
  isAnimating: boolean
  animationProgress: number
}>()

const emit = defineEmits<{
  (e: 'map-click', lat: number, lng: number): void
}>()

const mapContainer = ref<HTMLDivElement | null>(null)
let map: any = null
let tileLayer: any = null
let markersGroup: any = null
let routePolyline: any = null
let glowPolyline: any = null
let L: any = null

// Coordinates for Rhône-Alpes (centered on Lyon)
const defaultCenter: [number, number] = [45.7578137, 4.8320114]
const defaultZoom = 8

// Map press state
let pressTimer: any = null
let startLatLng: any = null

const startPress = (latlng: any) => {
  startLatLng = latlng
  if (pressTimer) clearTimeout(pressTimer)
  pressTimer = setTimeout(() => {
    if (startLatLng) {
      emit('map-click', startLatLng.lat, startLatLng.lng)
    }
    pressTimer = null
  }, 600)
}

const cancelPress = () => {
  if (pressTimer) {
    clearTimeout(pressTimer)
    pressTimer = null
  }
}

const checkMove = (latlng: any) => {
  if (pressTimer && startLatLng && map) {
    const p1 = map.latLngToContainerPoint(startLatLng)
    const p2 = map.latLngToContainerPoint(latlng)
    if (p1.distanceTo(p2) > 10) {
      cancelPress()
    }
  }
}

// Tile layers mapping
const getTileUrl = (t: 'light' | 'dark') => {
  return t === 'dark' 
    ? 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
    : 'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png'
}

const getTileAttribution = () => {
  return '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
}

// Markers helper
const createMarkerIcon = (order: number, isStart: boolean, isEnd: boolean) => {
  if (isStart) {
    return L.divIcon({
      className: 'custom-map-marker-flag start-flag',
      html: `
        <div class="flex items-center justify-center w-8 h-8 rounded-full bg-emerald-500 text-white shadow-md border-2 border-white transform transition-transform hover:scale-110">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"></path>
            <line x1="4" y1="22" x2="4" y2="15"></line>
          </svg>
        </div>
      `,
      iconSize: [32, 32],
      iconAnchor: [8, 32]
    })
  }

  if (isEnd) {
    return L.divIcon({
      className: 'custom-map-marker-flag end-flag',
      html: `
        <div class="flex items-center justify-center w-8 h-8 rounded-full bg-slate-900 text-white shadow-md border-2 border-white transform transition-transform hover:scale-110 overflow-hidden relative">
          <svg viewBox="0 0 24 24" class="w-5 h-5 text-white" fill="currentColor">
            <rect x="2" y="2" width="2" height="20" fill="#94a3b8"/>
            <rect x="4" y="4" width="3.5" height="3" fill="#ffffff"/>
            <rect x="7.5" y="4" width="3.5" height="3" fill="#000000"/>
            <rect x="11" y="4" width="3.5" height="3" fill="#ffffff"/>
            <rect x="14.5" y="4" width="3.5" height="3" fill="#000000"/>
            <rect x="4" y="7" width="3.5" height="3" fill="#000000"/>
            <rect x="7.5" y="7" width="3.5" height="3" fill="#ffffff"/>
            <rect x="11" y="7" width="3.5" height="3" fill="#000000"/>
            <rect x="14.5" y="7" width="3.5" height="3" fill="#ffffff"/>
            <rect x="4" y="10" width="3.5" height="3" fill="#ffffff"/>
            <rect x="7.5" y="10" width="3.5" height="3" fill="#000000"/>
            <rect x="11" y="10" width="3.5" height="3" fill="#ffffff"/>
            <rect x="14.5" y="10" width="3.5" height="3" fill="#000000"/>
          </svg>
        </div>
      `,
      iconSize: [32, 32],
      iconAnchor: [8, 32]
    })
  }

  // Intermediate steps
  return L.divIcon({
    className: 'custom-map-marker-step',
    html: `
      <div class="flex items-center justify-center w-7 h-7 rounded-full bg-brand-violet-600 text-white text-xs font-bold shadow-md border-2 border-white transform transition-all hover:scale-110">
        ${order}
      </div>
    `,
    iconSize: [28, 28],
    iconAnchor: [14, 14]
  })
}

const updateMarkers = () => {
  if (!map || !L) return
  
  // Clear existing markers
  markersGroup.clearLayers()
  
  props.locations.forEach((loc, idx) => {
    const isStart = idx === 0
    const isEnd = idx === props.locations.length - 1 && props.locations.length > 1
    
    const marker = L.marker([loc.lat, loc.lng], {
      icon: createMarkerIcon(loc.order, isStart, isEnd)
    })
    
    // Popup details
    marker.bindPopup(`
      <div class="text-sm">
        <div class="font-bold text-brand-violet-700 dark:text-brand-violet-400">Étape ${loc.order}</div>
        <div class="text-slate-700 dark:text-slate-300 font-medium mt-0.5 leading-tight">${loc.name.split(',')[0]}</div>
        <div class="text-slate-400 dark:text-slate-500 text-xs mt-1 font-mono">${loc.lat.toFixed(5)}, ${loc.lng.toFixed(5)}</div>
      </div>
    `)
    
    markersGroup.addLayer(marker)
  })
}

// Drawing route helper
const updateRoute = () => {
  if (!map || !L) return
  
  if (routePolyline) map.removeLayer(routePolyline)
  if (glowPolyline) map.removeLayer(glowPolyline)
  
  let routeCoords = toRaw(props.optimalRoute)
  
  if (props.isAnimating && routeCoords.length > 0) {
    const pointsToDraw = Math.ceil(routeCoords.length * props.animationProgress)
    routeCoords = routeCoords.slice(0, pointsToDraw)
  }
  
  if (routeCoords.length === 0) return
  
  // Outer glowing polyline (thick translucent pink)
  glowPolyline = L.polyline(routeCoords, {
    color: '#ec4899', // brand-pink-500
    weight: 8,
    opacity: 0.35,
    lineJoin: 'round',
    lineCap: 'round'
  }).addTo(map)
  
  // Inner core polyline (medium bright pink)
  routePolyline = L.polyline(routeCoords, {
    color: '#db2777', // brand-pink-600
    weight: 4,
    opacity: 0.95,
    lineJoin: 'round',
    lineCap: 'round'
  }).addTo(map)
  
  // Adjust bounds when new route is computed (and not animating)
  if (routeCoords.length > 0 && !props.isAnimating) {
    const bounds = L.latLngBounds(routeCoords)
    map.fitBounds(bounds, { padding: [50, 50], maxZoom: 14 })
  }
}

// Initializing map
onMounted(async () => {
  L = await import('leaflet')
  
  if (!mapContainer.value) return
  
  map = L.map(mapContainer.value, {
    zoomControl: false // custom position zoom control
  }).setView(defaultCenter, defaultZoom)
  
  // Add zoom control at top-left
  L.control.zoom({
    position: 'topleft'
  }).addTo(map)
  
  // Tile layer
  tileLayer = L.tileLayer(getTileUrl(props.theme), {
    attribution: getTileAttribution(),
    maxZoom: 19
  }).addTo(map)
  
  markersGroup = L.layerGroup().addTo(map)
  
  // Map interaction events for long press
  map.on('mousedown', (e: any) => {
    if (e.originalEvent && e.originalEvent.button !== 0) return
    startPress(e.latlng)
  })
  
  map.on('mouseup', () => cancelPress())
  map.on('mousemove', (e: any) => checkMove(e.latlng))
  map.on('dragstart', () => cancelPress())
  map.on('zoomstart', () => cancelPress())
  
  // Mobile touch support
  map.on('touchstart', (e: any) => startPress(e.latlng))
  map.on('touchend', () => cancelPress())
  map.on('touchmove', (e: any) => checkMove(e.latlng))
  
  // Center map on Rhône-Alpes if we have no stops yet
  if (props.locations.length === 0) {
    map.setView(defaultCenter, defaultZoom)
  } else {
    // Fit bounds of locations
    const bounds = L.latLngBounds(props.locations.map(loc => [loc.lat, loc.lng]))
    map.fitBounds(bounds, { padding: [50, 50] })
  }
  
  updateMarkers()
  updateRoute()
})

onBeforeUnmount(() => {
  if (pressTimer) clearTimeout(pressTimer)
  if (map) {
    map.remove()
    map = null
  }
})

// Observers
watch(() => props.theme, (newTheme) => {
  if (tileLayer && map) {
    tileLayer.setUrl(getTileUrl(newTheme))
  }
})

watch(() => props.locations, () => {
  updateMarkers()
  
  // Auto-adjust view when a location is added
  if (map && props.locations.length > 0 && props.optimalRoute.length === 0) {
    const bounds = L.latLngBounds(props.locations.map(loc => [loc.lat, loc.lng]))
    map.fitBounds(bounds, { padding: [50, 50], maxZoom: 12 })
  }
}, { deep: true })

watch(() => props.optimalRoute, () => {
  updateRoute()
}, { deep: true })

watch(() => props.animationProgress, () => {
  if (props.isAnimating) {
    updateRoute()
  }
})
</script>

<style>
/* Checkered SVG pole offset */
.custom-map-marker-flag {
  overflow: visible !important;
}

.custom-map-marker-flag div {
  transform: translateY(-8px);
}
</style>
