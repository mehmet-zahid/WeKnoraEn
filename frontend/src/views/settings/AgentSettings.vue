<template>
  <div class="agent-settings">
    <div v-if="activeSection === 'modes'">
      <div class="section-header">
        <h2>{{ $t('settings.conversationStrategy') }}</h2>
        <p class="section-description">{{ $t('conversationSettings.description') }}</p>
        <div class="global-config-notice">
          <t-icon name="info-circle" />
          <span>{{ $t('agentSettings.globalConfigNotice') }}</span>
        </div>
      </div>

      <t-tabs v-model="activeTab" class="conversation-tabs">
      <!-- Agent Mode Settings Tab -->
      <t-tab-panel value="agent" :label="$t('conversationSettings.agentMode')">
        <div class="tab-content">
          <!-- Agent Status Display -->
          <div class="agent-status-row">
        <div class="status-label">
          <label>{{ $t('agentSettings.status.label') }}</label>
        </div>
        <div class="status-control">
          <div class="status-badge" :class="{ ready: isAgentReady }">
            <t-icon 
              v-if="isAgentReady" 
              name="check-circle-filled" 
              class="status-icon"
            />
            <t-icon 
              v-else 
              name="error-circle-filled" 
              class="status-icon"
            />
            <span class="status-text">
              {{ isAgentReady ? $t('agentSettings.status.ready') : $t('agentSettings.status.notReady') }}
            </span>
          </div>
          <span v-if="!isAgentReady" class="status-hint">
            {{ agentStatusMessage }}
            <t-link v-if="needsModelConfig" @click="handleGoToModelSettings" theme="primary">
              {{ $t('agentSettings.status.goConfigureModels') }}
            </t-link>
          </span>
          <p v-if="!isAgentReady" class="status-tip">
            <t-icon name="info-circle" class="tip-icon" />
            {{ $t('agentSettings.status.hint') }}
          </p>
        </div>
      </div>

          <!-- Model Recommendation Hint -->
          <div class="model-recommendation-box">
            <div class="recommendation-header">
              <t-icon name="info-circle" class="recommendation-icon" />
              <span class="recommendation-title">{{ $t('agentSettings.modelRecommendation.title') }}</span>
            </div>
            <div class="recommendation-content">
              <p>{{ $t('agentSettings.modelRecommendation.content') }}</p>
            </div>
          </div>

          <div class="settings-group">

      <!-- Max Iterations -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ $t('agentSettings.maxIterations.label') }}</label>
          <p class="desc">{{ $t('agentSettings.maxIterations.desc') }}</p>
        </div>
        <div class="setting-control">
          <div class="slider-with-value">
          <t-slider 
            v-model="localMaxIterations" 
            :min="1" 
            :max="30" 
            :step="1"
            :marks="{ 1: '1', 5: '5', 10: '10', 15: '15', 20: '20', 25: '25', 30: '30' }"
            @change="handleMaxIterationsChangeDebounced"
              style="width: 200px;"
          />
            <span class="value-display">{{ localMaxIterations }}</span>
          </div>
        </div>
      </div>

      <!-- Temperature Parameter -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ $t('agentSettings.temperature.label') }}</label>
          <p class="desc">{{ $t('agentSettings.temperature.desc') }}</p>
        </div>
        <div class="setting-control">
          <div class="slider-with-value">
          <t-slider 
            v-model="localTemperature" 
            :min="0" 
            :max="1" 
            :step="0.1"
            :marks="{ 0: '0', 0.5: '0.5', 1: '1' }"
            @change="handleTemperatureChange"
              style="width: 200px;"
          />
            <span class="value-display">{{ localTemperature.toFixed(1) }}</span>
          </div>
        </div>
      </div>

      <!-- Allowed Tools -->
      <div class="setting-row vertical">
        <div class="setting-info">
          <label>{{ $t('agentSettings.allowedTools.label') }}</label>
          <p class="desc">{{ $t('agentSettings.allowedTools.desc') }}</p>
        </div>
        <div class="setting-control full-width allowed-tools-display">
          <div v-if="displayAllowedTools.length" class="allowed-tool-list">
            <div
              v-for="tool in displayAllowedTools"
              :key="tool.name"
              class="allowed-tool-chip"
            >
              <span class="allowed-tool-label">{{ tool.label }}</span>
              <span
                v-if="tool.description"
                class="allowed-tool-desc"
              >
                {{ tool.description }}
              </span>
            </div>
          </div>
          <p v-else class="allowed-tools-empty">
            {{ $t('agentSettings.allowedTools.empty') }}
          </p>
        </div>
      </div>

      <!-- System Prompt -->
      <div class="setting-row vertical">
        <div class="setting-info">
          <label>{{ $t('agentSettings.systemPrompt.label') }}</label>
          <p class="desc">{{ $t('agentSettings.systemPrompt.desc') }}</p>
          <div class="placeholder-hint">
            <p class="hint-title">{{ $t('agentSettings.systemPrompt.availablePlaceholders') }}</p>
            <ul class="placeholder-list">
              <li v-for="placeholder in availablePlaceholders" :key="placeholder.name">
                <code v-html="`{{${placeholder.name}}}`"></code> - {{ placeholder.label }} ({{ placeholder.description }})
              </li>
            </ul>
            <p class="hint-tip">{{ $t('agentSettings.systemPrompt.hintPrefix') }} <code>&#123;&#123;</code> {{ $t('agentSettings.systemPrompt.hintSuffix') }}</p>
          </div>
        </div>
        <div class="setting-control full-width" style="position: relative;">
          <div class="prompt-header">
            <t-button
              theme="default"
              variant="outline"
              size="small"
              @click="handleResetToDefault"
              :loading="isResettingPrompt"
            >
              {{ $t('common.resetToDefault') }}
            </t-button>
          </div>
          <p class="prompt-tab-hint">
            {{ $t('agentSettings.systemPrompt.tabHint') }} (Leave empty to use system default, use {{web_search_status}} placeholder to dynamically control web search behavior)
          </p>
          <div class="system-prompt-tabs">
            <div class="prompt-textarea-wrapper textarea-with-template">
              <t-textarea
                ref="promptTextareaRef"
                v-model="localSystemPrompt"
                :autosize="{ minRows: 15, maxRows: 30 }"
                :placeholder="$t('agentSettings.systemPrompt.placeholder')"
                @blur="handleSystemPromptChange"
                @input="handlePromptInput"
                @keydown="handlePromptKeydown"
                style="width: 100%; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; font-size: 13px;"
              />
              <PromptTemplateSelector 
                type="systemPrompt" 
                position="corner"
                :hasKnowledgeBase="true"
                @select="handleAgentSystemPromptTemplateSelect"
              />
            </div>
          </div>
          <!-- Placeholder Hint Dropdown -->
          <teleport to="body">
            <div
              v-if="showPlaceholderPopup && filteredPlaceholders.length > 0"
              class="placeholder-popup-wrapper"
              :style="popupStyle"
            >
              <div class="placeholder-popup">
              <div
                v-for="(placeholder, index) in filteredPlaceholders"
                :key="placeholder.name"
                class="placeholder-item"
                :class="{ active: selectedPlaceholderIndex === index }"
                @mousedown.prevent="insertPlaceholder(placeholder.name)"
                @mouseenter="selectedPlaceholderIndex = index"
              >
                  <div class="placeholder-name">
                    <code v-html="`{{${placeholder.name}}}`"></code>
                  </div>
                  <div class="placeholder-desc">{{ placeholder.description }}</div>
                </div>
              </div>
            </div>
          </teleport>
        </div>
      </div>
        </div>
      </div>
      </t-tab-panel>

      <!-- Normal Mode Settings Tab -->
      <t-tab-panel value="normal" :label="$t('conversationSettings.normalMode')">
        <div class="tab-content">
          <div class="settings-group">
            <!-- System Prompt (Normal Mode, Custom Toggle) -->
            <div class="setting-row vertical">
              <div class="setting-info">
                <label>{{ $t('conversationSettings.systemPrompt.label') }}</label>
                <p class="desc">{{ $t('conversationSettings.systemPrompt.desc') }} (Leave empty to use system default)</p>
              </div>
              <div class="setting-control full-width">
                <div class="prompt-textarea-wrapper textarea-with-template">
                  <t-textarea
                    v-model="localSystemPromptNormal"
                    :autosize="{ minRows: 10, maxRows: 20 }"
                    :placeholder="$t('conversationSettings.systemPrompt.placeholder')"
                    @blur="handleSystemPromptNormalChange"
                    style="width: 100%; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; font-size: 13px;"
                  />
                  <PromptTemplateSelector 
                    type="systemPrompt" 
                    position="corner"
                    :hasKnowledgeBase="true"
                    @select="handleNormalSystemPromptTemplateSelect"
                  />
                </div>
              </div>
            </div>

            <!-- Context Template (Normal Mode) -->
            <div class="setting-row vertical">
              <div class="setting-info">
                <label>{{ $t('conversationSettings.contextTemplate.label') }}</label>
                <p class="desc">{{ $t('conversationSettings.contextTemplate.desc') }} (Leave empty to use system default)</p>
              </div>
              <div class="setting-control full-width">
                <div class="prompt-textarea-wrapper textarea-with-template">
                  <t-textarea
                    v-model="localContextTemplate"
                    :autosize="{ minRows: 15, maxRows: 30 }"
                    :placeholder="$t('conversationSettings.contextTemplate.placeholder')"
                    @blur="handleContextTemplateChange"
                    style="width: 100%; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; font-size: 13px;"
                  />
                  <PromptTemplateSelector 
                    type="contextTemplate" 
                    position="corner"
                    :hasKnowledgeBase="true"
                    @select="handleContextTemplateTemplateSelect"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </t-tab-panel>
    </t-tabs>
    </div>

    <div v-else-if="activeSection === 'models'" class="section-block" data-conversation-section="models">
      <div class="section-header">
        <h2>{{ $t('conversationSettings.menus.models') }}</h2>
        <p class="section-description">{{ $t('conversationSettings.models.description') }}</p>
      </div>

      <div class="settings-group">
        <!-- Default LLM (Chat/Summary Model) -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.models.chatGroupLabel') }}</label>
            <p class="desc">{{ $t('conversationSettings.models.chatGroupDesc') }}</p>
          </div>
          <div class="setting-control">
            <t-select
              v-model="localSummaryModelId"
              :loading="loadingModels"
              filterable
              :placeholder="$t('conversationSettings.models.chatModel.placeholder')"
              style="width: 320px;"
              @focus="loadAllModels"
              @change="handleConversationSummaryModelChange"
            >
              <t-option
                v-for="model in chatModels"
                :key="model.id"
                :value="model.id"
                :label="model.name"
              />
              <t-option value="__add_model__" class="add-model-option">
                <div class="model-option add">
                  <t-icon name="add" class="add-icon" />
                  <span class="model-name">{{ $t('agentSettings.model.addChat') }}</span>
                </div>
              </t-option>
            </t-select>
          </div>
        </div>

        <!-- Default ReRank Model -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.models.rerankGroupLabel') }}</label>
            <p class="desc">{{ $t('conversationSettings.models.rerankGroupDesc') }}</p>
          </div>
          <div class="setting-control">
            <t-select
              v-model="localConversationRerankModelId"
              :loading="loadingModels"
              filterable
              :placeholder="$t('conversationSettings.models.rerankModel.placeholder')"
              style="width: 320px;"
              @focus="loadAllModels"
              @change="handleConversationRerankModelChange"
            >
              <t-option
                v-for="model in rerankModels"
                :key="model.id"
                :value="model.id"
                :label="model.name"
              />
              <t-option value="__add_model__" class="add-model-option">
                <div class="model-option add">
                  <t-icon name="add" class="add-icon" />
                  <span class="model-name">{{ $t('agentSettings.model.addRerank') }}</span>
                </div>
              </t-option>
            </t-select>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="activeSection === 'thresholds'" class="section-block">
      <div class="section-header">
        <h2>{{ $t('conversationSettings.menus.thresholds') }}</h2>
        <p class="section-description">{{ $t('conversationSettings.thresholds.description') }}</p>
      </div>

      <div class="settings-group">
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.maxRounds.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.maxRounds.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-input-number
              v-model="localMaxRounds"
              :min="1"
              :max="50"
              @change="handleMaxRoundsChange"
            />
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.embeddingTopK.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.embeddingTopK.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-input-number
              v-model="localEmbeddingTopK"
              :min="1"
              :max="50"
              @change="handleEmbeddingTopKChange"
            />
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.keywordThreshold.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.keywordThreshold.desc') }}</p>
          </div>
          <div class="setting-control slider-with-value">
            <t-slider
              v-model="localKeywordThreshold"
              :min="0"
              :max="1"
              :step="0.05"
              style="width: 240px;"
              @change="handleKeywordThresholdChange"
            />
            <span class="value-display">{{ localKeywordThreshold.toFixed(2) }}</span>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.vectorThreshold.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.vectorThreshold.desc') }}</p>
          </div>
          <div class="setting-control slider-with-value">
            <t-slider
              v-model="localVectorThreshold"
              :min="0"
              :max="1"
              :step="0.05"
              style="width: 240px;"
              @change="handleVectorThresholdChange"
            />
            <span class="value-display">{{ localVectorThreshold.toFixed(2) }}</span>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.rerankTopK.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.rerankTopK.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-input-number
              v-model="localRerankTopK"
              :min="1"
              :max="20"
              @change="handleRerankTopKChange"
            />
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.rerankThreshold.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.rerankThreshold.desc') }}</p>
          </div>
          <div class="setting-control slider-with-value">
            <t-slider
              v-model="localRerankThreshold"
              :min="0"
              :max="1"
              :step="0.05"
              style="width: 240px;"
              @change="handleRerankThresholdChange"
            />
            <span class="value-display">{{ localRerankThreshold.toFixed(2) }}</span>
          </div>
        </div>

      </div>
    </div>

    <div v-else-if="activeSection === 'advanced'" class="section-block">
      <div class="section-header">
        <h2>{{ $t('conversationSettings.menus.advanced') }}</h2>
        <p class="section-description">{{ $t('conversationSettings.advanced.description') }}</p>
      </div>

      <div class="settings-group">
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.enableQueryExpansion.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.enableQueryExpansion.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-switch
              v-model="localEnableQueryExpansion"
              :label="[$t('common.off'), $t('common.on')]"
              @change="handleEnableQueryExpansionChange"
            />
          </div>
        </div>
        <!-- Enable Question Rewrite -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.enableRewrite.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.enableRewrite.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-switch
              v-model="localEnableRewrite"
              :label="[$t('common.off'), $t('common.on')]"
              @change="handleEnableRewriteChange"
            />
          </div>
        </div>

        <!-- Rewrite Prompt: Only shown when rewrite is enabled -->
        <div v-if="localEnableRewrite" class="setting-row vertical">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.rewritePrompt.system') }}</label>
            <p class="desc">{{ $t('conversationSettings.rewritePrompt.desc') }}</p>
          </div>
          <div class="setting-control full-width">
            <div class="textarea-with-template">
              <t-textarea
                v-model="localRewritePromptSystem"
                :autosize="{ minRows: 8, maxRows: 16 }"
                @blur="handleRewritePromptSystemChange"
              />
              <PromptTemplateSelector 
                type="rewriteSystem" 
                position="corner"
                @select="handleRewriteSystemTemplateSelect"
              />
            </div>
          </div>
        </div>

        <div v-if="localEnableRewrite" class="setting-row vertical">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.rewritePrompt.user') }}</label>
            <p class="desc">{{ $t('conversationSettings.rewritePrompt.userDesc') }}</p>
          </div>
          <div class="setting-control full-width">
            <div class="textarea-with-template">
              <t-textarea
                v-model="localRewritePromptUser"
                :autosize="{ minRows: 8, maxRows: 16 }"
                @blur="handleRewritePromptUserChange"
              />
              <PromptTemplateSelector 
                type="rewriteUser" 
                position="corner"
                @select="handleRewriteUserTemplateSelect"
              />
            </div>
          </div>
        </div>

        <!-- Fallback Strategy -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.fallbackStrategy.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.fallbackStrategy.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-radio-group v-model="localFallbackStrategy" @change="handleFallbackStrategyChange">
              <t-radio value="fixed">{{ $t('conversationSettings.fallbackStrategy.fixed') }}</t-radio>
              <t-radio value="model">{{ $t('conversationSettings.fallbackStrategy.model') }}</t-radio>
            </t-radio-group>
          </div>
        </div>

        <!-- Fixed Fallback Reply: Only shown when fixed reply is selected -->
        <div v-if="localFallbackStrategy === 'fixed'" class="setting-row vertical">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.fallbackResponse.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.fallbackResponse.desc') }}</p>
          </div>
          <div class="setting-control full-width">
            <div class="textarea-with-template">
              <t-textarea
                v-model="localFallbackResponse"
                :autosize="{ minRows: 3, maxRows: 6 }"
                @blur="handleFallbackResponseChange"
              />
              <PromptTemplateSelector 
                type="fallback" 
                position="corner"
                @select="handleFallbackResponseTemplateSelect"
              />
            </div>
          </div>
        </div>

        <!-- Fallback Prompt: Only shown when "Let model continue generation" is selected -->
        <div v-else-if="localFallbackStrategy === 'model'" class="setting-row vertical">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.fallbackPrompt.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.fallbackPrompt.desc') }}</p>
          </div>
          <div class="setting-control full-width">
            <div class="textarea-with-template">
              <t-textarea
                v-model="localFallbackPrompt"
                :autosize="{ minRows: 8, maxRows: 16 }"
                @blur="handleFallbackPromptChange"
              />
              <PromptTemplateSelector 
                type="fallback" 
                position="corner"
                @select="handleFallbackPromptTemplateSelect"
              />
            </div>
          </div>
        </div>

        <!-- Normal mode generation parameter: Temperature -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.temperature.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.temperature.desc') }}</p>
          </div>
          <div class="setting-control">
            <div class="slider-with-value">
              <t-slider 
                v-model="localTemperatureNormal" 
                :min="0" 
                :max="1" 
                :step="0.1"
                :marks="{ 0: '0', 0.5: '0.5', 1: '1' }"
                @change="handleTemperatureNormalChange"
                style="width: 200px;"
              />
              <span class="value-display">{{ localTemperatureNormal.toFixed(1) }}</span>
            </div>
          </div>
        </div>

        <!-- Normal mode generation parameter: Max Tokens -->
        <div class="setting-row">
          <div class="setting-info">
            <label>{{ $t('conversationSettings.maxTokens.label') }}</label>
            <p class="desc">{{ $t('conversationSettings.maxTokens.desc') }}</p>
          </div>
          <div class="setting-control">
            <t-input-number
              v-model="localMaxCompletionTokens"
              :min="1"
              :max="100000"
              :step="100"
              @change="handleMaxCompletionTokensChange"
              style="width: 200px;"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, nextTick } from 'vue'
