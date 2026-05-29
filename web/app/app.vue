<template>
  <div class="flex flex-col h-screen w-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 overflow-hidden font-sans">
    <!-- Navbar / Header -->
    <header class="h-16 shrink-0 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 px-6 flex items-center justify-between shadow-sm z-50">
      <div class="flex items-center gap-3">
        <!-- Logo -->
        <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-brand-violet-600 to-brand-pink-500 flex items-center justify-center text-white shadow-md shadow-brand-violet-500/20 transform rotate-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 animate-pulse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="3 6 9 3 15 6 21 3 21 18 15 21 9 18 3 21"></polygon>
            <line x1="9" y1="3" x2="9" y2="18"></line>
            <line x1="15" y1="6" x2="15" y2="21"></line>
          </svg>
        </div>
        <div>
          <h1 class="text-lg font-black tracking-tight bg-gradient-to-r from-brand-violet-600 to-brand-pink-500 bg-clip-text text-transparent">
            METHIL MAPS
          </h1>
          <p class="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase tracking-wider -mt-1">
            Smart Route Optimizer
          </p>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <!-- Language selector -->
        <button
          @click="toggleLocale"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-xs font-semibold uppercase tracking-wider transition-colors"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m5 8 6 6"></path>
            <path d="m4 14 6-6 8 8"></path>
            <path d="M2 5h12"></path>
            <path d="M7 2h1"></path>
            <path d="m22 22-5-10-5 10"></path>
            <path d="M14 18h6"></path>
          </svg>
          <span>{{ locale }}</span>
        </button>

        <!-- Theme selector -->
        <button
          @click="toggleTheme"
          class="p-2 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 transition-colors"
        >
          <svg v-if="theme === 'light'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="4"></circle>
            <path d="M12 2v2"></path>
            <path d="M12 20v2"></path>
            <path d="m4.93 4.93 1.41 1.41"></path>
            <path d="m17.66 17.66 1.41 1.41"></path>
            <path d="M2 12h2"></path>
            <path d="M20 12h2"></path>
            <path d="m6.34 17.66-1.41 1.41"></path>
            <path d="m19.07 4.93-1.41 1.41"></path>
          </svg>
        </button>
      </div>
    </header>

    <!-- Main Workspace -->
    <main class="flex-1 flex flex-col md:flex-row overflow-hidden relative">
      <!-- Left Side Panel / Sidebar Drawer -->
      <aside class="w-full md:w-96 shrink-0 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col h-1/2 md:h-full z-10 shadow-lg md:shadow-none">
        <div class="p-5 flex-1 overflow-y-auto space-y-5">
          <!-- Search box -->
          <div class="space-y-2">
            <label class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
              {{ t('planner.search_placeholder') }}
            </label>
            <div class="relative">
              <input
                type="text"
                v-model="searchQuery"
                @input="handleSearch"
                :placeholder="t('planner.search_placeholder')"
                class="w-full h-10 pl-3 pr-10 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-950 text-sm focus:outline-none focus:ring-2 focus:ring-brand-violet-500/20 focus:border-brand-violet-500 transition-all placeholder:text-slate-400"
              />
              <div class="absolute right-3 top-3">
                <svg v-if="isSearching" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 animate-spin text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                </svg>
              </div>
            </div>
          </div>

          <!-- Search Results Dropdown -->
          <div v-if="searchResults.length > 0" class="border border-slate-100 dark:border-slate-800 rounded-xl overflow-hidden shadow-lg bg-white dark:bg-slate-900 max-h-40 overflow-y-auto">
            <div
              v-for="result in searchResults"
              :key="result.place_id"
              @click="addLocation(result)"
              class="p-2.5 hover:bg-slate-50 dark:hover:bg-slate-800 cursor-pointer flex items-start gap-2.5 transition-colors border-b border-slate-100 dark:border-slate-800 last:border-0"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 shrink-0 mt-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"></path>
                <circle cx="12" cy="10" r="3"></circle>
              </svg>
              <div class="min-w-0">
                <div class="font-bold text-xs text-slate-900 dark:text-white truncate">
                  {{ result.display_name.split(',')[0] }}
                </div>
                <div class="text-[10px] text-slate-400 dark:text-slate-500 truncate mt-0.5">
                  {{ result.display_name.split(',').slice(1).join(',').trim() }}
                </div>
              </div>
            </div>
          </div>

          <!-- Selected Stops -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
                {{ t('planner.selected_locations') }} ({{ locations.length }})
              </label>
              <button
                v-if="locations.length > 0"
                @click="clearLocations"
                class="text-[10px] font-bold text-rose-500 hover:text-rose-600 transition-colors uppercase tracking-wider"
              >
                Tout effacer
              </button>
            </div>

            <!-- Empty state -->
            <div
              v-if="locations.length === 0"
              class="p-5 border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-2xl text-center text-xs text-slate-400 dark:text-slate-500 leading-relaxed font-medium"
            >
              {{ t('planner.no_locations') }}
            </div>

            <!-- Stops list -->
            <div v-else class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div
                v-for="(loc, idx) in locations"
                :key="loc.id"
                draggable="true"
                @dragstart="onDragStart($event, idx)"
                @dragover.prevent="onDragOver(idx)"
                @dragend="onDragEnd"
                class="flex items-center gap-2 bg-slate-50 dark:bg-slate-950 border rounded-xl transition-all group cursor-grab active:cursor-grabbing p-1.5"
                :class="draggedIdx === idx 
                  ? 'opacity-40 border-dashed border-brand-violet-400 dark:border-brand-violet-800 bg-brand-violet-50/10'
                  : 'border-slate-200 dark:border-slate-800 hover:border-slate-300 dark:hover:border-slate-700'"
              >
                <!-- Drag Grip Icon -->
                <div class="text-slate-400 dark:text-slate-600 hover:text-slate-500 dark:hover:text-slate-400 shrink-0 select-none cursor-grab">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="9" cy="5" r="1"></circle>
                    <circle cx="9" cy="12" r="1"></circle>
                    <circle cx="9" cy="19" r="1"></circle>
                    <circle cx="15" cy="5" r="1"></circle>
                    <circle cx="15" cy="12" r="1"></circle>
                    <circle cx="15" cy="19" r="1"></circle>
                  </svg>
                </div>

                <!-- Order Badge -->
                <div
                  class="w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-extrabold text-white shrink-0 select-none"
                  :class="idx === 0 ? 'bg-emerald-500' : (idx === locations.length - 1 ? 'bg-slate-900' : 'bg-brand-violet-600')"
                >
                  {{ idx + 1 }}
                </div>

                <!-- Stop Name -->
                <div class="flex-1 min-w-0 select-none">
                  <div class="font-bold text-xs text-slate-800 dark:text-slate-200 truncate pr-1">
                    {{ loc.name.split(',')[0] }}
                  </div>
                </div>

                <!-- Move Buttons -->
                <div class="flex items-center gap-1 opacity-60 group-hover:opacity-100 transition-opacity">
                  <button
                    @click="moveLocation(idx, 'up')"
                    :disabled="idx === 0"
                    class="p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-800 disabled:opacity-30 disabled:hover:bg-transparent"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="18 15 12 9 6 15"></polyline>
                    </svg>
                  </button>
                  <button
                    @click="moveLocation(idx, 'down')"
                    :disabled="idx === locations.length - 1"
                    class="p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-800 disabled:opacity-30 disabled:hover:bg-transparent"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="6 9 12 15 18 9"></polyline>
                    </svg>
                  </button>
                  
                  <!-- Remove Stop -->
                  <button
                    @click="removeLocation(loc.id)"
                    class="p-1 rounded hover:bg-red-50 dark:hover:bg-red-950/30 text-rose-500 hover:text-rose-600 transition-colors ml-0.5"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Optimization Criteria Widget -->
          <div class="space-y-2">
            <label class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
              {{ t('settings.optimization_criterion') }}
            </label>
            <div class="grid grid-cols-5 gap-1">
              <button
                v-for="crit in ['distance', 'time', 'economy', 'smart', 'toll']"
                :key="crit"
                @click="optimizationCriterion = crit"
                class="py-2 px-0.5 text-center text-[10px] font-bold rounded-lg border transition-all truncate"
                :class="optimizationCriterion === crit
                  ? 'bg-brand-violet-600 border-brand-violet-600 text-white shadow-sm shadow-brand-violet-500/20'
                  : 'bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'"
              >
                {{ t(`settings.criterion_${crit}`) }}
              </button>
            </div>
          </div>

          <!-- Time Value Slider (Only shown at top level if Smart is selected) -->
          <div v-if="optimizationCriterion === 'smart'" class="p-4 border border-slate-200 dark:border-slate-800 rounded-2xl bg-white dark:bg-slate-900 space-y-1.5 shadow-sm">
            <div class="flex items-center justify-between text-xs font-bold">
              <span class="text-slate-600 dark:text-slate-400">{{ t('settings.time_value') }}</span>
              <span class="text-brand-violet-600 font-mono">{{ timeValue.toFixed(0) }} €/h</span>
            </div>
            <input
              type="range"
              min="5"
              max="50"
              step="1"
              v-model.number="timeValue"
              class="w-full h-1 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-brand-violet-600"
            />
          </div>

          <!-- Settings (Collapsible Card Widget) -->
          <div class="border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden bg-slate-50/50 dark:bg-slate-950/20">
            <button
              @click="isSettingsOpen = !isSettingsOpen"
              class="w-full p-4 flex items-center justify-between font-bold text-sm text-slate-800 dark:text-slate-200 bg-white dark:bg-slate-900 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              :class="isSettingsOpen ? 'border-b border-slate-200 dark:border-slate-800' : 'border-b border-transparent'"
            >
              <div class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-brand-violet-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                </svg>
                <span>{{ t('settings.title') }}</span>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform duration-200" :class="isSettingsOpen ? 'transform rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
            </button>

            <!-- Collapsible Settings Body -->
            <div v-show="isSettingsOpen" class="p-4 space-y-4">
              <!-- Fuel Type -->
              <div class="space-y-1.5">
                <label class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
                  {{ t('settings.fuel_type') }}
                </label>
                <select
                  v-model="fuelType"
                  class="w-full h-10 px-3 text-xs font-medium rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 focus:outline-none focus:ring-1 focus:ring-brand-violet-500"
                >
                  <option v-for="ft in ['Gazole', 'E10', 'SP98', 'SP95', 'E85', 'GPLc']" :key="ft" :value="ft">
                    {{ ft }}
                  </option>
                </select>
              </div>

              <!-- Fuel Consumption Slider -->
              <div class="space-y-1.5">
                <div class="flex items-center justify-between text-xs font-bold">
                  <span class="text-slate-600 dark:text-slate-400">{{ t('settings.fuel_consumption') }}</span>
                  <span class="text-brand-violet-600 font-mono">{{ fuelConsumption.toFixed(1) }} L</span>
                </div>
                <input
                  type="range"
                  min="3"
                  max="15"
                  step="0.5"
                  v-model.number="fuelConsumption"
                  class="w-full h-1 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-brand-violet-600"
                />
              </div>

              <!-- Fuel Price Input -->
              <div class="space-y-1.5">
                <label class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
                  {{ t('settings.fuel_price') }}
                </label>
                <div class="relative">
                  <input
                    type="number"
                    step="0.00001"
                    v-model.number="fuelPrice"
                    class="w-full h-10 pl-3 pr-12 text-xs font-mono font-bold rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 focus:outline-none"
                  />
                  <span class="absolute right-3 top-3 text-[10px] font-bold text-slate-400">€/L</span>
                </div>
              </div>

              <!-- Toll Price Slider -->
              <div class="space-y-1.5">
                <div class="flex items-center justify-between text-xs font-bold">
                  <span class="text-slate-600 dark:text-slate-400">{{ t('settings.toll_price') }}</span>
                  <span class="text-brand-violet-600 font-mono">{{ tollPrice.toFixed(2) }} €/km</span>
                </div>
                <input
                  type="range"
                  min="0.05"
                  max="0.25"
                  step="0.01"
                  v-model.number="tollPrice"
                  class="w-full h-1 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-brand-violet-600"
                />
                <p class="text-[10px] text-slate-400 dark:text-slate-500 leading-normal">
                  {{ t('settings.toll_price_hint') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Sticky Footer containing Route Action Button -->
        <div class="p-5 border-t border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shrink-0">
          <button
            @click="calculateRoute"
            :disabled="locations.length < 2 || isCalculating"
            class="w-full py-2.5 px-4 rounded-xl text-white font-semibold text-xs uppercase tracking-wider flex items-center justify-center gap-2 border border-transparent shadow-md transform transition-all select-none disabled:opacity-40 disabled:pointer-events-none active:scale-[0.98]"
            :class="isCalculating
              ? 'bg-slate-700'
              : 'bg-brand-violet-600 hover:bg-brand-violet-700 shadow-brand-violet-500/10 active:shadow-none'"
          >
            <span v-if="isCalculating" class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
              </svg>
              <span>{{ t('planner.calculating') }}</span>
            </span>
            <span v-else class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
              <span>{{ t('planner.find_route') }}</span>
            </span>
          </button>
        </div>
      </aside>

      <!-- Map Display Area -->
      <section class="flex-1 h-1/2 md:h-full relative z-0">
        <client-only>
          <LeafletMap
            :locations="locations"
            :optimalRoute="optimalRoute"
            :theme="theme"
            :isAnimating="isAnimating"
            :animationProgress="animationProgress"
            @map-click="addLocationFromCoordinates"
          />
        </client-only>

        <!-- Floating OSRM Fallback Warning -->
        <div v-if="isUsingOSRM && optimalRoute.length > 0" class="absolute top-4 left-4 right-4 md:left-auto md:right-4 md:w-80 bg-amber-500/90 backdrop-blur-md text-slate-950 p-3 rounded-xl shadow-lg border border-amber-400/20 z-[1000] flex gap-2.5 animate-pop-in">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0 mt-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          <div class="text-[11px] leading-snug font-medium">
            <div class="font-extrabold text-xs mb-0.5">{{ t('planner.osrm_warning_title') }}</div>
            <div>{{ t('planner.osrm_warning_desc') }}</div>
          </div>
        </div>

        <!-- Floating Route Results Card Overlay (Glow corner) -->
        <div
          v-if="optimalRoute.length > 0"
          class="absolute bottom-4 right-4 w-[calc(100vw-2rem)] md:w-80 bg-white/90 dark:bg-slate-900/90 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-5 rounded-2xl z-[1000] animate-pop-in shadow-[0_0_25px_rgba(109,40,217,0.25)] dark:shadow-[0_0_25px_rgba(109,40,217,0.18)]"
        >
          <!-- Card Header / Frozen Badge -->
          <div class="flex items-center justify-between mb-4">
            <h4 class="text-sm font-extrabold tracking-tight text-slate-950 dark:text-white uppercase">
              {{ t('planner.results') }}
            </h4>
            <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase tracking-wider bg-brand-violet-100 dark:bg-brand-violet-950 text-brand-violet-600 dark:text-brand-violet-400">
              {{ t(`settings.criterion_${calculatedCriterion}`) }}
            </span>
          </div>

          <!-- Cost breakdowns grid -->
          <div class="grid grid-cols-2 gap-3.5 mb-4 text-xs">
            <div class="p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800">
              <div class="text-slate-400 dark:text-slate-500 font-bold uppercase text-[9px] tracking-wider">
                {{ t('planner.distance') }}
              </div>
              <div class="text-sm font-extrabold text-slate-800 dark:text-slate-200 font-mono mt-0.5">
                {{ (totalDistance / 1000).toFixed(1) }} km
              </div>
            </div>

            <div class="p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800">
              <div class="text-slate-400 dark:text-slate-500 font-bold uppercase text-[9px] tracking-wider">
                {{ t('planner.duration') }}
              </div>
              <div class="text-sm font-extrabold text-slate-800 dark:text-slate-200 font-mono mt-0.5">
                {{ formatDuration(totalDuration) }}
              </div>
            </div>

            <div class="p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800">
              <div class="text-slate-400 dark:text-slate-500 font-bold uppercase text-[9px] tracking-wider">
                {{ t('planner.fuel_cost') }}
              </div>
              <div class="text-sm font-extrabold text-slate-800 dark:text-slate-200 font-mono mt-0.5">
                {{ totalFuelCost.toFixed(2) }} €
              </div>
            </div>

            <div class="p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800">
              <div class="text-slate-400 dark:text-slate-500 font-bold uppercase text-[9px] tracking-wider">
                {{ t('planner.toll_cost') }}
              </div>
              <div class="text-sm font-extrabold text-slate-800 dark:text-slate-200 font-mono mt-0.5">
                {{ totalTollCost.toFixed(2) }} €
              </div>
            </div>
          </div>

          <!-- Action buttons layout -->
          <div class="space-y-1.5">
            <div class="grid grid-cols-2 gap-1.5">
              <!-- Animation Toggle -->
              <button
                @click="toggleAnimation"
                class="flex items-center justify-center gap-1.5 py-2 rounded-lg border border-slate-200 dark:border-slate-800 text-xs font-bold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              >
                <svg v-if="isAnimating" xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-rose-500" viewBox="0 0 24 24" fill="currentColor">
                  <rect x="4" y="4" width="16" height="16" rx="2"></rect>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-500" viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
                <span>{{ isAnimating ? t('planner.stop') : t('planner.animate') }}</span>
              </button>

              <!-- Share Route -->
              <button
                @click="shareRoute"
                class="flex items-center justify-center gap-1.5 py-2 rounded-lg border border-slate-200 dark:border-slate-800 text-xs font-bold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-brand-violet-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="18" cy="5" r="3"></circle>
                  <circle cx="6" cy="12" r="3"></circle>
                  <circle cx="18" cy="19" r="3"></circle>
                  <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
                  <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
                </svg>
                <span>{{ t('planner.share') }}</span>
              </button>
            </div>

            <div class="grid grid-cols-2 gap-1.5">
              <!-- Export json -->
              <button
                @click="exportRouteData"
                class="flex items-center justify-center gap-1.5 py-2 rounded-lg border border-slate-200 dark:border-slate-800 text-xs font-bold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors whitespace-nowrap"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                  <polyline points="7 10 12 15 17 10"></polyline>
                  <line x1="12" y1="15" x2="12" y2="3"></line>
                </svg>
                <span>{{ t('planner.export') }}</span>
              </button>

              <!-- Open Detailed Analysis Modal -->
              <button
                @click="isAnalysisOpen = true"
                class="flex items-center justify-center gap-1.5 py-2 rounded-lg bg-brand-violet-50 hover:bg-brand-violet-100 dark:bg-brand-violet-950/40 dark:hover:bg-brand-violet-950/60 text-brand-violet-600 dark:text-brand-violet-400 text-xs font-extrabold transition-colors border border-brand-violet-200/20 whitespace-nowrap"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="20" x2="18" y2="10"></line>
                  <line x1="12" y1="20" x2="12" y2="4"></line>
                  <line x1="6" y1="20" x2="6" y2="14"></line>
                </svg>
                <span>{{ t('analysis.title') }}</span>
              </button>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Route Analysis Detailed Modal -->
    <RouteAnalysis
      :isOpen="isAnalysisOpen"
      :locations="locations"
      :totalDistance="totalDistance"
      :totalDuration="totalDuration"
      :totalFuelCost="totalFuelCost"
      :totalTollCost="totalTollCost"
      :calculatedCriterion="calculatedCriterion"
      :fuelType="fuelType"
      :fuelConsumption="fuelConsumption"
      :fuelPrice="fuelPrice"
      :tollPrice="tollPrice"
      :timeValue="timeValue"
      :distanceMatrix="routeOptimizer.distanceMatrix"
      :durationMatrix="routeOptimizer.durationMatrix"
      :fuelCostMatrix="routeOptimizer.fuelCostMatrix"
      :tollCostMatrix="routeOptimizer.tollCostMatrix"
      @close="isAnalysisOpen = false"
    />

    <!-- Toast Notifications -->
    <transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="notification && notification.visible"
        class="fixed top-6 right-6 z-[2000] w-full max-w-sm overflow-hidden bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border rounded-2xl shadow-xl transition-all"
        :class="notification.type === 'error' ? 'border-red-200 dark:border-red-900/50' : 'border-emerald-200 dark:border-emerald-900/50'"
      >
        <div class="p-4 flex items-start gap-3">
          <!-- Icon -->
          <div class="shrink-0 mt-0.5">
            <svg
              v-if="notification.type === 'error'"
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5 text-rose-500"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <svg
              v-else
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5 text-emerald-500"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
          </div>

          <!-- Content -->
          <div class="flex-1 min-w-0">
            <h5 class="text-xs font-black uppercase tracking-wider text-slate-900 dark:text-white">
              {{ notification.title }}
            </h5>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-snug">
              {{ notification.message }}
            </p>
          </div>

          <!-- Close Button -->
          <button
            @click="notification.visible = false"
            class="shrink-0 p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-slate-500 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, reactive } from 'vue'
import { useRoute, useRouter } from 'nuxt/app'
import { useTrans } from './composables/useTrans'
import { fetchAverageFuelPrice } from './utils/fuelApi'
import { RouteOptimizer } from './utils/routeOptimizer'
import type { Location, OptimizationCriterion } from './utils/routeOptimizer'
import type { FuelType } from './utils/fuelApi'

// Load translations and themes
const { t, locale, toggleLocale, setLocale } = useTrans()
const theme = ref<'light' | 'dark'>('light')

const toggleTheme = () => {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
  if (theme.value === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
  localStorage.setItem('methil_maps_theme', theme.value)
}

// Side drawers & modals open state
const isSettingsOpen = ref(false)
const isAnalysisOpen = ref(false)

// Notification system
const notification = ref<{
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message: string
  visible: boolean
} | null>(null)

let notificationTimeout: any = null

const showNotification = (type: 'success' | 'error' | 'warning' | 'info', title: string, message: string) => {
  if (notificationTimeout) clearTimeout(notificationTimeout)
  notification.value = { type, title, message, visible: true }
  notificationTimeout = setTimeout(() => {
    if (notification.value) {
      notification.value.visible = false
    }
  }, 5000)
}

// Routing states
const locations = ref<Location[]>([])
const optimalRoute = ref<[number, number][]>([])
const isCalculating = ref(false)

const totalDistance = ref(0)
const totalDuration = ref(0)
const totalFuelCost = ref(0)
const totalTollCost = ref(0)
const isUsingOSRM = ref(false)

let isInitializing = true

// Optimization configurations
const optimizationCriterion = ref<OptimizationCriterion>('distance')
const calculatedCriterion = ref<OptimizationCriterion>('distance')
const fuelType = ref<FuelType>('Gazole')
const fuelConsumption = ref(7.0)
const fuelPrice = ref(1.80)
const tollPrice = ref(0.13)
const timeValue = ref(15.0)

// Search inputs
const searchQuery = ref('')
const searchResults = ref<any[]>([])
const isSearching = ref(false)
let searchTimeout: any = null

// Animation variables
const isAnimating = ref(false)
const animationProgress = ref(0)
let animationFrameId: any = null
let animationStartTime = 0

// Route Optimizer Instance
const routeOptimizer = reactive(new RouteOptimizer())

// Nuxt URL Search Query Bindings
const route = useRoute()
const router = useRouter()

// Sync local state with URL Query Parameters
const syncStateToUrl = () => {
  const query: Record<string, string> = {}
  
  if (locations.value.length > 0) {
    try {
      // Base64 encoding compatible with original React implementation
      const locsData = btoa(unescape(encodeURIComponent(JSON.stringify(locations.value))))
      query.l = locsData
    } catch (e) {
      console.error(e)
    }
  }
  
  query.c = optimizationCriterion.value
  query.ft = fuelType.value
  query.fc = fuelConsumption.value.toString()
  query.fp = fuelPrice.value.toString()
  query.tp = tollPrice.value.toString()
  query.tv = timeValue.value.toString()
  
  // Set window location query parameter
  router.replace({ query })
}

// Read states from URL on mount
const loadStateFromUrl = () => {
  const query = route.query
  
  if (query.l && typeof query.l === 'string') {
    try {
      // Decode Base64 string compatible with original React implementation
      const decoded = JSON.parse(decodeURIComponent(escape(atob(query.l))))
      if (Array.isArray(decoded)) {
        locations.value = decoded
      }
    } catch (e) {
      console.error('Failed to parse stops from URL:', e)
    }
  }
  
  if (query.c && typeof query.c === 'string') {
    optimizationCriterion.value = query.c as OptimizationCriterion
    calculatedCriterion.value = query.c as OptimizationCriterion
  }
  if (query.ft && typeof query.ft === 'string') {
    fuelType.value = query.ft as FuelType
  }
  if (query.fc && typeof query.fc === 'string') {
    fuelConsumption.value = parseFloat(query.fc) || 7.0
  }
  if (query.fp && typeof query.fp === 'string') {
    fuelPrice.value = parseFloat(query.fp) || 1.80
  }
  if (query.tp && typeof query.tp === 'string') {
    tollPrice.value = parseFloat(query.tp) || 0.13
  }
  if (query.tv && typeof query.tv === 'string') {
    timeValue.value = parseFloat(query.tv) || 20.0
  }
}

// Nominatim Search
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  
  isSearching.value = true
  searchTimeout = setTimeout(async () => {
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery.value)}&limit=5`,
        {
          headers: {
            'Accept-Language': 'fr-FR,fr;q=0.9,en;q=0.8',
            'User-Agent': 'MethilMapsOptimizer/2.0'
          }
        }
      )
      if (response.ok) {
        searchResults.value = await response.json()
      }
    } catch (e) {
      console.error(e)
    } finally {
      isSearching.value = false
    }
  }, 500)
}

// Add location
const addLocation = (res: any) => {
  const newLoc: Location = {
    id: Date.now().toString() + Math.random().toString(36).substring(2, 5),
    name: res.display_name,
    lat: parseFloat(res.lat),
    lng: parseFloat(res.lon),
    order: locations.value.length + 1
  }
  
  locations.value.push(newLoc)
  searchQuery.value = ''
  searchResults.value = []
  clearRouteData()
}

// Click / Long Press Map adding
const addLocationFromCoordinates = async (lat: number, lng: number) => {
  const tempId = `click_${Date.now()}`
  const tempName = `${t('common.loading')} (${lat.toFixed(4)}, ${lng.toFixed(4)})`
  
  const newLoc: Location = {
    id: tempId,
    name: tempName,
    lat,
    lng,
    order: locations.value.length + 1
  }
  
  locations.value.push(newLoc)
  clearRouteData()
  
  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`,
      {
        headers: {
          'Accept-Language': 'fr-FR,fr;q=0.9,en;q=0.8',
          'User-Agent': 'MethilMapsOptimizer/2.0'
        }
      }
    )
    if (response.ok) {
      const data = await response.json()
      const finalName = data.display_name || `${lat.toFixed(5)}, ${lng.toFixed(5)}`
      locations.value = locations.value.map(loc => 
        loc.id === tempId ? { ...loc, name: finalName } : loc
      )
    } else {
      locations.value = locations.value.map(loc => 
        loc.id === tempId ? { ...loc, name: `${lat.toFixed(5)}, ${lng.toFixed(5)}` } : loc
      )
    }
  } catch (e) {
    locations.value = locations.value.map(loc => 
      loc.id === tempId ? { ...loc, name: `${lat.toFixed(5)}, ${lng.toFixed(5)}` } : loc
    )
  }
}

