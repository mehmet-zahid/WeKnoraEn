<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="settings-overlay">
        <div class="settings-modal">
          <!-- Close Button -->
          <button class="close-btn" @click="handleClose" :aria-label="$t('general.close')">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path d="M15 5L5 15M5 5L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>

          <div class="settings-container">
            <!-- Left Navigation -->
            <div class="settings-sidebar">
              <div class="sidebar-header">
                <h2 class="sidebar-title">{{ $t('general.settings') }}</h2>
              </div>
              <div class="settings-nav">
                <template v-for="(item, index) in navItems" :key="index">
                  <div 
                    :class="['nav-item', { 
                      'active': currentSection === (item as NavItem).key,
                      'has-submenu': (item as NavItem).children && (item as NavItem).children!.length > 0,
                      'expanded': expandedMenus.includes((item as NavItem).key)
                    }]"
                    @click="handleNavClick(item as NavItem)"
                  >
                    <!-- Web search uses custom SVG icon -->
                    <svg 
                      v-if="(item as NavItem).key === 'websearch'"
                      width="18" 
                      height="18" 
                      viewBox="0 0 18 18" 
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                      class="nav-icon"
                    >
                      <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="1.2" fill="none"/>
                      <path d="M 9 2 A 3.5 7 0 0 0 9 16" stroke="currentColor" stroke-width="1.2" fill="none"/>
                      <path d="M 9 2 A 3.5 7 0 0 1 9 16" stroke="currentColor" stroke-width="1.2" fill="none"/>
                      <line x1="2.94" y1="5.5" x2="15.06" y2="5.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/>
                      <line x1="2.94" y1="12.5" x2="15.06" y2="12.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/>
                    </svg>
                    <t-icon v-else :name="(item as NavItem).icon" class="nav-icon" />
                    <span class="nav-label">{{ (item as NavItem).label }}</span>
                    <t-icon 
                      v-if="(item as NavItem).children && (item as NavItem).children!.length > 0"
                      :name="expandedMenus.includes((item as NavItem).key) ? 'chevron-down' : 'chevron-right'"
                      class="expand-icon"
                    />
                  </div>
                  
                  <!-- Submenu -->
                  <Transition name="submenu">
                    <div 
                      v-if="(item as NavItem).children && expandedMenus.includes((item as NavItem).key)" 
                      class="submenu"
                    >
                      <div
                        v-for="(child, childIndex) in (item as NavItem).children"
                        :key="childIndex"
                        :class="['submenu-item', { 'active': currentSubSection === (child as NavItemChild).key }]"
                        @click.stop="handleSubMenuClick((item as NavItem).key, (child as NavItemChild).key)"
                      >
                        <span class="submenu-label">{{ (child as NavItemChild).label }}</span>
                      </div>
                    </div>
                  </Transition>
                </template>
              </div>
            </div>

            <!-- Right Content Area -->
            <div class="settings-content">
              <div class="content-wrapper">
                <!-- General Settings -->
                <div v-if="currentSection === 'general'" class="section">
                  <GeneralSettings />
                </div>

                <!-- Model Configuration -->
                <div v-if="currentSection === 'models'" class="section">
                  <ModelSettings />
                </div>

                <!-- Ollama Settings -->
                <div v-if="currentSection === 'ollama'" class="section">
                  <OllamaSettings />
                </div>

                <!-- Web Search Configuration -->
                <div v-if="currentSection === 'websearch'" class="section">
                  <WebSearchSettings />
                </div>

                <!-- System Information -->
                <div v-if="currentSection === 'system'" class="section">
                  <SystemInfo />
                </div>

                <!-- Tenant Information -->
                <div v-if="currentSection === 'tenant'" class="section">
                  <TenantInfo />
                </div>

                <!-- API Information -->
                <div v-if="currentSection === 'api'" class="section">
                  <ApiInfo />
                </div>

                <!-- MCP Service -->
                <div v-if="currentSection === 'mcp'" class="section">
                  <McpSettings />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { useI18n } from 'vue-i18n'