import type { Ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSettingsStore } from '@/stores/settings'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import { listModels, type ModelConfig } from '@/api/model'
import { getAgentConfig, updateAgentConfig, getConversationConfig, updateConversationConfig, type AgentConfig, type ConversationConfig, type ToolDefinition, type PlaceholderDefinition } from '@/api/system'
import PromptTemplateSelector from '@/components/PromptTemplateSelector.vue'

const props = defineProps<{
  // Submenu key from external settings modal: 'modes' | 'models' | 'thresholds' | 'advanced'
  activeSubSection?: string
}>()

// Current subpage (modes, models, thresholds, advanced)
const activeSection = computed(() => props.activeSubSection || 'modes')

const settingsStore = useSettingsStore()
const router = useRouter()
const { t } = useI18n()

// Tab state
const activeTab = ref('agent')

const getDefaultConversationConfig = (): ConversationConfig => ({
  prompt: '',
  context_template: '',
  temperature: 0.3,
  max_completion_tokens: 2048,
  max_rounds: 5,
  embedding_top_k: 10,
  keyword_threshold: 0.3,
  vector_threshold: 0.5,
  rerank_top_k: 5,
  rerank_threshold: 0.5,
  enable_rewrite: true,
  enable_query_expansion: true,
  fallback_strategy: 'fixed',
  fallback_response: '',
  fallback_prompt: '',
  summary_model_id: '',
  rerank_model_id: '',
  rewrite_prompt_system: '',
  rewrite_prompt_user: '',
})

