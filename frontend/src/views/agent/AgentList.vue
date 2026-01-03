<template>
  <div class="agent-list-container">
    <!-- Header -->
    <div class="header">
      <div class="header-title">
        <h2>{{ $t('agent.title') }}</h2>
        <p class="header-subtitle">{{ $t('agent.subtitle') }}</p>
      </div>
    </div>
    <div class="header-divider"></div>

    <!-- Card grid -->
    <div v-if="agents.length > 0" class="agent-card-wrap">
      <div 
        v-for="agent in agents" 
        :key="agent.id" 
        class="agent-card"
        :class="{ 
          'is-builtin': agent.is_builtin,
          'agent-mode-normal': agent.config?.agent_mode === 'quick-answer',
          'agent-mode-agent': agent.config?.agent_mode === 'smart-reasoning'
        }"
        @click="handleCardClick(agent)"
      >
        <!-- Decorative star -->
        <div class="card-decoration">
          <svg class="star-icon" width="24" height="24" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 3L10.8 6.2C10.9 6.7 11.3 7.1 11.8 7.2L15 8L11.8 8.8C11.3 8.9 10.9 9.3 10.8 9.8L10 13L9.2 9.8C9.1 9.3 8.7 8.9 8.2 8.8L5 8L8.2 7.2C8.7 7.1 9.1 6.7 9.2 6.2L10 3Z" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" fill="currentColor" fill-opacity="0.15"/>
          </svg>
          <svg class="star-icon small" width="14" height="14" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 3L10.8 6.2C10.9 6.7 11.3 7.1 11.8 7.2L15 8L11.8 8.8C11.3 8.9 10.9 9.3 10.8 9.8L10 13L9.2 9.8C9.1 9.3 8.7 8.9 8.2 8.8L5 8L8.2 7.2C8.7 7.1 9.1 6.7 9.2 6.2L10 3Z" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" fill="currentColor" fill-opacity="0.15"/>
          </svg>
        </div>
        
        <!-- Card header -->
        <div class="card-header">
          <div class="card-header-left">
            <!-- Built-in agents use simple icon -->
            <div v-if="agent.is_builtin" class="builtin-avatar" :class="agent.config?.agent_mode === 'smart-reasoning' ? 'agent' : 'normal'">
              <t-icon :name="agent.config?.agent_mode === 'smart-reasoning' ? 'control-platform' : 'chat'" size="18px" />
            </div>
            <!-- Custom agents use AgentAvatar -->
            <AgentAvatar v-else :name="getLocalizedAgentName(agent)" size="medium" />
            <span class="card-title" :title="getLocalizedAgentName(agent)">{{ getLocalizedAgentName(agent) }}</span>
          </div>
          <t-popup 
            v-model="agent.showMore" 
            overlayClassName="card-more-popup"
            :on-visible-change="(visible: boolean) => onVisibleChange(visible, agent)"
            trigger="click" 
            destroy-on-close 
            placement="bottom-right"
          >
            <div 
              class="more-wrap" 
              @click.stop
              :class="{ 'active-more': agent.showMore }"
            >
              <img class="more-icon" src="@/assets/img/more.png" alt="" />
            </div>
            <template #content>
              <div class="popup-menu">
                <div class="popup-menu-item" @click="handleEdit(agent)">
                  <t-icon class="menu-icon" name="edit" />
                  <span>{{ $t('common.edit') }}</span>
                </div>
                <div class="popup-menu-item" @click="handleCopy(agent)">
                  <t-icon class="menu-icon" name="file-copy" />
                  <span>{{ $t('common.copy') }}</span>
                </div>
                <div v-if="!agent.is_builtin" class="popup-menu-item delete" @click="handleDelete(agent)">
                  <t-icon class="menu-icon" name="delete" />
                  <span>{{ $t('common.delete') }}</span>
                </div>
              </div>
            </template>
          </t-popup>
        </div>

        <!-- Card content -->
        <div class="card-content">
          <div class="card-description">
            {{ getLocalizedAgentDescription(agent) }}
          </div>
        </div>

        <!-- Card footer -->
        <div class="card-bottom" @click.stop>
          <div class="bottom-left">
            <div class="feature-badges">
              <t-tooltip :content="agent.config?.agent_mode === 'smart-reasoning' ? $t('agent.mode.agent') : $t('agent.mode.normal')" placement="top">
                <div class="feature-badge" :class="{ 'mode-normal': agent.config?.agent_mode === 'quick-answer', 'mode-agent': agent.config?.agent_mode === 'smart-reasoning' }">
                  <t-icon :name="agent.config?.agent_mode === 'smart-reasoning' ? 'control-platform' : 'chat'" size="14px" />
                </div>
              </t-tooltip>
              <t-tooltip v-if="agent.config?.web_search_enabled" :content="$t('agent.features.webSearch')" placement="top">
                <div class="feature-badge web-search">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
                    <ellipse cx="8" cy="8" rx="2.5" ry="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
                    <line x1="2" y1="6" x2="14" y2="6" stroke="currentColor" stroke-width="1.2"/>
                    <line x1="2" y1="10" x2="14" y2="10" stroke="currentColor" stroke-width="1.2"/>
                  </svg>
                </div>
              </t-tooltip>
              <t-tooltip v-if="agent.config?.knowledge_bases?.length || agent.config?.kb_selection_mode === 'all'" :content="$t('agent.features.knowledgeBase')" placement="top">
                <div class="feature-badge knowledge">
                  <t-icon name="folder" size="16px" />
                </div>
              </t-tooltip>
              <t-tooltip v-if="agent.config?.mcp_services?.length || agent.config?.mcp_selection_mode === 'all'" :content="$t('agent.features.mcp')" placement="top">
                <div class="feature-badge mcp">
                  <t-icon name="extension" size="16px" />
                </div>
              </t-tooltip>
              <t-tooltip v-if="agent.config?.multi_turn_enabled" :content="$t('agent.features.multiTurn')" placement="top">
                <div class="feature-badge multi-turn">
                  <t-icon name="chat-bubble" size="16px" />
                </div>
              </t-tooltip>
            </div>
          </div>
          <div v-if="agent.is_builtin" class="builtin-badge">
            <t-icon name="lock-on" size="12px" />
            <span>{{ $t('agent.builtin') }}</span>
          </div>
          <span v-else-if="agent.updated_at" class="card-time">{{ formatDate(agent.updated_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="!loading" class="empty-state">
      <img class="empty-img" src="@/assets/img/upload.svg" alt="">
      <span class="empty-txt">{{ $t('agent.empty.title') }}</span>
      <span class="empty-desc">{{ $t('agent.empty.description') }}</span>
    </div>

    <!-- Delete confirmation dialog -->
    <t-dialog 
      v-model:visible="deleteVisible" 
      dialogClassName="del-agent-dialog" 
      :closeBtn="false" 
      :cancelBtn="null"
      :confirmBtn="null"
    >
      <div class="circle-wrap">
        <div class="dialog-header">
          <img class="circle-img" src="@/assets/img/circle.png" alt="">
          <span class="circle-title">{{ $t('agent.delete.confirmTitle') }}</span>
        </div>
        <span class="del-circle-txt">
          {{ $t('agent.delete.confirmMessage', { name: deletingAgent?.name ?? '' }) }}
        </span>
        <div class="circle-btn">
          <span class="circle-btn-txt" @click="deleteVisible = false">{{ $t('common.cancel') }}</span>
          <span class="circle-btn-txt confirm" @click="confirmDelete">{{ $t('agent.delete.confirmButton') }}</span>
        </div>
      </div>
    </t-dialog>

    <!-- Agent editor modal -->
    <AgentEditorModal 
      :visible="editorVisible"
      :mode="editorMode"
      :agent="editingAgent"
      :initialSection="editorInitialSection"
      @update:visible="editorVisible = $event"
      @success="handleEditorSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, Icon as TIcon } from 'tdesign-vue-next'
import { listAgents, deleteAgent, copyAgent, type CustomAgent, BUILTIN_QUICK_ANSWER_ID, BUILTIN_SMART_REASONING_ID } from '@/api/agent'
import { formatStringDate } from '@/utils/index'
import { useI18n } from 'vue-i18n'
import AgentEditorModal from './AgentEditorModal.vue'
import AgentAvatar from '@/components/AgentAvatar.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

interface AgentWithUI extends CustomAgent {
  showMore?: boolean
}

const agents = ref<AgentWithUI[]>([])
const loading = ref(false)
const deleteVisible = ref(false)
const deletingAgent = ref<AgentWithUI | null>(null)
const editorVisible = ref(false)
const editorMode = ref<'create' | 'edit'>('create')
const editingAgent = ref<CustomAgent | null>(null)
const editorInitialSection = ref<string>('basic')

// Get localized name and description for built-in agents
const getLocalizedAgentName = (agent: AgentWithUI): string => {
  if (agent.id === BUILTIN_QUICK_ANSWER_ID) {
    return t('agent.type.normal')
  } else if (agent.id === BUILTIN_SMART_REASONING_ID) {
    return t('agent.type.agent')
  }
  return agent.name
}

const getLocalizedAgentDescription = (agent: AgentWithUI): string => {
  if (agent.id === BUILTIN_QUICK_ANSWER_ID) {
    return t('agent.capabilities.normal')
  } else if (agent.id === BUILTIN_SMART_REASONING_ID) {
    return t('agent.capabilities.agent')
  }
  return agent.description || t('agent.noDescription')
}

const fetchList = () => {
  loading.value = true
  return listAgents().then((res: any) => {
    const data = res.data || []
    // Display all agents (including built-in agents)
    agents.value = data.map((agent: CustomAgent) => ({
      ...agent,
      showMore: false
    }))
    
    // Check if there is an edit parameter in the URL, if so, open the edit modal for the corresponding agent
    checkAndOpenEditModal()
  }).finally(() => loading.value = false)
}

// Check URL parameters and open edit modal
const checkAndOpenEditModal = () => {
  const editId = route.query.edit as string
  const section = route.query.section as string
  if (editId) {
    const agent = agents.value.find(a => a.id === editId)
    if (agent) {
      editingAgent.value = agent
      editorMode.value = 'edit'
      editorInitialSection.value = section || 'basic'
      editorVisible.value = true
    }
    // Clear parameters from URL
    router.replace({ path: route.path, query: {} })
  }
}

// Listen for menu create agent event
const handleOpenAgentEditor = (event: CustomEvent) => {
  if (event.detail?.mode === 'create') {
    openCreateModal()
  }
}

onMounted(() => {
  fetchList()
  window.addEventListener('openAgentEditor', handleOpenAgentEditor as EventListener)
})

onUnmounted(() => {
  window.removeEventListener('openAgentEditor', handleOpenAgentEditor as EventListener)
})

const onVisibleChange = (visible: boolean, agent: AgentWithUI) => {
  if (!visible) {
    agent.showMore = false
  }
}

const handleCardClick = (agent: AgentWithUI) => {
  // If popup is showing, don't trigger edit
  if (agent.showMore) {
    return
  }
  // Click card to edit (including built-in agents)
  try {
    handleEdit(agent)
  } catch (error) {
    console.error('Error opening agent editor:', error)
    MessagePlugin.error(t('agent.messages.openFailed') || 'Failed to open agent editor')
  }
}

const handleEdit = (agent: AgentWithUI) => {
  try {
    agent.showMore = false
    // Ensure agent data is properly structured
    if (!agent || !agent.id) {
      console.error('Invalid agent data:', agent)
      MessagePlugin.error(t('agent.messages.invalidAgent') || 'Invalid agent data')
      return
    }
    editingAgent.value = agent
    editorMode.value = 'edit'
    editorVisible.value = true
  } catch (error) {
    console.error('Error in handleEdit:', error)
    MessagePlugin.error(t('agent.messages.openFailed') || 'Failed to open agent editor')
  }
}

const handleDelete = (agent: AgentWithUI) => {
  agent.showMore = false
  deletingAgent.value = agent
  deleteVisible.value = true
}

const handleCopy = (agent: AgentWithUI) => {
  agent.showMore = false
  copyAgent(agent.id).then((res: any) => {
    if (res.data) {
      MessagePlugin.success(t('agent.messages.copied'))
      fetchList()
    } else {
      MessagePlugin.error(res.message || t('agent.messages.copyFailed'))
    }
  }).catch((e: any) => {
    MessagePlugin.error(e?.message || t('agent.messages.copyFailed'))
  })
}

const confirmDelete = () => {
  if (!deletingAgent.value) return
  
  deleteAgent(deletingAgent.value.id).then((res: any) => {
    if (res.success) {
      MessagePlugin.success(t('agent.messages.deleted'))
      deleteVisible.value = false
      deletingAgent.value = null
      fetchList()
    } else {
      MessagePlugin.error(res.message || t('agent.messages.deleteFailed'))
    }
  }).catch((e: any) => {
    MessagePlugin.error(e?.message || t('agent.messages.deleteFailed'))
  })
}

const handleEditorSuccess = () => {
  editorVisible.value = false
  editingAgent.value = null
  fetchList()
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return formatStringDate(new Date(dateStr))
}

// Expose create method for external calls
const openCreateModal = () => {
  editingAgent.value = null
  editorMode.value = 'create'
  editorVisible.value = true
}

defineExpose({
  openCreateModal
})
</script>

<style scoped lang="less">
.agent-list-container {
  padding: 24px 44px;
  margin: 0 20px;
  height: calc(100vh);
  overflow-y: auto;
  box-sizing: border-box;
  flex: 1;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;

  .header-title {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  h2 {
    margin: 0;
    color: #000000e6;
    font-family: "PingFang SC";
    font-size: 24px;
    font-weight: 600;
    line-height: 32px;
  }
}

.header-subtitle {
  margin: 0;
  color: #00000099;
  font-family: "PingFang SC";
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
}

.header-divider {
  height: 1px;
  background: #e7ebf0;
  margin-bottom: 20px;
}

.agent-card-wrap {
  display: grid;
  gap: 16px;
  grid-template-columns: 1fr;
}

.agent-card {
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  overflow: hidden;
  box-sizing: border-box;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  background: #fff;
  position: relative;
  cursor: pointer;
  transition: all 0.25s ease;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  min-height: 180px;

  &:hover {
    border-color: #07c05f;
    box-shadow: 0 4px 12px rgba(7, 192, 95, 0.12);
  }

  // Normal mode styles
  &.agent-mode-normal {
    background: linear-gradient(135deg, #ffffff 0%, #f8fcfa 100%);
    border-color: #e8f5ed;

    &:hover {
      border-color: #07c05f;
      background: linear-gradient(135deg, #ffffff 0%, #f0fdf4 100%);
    }

    .card-decoration {
      color: rgba(7, 192, 95, 0.35);
    }

    &:hover .card-decoration {
      color: rgba(7, 192, 95, 0.5);
    }
  }

  // Agent mode styles
  &.agent-mode-agent {
    background: linear-gradient(135deg, #ffffff 0%, #f8f5ff 100%);
    border-color: #ede8ff;

    &:hover {
      border-color: #7c4dff;
      box-shadow: 0 4px 12px rgba(124, 77, 255, 0.12);
      background: linear-gradient(135deg, #ffffff 0%, #f3efff 100%);
    }

    .card-decoration {
      color: rgba(124, 77, 255, 0.35);
    }

    &:hover .card-decoration {
      color: rgba(124, 77, 255, 0.5);
    }
  }

  // Ensure content is above decoration
  .card-header,
  .card-content,
  .card-bottom {
    position: relative;
    z-index: 1;
  }
}

.card-decoration {
  position: absolute;
  top: 12px;
  right: 50px;
  display: flex;
  align-items: flex-start;
  gap: 4px;
  pointer-events: none;
  z-index: 0;
  transition: color 0.25s ease;
  
  .star-icon {
    opacity: 0.9;
    
    &.small {
      margin-top: 14px;
      opacity: 0.7;
    }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.card-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.card-title {
  color: #1a1a1a;
  font-family: "PingFang SC";
  font-size: 15px;
  font-weight: 600;
  line-height: 22px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.builtin-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
  color: #666;
  font-family: "PingFang SC";
  font-size: 11px;
  font-weight: 500;
  flex-shrink: 0;
}

.builtin-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  flex-shrink: 0;
  
  &.normal {
    background: linear-gradient(135deg, rgba(7, 192, 95, 0.15) 0%, rgba(7, 192, 95, 0.08) 100%);
    color: #059669;
  }
  
  &.agent {
    background: linear-gradient(135deg, rgba(124, 77, 255, 0.15) 0%, rgba(124, 77, 255, 0.08) 100%);
    color: #7c4dff;
  }
}

.edit-btn {
  display: flex;
  width: 32px;
  height: 32px;
  justify-content: center;
  align-items: center;
  border-radius: 6px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s ease;
  color: #00000066;

  &:hover {
    background: rgba(0, 0, 0, 0.06);
    color: #07c05f;
  }
}

.more-wrap {
  display: flex;
  width: 28px;
  height: 28px;
  justify-content: center;
  align-items: center;
  border-radius: 6px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s ease;
  opacity: 0;

  .agent-card:hover & {
    opacity: 0.6;
  }

  &:hover {
    background: rgba(0, 0, 0, 0.05);
    opacity: 1 !important;
  }

  &.active-more {
    background: rgba(0, 0, 0, 0.06);
    opacity: 1 !important;
  }

  .more-icon {
    width: 16px;
    height: 16px;
  }
}

.card-content {
  flex: 1;
  margin-bottom: 12px;
  overflow: visible;
  min-height: 40px;
}

.card-description {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  overflow: hidden;
  color: #666;
  font-family: "PingFang SC";
  font-size: 13px;
  font-weight: 400;
  line-height: 20px;
  word-break: break-word;
}

.card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.bottom-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.feature-badges {
  display: flex;
  align-items: center;
  gap: 4px;
}

.feature-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  cursor: default;
  transition: background 0.2s ease;

  &.mode-normal {
    background: rgba(7, 192, 95, 0.08);
    color: #059669;

    &:hover {
      background: rgba(7, 192, 95, 0.12);
    }
  }

  &.mode-agent {
    background: rgba(124, 77, 255, 0.08);
    color: #7c4dff;

    &:hover {
      background: rgba(124, 77, 255, 0.12);
    }
  }

  &.web-search {
    background: rgba(255, 152, 0, 0.08);
    color: #f59e0b;

    &:hover {
      background: rgba(255, 152, 0, 0.12);
    }
  }

  &.knowledge {
    background: rgba(7, 192, 95, 0.08);
    color: #059669;

    &:hover {
      background: rgba(7, 192, 95, 0.12);
    }
  }

  &.mcp {
    background: rgba(236, 72, 153, 0.08);
    color: #ec4899;

    &:hover {
      background: rgba(236, 72, 153, 0.12);
    }
  }

  &.multi-turn {
    background: rgba(59, 130, 246, 0.08);
    color: #3b82f6;

    &:hover {
      background: rgba(59, 130, 246, 0.12);
    }
  }
}

.card-time {
  color: #999;
  font-family: "PingFang SC";
  font-size: 12px;
  font-weight: 400;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px 20px;

  .empty-img {
    width: 162px;
    height: 162px;
    margin-bottom: 20px;
  }

  .empty-txt {
    color: #00000099;
    font-family: "PingFang SC";
    font-size: 16px;
    font-weight: 600;
    line-height: 26px;
    margin-bottom: 8px;
  }

  .empty-desc {
    color: #00000066;
    font-family: "PingFang SC";
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
  }
}

// Responsive layout
@media (min-width: 900px) {
  .agent-card-wrap {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1250px) {
  .agent-card-wrap {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1600px) {
  .agent-card-wrap {
    grid-template-columns: repeat(4, 1fr);
  }
}

// Delete confirmation dialog styles
:deep(.del-agent-dialog) {
  padding: 0px !important;
  border-radius: 6px !important;

  .t-dialog__header {
    display: none;
  }

  .t-dialog__body {
    padding: 16px;
  }

  .t-dialog__footer {
    padding: 0;
  }
}

:deep(.t-dialog__position.t-dialog--top) {
  padding-top: 40vh !important;
}

.circle-wrap {
  .dialog-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }

  .circle-img {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }

  .circle-title {
    color: #000000e6;
    font-family: "PingFang SC";
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  .del-circle-txt {
    color: #00000099;
    font-family: "PingFang SC";
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    display: inline-block;
    margin-left: 29px;
    margin-bottom: 21px;
  }

  .circle-btn {
    height: 22px;
    width: 100%;
    display: flex;
    justify-content: flex-end;
  }

  .circle-btn-txt {
    color: #000000e6;
    font-family: "PingFang SC";
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    cursor: pointer;

    &:hover {
      opacity: 0.8;
    }
  }

  .confirm {
    color: #FA5151;
    margin-left: 40px;

    &:hover {
      opacity: 0.8;
    }
  }
}
</style>

<style lang="less">
// More actions popup styles
.card-more-popup {
  z-index: 99 !important;

  .t-popup__content {
    padding: 6px 0 !important;
    margin-top: 6px !important;
    min-width: 140px;
    border-radius: 6px !important;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1) !important;
    border: 1px solid #e7ebf0 !important;
  }
}

.popup-menu {
  display: flex;
  flex-direction: column;
}

.popup-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #000000e6;
  font-family: "PingFang SC";
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;

  .menu-icon {
    font-size: 16px;
    flex-shrink: 0;
    color: #00000099;
    transition: color 0.2s ease;
  }

  &:hover {
    background: #f7f9fc;
    
    .menu-icon {
      color: #000000e6;
    }
  }

  &.delete {
    color: #000000e6;
    
    &:hover {
      background: #fff1f0;
      color: #fa5151;

      .menu-icon {
        color: #fa5151;
      }
    }
  }
}
</style>