import SystemInfo from './SystemInfo.vue'
import TenantInfo from './TenantInfo.vue'
import ApiInfo from './ApiInfo.vue'
import GeneralSettings from './GeneralSettings.vue'
import ModelSettings from './ModelSettings.vue'
import OllamaSettings from './OllamaSettings.vue'
import McpSettings from './McpSettings.vue'
import WebSearchSettings from './WebSearchSettings.vue'

interface NavItemChild {
  key: string
  label: string
}

interface NavItem {
  key: string
  icon: string
  label: string
  children?: NavItemChild[]
}

const route = useRoute()
const router = useRouter()
const uiStore = useUIStore()
const { t } = useI18n()

const currentSection = ref<string>('general')
const currentSubSection = ref<string>('')
const expandedMenus = ref<string[]>([])

const navItems = computed<NavItem[]>(() => [
  { key: 'general', icon: 'setting', label: t('general.title') },
  { 
    key: 'models', 
    icon: 'control-platform', 
    label: t('settings.modelManagement'),
    children: [
      { key: 'chat', label: t('model.llmModel') },
      { key: 'embedding', label: t('model.embeddingModel') },
      { key: 'rerank', label: t('model.rerankModel') },
      { key: 'vllm', label: t('model.vlmModel') }
    ]
  },
  { key: 'ollama', icon: 'server', label: 'Ollama' },
  { key: 'websearch', icon: 'search', label: t('settings.webSearchConfig')  },
  { key: 'mcp', icon: 'tools', label: t('settings.mcpService') },
  { key: 'system', icon: 'info-circle', label: t('settings.systemSettings') },
  { key: 'tenant', icon: 'user-circle', label: t('settings.tenantInfo') },
  { key: 'api', icon: 'secured', label: t('settings.apiInfo') }
])

// Handle navigation item click
const handleNavClick = (item: NavItem) => {
  if (item.children && item.children.length > 0) {
    // Has submenu, toggle expand state
    const index = expandedMenus.value.indexOf(item.key)
    if (index > -1) {
      expandedMenus.value.splice(index, 1)
    } else {
      expandedMenus.value.push(item.key)
    }
    currentSubSection.value = item.children[0].key
  } else {
    currentSubSection.value = ''
  }
  
  // Switch to corresponding page
  currentSection.value = item.key
}

// Handle submenu click
const handleSubMenuClick = (parentKey: string, childKey: string) => {
  currentSection.value = parentKey
  currentSubSection.value = childKey
  
  // Scroll to corresponding model type area
  setTimeout(() => {
    const element = document.querySelector(`[data-model-type="${childKey}"]`)
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }, 100)
}

// Control modal display
const visible = computed(() => {
  return route.path === '/platform/settings' || uiStore.showSettingsModal
})

// Close modal
const handleClose = () => {
  uiStore.closeSettings()
  // If current route is settings page, go back
  if (route.path === '/platform/settings') {
    router.back()
  }
}

// Watch initial navigation settings
watch(() => uiStore.settingsInitialSection, (section) => {
  if (section && visible.value) {
    currentSection.value = section
    const navItem = navItems.value.find((item) => item.key === section)
    if (navItem && navItem.children && navItem.children.length > 0) {
      if (!expandedMenus.value.includes(section)) {
        expandedMenus.value.push(section)
      }
      currentSubSection.value = uiStore.settingsInitialSubSection || navItem.children[0].key
      if (uiStore.settingsInitialSubSection) {
        setTimeout(() => {
          const element = document.querySelector(`[data-model-type="${uiStore.settingsInitialSubSection}"]`)
          if (element) {
            element.scrollIntoView({ behavior: 'smooth', block: 'start' })
          }
        }, 300)
      }
    } else {
      currentSubSection.value = ''
    }
  }
}, { immediate: true })

// ESC key to close
const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && visible.value) {
    handleClose()
  }
}

// Handle quick navigation event
const handleSettingsNav = (e: CustomEvent) => {
  const { section, subsection } = e.detail
  if (section) {
    currentSection.value = section
    // If has submenu, auto expand
    const navItem = navItems.value.find((item) => item.key === section)
    if (navItem && navItem.children && navItem.children.length > 0) {
      if (!expandedMenus.value.includes(section)) {
        expandedMenus.value.push(section)
      }
      // If has subsection, select corresponding submenu item
      currentSubSection.value = subsection || navItem.children[0].key
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscape)
  window.addEventListener('settings-nav', handleSettingsNav as EventListener)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape)
  window.removeEventListener('settings-nav', handleSettingsNav as EventListener)
})
</script>