const normalizeConversationConfig = (config?: Partial<ConversationConfig>): ConversationConfig => ({
  ...getDefaultConversationConfig(),
  ...config,
})

const conversationConfig = ref<ConversationConfig>(getDefaultConversationConfig())
const conversationConfigLoaded = ref(false)
const conversationSaving = ref(false)

// Agent mode local state
const localMaxIterations = ref(5)
const localTemperature = ref(0.7)
const localAllowedTools = ref<string[]>([])

// Unified system prompt
const localSystemPrompt = ref('')
let savedSystemPrompt = ''

// Normal mode local state
const localContextTemplate = ref('')
const localSystemPromptNormal = ref('')
const localTemperatureNormal = ref(0.3)
const localMaxCompletionTokens = ref(2048)
let savedContextTemplate = ''
let savedSystemPromptNormal = ''
let savedTemperatureNormal = 0.3
let savedMaxCompletionTokens = 2048

const localMaxRounds = ref(5)
const localEmbeddingTopK = ref(10)
const localKeywordThreshold = ref(0.3)
const localVectorThreshold = ref(0.5)
const localRerankTopK = ref(5)
const localRerankThreshold = ref(0.5)
const localEnableRewrite = ref(true)
const localEnableQueryExpansion = ref(true)
const localFallbackStrategy = ref<'fixed' | 'model'>('fixed')
const localFallbackResponse = ref('')
const localFallbackPrompt = ref('')
const localRewritePromptSystem = ref('')
const localRewritePromptUser = ref('')
const localSummaryModelId = ref('')
const localConversationRerankModelId = ref('')

const syncConversationLocals = () => {
  const cfg = conversationConfig.value
  localContextTemplate.value = cfg.context_template ?? ''
  savedContextTemplate = localContextTemplate.value
  localSystemPromptNormal.value = cfg.prompt ?? ''
  savedSystemPromptNormal = localSystemPromptNormal.value
  localTemperatureNormal.value = cfg.temperature ?? 0.3
  savedTemperatureNormal = localTemperatureNormal.value
  localMaxCompletionTokens.value = cfg.max_completion_tokens ?? 2048
  savedMaxCompletionTokens = localMaxCompletionTokens.value

  localMaxRounds.value = cfg.max_rounds ?? 5
  localEmbeddingTopK.value = cfg.embedding_top_k ?? 10
  localKeywordThreshold.value = cfg.keyword_threshold ?? 0.3
  localVectorThreshold.value = cfg.vector_threshold ?? 0.5
  localRerankTopK.value = cfg.rerank_top_k ?? 5
  localRerankThreshold.value = cfg.rerank_threshold ?? 0.5
  localEnableRewrite.value = cfg.enable_rewrite ?? true
  localEnableQueryExpansion.value = cfg.enable_query_expansion ?? true
  localFallbackStrategy.value = (cfg.fallback_strategy as 'fixed' | 'model') || 'fixed'
  localFallbackResponse.value = cfg.fallback_response ?? ''
  localFallbackPrompt.value = cfg.fallback_prompt ?? ''
  localRewritePromptSystem.value = cfg.rewrite_prompt_system ?? ''
  localRewritePromptUser.value = cfg.rewrite_prompt_user ?? ''
  localSummaryModelId.value = cfg.summary_model_id ?? ''
  localConversationRerankModelId.value = cfg.rerank_model_id ?? ''

  settingsStore.updateConversationModels({
    summaryModelId: localSummaryModelId.value || '',
    rerankModelId: localConversationRerankModelId.value || '',
  })
}