// Remove stop
const removeLocation = (id: string) => {
  locations.value = locations.value
    .filter(loc => loc.id !== id)
    .map((loc, idx) => ({ ...loc, order: idx + 1 }))
  clearRouteData()
}

// Re-order stops
const moveLocation = (index: number, direction: 'up' | 'down') => {
  const targetIndex = direction === 'up' ? index - 1 : index + 1
  if (targetIndex < 0 || targetIndex >= locations.value.length) return
  
  const updated = [...locations.value]
  const temp = updated[index]
  updated[index] = updated[targetIndex]
  updated[targetIndex] = temp
  
  locations.value = updated.map((loc, idx) => ({ ...loc, order: idx + 1 }))
  clearRouteData()
}

// Drag and drop states & handlers
const draggedIdx = ref<number | null>(null)

const onDragStart = (event: DragEvent, index: number) => {
  draggedIdx.value = index
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', index.toString())
  }
}

const onDragOver = (index: number) => {
  if (draggedIdx.value === null || draggedIdx.value === index) return
  
  const items = [...locations.value]
  const draggedItem = items[draggedIdx.value]
  
  // Swap items dynamically in the list for real-time visual feedback
  items.splice(draggedIdx.value, 1)
  items.splice(index, 0, draggedItem)
  
  locations.value = items.map((loc, i) => ({ ...loc, order: i + 1 }))
  draggedIdx.value = index
}

