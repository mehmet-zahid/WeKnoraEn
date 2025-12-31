<template>
  <div class="websearch-settings">
    <div class="section-header">
      <h2>{{ t('webSearchSettings.title') }}</h2>
      <p class="section-description">{{ t('webSearchSettings.description') }}</p>
    </div>

    <div class="settings-group">
      <!-- Search Engine Provider -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.providerLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.providerDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-select
            v-model="localProvider"
            :loading="loadingProviders"
            filterable
            :placeholder="t('webSearchSettings.providerPlaceholder')"
            @change="handleProviderChange"
            @focus="loadProviders"
            style="width: 280px;"
          >
            <t-option
              v-for="provider in providers"
              :key="provider.id"
              :value="provider.id"
              :label="provider.name"
            >
              <div class="provider-option-wrapper">
                <div class="provider-option">
                  <span class="provider-name">{{ provider.name }}</span>
                </div>
              </div>
            </t-option>
          </t-select>
        </div>
      </div>

      <!-- API Key -->
      <div v-if="selectedProvider && selectedProvider.requires_api_key" class="setting-row">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.apiKeyLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.apiKeyDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-input
            v-model="localAPIKey"
            type="password"
            :placeholder="t('webSearchSettings.apiKeyPlaceholder')"
            @change="handleAPIKeyChange"
            style="width: 400px;"
            :show-password="true"
          />
        </div>
      </div>

      <!-- Max Results -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.maxResultsLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.maxResultsDescription') }}</p>
        </div>
        <div class="setting-control">
          <div class="slider-with-value">
            <t-slider 
              v-model="localMaxResults" 
              :min="1" 
              :max="50" 
              :step="1"
              :marks="{ 1: '1', 10: '10', 20: '20', 30: '30', 40: '40', 50: '50' }"
              @change="handleMaxResultsChange"
              style="width: 200px;"
            />
            <span class="value-display">{{ localMaxResults }}</span>
          </div>
        </div>
      </div>

      <!-- Include Date -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.includeDateLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.includeDateDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-switch
            v-model="localIncludeDate"
            @change="handleIncludeDateChange"
          />
        </div>
      </div>

      <!-- Compression Method -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.compressionLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.compressionDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-select
            v-model="localCompressionMethod"
            @change="handleCompressionMethodChange"
            style="width: 280px;"
          >
            <t-option value="none" :label="t('webSearchSettings.compressionNone')">
              {{ t('webSearchSettings.compressionNone') }}
            </t-option>
            <t-option value="llm_summary" :label="t('webSearchSettings.compressionSummary')">
              {{ t('webSearchSettings.compressionSummary') }}
            </t-option>
          </t-select>
        </div>
      </div>

      <!-- Blacklist -->
      <div class="setting-row vertical">
        <div class="setting-info">
          <label>{{ t('webSearchSettings.blacklistLabel') }}</label>
          <p class="desc">{{ t('webSearchSettings.blacklistDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-textarea
            v-model="localBlacklistText"
            :placeholder="t('webSearchSettings.blacklistPlaceholder')"
            :autosize="{ minRows: 4, maxRows: 8 }"
            @change="handleBlacklistChange"
            style="width: 500px;"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import { getWebSearchProviders, getTenantWebSearchConfig, updateTenantWebSearchConfig, type WebSearchProviderConfig, type WebSearchConfig } from '@/api/web-search'

const { t } = useI18n()

// Local state
const loadingProviders = ref(false)
const providers = ref<WebSearchProviderConfig[]>([])
const localProvider = ref<string>('')
const localAPIKey = ref<string>('')
const localMaxResults = ref<number>(5)
const localIncludeDate = ref<boolean>(true)
const localCompressionMethod = ref<string>('none')
const localBlacklistText = ref<string>('')
const isInitializing = ref(true) // Flag to track initialization, auto-save is disabled during initialization
const initialConfig = ref<WebSearchConfig | null>(null) // Save initial config for comparison

// Computed property: currently selected provider
const selectedProvider = computed(() => {
  return providers.value.find(p => p.id === localProvider.value)
})

// Load provider list
const loadProviders = async () => {
  if (providers.value.length > 0) {
    return // Already loaded
  }
  
  loadingProviders.value = true
  try {
    const response = await getWebSearchProviders()
    // Request interceptor has already processed the response, use data field directly
    if (response.data && Array.isArray(response.data)) {
      providers.value = response.data
    }
  } catch (error: any) {
    console.error('Failed to load web search providers:', error)
    const errorMessage = error?.message || t('webSearchSettings.errors.unknown')
    MessagePlugin.error(t('webSearchSettings.toasts.loadProvidersFailed', { message: errorMessage }))
  } finally {
    loadingProviders.value = false
  }
}

// Load tenant configuration
const loadTenantConfig = async () => {
  try {
    const response = await getTenantWebSearchConfig()
    // Request interceptor has already processed the response, use data field directly
    if (response.data) {
      const config = response.data
      // Disable auto-save when setting initial values
      isInitializing.value = true
      
      // Save a copy of initial config (for later comparison)
      const blacklist = (config.blacklist || []).join('\n')
      initialConfig.value = {
        provider: config.provider || '',
        api_key: config.api_key === '***' ? '***' : config.api_key || '',
        max_results: config.max_results || 5,
        include_date: config.include_date !== undefined ? config.include_date : true,
        compression_method: config.compression_method || 'none',
        blacklist: config.blacklist || []
      }
      
      // Set local state values
      localProvider.value = config.provider || ''
      // API key is hidden in response, if it's "***", it means configured but actual value not returned
      localAPIKey.value = config.api_key === '***' ? '***' : config.api_key || ''
      localMaxResults.value = config.max_results || 5
      localIncludeDate.value = config.include_date !== undefined ? config.include_date : true
      localCompressionMethod.value = config.compression_method || 'none'
      localBlacklistText.value = blacklist
      
      // Wait for all reactive updates to complete before enabling auto-save
      await nextTick()
      await nextTick()
      // Use setTimeout to ensure all events have been processed
      setTimeout(() => {
        isInitializing.value = false
      }, 100)
    } else {
      // If no config data, save default config
      initialConfig.value = {
        provider: '',
        api_key: '',
        max_results: 5,
        include_date: true,
        compression_method: 'none',
        blacklist: []
      }
      await nextTick()
      setTimeout(() => {
        isInitializing.value = false
      }, 100)
    }
  } catch (error: any) {
    console.error('Failed to load tenant web search config:', error)
    // If config doesn't exist, use default values (don't show error)
    initialConfig.value = {
      provider: '',
      api_key: '',
      max_results: 5,
      include_date: true,
      compression_method: 'none',
      blacklist: []
    }
    await nextTick()
    setTimeout(() => {
      isInitializing.value = false
    }, 100)
  }
}

// Check if configuration has changed
const hasConfigChanged = (): boolean => {
  if (!initialConfig.value) {
    return true // If no initial config, consider it changed
  }
  
  const blacklist = localBlacklistText.value
    .split('\n')
    .map(line => line.trim())
    .filter(line => line.length > 0)
  
  const currentConfig: WebSearchConfig = {
    provider: localProvider.value,
    api_key: localAPIKey.value,
    max_results: localMaxResults.value,
    include_date: localIncludeDate.value,
    compression_method: localCompressionMethod.value,
    blacklist: blacklist
  }
  
  // Compare if config has changed (ignore '***' placeholder for API key)
  const initial = initialConfig.value
  if (currentConfig.provider !== initial.provider) return true
  if (currentConfig.api_key !== initial.api_key && 
      !(currentConfig.api_key === '***' && initial.api_key === '***')) return true
  if (currentConfig.max_results !== initial.max_results) return true
  if (currentConfig.include_date !== initial.include_date) return true
  if (currentConfig.compression_method !== initial.compression_method) return true
  
  // Compare blacklist arrays
  const currentBlacklist = blacklist.sort().join(',')
  const initialBlacklist = (initial.blacklist || []).sort().join(',')
  if (currentBlacklist !== initialBlacklist) return true
  
  return false
}

// Save configuration
const saveConfig = async () => {
  // If config hasn't changed, don't save
  if (!hasConfigChanged()) {
    return
  }
  
  try {
    const blacklist = localBlacklistText.value
      .split('\n')
      .map(line => line.trim())
      .filter(line => line.length > 0)
    
    const config: WebSearchConfig = {
      provider: localProvider.value,
      api_key: localAPIKey.value,
      max_results: localMaxResults.value,
      include_date: localIncludeDate.value,
      compression_method: localCompressionMethod.value,
      blacklist: blacklist
    }
    
    await updateTenantWebSearchConfig(config)
    
    // Update initial config to avoid duplicate saves
    initialConfig.value = {
      provider: config.provider,
      api_key: config.api_key,
      max_results: config.max_results,
      include_date: config.include_date,
      compression_method: config.compression_method,
      blacklist: [...config.blacklist]
    }
    
    MessagePlugin.success(t('webSearchSettings.toasts.saveSuccess'))
  } catch (error: any) {
    console.error('Failed to save web search config:', error)
    const errorMessage = error?.message || t('webSearchSettings.errors.unknown')
    MessagePlugin.error(t('webSearchSettings.toasts.saveFailed', { message: errorMessage }))
    throw error
  }
}

// Debounced save
let saveTimer: number | null = null
const debouncedSave = () => {
  // Don't trigger auto-save during initialization
  if (isInitializing.value) {
    return
  }
  if (saveTimer) {
    clearTimeout(saveTimer)
  }
  saveTimer = window.setTimeout(() => {
    saveConfig().catch(() => {
      // Error already handled in saveConfig
    })
  }, 500)
}

// Handle changes
const handleProviderChange = () => {
  debouncedSave()
}

const handleAPIKeyChange = () => {
  debouncedSave()
}

const handleMaxResultsChange = () => {
  debouncedSave()
}

const handleIncludeDateChange = () => {
  debouncedSave()
}

const handleCompressionMethodChange = () => {
  debouncedSave()
}

const handleBlacklistChange = () => {
  debouncedSave()
}

// Initialize
onMounted(async () => {
  isInitializing.value = true
  await loadProviders()
  await loadTenantConfig()
  // loadTenantConfig already handles isInitializing internally, no need to set it here
})
</script>

<style lang="less" scoped>
.websearch-settings {
  width: 100%;
}

.section-header {
  margin-bottom: 32px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #333333;
    margin: 0 0 8px 0;
  }

  .section-description {
    font-size: 14px;
    color: #666666;
    margin: 0;
    line-height: 1.5;
  }
}

.settings-group {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.setting-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid #e5e7eb;

  &:last-child {
    border-bottom: none;
  }

  &.vertical {
    flex-direction: column;
    gap: 12px;

    .setting-control {
      width: 100%;
      max-width: 100%;
    }
  }
}

.setting-info {
  flex: 1;
  max-width: 65%;
  padding-right: 24px;

  label {
    font-size: 15px;
    font-weight: 500;
    color: #333333;
    display: block;
    margin-bottom: 4px;
  }

  .desc {
    font-size: 13px;
    color: #666666;
    margin: 0;
    line-height: 1.5;
  }
}

.setting-control {
  flex-shrink: 0;
  min-width: 280px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.slider-with-value {
  display: flex;
  align-items: center;
  gap: 12px;
}

.value-display {
  min-width: 40px;
  text-align: right;
  font-size: 14px;
  font-weight: 500;
  color: #333333;
}

.provider-option-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 2px 0;
}

.provider-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  flex-wrap: wrap;
}

.provider-name {
  font-weight: 500;
  font-size: 14px;
  color: #333;
  flex-shrink: 0;
}

.provider-tags {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.provider-desc {
  font-size: 12px;
  color: #999;
  line-height: 1.4;
  margin-top: 2px;
}

/* Fix dropdown item description overlapping with entries: make options support multi-line adaptive height */
:deep(.t-select-option) {
  height: auto;
  align-items: flex-start;
  padding-top: 6px;
  padding-bottom: 6px;
}

:deep(.t-select-option__content) {
  white-space: normal;
}

</style>
<style lang="less">
.t-select__dropdown .t-select-option {
  height: auto;
  align-items: flex-start;
  padding-top: 6px;
  padding-bottom: 6px;
}
.t-select__dropdown .t-select-option__content {
  white-space: normal;
}
.t-select__dropdown .provider-option-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 2px 0;
}
</style>