const saveConversationConfig = async (partial: Partial<ConversationConfig>, toastMessage?: string) => {
  if (!conversationConfigLoaded.value) return

  const payload = normalizeConversationConfig({
    ...conversationConfig.value,
    ...partial,
  })

  try {
    conversationSaving.value = true
    const res = await updateConversationConfig(payload)
    conversationConfig.value = normalizeConversationConfig(res.data ?? payload)
    syncConversationLocals()
    if (toastMessage) {
      MessagePlugin.success(toastMessage)
    }
  } catch (error) {
    console.error('Failed to save conversation config:', error)
    MessagePlugin.error(getErrorMessage(error))
    throw error
  } finally {
    conversationSaving.value = false
  }
}

// Calculate if Agent is ready
const isAgentReady = computed(() => {
  return (
    localAllowedTools.value.length > 0 &&
    localSummaryModelId.value &&
    localSummaryModelId.value.trim() !== '' &&
    localConversationRerankModelId.value &&
    localConversationRerankModelId.value.trim() !== ''
  )
})

const buildAgentConfigPayload = (overrides: Partial<AgentConfig> = {}): AgentConfig => ({
  max_iterations: localMaxIterations.value,
  reflection_enabled: false,
  allowed_tools: localAllowedTools.value,
  temperature: localTemperature.value,
  system_prompt: localSystemPrompt.value,
  ...overrides,
})

// Whether model configuration is missing
const needsModelConfig = computed(() => {
  return (
    (!localSummaryModelId.value || localSummaryModelId.value.trim() === '') ||
    (!localConversationRerankModelId.value || localConversationRerankModelId.value.trim() === '')
  )
})

// Agent status message
const agentStatusMessage = computed(() => {
  const missing: string[] = []
  
  if (localAllowedTools.value.length === 0) {
    missing.push(t('agentSettings.status.missingAllowedTools'))
  }
  
  if (!localSummaryModelId.value || localSummaryModelId.value.trim() === '') {
    missing.push(t('agentSettings.status.missingSummaryModel'))
  }
  
  if (!localConversationRerankModelId.value || localConversationRerankModelId.value.trim() === '') {
    missing.push(t('agentSettings.status.missingRerankModel'))
  }
  
  if (missing.length === 0) {
    return ''
  }
  
  return t('agentSettings.status.pleaseConfigure', { items: missing.join('、') })
})