const onDragEnd = () => {
  draggedIdx.value = null
  clearRouteData()
}

const clearLocations = () => {
  locations.value = []
  clearRouteData()
  syncStateToUrl()
}

const clearRouteData = () => {
  optimalRoute.value = []
  totalDistance.value = 0
  totalDuration.value = 0
  totalFuelCost.value = 0
  totalTollCost.value = 0
  isUsingOSRM.value = false
  stopAnimation()
}

// Optimization routing execution
const calculateRoute = async () => {
  if (locations.value.length < 2) return
  
  isCalculating.value = true
  stopAnimation()
  
  try {
    // Configure Optimizer
    routeOptimizer.optimizationCriterion = optimizationCriterion.value
    routeOptimizer.fuelConsumptionLitersPer100Km = fuelConsumption.value
    routeOptimizer.fuelPricePerLiter = fuelPrice.value
    routeOptimizer.tollPricePerKm = tollPrice.value
    routeOptimizer.timeValuePerHour = timeValue.value
    
    await routeOptimizer.buildGraph(locations.value)
    
    const res = routeOptimizer.findOptimalRoute(locations.value[0].id)
    
    optimalRoute.value = res.path
    totalDistance.value = res.totalDistance
    totalDuration.value = res.totalDuration
    totalFuelCost.value = res.totalFuelCost
    totalTollCost.value = res.totalTollCost
    isUsingOSRM.value = routeOptimizer.isUsingOSRM
    
    // Freeze calculated criterion badge
    calculatedCriterion.value = optimizationCriterion.value
    
    // Sync state to URL upon successful calculation
    syncStateToUrl()
  } catch (e: any) {
    console.error('Calculation error:', e)
    showNotification(
      'error',
      t('notifications.calculation_error'),
      e?.message || t('notifications.calculation_error_desc')
    )
  } finally {
    isCalculating.value = false
  }
}