<style lang="less" scoped>
/* Overlay */
.settings-overlay {
  position: fixed;
  inset: 0;
  z-index: 1100;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  backdrop-filter: blur(4px);
}

/* Modal Container */
.settings-modal {
  position: relative;
  width: 100%;
  max-width: 900px;
  height: 700px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 6px 28px rgba(15, 23, 42, 0.08);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Close Button */
.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: #666666;
  cursor: pointer;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  z-index: 10;

  &:hover {
    background: #f5f5f5;
    color: #333333;
  }
}

.settings-container {
  display: flex;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

/* Left Navigation Bar */
.settings-sidebar {
  width: 220px;
  background-color: #f8f9fa;
  border-right: 1px solid #e5e7eb;
  flex-shrink: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 24px 16px 16px;
  border-bottom: 1px solid #e5e7eb;
}

.sidebar-title {
  font-size: 18px;
  font-weight: 600;
  color: #333333;
  margin: 0;
}

.settings-nav {
  padding: 16px 8px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  margin-bottom: 4px;
  border-radius: 6px;
  cursor: pointer;
  color: #666666;
  font-size: 14px;
  transition: all 0.2s ease;
  user-select: none;

  &:hover {
    background-color: #e8f5ed;
    color: #333333;
  }

  &.active {
    background-color: rgba(7, 192, 95, 0.1);
    color: #07C05F;
    font-weight: 500;
  }
}

.nav-icon {
  margin-right: 12px;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: inherit;
}

.nav-label {
  flex: 1;
}

.expand-icon {
  margin-left: 4px;
  font-size: 14px;
  transition: transform 0.2s ease;
}

/* Submenu */
.submenu {
  margin-left: 32px;
  margin-bottom: 4px;
  overflow: hidden;
}

.submenu-item {
  padding: 8px 16px;
  margin-bottom: 2px;
  border-radius: 4px;
  cursor: pointer;
  color: #666666;
  font-size: 13px;
  transition: all 0.2s ease;
  user-select: none;

  &:hover {
    background-color: #f5f7fa;
    color: #333333;
  }

  &.active {
    background-color: rgba(7, 192, 95, 0.08);
    color: #07C05F;
    font-weight: 500;
  }
}

.submenu-label {
  display: block;
}

/* Submenu Animation */
.submenu-enter-active,
.submenu-leave-active {
  transition: all 0.2s ease;
}

.submenu-enter-from {
  opacity: 0;
  max-height: 0;
}

.submenu-enter-to {
  opacity: 1;
  max-height: 300px;
}

.submenu-leave-from {
  opacity: 1;
  max-height: 300px;
}

.submenu-leave-to {
  opacity: 0;
  max-height: 0;
}

/* Right Content Area */
.settings-content {
  flex: 1;
  overflow-y: auto;
  background-color: #ffffff;
}

.content-wrapper {
  max-width: 600px;
  padding: 40px 48px;
}

.section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Modal Animation */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .settings-modal,
.modal-leave-active .settings-modal {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .settings-modal,
.modal-leave-to .settings-modal {
  transform: scale(0.95);
  opacity: 0;
}

/* Scrollbar Styles */
.settings-sidebar::-webkit-scrollbar,
.settings-content::-webkit-scrollbar {
  width: 6px;
}

.settings-sidebar::-webkit-scrollbar-track {
  background: #f8f9fa;
}

.settings-sidebar::-webkit-scrollbar-thumb {
  background: #d0d0d0;
  border-radius: 3px;
}

.settings-sidebar::-webkit-scrollbar-thumb:hover {
  background: #b0b0b0;
}

.settings-content::-webkit-scrollbar-track {
  background: #ffffff;
}

.settings-content::-webkit-scrollbar-thumb {
  background: #d0d0d0;
  border-radius: 3px;
}

.settings-content::-webkit-scrollbar-thumb:hover {
  background: #b0b0b0;
}
</style>