// Navigate to model configuration
const handleGoToModelSettings = () => {
  router.push('/platform/settings')

  setTimeout(() => {
    const event = new CustomEvent('settings-nav', {
      detail: { section: 'agent', subsection: 'models' }
    })
    window.dispatchEvent(event)

    setTimeout(() => {
      const sectionEl = document.querySelector('[data-conversation-section="models"]')
      if (sectionEl) {
        sectionEl.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }, 150)
  }, 100)
}

// Model list state
const chatModels = ref<ModelConfig[]>([])
const rerankModels = ref<ModelConfig[]>([])
const loadingModels = ref(false)

// Available tools list
const availableTools = ref<ToolDefinition[]>([])
// Available placeholders list
const availablePlaceholders = ref<PlaceholderDefinition[]>([])
const displayAllowedTools = computed(() => {
  return localAllowedTools.value.map(name => {
    const detail = availableTools.value.find(tool => tool.name === name)
    return {
      name,
      label: detail?.label || name,
      description: detail?.description || ''
    }
  })
})

// Configuration loading state
const loadingConfig = ref(false)
const configLoaded = ref(false) // Prevent duplicate loading
const isInitializing = ref(true) // Flag to track initialization, prevent save during initialization

// Reset to default Prompt loading state
const isResettingPrompt = ref(false)

// Placeholder hint related state
const promptTextareaRef = ref<any>(null)
const showPlaceholderPopup = ref(false)
const selectedPlaceholderIndex = ref(0)
let placeholderPopupTimer: any = null
const placeholderPrefix = ref('') // Current input prefix for filtering
const popupStyle = ref({ top: '0px', left: '0px' }) // Popup position

// Setup textarea native event listeners
const setupTextareaEventListeners = () => {
  nextTick(() => {
    const textarea = getTextareaElement()
    if (textarea) {
      // Add native keydown event listener (using capture phase to ensure priority handling)
      textarea.addEventListener('keydown', (e: KeyboardEvent) => {
        // If placeholder hint is showing, prioritize placeholder-related keys
        if (showPlaceholderPopup.value && filteredPlaceholders.value.length > 0) {
          if (e.key === 'ArrowDown') {
            // Arrow down: select next
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()
            if (selectedPlaceholderIndex.value < filteredPlaceholders.value.length - 1) {
              selectedPlaceholderIndex.value++
            } else {
              selectedPlaceholderIndex.value = 0 // Loop to first
            }
            return
          } else if (e.key === 'ArrowUp') {
            // Arrow up: select previous
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()
            if (selectedPlaceholderIndex.value > 0) {
              selectedPlaceholderIndex.value--
            } else {
              selectedPlaceholderIndex.value = filteredPlaceholders.value.length - 1 // Loop to last
            }
            return
          } else if (e.key === 'Enter') {
            // Enter key: insert selected placeholder
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()
            const selected = filteredPlaceholders.value[selectedPlaceholderIndex.value]
            if (selected) {
              insertPlaceholder(selected.name)
            }
            return
          } else if (e.key === 'Escape') {
            // ESC key: close hint
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()
            showPlaceholderPopup.value = false
            placeholderPrefix.value = ''
            return
          }
        }
        
        // If { key is pressed
        if (e.key === '{') {
          // Clear previous timer
          if (placeholderPopupTimer) {
            clearTimeout(placeholderPopupTimer)
          }
          
          // Delay check, wait for input to complete (two consecutive {)
          placeholderPopupTimer = setTimeout(() => {
            checkAndShowPlaceholderPopup()
          }, 150)
        }
      }, true) // Use capture phase
      
      // Add native input event listener (as backup)
      textarea.addEventListener('input', () => {
        if (placeholderPopupTimer) {
          clearTimeout(placeholderPopupTimer)
        }
        placeholderPopupTimer = setTimeout(() => {
          checkAndShowPlaceholderPopup()
        }, 50)
      })
    }
  })
}

// Helper function to get textarea element
const getTextareaElement = (): HTMLTextAreaElement | null => {
  if (promptTextareaRef.value) {
    if (promptTextareaRef.value.$el) {
      return promptTextareaRef.value.$el.querySelector('textarea')
    } else if (promptTextareaRef.value instanceof HTMLTextAreaElement) {
      return promptTextareaRef.value
    }
  }
  
  // If still not found, try to find via DOM
  const wrapper = document.querySelector('.setting-control.full-width')
  return wrapper?.querySelector('textarea') || null
}

// Initialize and load
onMounted(async () => {
  // Prevent duplicate loading
  if (configLoaded.value) return
  
  loadingConfig.value = true
  configLoaded.value = true
  isInitializing.value = true
  
  try {
    // Load config from backend
    const res = await getAgentConfig()
    const config = res.data
    
    // Update local state (won't trigger save during initialization)
    localMaxIterations.value = config.max_iterations
    lastSavedValue = config.max_iterations // Record saved value during initialization
    localTemperature.value = config.temperature
    localAllowedTools.value = config.allowed_tools || []
    const systemPrompt = config.system_prompt || ''
    localSystemPrompt.value = systemPrompt
    savedSystemPrompt = systemPrompt
    availableTools.value = config.available_tools || []
    availablePlaceholders.value = config.available_placeholders || []
    
    // Debug info
    console.log('Loaded placeholder list:', availablePlaceholders.value)
    
    // Unified load all models (only call API once)
      await loadAllModels()
    
    // Sync to store (only update local storage, don't trigger API save)
    // Note: Don't automatically set isAgentEnabled, keep user's previous choice
    // enabled state should be manually controlled by user, not automatically set based on config
    settingsStore.updateAgentConfig({
      maxIterations: config.max_iterations,
      temperature: config.temperature,
      allowedTools: config.allowed_tools || [],
      system_prompt: systemPrompt,
    })

    // Load normal mode configuration
    if (!conversationConfigLoaded.value) {
      try {
        const convRes = await getConversationConfig()
        conversationConfig.value = normalizeConversationConfig(convRes.data)
        conversationConfigLoaded.value = true
        syncConversationLocals()
      } catch (error) {
        console.error('Failed to load normal mode configuration:', error)
        // Use default values
        conversationConfigLoaded.value = true
      }
    }
    
    // Wait for next tick to ensure all reactive updates complete
    await nextTick()
    // Wait one more frame to ensure all event listeners are set up
    requestAnimationFrame(() => {
      // Initialization complete, now allow save operations
      isInitializing.value = false
      
      // Setup native event listeners (as backup)
      setupTextareaEventListeners()
    })
  } catch (error) {
    console.error('Failed to load Agent config:', error)
    MessagePlugin.error(t('agentSettings.toasts.loadFailed'))
    configLoaded.value = false // Reset flag on load failure to allow retry
    
    // On failure, load from store
    localMaxIterations.value = settingsStore.agentConfig.maxIterations
    localTemperature.value = settingsStore.agentConfig.temperature
  } finally {
    loadingConfig.value = false
    isInitializing.value = false // Ensure initialization completes, even on failure allow subsequent operations
  }
})

// Error code to error message mapping
const getErrorMessage = (error: any): string => {
  const errorCode = error?.response?.data?.error?.code
  const errorMessage = error?.response?.data?.error?.message
  
  switch (errorCode) {
    case 2100:
      return t('agentSettings.errors.selectThinkingModel')
    case 2101:
      return t('agentSettings.errors.selectAtLeastOneTool')
    case 2102:
      return t('agentSettings.errors.iterationsRange')
    case 2103:
      return t('agentSettings.errors.temperatureRange')
    case 1010:
      return errorMessage || t('agentSettings.errors.validationFailed')
    default:
      return errorMessage || t('common.saveFailed')
  }
}

// Debounce timer
let maxIterationsDebounceTimer: any = null
// Last saved value, used to avoid saving duplicate values
let lastSavedValue: number | null = null

// Handle max iterations change (debounced version, used for both click and drag)
const handleMaxIterationsChangeDebounced = (value: number) => {
  // If initializing, don't trigger save
  if (isInitializing.value) return
  
  // Ensure value is a number
  const numValue = typeof value === 'number' ? value : Number(value)
  if (isNaN(numValue)) {
    console.error('Invalid max_iterations value:', value)
    return
  }
  
  // If value hasn't changed, don't save
  if (lastSavedValue === numValue) {
    return
  }
  
  // Clear previous timer
  if (maxIterationsDebounceTimer) {
    clearTimeout(maxIterationsDebounceTimer)
}

  // Set new timer, save after 300ms (reduce delay, improve responsiveness)
  maxIterationsDebounceTimer = setTimeout(async () => {
    // Check again if value has changed (may have changed during wait)
    if (lastSavedValue === numValue) {
      maxIterationsDebounceTimer = null
      return
    }
  
  try {
    const config = buildAgentConfigPayload({ max_iterations: numValue })
    await updateAgentConfig(config)
      settingsStore.updateAgentConfig({ maxIterations: numValue })
      lastSavedValue = numValue // Record saved value
    MessagePlugin.success(t('agentSettings.toasts.iterationsSaved'))
  } catch (error) {
    console.error('Failed to save:', error)
    MessagePlugin.error(getErrorMessage(error))
    } finally {
      maxIterationsDebounceTimer = null
  }
  }, 300)
}

// Unified load all models (only call API once)
const loadAllModels = async () => {
  if (chatModels.value.length > 0 && rerankModels.value.length > 0) return // Already loaded
  
  loadingModels.value = true
  try {
    const allModels = await listModels()
    // Filter by type to avoid duplicate calls
    chatModels.value = allModels.filter(m => m.type === 'KnowledgeQA')
    rerankModels.value = allModels.filter(m => m.type === 'Rerank')
  } catch (error) {
    console.error('Failed to load model list:', error)
    MessagePlugin.error('Failed to load model list')
  } finally {
    loadingModels.value = false
  }
}

// Load chat model list (deprecated, use loadAllModels)
const loadChatModels = async () => {
  await loadAllModels()
}

// Load Rerank model list (deprecated, use loadAllModels)
const loadRerankModels = async () => {
  await loadAllModels()
}

// Handle temperature parameter change
const handleTemperatureChange = async (value: number) => {
  // If initializing, don't trigger save
  if (isInitializing.value) return
  
  try {
    const config = buildAgentConfigPayload({ temperature: value })
    await updateAgentConfig(config)
    settingsStore.updateAgentConfig({ temperature: value })
    MessagePlugin.success(t('agentSettings.toasts.temperatureSaved'))
  } catch (error) {
    console.error('Failed to save:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

// Handle system Prompt keyboard events (as backup, main logic in native event listeners)
const handlePromptKeydown = (e: KeyboardEvent) => {
  // If placeholder hint is showing and input is letter, number or underscore, update filter in real-time
  if (showPlaceholderPopup.value && /^[a-zA-Z0-9_]$/.test(e.key)) {
    // Delay check, wait for character input to complete
    if (placeholderPopupTimer) {
      clearTimeout(placeholderPopupTimer)
    }
    placeholderPopupTimer = setTimeout(() => {
      checkAndShowPlaceholderPopup()
    }, 50)
  }
}

// Filtered placeholder list (based on prefix match)
const filteredPlaceholders = computed(() => {
  if (!placeholderPrefix.value) {
    return availablePlaceholders.value
  }
  
  const prefix = placeholderPrefix.value.toLowerCase()
  return availablePlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  )
})

// Calculate cursor position in pixels within textarea
const calculateCursorPosition = (textarea: HTMLTextAreaElement) => {
  const cursorPos = textarea.selectionStart
  const activePromptValue = getActivePromptRef().value
  const textBeforeCursor = activePromptValue.substring(0, cursorPos)
  
  // Get textarea style and position
  const style = window.getComputedStyle(textarea)
  const textareaRect = textarea.getBoundingClientRect()
  
  // Calculate line number and current line text
  const lines = textBeforeCursor.split('\n')
  const currentLine = lines.length - 1
  const lineText = lines[currentLine] || ''
  
  // Get line height
  const lineHeight = parseFloat(style.lineHeight) || parseFloat(style.fontSize) * 1.2
  
  // Get padding
  const paddingTop = parseFloat(style.paddingTop) || 0
  const paddingLeft = parseFloat(style.paddingLeft) || 0
  
  // Use canvas to measure current line text width (more accurate)
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  let textWidth = 0
  
  if (context) {
    context.font = `${style.fontSize} ${style.fontFamily}`
    textWidth = context.measureText(lineText).width
  } else {
    // Fallback: use monospace font estimation (Monaco/Menlo are monospace fonts)
    const charWidth = parseFloat(style.fontSize) * 0.6 // Monospace character width is approximately 0.6 times font size
    textWidth = lineText.length * charWidth
  }
  
  // Calculate cursor position top (considering scroll)
  const scrollTop = textarea.scrollTop
  const top = textareaRect.top + paddingTop + (currentLine * lineHeight) - scrollTop + lineHeight + 4
  
  // Calculate cursor position left (considering scroll)
  const scrollLeft = textarea.scrollLeft
  const left = textareaRect.left + paddingLeft + textWidth - scrollLeft
  
  return { top, left }
}

// Check and show placeholder hint
const checkAndShowPlaceholderPopup = () => {
  const textarea = getTextareaElement()
  
  if (!textarea) {
    return
  }
  
  const cursorPos = textarea.selectionStart
  const textBeforeCursor = getActivePromptRef().value.substring(0, cursorPos)
  
  // Check if {{ was entered (find nearest {{ before cursor position)
  // Need to find the nearest {{ before cursor, with no }} in between
  let lastOpenPos = -1
  for (let i = cursorPos - 1; i >= 0; i--) {
    if (i > 0 && textBeforeCursor[i - 1] === '{' && textBeforeCursor[i] === '{') {
      // Found {{
      const textAfterOpen = textBeforeCursor.substring(i + 1)
      // Check if }} is already included (placeholder is complete)
      if (!textAfterOpen.includes('}}')) {
        lastOpenPos = i - 1
        break
      }
    }
  }
  
  if (lastOpenPos === -1) {
    // No valid {{ found, hide hint
    showPlaceholderPopup.value = false
    placeholderPrefix.value = ''
    return
  }
  
  // Get content from {{ to cursor position as prefix
  const textAfterOpen = textBeforeCursor.substring(lastOpenPos + 2)
  
  // Update prefix
  placeholderPrefix.value = textAfterOpen
  
  // Filter placeholders based on prefix
  const filtered = filteredPlaceholders.value
  
  if (filtered.length > 0) {
    // Has matching placeholders, show hint
    // Calculate cursor position
    nextTick(() => {
      const position = calculateCursorPosition(textarea)
      popupStyle.value = {
        top: `${position.top}px`,
        left: `${position.left}px`
      }
      showPlaceholderPopup.value = true
      // Reset selected index to first (default select first)
      selectedPlaceholderIndex.value = 0
    })
  } else {
    // No matching placeholders, hide hint
    showPlaceholderPopup.value = false
  }
}

// Handle system Prompt input
const handlePromptInput = () => {
  // Clear previous timer
  if (placeholderPopupTimer) {
    clearTimeout(placeholderPopupTimer)
  }
  
  // Delay check to avoid frequent triggers
  placeholderPopupTimer = setTimeout(() => {
    checkAndShowPlaceholderPopup()
  }, 50)
}

// Insert placeholder
const insertPlaceholder = (placeholderName: string) => {
  const textarea = getTextareaElement()
  if (!textarea) {
    return
  }
  
  // Close hint first to avoid triggering blur event
  showPlaceholderPopup.value = false
  placeholderPrefix.value = ''
  selectedPlaceholderIndex.value = 0
  
  // Delay execution to ensure hint box is closed
  nextTick(() => {
    const cursorPos = textarea.selectionStart
    const promptRef = getActivePromptRef()
    const currentValue = promptRef.value
    const textBeforeCursor = currentValue.substring(0, cursorPos)
    const textAfterCursor = currentValue.substring(cursorPos)
    
    // Find position of last {{
    const lastOpenPos = textBeforeCursor.lastIndexOf('{{')
    if (lastOpenPos === -1) {
      // If {{ not found, directly insert complete placeholder
      const placeholder = `{{${placeholderName}}}`
      promptRef.value = textBeforeCursor + placeholder + textAfterCursor
      // Set cursor position
      nextTick(() => {
        const newPos = cursorPos + placeholder.length
        textarea.setSelectionRange(newPos, newPos)
        textarea.focus()
      })
    } else {
      // Replace content from {{ to cursor position with complete placeholder
      const beforePlaceholder = textBeforeCursor.substring(0, lastOpenPos)
      const placeholder = `{{${placeholderName}}}`
      promptRef.value = beforePlaceholder + placeholder + textAfterCursor
      // Set cursor position
      nextTick(() => {
        const newPos = lastOpenPos + placeholder.length
        textarea.setSelectionRange(newPos, newPos)
        textarea.focus()
      })
    }
  })
}

// Reset to default Prompt
const handleResetToDefault = async () => {
  const confirmDialog = DialogPlugin.confirm({
    header: t('agentSettings.reset.header'),
    body: t('agentSettings.reset.body'),
    confirmBtn: t('common.confirm'),
    cancelBtn: t('common.cancel'),
    onConfirm: async () => {
      try {
        isResettingPrompt.value = true
        
        // Get default value by setting system_prompt to empty string
        // Backend returns default value when field is empty
        const tempConfig = buildAgentConfigPayload({
          system_prompt: '',
        })
        
        await updateAgentConfig(tempConfig)
        
        // Reload config to get complete default Prompt content
        const res = await getAgentConfig()
        const defaultPrompt = res.data.system_prompt || ''
        
        // Set to default Prompt content
        localSystemPrompt.value = defaultPrompt
        savedSystemPrompt = defaultPrompt
        
        MessagePlugin.success(t('agentSettings.toasts.resetToDefault'))
        confirmDialog.hide()
      } catch (error) {
        console.error('Failed to reset to default Prompt:', error)
        MessagePlugin.error(getErrorMessage(error))
      } finally {
        isResettingPrompt.value = false
      }
    }
  })
}

// Handle system Prompt change
const handleSystemPromptChange = async (e?: FocusEvent) => {
  // If clicked on placeholder hint box, don't trigger save
  if (e?.relatedTarget) {
    const target = e.relatedTarget as HTMLElement
    if (target.closest('.placeholder-popup-wrapper')) {
      return
    }
  }
  
  // Delay check to avoid immediate trigger when clicking placeholder
  await nextTick()
  
  // If placeholder hint box is still showing, user clicked placeholder, don't trigger save
  if (showPlaceholderPopup.value) {
    return
  }
  
  // Hide placeholder hint
  placeholderPrefix.value = ''
  
  // If initializing, don't trigger save
  if (isInitializing.value) return

  // Check if content has changed
  if (localSystemPrompt.value === savedSystemPrompt) {
    return // Content unchanged, don't call API
  }
  
  try {
    const config = buildAgentConfigPayload()
    await updateAgentConfig(config)
    savedSystemPrompt = localSystemPrompt.value // Update saved value
    MessagePlugin.success(t('agentSettings.toasts.systemPromptSaved'))
  } catch (error) {
    console.error('Failed to save system Prompt:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

// Watch Agent ready state changes, sync to store
watch(isAgentReady, (newValue, oldValue) => {
  if (!isInitializing.value) {
    // If config changes from "ready" to "not ready", and Agent is currently enabled, auto-disable
    if (!newValue && oldValue && settingsStore.isAgentEnabled) {
      settingsStore.toggleAgent(false)
      MessagePlugin.warning(t('agentSettings.toasts.autoDisabled'))
    }
    // Note: When config changes from "not ready" to "ready", don't auto-enable (let user decide)
  }
})

// Normal mode config handler functions
const handleContextTemplateChange = async () => {
  if (!conversationConfigLoaded.value) return
  
  if (localContextTemplate.value === savedContextTemplate) {
    return
  }
  
  try {
    await saveConversationConfig(
      {
        context_template: localContextTemplate.value,
      },
      t('conversationSettings.toasts.contextTemplateSaved')
    )
    savedContextTemplate = localContextTemplate.value
  } catch (error) {
    console.error('Failed to save Context Template:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

const reloadConversationConfig = async () => {
  const convRes = await getConversationConfig()
  conversationConfig.value = normalizeConversationConfig(convRes.data)
  syncConversationLocals()
}

const handleSystemPromptNormalChange = async () => {
  if (!conversationConfigLoaded.value) return
  
  if (localSystemPromptNormal.value === savedSystemPromptNormal) {
    return
  }
  
  try {
    await saveConversationConfig(
      {
        prompt: localSystemPromptNormal.value,
      },
      t('conversationSettings.toasts.systemPromptSaved')
    )
    savedSystemPromptNormal = localSystemPromptNormal.value
  } catch (error) {
    console.error('Failed to save System Prompt:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

const handleTemperatureNormalChange = async (value: number) => {
  if (!conversationConfigLoaded.value) return
  if (value === savedTemperatureNormal) return
  
  try {
    await saveConversationConfig(
      { temperature: value },
      t('conversationSettings.toasts.temperatureSaved')
    )
    savedTemperatureNormal = value
  } catch (error) {
    console.error('Failed to save Temperature:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

const handleMaxCompletionTokensChange = async (value: number) => {
  if (!conversationConfigLoaded.value) return
  
  try {
    await saveConversationConfig(
      { max_completion_tokens: value },
      t('conversationSettings.toasts.maxTokensSaved')
    )
    savedMaxCompletionTokens = value
  } catch (error) {
    console.error('Failed to save Max Tokens:', error)
    MessagePlugin.error(getErrorMessage(error))
  }
}

const handleMaxRoundsChange = async (value: number) => {
  try {
    await saveConversationConfig({ max_rounds: value }, t('conversationSettings.toasts.maxRoundsSaved'))
  } catch (error) {
    console.error('Failed to save max_rounds:', error)
    localMaxRounds.value = conversationConfig.value.max_rounds
  }
}

const handleEmbeddingTopKChange = async (value: number) => {
  try {
    await saveConversationConfig({ embedding_top_k: value }, t('conversationSettings.toasts.embeddingSaved'))
  } catch (error) {
    console.error('Failed to save embedding_top_k:', error)
    localEmbeddingTopK.value = conversationConfig.value.embedding_top_k
  }
}

const handleKeywordThresholdChange = async (value: number) => {
  try {
    await saveConversationConfig({ keyword_threshold: value }, t('conversationSettings.toasts.keywordThresholdSaved'))
  } catch (error) {
    console.error('Failed to save keyword_threshold:', error)
    localKeywordThreshold.value = conversationConfig.value.keyword_threshold
  }
}

const handleVectorThresholdChange = async (value: number) => {
  try {
    await saveConversationConfig({ vector_threshold: value }, t('conversationSettings.toasts.vectorThresholdSaved'))
  } catch (error) {
    console.error('Failed to save vector_threshold:', error)
    localVectorThreshold.value = conversationConfig.value.vector_threshold
  }
}

const handleRerankTopKChange = async (value: number) => {
  try {
    await saveConversationConfig({ rerank_top_k: value }, t('conversationSettings.toasts.rerankTopKSaved'))
  } catch (error) {
    console.error('Failed to save rerank_top_k:', error)
    localRerankTopK.value = conversationConfig.value.rerank_top_k
  }
}

const handleRerankThresholdChange = async (value: number) => {
  try {
    await saveConversationConfig({ rerank_threshold: value }, t('conversationSettings.toasts.rerankThresholdSaved'))
  } catch (error) {
    console.error('Failed to save rerank_threshold:', error)
    localRerankThreshold.value = conversationConfig.value.rerank_threshold
  }
}

const handleEnableRewriteChange = async (value: boolean) => {
  try {
    await saveConversationConfig({ enable_rewrite: value }, t('conversationSettings.toasts.enableRewriteSaved'))
  } catch (error) {
    console.error('Failed to save enable_rewrite:', error)
    localEnableRewrite.value = conversationConfig.value.enable_rewrite
  }
}

const handleEnableQueryExpansionChange = async (value: boolean) => {
  try {
    await saveConversationConfig(
      { enable_query_expansion: value },
      t('conversationSettings.toasts.enableQueryExpansionSaved')
    )
  } catch (error) {
    console.error('Failed to save enable_query_expansion:', error)
    localEnableQueryExpansion.value = conversationConfig.value.enable_query_expansion ?? true
  }
}

const handleFallbackStrategyChange = async (value: 'fixed' | 'model') => {
  try {
    await saveConversationConfig({ fallback_strategy: value }, t('conversationSettings.toasts.fallbackStrategySaved'))
  } catch (error) {
    console.error('Failed to save fallback_strategy:', error)
    localFallbackStrategy.value = (conversationConfig.value.fallback_strategy as 'fixed' | 'model') || 'fixed'
  }
}

const handleFallbackResponseChange = async () => {
  if (localFallbackResponse.value === (conversationConfig.value.fallback_response ?? '')) return
  try {
    await saveConversationConfig({ fallback_response: localFallbackResponse.value }, t('conversationSettings.toasts.fallbackResponseSaved'))
  } catch (error) {
    console.error('Failed to save fallback_response:', error)
    localFallbackResponse.value = conversationConfig.value.fallback_response ?? ''
  }
}

const handleRewritePromptSystemChange = async () => {
  if (localRewritePromptSystem.value === (conversationConfig.value.rewrite_prompt_system ?? '')) return
  try {
    await saveConversationConfig({ rewrite_prompt_system: localRewritePromptSystem.value }, t('conversationSettings.toasts.rewritePromptSystemSaved'))
  } catch (error) {
    console.error('Failed to save rewrite_prompt_system:', error)
    localRewritePromptSystem.value = conversationConfig.value.rewrite_prompt_system ?? ''
  }
}

const handleRewritePromptUserChange = async () => {
  if (localRewritePromptUser.value === (conversationConfig.value.rewrite_prompt_user ?? '')) return
  try {
    await saveConversationConfig({ rewrite_prompt_user: localRewritePromptUser.value }, t('conversationSettings.toasts.rewritePromptUserSaved'))
  } catch (error) {
    console.error('Failed to save rewrite_prompt_user:', error)
    localRewritePromptUser.value = conversationConfig.value.rewrite_prompt_user ?? ''
  }
}

const handleFallbackPromptChange = async () => {
  if (localFallbackPrompt.value === (conversationConfig.value.fallback_prompt ?? '')) return
  try {
    await saveConversationConfig({ fallback_prompt: localFallbackPrompt.value }, t('conversationSettings.toasts.fallbackPromptSaved'))
  } catch (error) {
    console.error('Failed to save fallback_prompt:', error)
    localFallbackPrompt.value = conversationConfig.value.fallback_prompt ?? ''
  }
}

// Template selection handler functions
const handleAgentSystemPromptTemplateSelect = (template: string) => {
  localSystemPrompt.value = template
}

const handleNormalSystemPromptTemplateSelect = (template: string) => {
  localSystemPromptNormal.value = template
}

const handleContextTemplateTemplateSelect = (template: string) => {
  localContextTemplate.value = template
}

const handleRewriteSystemTemplateSelect = (template: string) => {
  localRewritePromptSystem.value = template
}

const handleRewriteUserTemplateSelect = (template: string) => {
  localRewritePromptUser.value = template
}

const handleFallbackResponseTemplateSelect = (template: string) => {
  localFallbackResponse.value = template
}

const handleFallbackPromptTemplateSelect = (template: string) => {
  localFallbackPrompt.value = template
}

const navigateToModelSettings = (subsection: 'chat' | 'rerank') => {
  router.push('/platform/settings')

  setTimeout(() => {
    const event = new CustomEvent('settings-nav', {
      detail: { section: 'models', subsection },
    })
    window.dispatchEvent(event)

    setTimeout(() => {
      const selector = subsection === 'rerank' ? '[data-model-type="rerank"]' : '[data-model-type="chat"]'
      const element = document.querySelector(selector)
      if (element) {
        element.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }, 200)
  }, 100)
}

const handleConversationSummaryModelChange = async (value: string) => {
  if (value === '__add_model__') {
    localSummaryModelId.value = conversationConfig.value.summary_model_id ?? ''
    navigateToModelSettings('chat')
    return
  }

  try {
    await saveConversationConfig({ summary_model_id: value }, t('conversationSettings.toasts.chatModelSaved'))
  } catch (error) {
    console.error('Failed to save summary_model_id:', error)
    localSummaryModelId.value = conversationConfig.value.summary_model_id ?? ''
  }
}

const handleConversationRerankModelChange = async (value: string) => {
  if (value === '__add_model__') {
    localConversationRerankModelId.value = conversationConfig.value.rerank_model_id ?? ''
    navigateToModelSettings('rerank')
    return
  }

  try {
    await saveConversationConfig({ rerank_model_id: value }, t('conversationSettings.toasts.rerankModelSaved'))
  } catch (error) {
    console.error('Failed to save rerank_model_id:', error)
    localConversationRerankModelId.value = conversationConfig.value.rerank_model_id ?? ''
  }
}
</script>

<style lang="less" scoped>
.agent-settings {
  width: 100%;
}


.section-header {

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #333333;
    margin: 0 0 8px 0;
  }

  .section-description {
    font-size: 14px;
    color: #666666;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }

  .global-config-notice {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 12px 16px;
    background: #f0f9ff;
    border: 1px solid #bae0ff;
    border-radius: 8px;
    margin-bottom: 20px;
    color: #0958d9;
    font-size: 13px;
    line-height: 1.5;

    .t-icon {
      font-size: 16px;
      flex-shrink: 0;
      margin-top: 2px;
    }
  }
}

.agent-status-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid #e5e7eb;
  margin-top: 8px;

  .status-label {
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
  }

  .status-control {
    flex-shrink: 0;
    min-width: 280px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 4px 12px;
      border-radius: 4px;
      font-size: 14px;
      font-weight: 500;

      &.ready {
        background: #f0fdf4;
        color: #16a34a;
        
        .status-icon {
          color: #16a34a;
          font-size: 16px;
        }
      }

      &:not(.ready) {
        background: #fff7ed;
        color: #ea580c;
        
        .status-icon {
          color: #ea580c;
          font-size: 16px;
        }
      }

      .status-text {
        line-height: 1.4;
      }
    }

    .status-hint {
      font-size: 13px;
      color: #666666;
      text-align: right;
      line-height: 1.5;
      max-width: 280px;
    }

    .status-tip {
      margin: 8px 0 0 0;
      font-size: 12px;
      color: #999999;
      text-align: right;
      line-height: 1.5;
      max-width: 280px;
      display: flex;
      align-items: flex-start;
      gap: 4px;
      justify-content: flex-end;

      .tip-icon {
        font-size: 14px;
        color: #999999;
        flex-shrink: 0;
        margin-top: 2px;
      }
    }
  }
}

.model-recommendation-box {
  margin: 20px 0;
  background: #f0fdf6;
  border: 1px solid #d1fae5;
  border-left: 3px solid #07C05F;
  border-radius: 6px;
  padding: 16px;

  .recommendation-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;

    .recommendation-icon {
      font-size: 16px;
      color: #07C05F;
      flex-shrink: 0;
    }

    .recommendation-title {
      font-size: 14px;
      font-weight: 500;
      color: #059669;
    }
  }

  .recommendation-content {
    font-size: 13px;
    line-height: 1.6;
    color: #065f46;

    p {
      margin: 0;
    }
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
    align-items: flex-start;

    .setting-info {
      margin-bottom: 12px;
      max-width: 100%;
    }

    .setting-control.full-width {
      width: 100%;
    }
  }
}

.setting-info {
  flex: 1;
  max-width: 55%;
  word-break: keep-all;
  white-space: normal;

  .setting-info-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 4px;
    
    label {
      margin-bottom: 0;
    }
  }

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

  .hint-tip {
    margin: 8px 0 0 0;
    font-size: 12px;
    color: #999999;
    line-height: 1.5;
    display: flex;
    align-items: flex-start;
    gap: 4px;

    .tip-icon {
      font-size: 14px;
      color: #999999;
      flex-shrink: 0;
      margin-top: 2px;
    }
  }
}

.model-row {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.model-column {
  min-width: 260px;
  flex: 1;
}

.model-column-label {
  font-size: 13px;
  font-weight: 500;
  color: #555;
  margin-bottom: 4px;
}

.model-column-desc {
  margin: 0 0 8px 0;
  font-size: 12px;
  color: #888;
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
  gap: 16px;
  justify-content: flex-end;

  .value-display {
    font-size: 14px;
    font-weight: 500;
    color: #333333;
    min-width: 40px;
    text-align: right;
  }
}

// Model selector styles
.model-option {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .model-icon {
    font-size: 14px;
    color: #07C05F;
  }
  
  .add-icon {
    font-size: 14px;
    color: #07C05F;
  }
  
  .model-name {
    flex: 1;
    font-size: 13px;
  }
  
  &.add {
    .model-name {
      color: #07C05F;
      font-weight: 500;
    }
  }
}

.prompt-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  width: 100%;
}

.prompt-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
}

.prompt-toggle-label {
  font-size: 13px !important;
  color: #555;
}

.prompt-toggle :deep(.t-switch) {
  font-size: 0;
}

.prompt-toggle :deep(.t-switch__label),
.prompt-toggle :deep(.t-switch__content) {
  font-size: 12px !important;
  line-height: 18px;
  color: #666;
}

.prompt-toggle :deep(.t-switch__label--off),
.prompt-toggle :deep(.t-switch__content) {
  color: #fafafa !important;
}

.prompt-disabled-hint {
  margin: 0 0 8px;
  color: #666;
  font-size: 12px;
}

.prompt-tab-hint {
  margin: 0 0 12px;
  color: #666;
  font-size: 12px;
}

.system-prompt-tabs {
  width: 100%;
}

.allowed-tools-display {
  width: 100%;
}

.allowed-tool-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.allowed-tool-chip {
  background: #f5f7fa;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 10px 12px;
  min-width: 180px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.allowed-tool-label {
  font-size: 13px;
  font-weight: 600;
  color: #1d2129;
}

.allowed-tool-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.allowed-tools-empty {
  margin: 0;
  font-size: 12px;
  color: #999;
}

.prompt-textarea-readonly {
  background-color: #fafafa;
}

.prompt-textarea-wrapper {
  width: 100%;
}

.textarea-with-template {
  position: relative;
  width: 100%;
}

.setting-control.full-width {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.placeholder-hint {
  margin-top: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.6;

  .hint-title {
    font-weight: 500;
    color: #333;
    margin: 0 0 8px 0;
  }

  .placeholder-list {
    margin: 8px 0;
    padding-left: 20px;
    color: #666;

    li {
      margin: 4px 0;

      code {
        background: #fff;
        padding: 2px 6px;
        border-radius: 3px;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 11px;
        color: #e83e8c;
        border: 1px solid #e1e8ed;
      }
    }
  }

  .hint-tip {
    margin: 8px 0 0 0;
    color: #999;
    font-style: italic;
  }
}

.placeholder-popup-wrapper {
  position: fixed;
  z-index: 10001;
  pointer-events: auto;
}

.placeholder-popup {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  max-width: 400px;
  max-height: 300px;
  overflow-y: auto;
  padding: 4px 0;
}

.placeholder-item {
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover,
  &.active {
    background-color: #f5f7fa;
  }

  .placeholder-name {
    font-weight: 500;
    margin-bottom: 4px;

    code {
      background: #f5f7fa;
      padding: 2px 6px;
      border-radius: 3px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 12px;
      color: #e83e8c;
    }
  }

  .placeholder-desc {
    font-size: 12px;
    color: #666;
    line-height: 1.4;
  }
}

</style>