// Animation
const toggleAnimation = () => {
  if (isAnimating.value) {
    stopAnimation()
  } else {
    startAnimation()
  }
}

const startAnimation = () => {
  if (optimalRoute.value.length === 0) return
  isAnimating.value = true
  animationProgress.value = 0
  animationStartTime = performance.now()
  
  const animate = (timestamp: number) => {
    const duration = 4000 // duration of animation (4s)
    const elapsed = timestamp - animationStartTime
    const progress = Math.min(elapsed / duration, 1)
    
    animationProgress.value = progress
    
    if (progress < 1) {
      animationFrameId = requestAnimationFrame(animate)
    } else {
      isAnimating.value = false
    }
  }
  
  animationFrameId = requestAnimationFrame(animate)
}

const stopAnimation = () => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
  isAnimating.value = false
  animationProgress.value = 0
}

// Format duration
const formatDuration = (sec: number) => {
  const hours = Math.floor(sec / 3600)
  const minutes = Math.floor((sec % 3600) / 60)
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes} min`
}

// Export route data as JSON file download
const exportRouteData = () => {
  if (optimalRoute.value.length === 0) return
  
  const data = {
    locations: locations.value,
    totalDistance: totalDistance.value,
    totalDuration: totalDuration.value,
    totalFuelCost: totalFuelCost.value,
    totalTollCost: totalTollCost.value,
    optimalRoute: optimalRoute.value
  }
  
  const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(data, null, 2))
  const downloadAnchor = document.createElement('a')
  downloadAnchor.setAttribute("href", dataStr)
  downloadAnchor.setAttribute("download", `methil-route-${Date.now()}.json`)
  document.body.appendChild(downloadAnchor)
  downloadAnchor.click()
  downloadAnchor.remove()
}

// Share route URL link
const shareRoute = () => {
  // Generate link containing state parameters
  const currentUrl = new URL(window.location.href)
  navigator.clipboard.writeText(currentUrl.toString())
  
  showNotification(
    'success',
    t('notifications.link_copied'),
    t('notifications.link_copied_desc')
  )
}

// Observers
watch(fuelType, async (newType) => {
  if (isInitializing) return
  const price = await fetchAverageFuelPrice(newType)
  if (price !== null) {
    fuelPrice.value = parseFloat(price.toFixed(5))
  }
})

// Initialize
onMounted(() => {
  // Check saved theme or OS preference
  const savedTheme = localStorage.getItem('methil_maps_theme')
  if (savedTheme === 'dark' || savedTheme === 'light') {
    theme.value = savedTheme
    if (savedTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  } else {
    const isDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    if (isDark) {
      theme.value = 'dark'
      document.documentElement.classList.add('dark')
      localStorage.setItem('methil_maps_theme', 'dark')
    } else {
      theme.value = 'light'
      document.documentElement.classList.remove('dark')
      localStorage.setItem('methil_maps_theme', 'light')
    }
  }
  
  // Check saved locale or browser preference
  const savedLocale = localStorage.getItem('methil_maps_locale')
  if (savedLocale === 'en' || savedLocale === 'fr') {
    setLocale(savedLocale)
  } else {
    const lang = navigator.language || (navigator as any).userLanguage
    if (lang && lang.startsWith('en')) {
      setLocale('en')
    } else {
      setLocale('fr')
    }
  }

  loadStateFromUrl()
  
  if (!route.query.fp) {
    fetchAverageFuelPrice(fuelType.value)
      .then((price) => {
        if (price !== null) {
          fuelPrice.value = parseFloat(price.toFixed(5))
        }
      })
      .finally(() => {
        isInitializing = false
        // Auto-run computation if locations are fully loaded from URL parameters
        if (locations.value.length >= 2) {
          calculateRoute()
        }
      })
  } else {
    isInitializing = false
    // Auto-run computation if locations are fully loaded from URL parameters
    if (locations.value.length >= 2) {
      calculateRoute()
    }
  }
})
</script>

<style>
/* CSS animations and custom properties injected dynamically */
</style>
