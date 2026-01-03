<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="settings-overlay" @click.self="handleClose">
        <div class="settings-modal">
          <!-- Close button -->
          <button class="close-btn" @click="handleClose" :aria-label="$t('common.close')">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path d="M15 5L5 15M5 5L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>

          <div class="settings-container">
            <!-- Left navigation -->
            <div class="settings-sidebar">
              <div class="sidebar-header">
                <h2 class="sidebar-title">{{ mode === 'create' ? $t('agent.editor.createTitle') : $t('agent.editor.editTitle') }}</h2>
              </div>
              <div class="settings-nav">
                <div 
                  v-for="(item, index) in navItems" 
                  :key="index"
                  :class="['nav-item', { 'active': currentSection === item.key }]"
                  @click="currentSection = item.key"
                >
                  <t-icon :name="item.icon" class="nav-icon" />
                  <span class="nav-label">{{ item.label }}</span>
                </div>
              </div>
            </div>

            <!-- Right content area -->
            <div class="settings-content">
              <div class="content-wrapper">
                <!-- Basic settings -->
                <div v-show="currentSection === 'basic'" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.basicInfo') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.basicInfoDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Built-in agent notice -->
                    <div v-if="isBuiltinAgent" class="builtin-agent-notice">
                      <t-icon name="info-circle" />
                      <span>{{ $t('agent.editor.builtinAgentNotice') }}</span>
                    </div>

                    <!-- Running mode (select first) -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.mode') }} <span class="required">*</span></label>
                        <p class="desc">{{ agentMode === 'smart-reasoning' ? $t('agent.editor.agentDesc') : $t('agent.editor.normalDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-radio-group v-model="agentMode" :disabled="isBuiltinAgent">
                          <t-radio-button value="quick-answer">
                            {{ $t('agent.type.normal') }}
                          </t-radio-button>
                          <t-radio-button value="smart-reasoning">
                            {{ $t('agent.type.agent') }}
                          </t-radio-button>
                        </t-radio-group>
                      </div>
                    </div>

                    <!-- Name -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.name') }} <span v-if="!isBuiltinAgent" class="required">*</span></label>
                        <p class="desc">{{ $t('agent.editor.nameDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="name-input-wrapper">
                          <!-- Built-in agents use simple icon -->
                          <div v-if="isBuiltinAgent" class="builtin-avatar" :class="isAgentMode ? 'agent' : 'normal'">
                            <t-icon :name="isAgentMode ? 'control-platform' : 'chat'" size="24px" />
                          </div>
                          <!-- Custom agents use AgentAvatar -->
                          <AgentAvatar v-else :name="formData.name || '?'" size="large" />
                          <t-input 
                            v-model="formData.name" 
                            :placeholder="$t('agent.editor.namePlaceholder')" 
                            class="name-input"
                            :disabled="isBuiltinAgent"
                          />
                        </div>
                      </div>
                    </div>

                    <!-- Description -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.description') }}</label>
                        <p class="desc">{{ $t('agent.editor.descriptionDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-textarea 
                          v-model="formData.description" 
                          :placeholder="$t('agent.editor.descriptionPlaceholder')"
                          :autosize="{ minRows: 2, maxRows: 4 }"
                          :disabled="isBuiltinAgent"
                        />
                      </div>
                    </div>

                    <!-- System prompt -->
                    <div class="setting-row setting-row-vertical">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.systemPrompt') }} <span v-if="!isBuiltinAgent" class="required">*</span></label>
                        <p class="desc">{{ isBuiltinAgent ? $t('agent.editor.systemPromptDescBuiltin') : $t('agent.editor.systemPromptDesc') }}</p>
                        <div class="placeholder-tags">
                          <span class="placeholder-label">{{ $t('agent.editor.availableVariables') }}</span>
                          <t-tooltip 
                            v-for="placeholder in availablePlaceholders" 
                            :key="placeholder.name"
                            :content="placeholder.description + ' ' + $t('agent.editor.clickToInsert')"
                            placement="top"
                          >
                            <span 
                              class="placeholder-tag"
                              @click="handlePlaceholderClick('system', placeholder.name)"
                              v-text="formatPlaceholder(placeholder.name)"
                            ></span>
                          </t-tooltip>
                          <span class="placeholder-hint">{{ $t('agent.editor.clickToInsertOrType') }}</span>
                        </div>
                      </div>
                      <div class="setting-control setting-control-full" style="position: relative;">
                        <!-- Agent mode: unified prompt (use {{web_search_status}} placeholder to dynamically control behavior) -->
                        <div v-if="isAgentMode" class="textarea-with-template">
                          <t-textarea 
                            ref="promptTextareaRef"
                            v-model="formData.config.system_prompt" 
                            :placeholder="systemPromptPlaceholder"
                            :autosize="{ minRows: 10, maxRows: 25 }"
                            @input="handlePromptInput"
                            class="system-prompt-textarea"
                          />
                          <PromptTemplateSelector 
                            type="systemPrompt" 
                            position="corner"
                            :hasKnowledgeBase="hasKnowledgeBase"
                            @select="handleSystemPromptTemplateSelect"
                          />
                        </div>
                        <!-- Normal mode: single prompt -->
                        <div v-else class="textarea-with-template">
                          <t-textarea 
                            ref="promptTextareaRef"
                            v-model="formData.config.system_prompt" 
                            :placeholder="systemPromptPlaceholder"
                            :autosize="{ minRows: 10, maxRows: 25 }"
                            @input="handlePromptInput"
                            class="system-prompt-textarea"
                          />
                          <PromptTemplateSelector 
                            type="systemPrompt" 
                            position="corner"
                            :hasKnowledgeBase="hasKnowledgeBase"
                            @select="handleSystemPromptTemplateSelect"
                          />
                        </div>
                        <!-- Placeholder autocomplete dropdown -->
                        <Teleport to="body">
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
                                @mousedown.prevent="insertPlaceholder(placeholder.name, true)"
                                @mouseenter="selectedPlaceholderIndex = index"
                              >
                                <div class="placeholder-name">
                                  <code>{{ formatPlaceholder(placeholder.name) }}</code>
                                </div>
                                <div class="placeholder-desc">{{ placeholder.description }}</div>
                              </div>
                            </div>
                          </div>
                        </Teleport>
                      </div>
                    </div>

                    <!-- Context template (normal mode only) -->
                    <div v-if="!isAgentMode" class="setting-row setting-row-vertical">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.contextTemplate') }} <span v-if="!isBuiltinAgent" class="required">*</span></label>
                        <p class="desc">{{ isBuiltinAgent ? $t('agent.editor.contextTemplateDescBuiltin') : $t('agent.editor.contextTemplateDesc') }}</p>
                        <div class="placeholder-tags">
                          <span class="placeholder-label">{{ $t('agent.editor.availableVariables') }}</span>
                          <t-tooltip 
                            v-for="placeholder in contextTemplatePlaceholders" 
                            :key="placeholder.name"
                            :content="placeholder.description + ' ' + $t('agent.editor.clickToInsert')"
                            placement="top"
                          >
                            <span 
                              class="placeholder-tag"
                              @click="handlePlaceholderClick('context', placeholder.name)"
                              v-text="formatPlaceholder(placeholder.name)"
                            ></span>
                          </t-tooltip>
                          <span class="placeholder-hint">{{ $t('agent.editor.clickToInsertOrType') }}</span>
                        </div>
                      </div>
                      <div class="setting-control setting-control-full" style="position: relative;">
                        <div class="textarea-with-template">
                          <t-textarea 
                            ref="contextTemplateTextareaRef"
                            v-model="formData.config.context_template" 
                            :placeholder="contextTemplatePlaceholder"
                            :autosize="{ minRows: 8, maxRows: 20 }"
                            @input="handleContextTemplateInput"
                            class="system-prompt-textarea"
                          />
                          <PromptTemplateSelector 
                            type="contextTemplate" 
                            position="corner"
                            :hasKnowledgeBase="hasKnowledgeBase"
                            @select="handleContextTemplateSelect"
                          />
                        </div>
                        <!-- Context template placeholder autocomplete dropdown -->
                        <Teleport to="body">
                          <div
                            v-if="showContextPlaceholderPopup && filteredContextPlaceholders.length > 0"
                            class="placeholder-popup-wrapper"
                            :style="contextPopupStyle"
                          >
                            <div class="placeholder-popup">
                              <div
                                v-for="(placeholder, index) in filteredContextPlaceholders"
                                :key="placeholder.name"
                                class="placeholder-item"
                                :class="{ active: selectedContextPlaceholderIndex === index }"
                                @mousedown.prevent="insertContextPlaceholder(placeholder.name, true)"
                                @mouseenter="selectedContextPlaceholderIndex = index"
                              >
                                <div class="placeholder-name">
                                  <code>{{ formatPlaceholder(placeholder.name) }}</code>
                                </div>
                                <div class="placeholder-desc">{{ placeholder.description }}</div>
                              </div>
                            </div>
                          </div>
                        </Teleport>
                      </div>
                    </div>

                  </div>
                </div>

                <!-- Model configuration -->
                <div v-show="currentSection === 'model'" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.modelConfig') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.modelConfigDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Model selection -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.model') }} <span class="required">*</span></label>
                        <p class="desc">{{ $t('agent.editor.modelDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <ModelSelector
                          model-type="KnowledgeQA"
                          :selected-model-id="formData.config.model_id"
                          :all-models="allModels"
                          @update:selected-model-id="(val: string) => formData.config.model_id = val"
                          @add-model="handleAddModel('llm')"
                          :placeholder="$t('agent.editor.modelPlaceholder')"
                        />
                      </div>
                    </div>

                    <!-- Temperature -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.temperature') }}</label>
                        <p class="desc">{{ $t('agent.editor.temperatureDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="slider-wrapper">
                          <t-slider v-model="formData.config.temperature" :min="0" :max="1" :step="0.1" />
                          <span class="slider-value">{{ formData.config.temperature }}</span>
                        </div>
                      </div>
                    </div>

                    <!-- Max completion tokens (normal mode only) -->
                    <div v-if="!isAgentMode" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.maxCompletionTokens') }}</label>
                        <p class="desc">{{ $t('agent.editor.maxCompletionTokensDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-input-number v-model="formData.config.max_completion_tokens" :min="100" :max="100000" :step="100" theme="column" />
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Multi-turn conversation (normal mode only, Agent mode controls internally) -->
                <div v-show="currentSection === 'conversation' && !isAgentMode" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.conversationSettings') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.conversationSettingsDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Multi-turn conversation -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.multiTurn') }}</label>
                        <p class="desc">{{ $t('agent.editor.multiTurnDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-switch v-model="formData.config.multi_turn_enabled" />
                      </div>
                    </div>

                    <!-- History turns -->
                    <div v-if="formData.config.multi_turn_enabled" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.historyTurns') }}</label>
                        <p class="desc">{{ $t('agent.editor.historyTurnsDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-input-number v-model="formData.config.history_turns" :min="1" :max="20" theme="column" />
                      </div>
                    </div>

                    <!-- Query rewrite (only shown when multi-turn is enabled and in normal mode) -->
                    <div v-if="formData.config.multi_turn_enabled && !isAgentMode" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.enableRewrite') }}</label>
                        <p class="desc">{{ $t('agent.editor.enableRewriteDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-switch v-model="formData.config.enable_rewrite" />
                      </div>
                    </div>

                    <!-- Rewrite system prompt -->
                    <div v-if="formData.config.multi_turn_enabled && !isAgentMode && formData.config.enable_rewrite" class="setting-row setting-row-vertical">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.rewritePromptSystem') }}</label>
                        <p class="desc">{{ $t('agent.editor.rewriteSystemPromptDesc') }}</p>
                        <div class="placeholder-tags" v-if="rewriteSystemPlaceholders.length > 0">
                          <span class="placeholder-label">{{ $t('agent.editor.availableVariables') }}</span>
                          <t-tooltip 
                            v-for="placeholder in rewriteSystemPlaceholders" 
                            :key="placeholder.name"
                            :content="placeholder.description + ' ' + $t('agent.editor.clickToInsert')"
                            placement="top"
                          >
                            <span 
                              class="placeholder-tag"
                              @click="handlePlaceholderClick('rewriteSystem', placeholder.name)"
                              v-text="formatPlaceholder(placeholder.name)"
                            ></span>
                          </t-tooltip>
                          <span class="placeholder-hint">{{ $t('agent.editor.clickToInsertOrType') }}</span>
                        </div>
                      </div>
                      <div class="setting-control setting-control-full" style="position: relative;">
                        <div class="textarea-with-template">
                          <t-textarea 
                            ref="rewriteSystemTextareaRef"
                            v-model="formData.config.rewrite_prompt_system" 
                            :placeholder="defaultRewritePromptSystem || $t('agent.editor.rewritePromptSystemPlaceholder')"
                            :autosize="{ minRows: 4, maxRows: 10 }"
                            @input="handleRewriteSystemInput"
                          />
                          <PromptTemplateSelector 
                            type="rewriteSystem" 
                            position="corner"
                            @select="handleRewriteSystemTemplateSelect"
                          />
                        </div>
                        <Teleport to="body">
                          <div
                            v-if="rewriteSystemPopup.show && filteredRewriteSystemPlaceholders.length > 0"
                            class="placeholder-popup-wrapper"
                            :style="rewriteSystemPopup.style"
                          >
                            <div class="placeholder-popup">
                              <div
                                v-for="(placeholder, index) in filteredRewriteSystemPlaceholders"
                                :key="placeholder.name"
                                class="placeholder-item"
                                :class="{ active: rewriteSystemPopup.selectedIndex === index }"
                                @mousedown.prevent="insertGenericPlaceholder('rewriteSystem', placeholder.name, true)"
                                @mouseenter="rewriteSystemPopup.selectedIndex = index"
                              >
                                <div class="placeholder-name">
                                  <code>{{ formatPlaceholder(placeholder.name) }}</code>
                                </div>
                                <div class="placeholder-desc">{{ placeholder.description }}</div>
                              </div>
                            </div>
                          </div>
                        </Teleport>
                      </div>
                    </div>

                    <!-- Rewrite user prompt -->
                    <div v-if="formData.config.multi_turn_enabled && !isAgentMode && formData.config.enable_rewrite" class="setting-row setting-row-vertical">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.rewritePromptUser') }}</label>
                        <p class="desc">{{ $t('agent.editor.rewriteUserPromptDesc') }}</p>
                        <div class="placeholder-tags" v-if="rewritePlaceholders.length > 0">
                          <span class="placeholder-label">{{ $t('agent.editor.availableVariables') }}</span>
                          <t-tooltip 
                            v-for="placeholder in rewritePlaceholders" 
                            :key="placeholder.name"
                            :content="placeholder.description + ' ' + $t('agent.editor.clickToInsert')"
                            placement="top"
                          >
                            <span 
                              class="placeholder-tag"
                              @click="handlePlaceholderClick('rewriteUser', placeholder.name)"
                              v-text="formatPlaceholder(placeholder.name)"
                            ></span>
                          </t-tooltip>
                          <span class="placeholder-hint">{{ $t('agent.editor.clickToInsertOrType') }}</span>
                        </div>
                      </div>
                      <div class="setting-control setting-control-full" style="position: relative;">
                        <div class="textarea-with-template">
                          <t-textarea 
                            ref="rewriteUserTextareaRef"
                            v-model="formData.config.rewrite_prompt_user" 
                            :placeholder="defaultRewritePromptUser || $t('agent.editor.rewritePromptUserPlaceholder')"
                            :autosize="{ minRows: 4, maxRows: 10 }"
                            @input="handleRewriteUserInput"
                          />
                          <PromptTemplateSelector 
                            type="rewriteUser" 
                            position="corner"
                            @select="handleRewriteUserTemplateSelect"
                          />
                        </div>
                        <Teleport to="body">
                          <div
                            v-if="rewriteUserPopup.show && filteredRewriteUserPlaceholders.length > 0"
                            class="placeholder-popup-wrapper"
                            :style="rewriteUserPopup.style"
                          >
                            <div class="placeholder-popup">
                              <div
                                v-for="(placeholder, index) in filteredRewriteUserPlaceholders"
                                :key="placeholder.name"
                                class="placeholder-item"
                                :class="{ active: rewriteUserPopup.selectedIndex === index }"
                                @mousedown.prevent="insertGenericPlaceholder('rewriteUser', placeholder.name, true)"
                                @mouseenter="rewriteUserPopup.selectedIndex = index"
                              >
                                <div class="placeholder-name">
                                  <code>{{ formatPlaceholder(placeholder.name) }}</code>
                                </div>
                                <div class="placeholder-desc">{{ placeholder.description }}</div>
                              </div>
                            </div>
                          </div>
                        </Teleport>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Tools configuration (Agent mode only) -->
                <div v-show="currentSection === 'tools' && isAgentMode" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.toolsConfig') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.toolsConfigDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Allowed tools -->
                    <div class="setting-row setting-row-vertical">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.allowedTools') }}</label>
                        <p class="desc">{{ $t('agent.editor.allowedToolsDesc') }}</p>
                      </div>
                      <div class="setting-control setting-control-full">
                        <t-checkbox-group v-model="formData.config.allowed_tools" class="tools-checkbox-group">
                          <t-checkbox 
                            v-for="tool in availableTools" 
                            :key="tool.value" 
                            :value="tool.value"
                            :disabled="tool.disabled"
                            :class="['tool-checkbox-item', { 'tool-disabled': tool.disabled }]"
                          >
                            <div class="tool-item-content">
                              <span class="tool-name">{{ tool.label }}</span>
                              <span v-if="tool.description" class="tool-desc">{{ tool.description }}</span>
                              <span v-if="tool.disabled" class="tool-disabled-hint">{{ $t('agent.editor.toolDisabledHint') }}</span>
                            </div>
                          </t-checkbox>
                        </t-checkbox-group>
                      </div>
                    </div>

                    <!-- Max iterations -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.maxIterations') }}</label>
                        <p class="desc">{{ $t('agent.editor.maxIterationsDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-input-number v-model="formData.config.max_iterations" :min="1" :max="50" theme="column" />
                      </div>
                    </div>

                    <!-- MCP service selection -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.mcpService') }}</label>
                        <p class="desc">{{ $t('agent.editor.mcpServiceDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-radio-group v-model="mcpSelectionMode">
                          <t-radio-button value="all">{{ $t('agent.editor.allMcp') }}</t-radio-button>
                          <t-radio-button value="selected">{{ $t('agent.editor.selectedMcp') }}</t-radio-button>
                          <t-radio-button value="none">{{ $t('agent.editor.noneMcp') }}</t-radio-button>
                        </t-radio-group>
                      </div>
                    </div>

                    <!-- Select specific MCP services -->
                    <div v-if="mcpSelectionMode === 'selected' && mcpOptions.length > 0" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.selectMcpService') }}</label>
                        <p class="desc">{{ $t('agent.editor.selectMcpServiceDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-select 
                          v-model="formData.config.mcp_services" 
                          multiple 
                          :placeholder="$t('agent.editor.selectMcpServicePlaceholder')"
                          filterable
                        >
                          <t-option 
                            v-for="mcp in mcpOptions" 
                            :key="mcp.value" 
                            :value="mcp.value" 
                            :label="mcp.label" 
                          />
                        </t-select>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Knowledge base configuration -->
                <div v-show="currentSection === 'knowledge'" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.knowledgeConfig') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.knowledgeConfigDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Associate knowledge bases -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.knowledgeBases') }}</label>
                        <p class="desc">{{ $t('agent.editor.knowledgeBasesDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-radio-group v-model="kbSelectionMode">
                          <t-radio-button value="all">{{ $t('agent.editor.allKnowledgeBases') }}</t-radio-button>
                          <t-radio-button value="selected">{{ $t('agent.editor.selectedKnowledgeBases') }}</t-radio-button>
                          <t-radio-button value="none">{{ $t('agent.editor.noKnowledgeBase') }}</t-radio-button>
                        </t-radio-group>
                      </div>
                    </div>

                    <!-- Select specific knowledge bases (only shown when "selected" is chosen) -->
                    <div v-if="kbSelectionMode === 'selected'" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.selectKnowledgeBases') }}</label>
                        <p class="desc">{{ $t('agent.editor.selectKnowledgeBasesDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-select 
                          v-model="formData.config.knowledge_bases" 
                          multiple 
                          :placeholder="$t('agent.editor.selectKnowledgeBases')"
                          filterable
                        >
                          <t-option 
                            v-for="kb in kbOptions" 
                            :key="kb.value" 
                            :value="kb.value" 
                            :label="kb.label"
                          >
                            <div class="kb-option-item">
                              <span class="kb-option-icon" :class="kb.type === 'faq' ? 'faq-icon' : 'doc-icon'">
                                <t-icon :name="kb.type === 'faq' ? 'chat-bubble-help' : 'folder'" />
                              </span>
                              <span class="kb-option-label">{{ kb.label }}</span>
                              <span class="kb-option-count">({{ kb.count || 0 }})</span>
                            </div>
                          </t-option>
                        </t-select>
                      </div>
                    </div>

                    <!-- Supported file types (limit user-selectable file types) -->
                    <div v-if="hasKnowledgeBase" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.supportedFileTypes') }}</label>
                        <p class="desc">{{ $t('agent.editor.supportedFileTypesDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-select 
                          v-model="formData.config.supported_file_types" 
                          multiple 
                          :placeholder="$t('agent.editor.allTypes')"
                          :min-collapsed-num="3"
                          clearable
                        >
                          <t-option 
                            v-for="ft in availableFileTypes" 
                            :key="ft.value" 
                            :value="ft.value" 
                            :label="ft.label"
                          />
                        </t-select>
                      </div>
                    </div>

                    <!-- ReRank model (shown when knowledge base is configured) -->
                    <div v-if="needsRerankModel" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.rerankModel') }} <span class="required">*</span></label>
                        <p class="desc">{{ $t('agent.editor.rerankModelDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <ModelSelector
                          model-type="Rerank"
                          :selected-model-id="formData.config.rerank_model_id"
                          :all-models="allModels"
                          @update:selected-model-id="(val: string) => formData.config.rerank_model_id = val"
                          @add-model="handleAddModel('rerank')"
                          :placeholder="$t('agent.editor.rerankModelPlaceholder')"
                        />
                      </div>
                    </div>

                    <!-- FAQ strategy settings (only shown when FAQ type knowledge base is selected) -->
                    <div v-if="hasFaqKnowledgeBase" class="faq-strategy-section">
                      <div class="faq-strategy-header">
                        <t-icon name="chat-bubble-help" class="faq-icon" />
                        <span>{{ $t('agent.editor.faqPriorityStrategy') }}</span>
                        <t-tooltip :content="$t('agent.editor.faqPriorityStrategyDesc')">
                          <t-icon name="help-circle" class="help-icon" />
                        </t-tooltip>
                      </div>

                      <!-- FAQ priority toggle -->
                      <div class="setting-row">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.enableFaqPriority') }}</label>
                          <p class="desc">{{ $t('agent.editor.enableFaqPriorityDesc') }}</p>
                        </div>
                        <div class="setting-control">
                          <t-switch v-model="formData.config.faq_priority_enabled" />
                        </div>
                      </div>

                      <!-- FAQ direct answer threshold -->
                      <div v-if="formData.config.faq_priority_enabled" class="setting-row">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.directAnswerThreshold') }}</label>
                          <p class="desc">{{ $t('agent.editor.directAnswerThresholdDesc') }}</p>
                        </div>
                        <div class="setting-control">
                          <div class="slider-wrapper">
                            <t-slider v-model="formData.config.faq_direct_answer_threshold" :min="0.7" :max="1" :step="0.05" />
                            <span class="slider-value">{{ formData.config.faq_direct_answer_threshold?.toFixed(2) }}</span>
                          </div>
                        </div>
                      </div>

                      <!-- FAQ score boost -->
                      <div v-if="formData.config.faq_priority_enabled" class="setting-row">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.faqScoreBoost') }}</label>
                          <p class="desc">{{ $t('agent.editor.faqScoreBoostDesc') }}</p>
                        </div>
                        <div class="setting-control">
                          <div class="slider-wrapper">
                            <t-slider v-model="formData.config.faq_score_boost" :min="1" :max="2" :step="0.1" />
                            <span class="slider-value">{{ formData.config.faq_score_boost?.toFixed(1) }}x</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Web search configuration -->
                <div v-show="currentSection === 'websearch'" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.webSearchConfig') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.webSearchConfigDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Web search -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.webSearch') }}</label>
                        <p class="desc">{{ $t('agent.editor.webSearchDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-switch v-model="formData.config.web_search_enabled" />
                      </div>
                    </div>

                    <!-- Web search max results -->
                    <div v-if="formData.config.web_search_enabled" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.webSearchMaxResults') }}</label>
                        <p class="desc">{{ $t('agent.editor.webSearchMaxResultsDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="slider-wrapper">
                          <t-slider v-model="formData.config.web_search_max_results" :min="1" :max="10" />
                          <span class="slider-value">{{ formData.config.web_search_max_results }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Retrieval strategy (only shown when knowledge base capability is available) -->
                <div v-show="currentSection === 'retrieval' && hasKnowledgeBase" class="section">
                  <div class="section-header">
                    <h2>{{ $t('agent.editor.retrievalStrategy') }}</h2>
                    <p class="section-description">{{ $t('agent.editor.retrievalStrategyDesc') }}</p>
                  </div>
                  
                  <div class="settings-group">
                    <!-- Query expansion (normal mode only) -->
                    <div v-if="!isAgentMode" class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.enableQueryExpansion') }}</label>
                        <p class="desc">{{ $t('agent.editor.enableQueryExpansionDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-switch v-model="formData.config.enable_query_expansion" />
                      </div>
                    </div>

                    <!-- Vector recall TopK -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.embeddingTopK') }}</label>
                        <p class="desc">{{ $t('agent.editor.embeddingTopKDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-input-number v-model="formData.config.embedding_top_k" :min="1" :max="50" theme="column" />
                      </div>
                    </div>

                    <!-- Keyword threshold -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.keywordThreshold') }}</label>
                        <p class="desc">{{ $t('agent.editor.keywordThresholdDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="slider-wrapper">
                          <t-slider v-model="formData.config.keyword_threshold" :min="0" :max="1" :step="0.05" />
                          <span class="slider-value">{{ formData.config.keyword_threshold?.toFixed(2) }}</span>
                        </div>
                      </div>
                    </div>

                    <!-- Vector threshold -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.vectorThreshold') }}</label>
                        <p class="desc">{{ $t('agent.editor.vectorThresholdDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="slider-wrapper">
                          <t-slider v-model="formData.config.vector_threshold" :min="0" :max="1" :step="0.05" />
                          <span class="slider-value">{{ formData.config.vector_threshold?.toFixed(2) }}</span>
                        </div>
                      </div>
                    </div>

                    <!-- Rerank TopK -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.rerankTopK') }}</label>
                        <p class="desc">{{ $t('agent.editor.rerankTopKDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <t-input-number v-model="formData.config.rerank_top_k" :min="1" :max="20" theme="column" />
                      </div>
                    </div>

                    <!-- Rerank threshold -->
                    <div class="setting-row">
                      <div class="setting-info">
                        <label>{{ $t('agent.editor.rerankThreshold') }}</label>
                        <p class="desc">{{ $t('agent.editor.rerankThresholdDesc') }}</p>
                      </div>
                      <div class="setting-control">
                        <div class="slider-wrapper">
                          <t-slider v-model="formData.config.rerank_threshold" :min="0" :max="1" :step="0.05" />
                          <span class="slider-value">{{ formData.config.rerank_threshold?.toFixed(2) }}</span>
                        </div>
                      </div>
                    </div>

                    <!-- Fallback strategy (normal mode only) -->
                    <template v-if="!isAgentMode">
                      <div class="setting-row">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.fallbackStrategy') }}</label>
                          <p class="desc">{{ $t('agent.editor.fallbackStrategyDesc') }}</p>
                        </div>
                        <div class="setting-control">
                          <t-radio-group v-model="formData.config.fallback_strategy">
                            <t-radio-button value="fixed">{{ $t('agent.editor.fixedResponse') }}</t-radio-button>
                            <t-radio-button value="model">{{ $t('agent.editor.modelGeneration') }}</t-radio-button>
                          </t-radio-group>
                        </div>
                      </div>

                      <!-- Fixed fallback response -->
                      <div v-if="formData.config.fallback_strategy === 'fixed'" class="setting-row setting-row-vertical">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.fallbackResponse') }}</label>
                          <p class="desc">{{ $t('agent.editor.fallbackResponseDesc') }}</p>
                        </div>
                        <div class="setting-control setting-control-full">
                          <div class="textarea-with-template">
                            <t-textarea 
                              v-model="formData.config.fallback_response" 
                              :placeholder="defaultFallbackResponse || $t('agent.editor.fallbackResponsePlaceholder')"
                              :autosize="{ minRows: 2, maxRows: 6 }"
                            />
                            <PromptTemplateSelector 
                              type="fallback" 
                              position="corner"
                              @select="handleFallbackResponseTemplateSelect"
                            />
                          </div>
                        </div>
                      </div>

                      <!-- Fallback prompt -->
                      <div v-if="formData.config.fallback_strategy === 'model'" class="setting-row setting-row-vertical">
                        <div class="setting-info">
                          <label>{{ $t('agent.editor.fallbackPrompt') }}</label>
                          <p class="desc">{{ $t('agent.editor.fallbackPromptDesc') }}</p>
                          <div class="placeholder-tags" v-if="fallbackPlaceholders.length > 0">
                            <span class="placeholder-label">{{ $t('agent.editor.availableVariables') }}</span>
                            <t-tooltip 
                              v-for="placeholder in fallbackPlaceholders" 
                              :key="placeholder.name"
                              :content="placeholder.description + ' ' + $t('agent.editor.clickToInsert')"
                              placement="top"
                            >
                              <span 
                                class="placeholder-tag"
                                @click="handlePlaceholderClick('fallback', placeholder.name)"
                                v-text="formatPlaceholder(placeholder.name)"
                              ></span>
                            </t-tooltip>
                            <span class="placeholder-hint">{{ $t('agent.editor.clickToInsertOrType') }}</span>
                          </div>
                        </div>
                        <div class="setting-control setting-control-full" style="position: relative;">
                          <div class="textarea-with-template">
                            <t-textarea 
                              ref="fallbackPromptTextareaRef"
                              v-model="formData.config.fallback_prompt" 
                              :placeholder="defaultFallbackPrompt || $t('agent.editor.fallbackPromptPlaceholder')"
                              :autosize="{ minRows: 4, maxRows: 10 }"
                              @input="handleFallbackPromptInput"
                            />
                            <PromptTemplateSelector 
                              type="fallback" 
                              position="corner"
                              @select="handleFallbackPromptTemplateSelect"
                            />
                          </div>
                          <Teleport to="body">
                            <div
                              v-if="fallbackPromptPopup.show && filteredFallbackPlaceholders.length > 0"
                              class="placeholder-popup-wrapper"
                              :style="fallbackPromptPopup.style"
                            >
                              <div class="placeholder-popup">
                                <div
                                  v-for="(placeholder, index) in filteredFallbackPlaceholders"
                                  :key="placeholder.name"
                                  class="placeholder-item"
                                  :class="{ active: fallbackPromptPopup.selectedIndex === index }"
                                  @mousedown.prevent="insertGenericPlaceholder('fallback', placeholder.name, true)"
                                  @mouseenter="fallbackPromptPopup.selectedIndex = index"
                                >
                                  <div class="placeholder-name">
                                    <code>{{ formatPlaceholder(placeholder.name) }}</code>
                                  </div>
                                  <div class="placeholder-desc">{{ placeholder.description }}</div>
                                </div>
                              </div>
                            </div>
                          </Teleport>
                        </div>
                      </div>
                    </template>
                  </div>
                </div>
              </div>

              <!-- Bottom action bar -->
              <div class="settings-footer">
                <t-button variant="outline" @click="handleClose">{{ $t('common.cancel') }}</t-button>
                <t-button theme="primary" :loading="saving" @click="handleSave">{{ $t('common.confirm') }}</t-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { MessagePlugin } from 'tdesign-vue-next';
import { createAgent, updateAgent, getPlaceholders, type CustomAgent, type PlaceholderDefinition } from '@/api/agent';
import { listModels, type ModelConfig } from '@/api/model';
import { listKnowledgeBases } from '@/api/knowledge-base';
import { listMCPServices, type MCPService } from '@/api/mcp-service';
import { getAgentConfig, getConversationConfig } from '@/api/system';
import { useUIStore } from '@/stores/ui';
import AgentAvatar from '@/components/AgentAvatar.vue';
import PromptTemplateSelector from '@/components/PromptTemplateSelector.vue';
import ModelSelector from '@/components/ModelSelector.vue';

const uiStore = useUIStore();

const { t } = useI18n();

const props = defineProps<{
  visible: boolean;
  mode: 'create' | 'edit';
  agent?: CustomAgent | null;
  initialSection?: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void;
  (e: 'success'): void;
}>();

const currentSection = ref(props.initialSection || 'basic');
const saving = ref(false);
const allModels = ref<ModelConfig[]>([]);
const kbOptions = ref<{ label: string; value: string; type?: 'document' | 'faq'; count?: number }[]>([]);
const mcpOptions = ref<{ label: string; value: string }[]>([]);

// System default configuration (for built-in agents to display default prompts)
const defaultAgentSystemPrompt = ref('');  // Default system prompt for Agent mode (from agent-config)
const defaultNormalSystemPrompt = ref('');  // Default system prompt for normal mode (from conversation-config)
const defaultContextTemplate = ref('');
const defaultRewritePromptSystem = ref('');
const defaultRewritePromptUser = ref('');
const defaultFallbackPrompt = ref('');
const defaultFallbackResponse = ref('');
// Default retrieval parameters
const defaultEmbeddingTopK = ref(10);
const defaultKeywordThreshold = ref(0.3);
const defaultVectorThreshold = ref(0.5);
const defaultRerankTopK = ref(5);
const defaultRerankThreshold = ref(0.5);
const defaultMaxCompletionTokens = ref(2048);
const defaultTemperature = ref(0.7);

// Knowledge base related tools list
const knowledgeBaseTools = ['grep_chunks', 'knowledge_search', 'list_knowledge_chunks', 'query_knowledge_graph', 'get_document_info', 'database_query'];

// Initialization flag to prevent watch from automatically adding tools during initialization
const isInitializing = ref(false);

// Knowledge base selection mode: all=all, selected=specified, none=not used
const kbSelectionMode = ref<'all' | 'selected' | 'none'>('none');

// MCP service selection mode: all=all, selected=specified, none=not used
const mcpSelectionMode = ref<'all' | 'selected' | 'none'>('none');

// Available tools list (keep consistent with backend definitions.go)
const allTools = [
  { value: 'thinking', label: 'Thinking', description: 'Dynamic and reflective problem-solving thinking tool', requiresKB: false },
  { value: 'todo_write', label: 'Create Plan', description: 'Create structured research plans', requiresKB: false },
  { value: 'grep_chunks', label: 'Keyword Search', description: 'Quickly locate documents and chunks containing specific keywords', requiresKB: true },
  { value: 'knowledge_search', label: 'Semantic Search', description: 'Understand questions and find semantically related content', requiresKB: true },
  { value: 'list_knowledge_chunks', label: 'View Document Chunks', description: 'Get complete chunk content of documents', requiresKB: true },
  { value: 'query_knowledge_graph', label: 'Query Knowledge Graph', description: 'Query relationships from knowledge graph', requiresKB: true },
  { value: 'get_document_info', label: 'Get Document Info', description: 'View document metadata', requiresKB: true },
  { value: 'database_query', label: 'Query Database', description: 'Query information from database', requiresKB: true },
  { value: 'data_analysis', label: 'Data Analysis', description: 'Understand data files and perform data analysis', requiresKB: true },
  { value: 'data_schema', label: 'View Data Schema', description: 'Get metadata of table files', requiresKB: true },
];

// Dynamically calculate whether there is knowledge base capability based on knowledge base configuration
const hasKnowledgeBase = computed(() => {
  return kbSelectionMode.value !== 'none';
});

// Detect if selected knowledge bases contain FAQ type
const hasFaqKnowledgeBase = computed(() => {
  if (kbSelectionMode.value === 'none') return false;
  if (kbSelectionMode.value === 'all') {
    // All knowledge bases mode, check if there are any FAQ type knowledge bases
    return kbOptions.value.some(kb => kb.type === 'faq');
  }
  // Selected knowledge bases mode, check if selected knowledge bases contain FAQ type
  const selectedKbIds = formData.value.config.knowledge_bases || [];
  return kbOptions.value.some(kb => selectedKbIds.includes(kb.value) && kb.type === 'faq');
});

const availableTools = computed(() => {
  return allTools.map(tool => ({
    ...tool,
    disabled: tool.requiresKB && !hasKnowledgeBase.value
  }));
});

// Available file types list
const availableFileTypes = [
  { value: 'pdf', label: 'PDF', description: 'PDF document' },
  { value: 'docx', label: 'Word', description: 'Word document (.docx/.doc)' },
  { value: 'txt', label: 'Text', description: 'Plain text file (.txt)' },
  { value: 'md', label: 'Markdown', description: 'Markdown document' },
  { value: 'csv', label: 'CSV', description: 'Comma-separated values file' },
  { value: 'xlsx', label: 'Excel', description: 'Excel spreadsheet (.xlsx/.xls)' },
  { value: 'jpg', label: 'Image', description: 'Image file (.jpg/.jpeg/.png)' },
];

// Placeholder related - fetched from API
const placeholderData = ref<{
  system_prompt: PlaceholderDefinition[];
  agent_system_prompt: PlaceholderDefinition[];
  context_template: PlaceholderDefinition[];
  rewrite_system_prompt: PlaceholderDefinition[];
  rewrite_prompt: PlaceholderDefinition[];
  fallback_prompt: PlaceholderDefinition[];
}>({
  system_prompt: [],
  agent_system_prompt: [],
  context_template: [],
  rewrite_system_prompt: [],
  rewrite_prompt: [],
  fallback_prompt: [],
});

// System prompt placeholders (dynamically selected based on mode)
const availablePlaceholders = computed(() => {
  return isAgentMode.value ? placeholderData.value.agent_system_prompt : placeholderData.value.system_prompt;
});

// Context template placeholders
const contextTemplatePlaceholders = computed(() => placeholderData.value.context_template);

// Rewrite system prompt placeholders
const rewriteSystemPlaceholders = computed(() => placeholderData.value.rewrite_system_prompt);

// Rewrite user prompt placeholders
const rewritePlaceholders = computed(() => placeholderData.value.rewrite_prompt);

// Fallback prompt placeholders
const fallbackPlaceholders = computed(() => placeholderData.value.fallback_prompt);

// Format placeholder name with curly braces
const formatPlaceholder = (name: string): string => {
  return `{{${name}}}`;
};

const promptTextareaRef = ref<any>(null);
const showPlaceholderPopup = ref(false);
const selectedPlaceholderIndex = ref(0);
const placeholderPrefix = ref('');
const popupStyle = ref({ top: '0px', left: '0px' });
let placeholderPopupTimer: any = null;

// Context template placeholder related
const contextTemplateTextareaRef = ref<any>(null);
const showContextPlaceholderPopup = ref(false);
const selectedContextPlaceholderIndex = ref(0);
const contextPlaceholderPrefix = ref('');
const contextPopupStyle = ref({ top: '0px', left: '0px' });
let contextPlaceholderPopupTimer: any = null;

// Generic placeholder popup related (for rewrite prompts and fallback prompts)
interface PlaceholderPopupState {
  show: boolean;
  selectedIndex: number;
  prefix: string;
  style: { top: string; left: string };
  timer: any;
  fieldKey: string;
  placeholders: PlaceholderDefinition[];
}

const rewriteSystemPopup = ref<PlaceholderPopupState>({
  show: false, selectedIndex: 0, prefix: '', style: { top: '0px', left: '0px' }, timer: null, fieldKey: 'rewrite_prompt_system', placeholders: []
});
const rewriteUserPopup = ref<PlaceholderPopupState>({
  show: false, selectedIndex: 0, prefix: '', style: { top: '0px', left: '0px' }, timer: null, fieldKey: 'rewrite_prompt_user', placeholders: []
});
const fallbackPromptPopup = ref<PlaceholderPopupState>({
  show: false, selectedIndex: 0, prefix: '', style: { top: '0px', left: '0px' }, timer: null, fieldKey: 'fallback_prompt', placeholders: []
});

const rewriteSystemTextareaRef = ref<any>(null);
const rewriteUserTextareaRef = ref<any>(null);
const fallbackPromptTextareaRef = ref<any>(null);

const navItems = computed(() => {
  const items: { key: string; icon: string; label: string }[] = [
    { key: 'basic', icon: 'info-circle', label: t('agent.editor.basicInfo') },
    { key: 'model', icon: 'control-platform', label: t('agent.editor.modelConfig') },
  ];
  // Knowledge base configuration (placed above tools)
  items.push({ key: 'knowledge', icon: 'folder', label: t('agent.editor.knowledgeConfig') });
  // Only show tools configuration in Agent mode
  if (isAgentMode.value) {
    items.push({ key: 'tools', icon: 'tools', label: t('agent.editor.toolsConfig') });
  }
  // Only show retrieval strategy when knowledge base capability is available
  if (hasKnowledgeBase.value) {
    items.push({ key: 'retrieval', icon: 'search', label: t('agent.editor.retrievalStrategy') });
  }
  // Web search (independent menu)
  items.push({ key: 'websearch', icon: 'internet', label: t('agent.editor.webSearchConfig') });
  // Multi-turn conversation (only shown in normal mode, Agent mode controls internally)
  if (!isAgentMode.value) {
    items.push({ key: 'conversation', icon: 'chat', label: t('agent.editor.conversationSettings') });
  }
  return items;
});

// Initial data
const defaultFormData = {
  name: '',
  description: '',
  is_builtin: false,
  config: {
    // Basic settings
    agent_mode: 'quick-answer' as 'quick-answer' | 'smart-reasoning',
    system_prompt: '',
    context_template: '{{query}}',
    // Model settings
    model_id: '',
    rerank_model_id: '',
    temperature: 0.7,
    max_completion_tokens: 2048,
    // Agent mode settings
    max_iterations: 10,
    allowed_tools: [] as string[],
    reflection_enabled: false,
    // MCP service settings
    mcp_selection_mode: 'none' as 'all' | 'selected' | 'none',
    mcp_services: [] as string[],
    // Knowledge base settings
    kb_selection_mode: 'none' as 'all' | 'selected' | 'none',
    knowledge_bases: [] as string[],
    // File type restrictions
    supported_file_types: [] as string[],
    // FAQ strategy settings
    faq_priority_enabled: true, // Whether to enable FAQ priority strategy
    faq_direct_answer_threshold: 0.9, // FAQ direct answer threshold (use FAQ answer directly when similarity exceeds this value)
    faq_score_boost: 1.2, // FAQ score boost coefficient
    // Web search settings
    web_search_enabled: false,
    web_search_max_results: 5,
    // Multi-turn conversation settings
    multi_turn_enabled: false,
    history_turns: 5,
    // Retrieval strategy settings
    embedding_top_k: 10,
    keyword_threshold: 0.3,
    vector_threshold: 0.5,
    rerank_top_k: 5,
    rerank_threshold: 0.5,
    // Advanced settings (normal mode)
    enable_query_expansion: true,
    enable_rewrite: true,
    rewrite_prompt_system: '',
    rewrite_prompt_user: '',
    fallback_strategy: 'model' as 'fixed' | 'model',
    fallback_response: '',
    fallback_prompt: '',
    // Deprecated fields (kept for compatibility)
    welcome_message: '',
    suggested_prompts: [] as string[],
  }
};

const formData = ref(JSON.parse(JSON.stringify(defaultFormData)));
const agentMode = computed({
  get: () => formData.value.config.agent_mode,
  set: (val: 'quick-answer' | 'smart-reasoning') => { formData.value.config.agent_mode = val; }
});

const isAgentMode = computed(() => agentMode.value === 'smart-reasoning');

// Whether it is a built-in agent
const isBuiltinAgent = computed(() => {
  return formData.value.is_builtin === true;
});

// System prompt placeholder
const systemPromptPlaceholder = computed(() => {
  return t('agent.editor.systemPromptPlaceholder');
});

// Context template placeholder
const contextTemplatePlaceholder = computed(() => {
  return t('agent.editor.contextTemplatePlaceholder');
});

// Whether ReRank model configuration is needed (required when knowledge base capability is available)
const needsRerankModel = computed(() => {
  return hasKnowledgeBase.value;
});

// Watch visibility changes, reset form
watch(() => props.visible, async (val) => {
  if (val) {
    try {
      currentSection.value = props.initialSection || 'basic';
      // Load dependency data first (including default configuration)
      await loadDependencies();
      
      if (props.mode === 'edit' && props.agent) {
        // Deep copy object to avoid reference issues
        let agentData;
        try {
          agentData = JSON.parse(JSON.stringify(props.agent));
        } catch (e) {
          console.error('Error serializing agent data:', e);
          // Fallback: use shallow copy if JSON serialization fails
          agentData = { ...props.agent };
          if (props.agent.config) {
            agentData.config = { ...props.agent.config };
          }
        }
      
      // Ensure config object exists
      if (!agentData.config) {
        agentData.config = JSON.parse(JSON.stringify(defaultFormData.config));
      }
      
      // Fill in potentially missing fields
      agentData.config = { ...defaultFormData.config, ...agentData.config };
      
      // Ensure array fields exist
      if (!agentData.config.suggested_prompts) agentData.config.suggested_prompts = [];
      if (!agentData.config.knowledge_bases) agentData.config.knowledge_bases = [];
      if (!agentData.config.allowed_tools) agentData.config.allowed_tools = [];
      if (!agentData.config.mcp_services) agentData.config.mcp_services = [];
      if (!agentData.config.supported_file_types) agentData.config.supported_file_types = [];

      // Compatibility with old data: if agent_mode field is missing, infer from allowed_tools
      if (!agentData.config.agent_mode) {
        const isAgent = agentData.config.max_iterations > 1 || (agentData.config.allowed_tools && agentData.config.allowed_tools.length > 0);
        agentData.config.agent_mode = isAgent ? 'smart-reasoning' : 'quick-answer';
      }

      // Set initialization flag to prevent watch from automatically adding tools
      isInitializing.value = true;
      formData.value = agentData;
      // Initialize knowledge base selection mode
      initKbSelectionMode();
      initMcpSelectionMode();
      // Reset flag after initialization is complete
      nextTick(() => {
        isInitializing.value = false;
      });
      // Built-in agent: if prompt is empty, fill in system default values
      if (agentData.is_builtin) {
        fillBuiltinAgentDefaults();
      }
    } else {
      // Create new agent, use system default values
      const newFormData = JSON.parse(JSON.stringify(defaultFormData));
      // Apply system default retrieval parameters
      newFormData.config.embedding_top_k = defaultEmbeddingTopK.value;
      newFormData.config.keyword_threshold = defaultKeywordThreshold.value;
      newFormData.config.vector_threshold = defaultVectorThreshold.value;
      newFormData.config.rerank_top_k = defaultRerankTopK.value;
      newFormData.config.rerank_threshold = defaultRerankThreshold.value;
      newFormData.config.max_completion_tokens = defaultMaxCompletionTokens.value;
      newFormData.config.temperature = defaultTemperature.value;
      // Apply system default context template
      if (defaultContextTemplate.value) {
        newFormData.config.context_template = defaultContextTemplate.value;
      }
      formData.value = newFormData;
      kbSelectionMode.value = 'none';
      mcpSelectionMode.value = 'none';
    }
    } catch (error) {
      console.error('Error initializing agent editor:', error);
      MessagePlugin.error(t('agent.messages.openFailed') || 'Failed to initialize agent editor');
      // Close modal on error
      emit('update:visible', false);
    }
  } else {
    // Modal closed, reset form
    try {
      formData.value = JSON.parse(JSON.stringify(defaultFormData));
      currentSection.value = 'basic';
    } catch (error) {
      console.error('Error resetting form:', error);
      // Fallback: use shallow copy
      formData.value = { ...defaultFormData };
      if (defaultFormData.config) {
        formData.value.config = { ...defaultFormData.config };
      }
      currentSection.value = 'basic';
    }
  }
});

// Initialize knowledge base selection mode
const initKbSelectionMode = () => {
  if (formData.value.config.kb_selection_mode) {
    // If there is a saved mode, use it directly
    kbSelectionMode.value = formData.value.config.kb_selection_mode;
  } else if (formData.value.config.knowledge_bases?.length > 0) {
    // Has specified knowledge bases
    kbSelectionMode.value = 'selected';
  } else {
    kbSelectionMode.value = 'none';
  }
};

// Initialize MCP selection mode
const initMcpSelectionMode = () => {
  if (formData.value.config.mcp_selection_mode) {
    // If there is a saved mode, use it directly
    mcpSelectionMode.value = formData.value.config.mcp_selection_mode;
  } else if (formData.value.config.mcp_services?.length > 0) {
    // Has specified MCP services
    mcpSelectionMode.value = 'selected';
  } else {
    mcpSelectionMode.value = 'none';
  }
};

// Built-in agent: fill in system default values
const fillBuiltinAgentDefaults = () => {
  const config = formData.value.config;
  const isAgent = config.agent_mode === 'smart-reasoning';
  
  if (isAgent) {
    // Agent mode: use default prompt from agent-config
    if (!config.system_prompt && defaultAgentSystemPrompt.value) {
      config.system_prompt = defaultAgentSystemPrompt.value;
    }
  } else {
    // Normal mode: use default system prompt and context template from conversation-config
    if (!config.system_prompt && defaultNormalSystemPrompt.value) {
      config.system_prompt = defaultNormalSystemPrompt.value;
    }
    if (!config.context_template && defaultContextTemplate.value) {
      config.context_template = defaultContextTemplate.value;
    }
  }
  
  // Common default values
  if (!config.rewrite_prompt_system && defaultRewritePromptSystem.value) {
    config.rewrite_prompt_system = defaultRewritePromptSystem.value;
  }
  if (!config.rewrite_prompt_user && defaultRewritePromptUser.value) {
    config.rewrite_prompt_user = defaultRewritePromptUser.value;
  }
  if (!config.fallback_prompt && defaultFallbackPrompt.value) {
    config.fallback_prompt = defaultFallbackPrompt.value;
  }
  if (!config.fallback_response && defaultFallbackResponse.value) {
    config.fallback_response = defaultFallbackResponse.value;
  }
};

// Watch knowledge base selection mode changes
watch(kbSelectionMode, (mode) => {
  formData.value.config.kb_selection_mode = mode;
  if (mode === 'none') {
    // Not using knowledge base, clear related configuration
    formData.value.config.knowledge_bases = [];
  } else if (mode === 'all') {
    // All knowledge bases, clear specified list
    formData.value.config.knowledge_bases = [];
  }
  // selected mode keeps knowledge_bases unchanged
});

// Watch MCP selection mode changes
watch(mcpSelectionMode, (mode) => {
  formData.value.config.mcp_selection_mode = mode;
  if (mode === 'none') {
    // Not using MCP, clear related configuration
    formData.value.config.mcp_services = [];
  } else if (mode === 'all') {
    // All MCP, clear specified list
    formData.value.config.mcp_services = [];
  }
  // selected mode keeps mcp_services unchanged
});

// Watch mode changes, automatically adjust configuration
watch(agentMode, (val) => {
  if (val === 'smart-reasoning') {
    // Switch to Agent mode, enable tools based on knowledge base configuration
    if (formData.value.config.allowed_tools.length === 0) {
      if (hasKnowledgeBase.value) {
        // When knowledge base is available, enable all tools
        formData.value.config.allowed_tools = [
          'thinking',
          'todo_write',
          'knowledge_search',
          'grep_chunks',
          'list_knowledge_chunks',
          'query_knowledge_graph',
          'get_document_info',
          'database_query',
        ];
      } else {
        // When no knowledge base, only enable non-knowledge base tools
        formData.value.config.allowed_tools = ['thinking', 'todo_write'];
      }
    }
    if (formData.value.config.max_iterations <= 1) {
      formData.value.config.max_iterations = 10;
    }
  } else {
    // Switch to normal mode, clear tools
    formData.value.config.allowed_tools = [];
    formData.value.config.max_iterations = 1; // Set to 1 to indicate single-turn RAG
  }
});

// Watch knowledge base configuration changes, automatically remove/add knowledge base related tools
watch(hasKnowledgeBase, (hasKB, oldHasKB) => {
  // If currently on retrieval strategy page but no knowledge base capability, switch to basic settings
  if (!hasKB && currentSection.value === 'retrieval') {
    currentSection.value = 'basic';
  }
  
  // Do not automatically adjust tools during initialization or in non-Agent mode
  if (isInitializing.value || !isAgentMode.value) return;
  
  if (hasKB && !oldHasKB) {
    // Changed from no knowledge base to having knowledge base, automatically add knowledge base related tools
    const currentTools = formData.value.config.allowed_tools || [];
    const toolsToAdd = knowledgeBaseTools.filter((tool: string) => !currentTools.includes(tool));
    formData.value.config.allowed_tools = [...currentTools, ...toolsToAdd];
  } else if (!hasKB && oldHasKB) {
    // Changed from having knowledge base to no knowledge base, remove knowledge base related tools
    formData.value.config.allowed_tools = formData.value.config.allowed_tools.filter(
      (tool: string) => !knowledgeBaseTools.includes(tool)
    );
  }
});

// Watch running mode changes, automatically switch pages
watch(isAgentMode, (isAgent) => {
  // If currently on advanced settings page but switched to Agent mode, switch to basic settings
  if (isAgent && currentSection.value === 'advanced') {
    currentSection.value = 'basic';
  }
  // If currently on multi-turn conversation page but switched to Agent mode, switch to basic settings (Agent mode controls multi-turn conversation internally)
  if (isAgent && currentSection.value === 'conversation') {
    currentSection.value = 'basic';
  }
});

// Watch settings modal close, refresh model list
watch(() => uiStore.showSettingsModal, async (visible, prevVisible) => {
  // When returning from settings page (modal closed), refresh model list
  if (prevVisible && !visible && props.visible) {
    try {
      const models = await listModels();
      if (models && models.length > 0) {
        allModels.value = models;
      }
    } catch (e) {
      console.warn('Failed to refresh models after settings closed', e);
    }
  }
});

// Load dependency data
const loadDependencies = async () => {
  try {
    // Load all model list (ModelSelector component will automatically filter by type)
    const models = await listModels();
    if (models && models.length > 0) {
      allModels.value = models;
    }

    // Load knowledge base list
    const kbRes: any = await listKnowledgeBases();
    if (kbRes.data) {
      kbOptions.value = kbRes.data.map((kb: any) => ({ 
        label: kb.name, 
        value: kb.id,
        type: kb.type || 'document',
        count: kb.type === 'faq' ? (kb.chunk_count || 0) : (kb.knowledge_count || 0)
      }));
    }

    // Load MCP service list (only load enabled ones)
    try {
      const mcpList = await listMCPServices();
      if (mcpList && mcpList.length > 0) {
        mcpOptions.value = mcpList
          .filter((mcp: MCPService) => mcp.enabled)
          .map((mcp: MCPService) => ({ label: mcp.name, value: mcp.id }));
      }
    } catch (e) {
      console.warn('Failed to load MCP services', e);
    }

    // Load placeholder definitions (from unified API)
    try {
      const placeholdersRes = await getPlaceholders();
      if (placeholdersRes.data) {
        placeholderData.value = placeholdersRes.data;
      }
    } catch (e) {
      console.warn('Failed to load placeholders', e);
    }

    // Load Agent mode default prompt (from agent-config, for smart-reasoning mode)
    const agentConfig = await getAgentConfig();
    if (agentConfig.data?.system_prompt) {
      defaultAgentSystemPrompt.value = agentConfig.data.system_prompt;
    }

    // Load system default configuration (from conversation-config, for normal mode quick-answer)
    const conversationConfig = await getConversationConfig();
    if (conversationConfig.data?.prompt) {
      defaultNormalSystemPrompt.value = conversationConfig.data.prompt;
    }
    if (conversationConfig.data?.context_template) {
      defaultContextTemplate.value = conversationConfig.data.context_template;
    }
    if (conversationConfig.data?.rewrite_prompt_system) {
      defaultRewritePromptSystem.value = conversationConfig.data.rewrite_prompt_system;
    }
    if (conversationConfig.data?.rewrite_prompt_user) {
      defaultRewritePromptUser.value = conversationConfig.data.rewrite_prompt_user;
    }
    if (conversationConfig.data?.fallback_prompt) {
      defaultFallbackPrompt.value = conversationConfig.data.fallback_prompt;
    }
    if (conversationConfig.data?.fallback_response) {
      defaultFallbackResponse.value = conversationConfig.data.fallback_response;
    }
    // Load default retrieval parameters
    if (conversationConfig.data?.embedding_top_k) {
      defaultEmbeddingTopK.value = conversationConfig.data.embedding_top_k;
    }
    if (conversationConfig.data?.keyword_threshold !== undefined) {
      defaultKeywordThreshold.value = conversationConfig.data.keyword_threshold;
    }
    if (conversationConfig.data?.vector_threshold !== undefined) {
      defaultVectorThreshold.value = conversationConfig.data.vector_threshold;
    }
    if (conversationConfig.data?.rerank_top_k) {
      defaultRerankTopK.value = conversationConfig.data.rerank_top_k;
    }
    if (conversationConfig.data?.rerank_threshold !== undefined) {
      defaultRerankThreshold.value = conversationConfig.data.rerank_threshold;
    }
    if (conversationConfig.data?.max_completion_tokens) {
      defaultMaxCompletionTokens.value = conversationConfig.data.max_completion_tokens;
    }
    if (conversationConfig.data?.temperature !== undefined) {
      defaultTemperature.value = conversationConfig.data.temperature;
    }
  } catch (e) {
    console.error('Failed to load dependencies', e);
  }
};

// Navigate to model management page to add model
const handleAddModel = (subSection: string) => {
  uiStore.openSettings('models', subSection);
};

const handleClose = () => {
  showPlaceholderPopup.value = false;
  showContextPlaceholderPopup.value = false;
  rewriteSystemPopup.value.show = false;
  rewriteUserPopup.value.show = false;
  fallbackPromptPopup.value.show = false;
  emit('update:visible', false);
};

// Filtered placeholder list
const filteredPlaceholders = computed(() => {
  if (!placeholderPrefix.value) {
    return availablePlaceholders.value;
  }
  const prefix = placeholderPrefix.value.toLowerCase();
  return availablePlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  );
});

// Filtered context template placeholder list
const filteredContextPlaceholders = computed(() => {
  if (!contextPlaceholderPrefix.value) {
    return contextTemplatePlaceholders.value;
  }
  const prefix = contextPlaceholderPrefix.value.toLowerCase();
  return contextTemplatePlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  );
});

// Filtered rewrite system prompt placeholder list
const filteredRewriteSystemPlaceholders = computed(() => {
  if (!rewriteSystemPopup.value.prefix) {
    return rewriteSystemPlaceholders.value;
  }
  const prefix = rewriteSystemPopup.value.prefix.toLowerCase();
  return rewriteSystemPlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  );
});

// Filtered rewrite user prompt placeholder list
const filteredRewriteUserPlaceholders = computed(() => {
  if (!rewriteUserPopup.value.prefix) {
    return rewritePlaceholders.value;
  }
  const prefix = rewriteUserPopup.value.prefix.toLowerCase();
  return rewritePlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  );
});

// Filtered fallback prompt placeholder list
const filteredFallbackPlaceholders = computed(() => {
  if (!fallbackPromptPopup.value.prefix) {
    return fallbackPlaceholders.value;
  }
  const prefix = fallbackPromptPopup.value.prefix.toLowerCase();
  return fallbackPlaceholders.value.filter(p => 
    p.name.toLowerCase().startsWith(prefix)
  );
});

// Get textarea element
const getTextareaElement = (): HTMLTextAreaElement | null => {
  if (promptTextareaRef.value) {
    if (promptTextareaRef.value.$el) {
      return promptTextareaRef.value.$el.querySelector('textarea');
    }
    if (promptTextareaRef.value instanceof HTMLTextAreaElement) {
      return promptTextareaRef.value;
    }
  }
  return null;
};

// Calculate cursor position
const calculateCursorPosition = (textarea: HTMLTextAreaElement) => {
  const cursorPos = textarea.selectionStart;
  const textBeforeCursor = formData.value.config.system_prompt.substring(0, cursorPos);
  
  const style = window.getComputedStyle(textarea);
  const textareaRect = textarea.getBoundingClientRect();
  
  const lineHeight = parseFloat(style.lineHeight) || 20;
  const paddingTop = parseFloat(style.paddingTop) || 0;
  const paddingLeft = parseFloat(style.paddingLeft) || 0;
  
  // Calculate current line number
  const lines = textBeforeCursor.split('\n');
  const currentLine = lines.length - 1;
  const currentLineText = lines[currentLine];
  
  // Create temporary span to calculate text width
  const span = document.createElement('span');
  span.style.font = style.font;
  span.style.visibility = 'hidden';
  span.style.position = 'absolute';
  span.style.whiteSpace = 'pre';
  span.textContent = currentLineText;
  document.body.appendChild(span);
  const textWidth = span.offsetWidth;
  document.body.removeChild(span);
  
  const scrollTop = textarea.scrollTop;
  const top = textareaRect.top + paddingTop + (currentLine * lineHeight) - scrollTop + lineHeight + 4;
  const scrollLeft = textarea.scrollLeft;
  const left = textareaRect.left + paddingLeft + textWidth - scrollLeft;
  
  return { top, left };
};

// Check and show placeholder popup
const checkAndShowPlaceholderPopup = () => {
  const textarea = getTextareaElement();
  if (!textarea) return;
  
  const cursorPos = textarea.selectionStart;
  const textBeforeCursor = formData.value.config.system_prompt.substring(0, cursorPos);
  
  // Find the nearest {{ position
  let lastOpenPos = -1;
  for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
    if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
      const textAfterOpen = textBeforeCursor.substring(i + 1);
      if (!textAfterOpen.includes('}}')) {
        lastOpenPos = i - 1;
        break;
      }
    }
  }
  
  if (lastOpenPos === -1) {
    showPlaceholderPopup.value = false;
    placeholderPrefix.value = '';
    return;
  }
  
  const textAfterOpen = textBeforeCursor.substring(lastOpenPos + 2);
  placeholderPrefix.value = textAfterOpen;
  
  const filtered = filteredPlaceholders.value;
  if (filtered.length > 0) {
    nextTick(() => {
      const position = calculateCursorPosition(textarea);
      popupStyle.value = {
        top: `${position.top}px`,
        left: `${position.left}px`
      };
      showPlaceholderPopup.value = true;
      selectedPlaceholderIndex.value = 0;
    });
  } else {
    showPlaceholderPopup.value = false;
  }
};

// Handle input
const handlePromptInput = () => {
  if (placeholderPopupTimer) {
    clearTimeout(placeholderPopupTimer);
  }
  placeholderPopupTimer = setTimeout(() => {
    checkAndShowPlaceholderPopup();
  }, 50);
};

// Insert placeholder
const insertPlaceholder = (placeholderName: string, fromPopup: boolean = false) => {
  const textarea = getTextareaElement();
  if (!textarea) return;
  
  showPlaceholderPopup.value = false;
  placeholderPrefix.value = '';
  selectedPlaceholderIndex.value = 0;
  
  nextTick(() => {
    const cursorPos = textarea.selectionStart;
    const currentValue = formData.value.config.system_prompt || '';
    const textBeforeCursor = currentValue.substring(0, cursorPos);
    const textAfterCursor = currentValue.substring(cursorPos);
    
    // Only search for {{ and replace when selecting from dropdown list
    if (fromPopup) {
      let lastOpenPos = -1;
      for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
        if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
          lastOpenPos = i - 1;
          break;
        }
      }
      
      if (lastOpenPos !== -1) {
        const textBeforeOpen = currentValue.substring(0, lastOpenPos);
        const newValue = textBeforeOpen + `{{${placeholderName}}}` + textAfterCursor;
        formData.value.config.system_prompt = newValue;
        
        nextTick(() => {
          const newCursorPos = textBeforeOpen.length + placeholderName.length + 4;
          textarea.setSelectionRange(newCursorPos, newCursorPos);
          textarea.focus();
        });
        return;
      }
    }
    
    // Directly insert complete placeholder at cursor position
    const newValue = textBeforeCursor + `{{${placeholderName}}}` + textAfterCursor;
    formData.value.config.system_prompt = newValue;
    
    nextTick(() => {
      const newCursorPos = cursorPos + placeholderName.length + 4;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
      textarea.focus();
    });
  });
};

// Get context template textarea element
const getContextTemplateTextareaElement = (): HTMLTextAreaElement | null => {
  if (contextTemplateTextareaRef.value) {
    if (contextTemplateTextareaRef.value.$el) {
      return contextTemplateTextareaRef.value.$el.querySelector('textarea');
    }
    if (contextTemplateTextareaRef.value instanceof HTMLTextAreaElement) {
      return contextTemplateTextareaRef.value;
    }
  }
  return null;
};

// Calculate context template cursor position
const calculateContextCursorPosition = (textarea: HTMLTextAreaElement) => {
  const cursorPos = textarea.selectionStart;
  const textBeforeCursor = formData.value.config.context_template.substring(0, cursorPos);
  
  const style = window.getComputedStyle(textarea);
  const textareaRect = textarea.getBoundingClientRect();
  
  const lineHeight = parseFloat(style.lineHeight) || 20;
  const paddingTop = parseFloat(style.paddingTop) || 0;
  const paddingLeft = parseFloat(style.paddingLeft) || 0;
  
  const lines = textBeforeCursor.split('\n');
  const currentLine = lines.length - 1;
  const currentLineText = lines[currentLine];
  
  const span = document.createElement('span');
  span.style.font = style.font;
  span.style.visibility = 'hidden';
  span.style.position = 'absolute';
  span.style.whiteSpace = 'pre';
  span.textContent = currentLineText;
  document.body.appendChild(span);
  const textWidth = span.offsetWidth;
  document.body.removeChild(span);
  
  const scrollTop = textarea.scrollTop;
  const top = textareaRect.top + paddingTop + (currentLine * lineHeight) - scrollTop + lineHeight + 4;
  const scrollLeft = textarea.scrollLeft;
  const left = textareaRect.left + paddingLeft + textWidth - scrollLeft;
  
  return { top, left };
};

// Check and show context template placeholder popup
const checkAndShowContextPlaceholderPopup = () => {
  const textarea = getContextTemplateTextareaElement();
  if (!textarea) return;
  
  const cursorPos = textarea.selectionStart;
  const textBeforeCursor = formData.value.config.context_template.substring(0, cursorPos);
  
  let lastOpenPos = -1;
  for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
    if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
      const textAfterOpen = textBeforeCursor.substring(i + 1);
      if (!textAfterOpen.includes('}}')) {
        lastOpenPos = i - 1;
        break;
      }
    }
  }
  
  if (lastOpenPos === -1) {
    showContextPlaceholderPopup.value = false;
    contextPlaceholderPrefix.value = '';
    return;
  }
  
  const textAfterOpen = textBeforeCursor.substring(lastOpenPos + 2);
  contextPlaceholderPrefix.value = textAfterOpen;
  
  const filtered = filteredContextPlaceholders.value;
  if (filtered.length > 0) {
    nextTick(() => {
      const position = calculateContextCursorPosition(textarea);
      contextPopupStyle.value = {
        top: `${position.top}px`,
        left: `${position.left}px`
      };
      showContextPlaceholderPopup.value = true;
      selectedContextPlaceholderIndex.value = 0;
    });
  } else {
    showContextPlaceholderPopup.value = false;
  }
};

// Handle context template input
const handleContextTemplateInput = () => {
  if (contextPlaceholderPopupTimer) {
    clearTimeout(contextPlaceholderPopupTimer);
  }
  contextPlaceholderPopupTimer = setTimeout(() => {
    checkAndShowContextPlaceholderPopup();
  }, 50);
};

// Insert context template placeholder
const insertContextPlaceholder = (placeholderName: string, fromPopup: boolean = false) => {
  const textarea = getContextTemplateTextareaElement();
  if (!textarea) return;
  
  showContextPlaceholderPopup.value = false;
  contextPlaceholderPrefix.value = '';
  selectedContextPlaceholderIndex.value = 0;
  
  nextTick(() => {
    const cursorPos = textarea.selectionStart;
    const currentValue = formData.value.config.context_template || '';
    const textBeforeCursor = currentValue.substring(0, cursorPos);
    const textAfterCursor = currentValue.substring(cursorPos);
    
    // Only search for {{ and replace when selecting from dropdown list
    if (fromPopup) {
      let lastOpenPos = -1;
      for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
        if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
          lastOpenPos = i - 1;
          break;
        }
      }
      
      if (lastOpenPos !== -1) {
        const textBeforeOpen = currentValue.substring(0, lastOpenPos);
        const newValue = textBeforeOpen + `{{${placeholderName}}}` + textAfterCursor;
        formData.value.config.context_template = newValue;
        
        nextTick(() => {
          const newCursorPos = textBeforeOpen.length + placeholderName.length + 4;
          textarea.setSelectionRange(newCursorPos, newCursorPos);
          textarea.focus();
        });
        return;
      }
    }
    
    // Directly insert complete placeholder at cursor position
    const newValue = textBeforeCursor + `{{${placeholderName}}}` + textAfterCursor;
    formData.value.config.context_template = newValue;
    
    nextTick(() => {
      const newCursorPos = cursorPos + placeholderName.length + 4;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
      textarea.focus();
    });
  });
};

// Generic get textarea element
const getGenericTextareaElement = (type: 'rewriteSystem' | 'rewriteUser' | 'fallback'): HTMLTextAreaElement | null => {
  const refMap = {
    rewriteSystem: rewriteSystemTextareaRef,
    rewriteUser: rewriteUserTextareaRef,
    fallback: fallbackPromptTextareaRef,
  };
  const ref = refMap[type];
  if (ref.value) {
    if (ref.value.$el) {
      return ref.value.$el.querySelector('textarea');
    }
    if (ref.value instanceof HTMLTextAreaElement) {
      return ref.value;
    }
  }
  return null;
};

// Generic calculate cursor position
const calculateGenericCursorPosition = (textarea: HTMLTextAreaElement, fieldValue: string) => {
  const cursorPos = textarea.selectionStart;
  const textBeforeCursor = fieldValue.substring(0, cursorPos);
  const lines = textBeforeCursor.split('\n');
  const currentLine = lines.length - 1;
  const currentLineText = lines[currentLine];
  
  const textareaRect = textarea.getBoundingClientRect();
  const style = window.getComputedStyle(textarea);
  const lineHeight = parseFloat(style.lineHeight) || 20;
  const paddingTop = parseFloat(style.paddingTop) || 0;
  const paddingLeft = parseFloat(style.paddingLeft) || 0;
  
  const span = document.createElement('span');
  span.style.font = style.font;
  span.style.visibility = 'hidden';
  span.style.position = 'absolute';
  span.style.whiteSpace = 'pre';
  span.textContent = currentLineText;
  document.body.appendChild(span);
  const textWidth = span.offsetWidth;
  document.body.removeChild(span);
  
  const scrollTop = textarea.scrollTop;
  const top = textareaRect.top + paddingTop + (currentLine * lineHeight) - scrollTop + lineHeight + 4;
  const scrollLeft = textarea.scrollLeft;
  const left = textareaRect.left + paddingLeft + textWidth - scrollLeft;
  
  return { top, left };
};

// Generic check and show placeholder popup
const checkAndShowGenericPlaceholderPopup = (
  type: 'rewriteSystem' | 'rewriteUser' | 'fallback',
  popup: typeof rewriteSystemPopup,
  fieldKey: keyof typeof formData.value.config,
  filteredPlaceholders: PlaceholderDefinition[]
) => {
  const textarea = getGenericTextareaElement(type);
  if (!textarea) return;
  
  const cursorPos = textarea.selectionStart;
  const fieldValue = String(formData.value.config[fieldKey] || '');
  const textBeforeCursor = fieldValue.substring(0, cursorPos);
  
  let lastOpenPos = -1;
  for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
    if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
      const textAfterOpen = textBeforeCursor.substring(i + 1);
      if (!textAfterOpen.includes('}}')) {
        lastOpenPos = i - 1;
        break;
      }
    }
  }
  
  if (lastOpenPos === -1) {
    popup.value.show = false;
    popup.value.prefix = '';
    return;
  }
  
  const textAfterOpen = textBeforeCursor.substring(lastOpenPos + 2);
  popup.value.prefix = textAfterOpen;
  
  if (filteredPlaceholders.length > 0) {
    nextTick(() => {
      const position = calculateGenericCursorPosition(textarea, fieldValue);
      popup.value.style = {
        top: `${position.top}px`,
        left: `${position.left}px`
      };
      popup.value.show = true;
      popup.value.selectedIndex = 0;
    });
  } else {
    popup.value.show = false;
  }
};

// Handle rewrite system prompt input
const handleRewriteSystemInput = () => {
  if (rewriteSystemPopup.value.timer) {
    clearTimeout(rewriteSystemPopup.value.timer);
  }
  rewriteSystemPopup.value.timer = setTimeout(() => {
    checkAndShowGenericPlaceholderPopup('rewriteSystem', rewriteSystemPopup, 'rewrite_prompt_system', filteredRewriteSystemPlaceholders.value);
  }, 50);
};

// Handle rewrite user prompt input
const handleRewriteUserInput = () => {
  if (rewriteUserPopup.value.timer) {
    clearTimeout(rewriteUserPopup.value.timer);
  }
  rewriteUserPopup.value.timer = setTimeout(() => {
    checkAndShowGenericPlaceholderPopup('rewriteUser', rewriteUserPopup, 'rewrite_prompt_user', filteredRewriteUserPlaceholders.value);
  }, 50);
};

// Handle fallback prompt input
const handleFallbackPromptInput = () => {
  if (fallbackPromptPopup.value.timer) {
    clearTimeout(fallbackPromptPopup.value.timer);
  }
  fallbackPromptPopup.value.timer = setTimeout(() => {
    checkAndShowGenericPlaceholderPopup('fallback', fallbackPromptPopup, 'fallback_prompt', filteredFallbackPlaceholders.value);
  }, 50);
};

// Generic insert placeholder
const insertGenericPlaceholder = (type: 'rewriteSystem' | 'rewriteUser' | 'fallback', placeholderName: string, fromPopup: boolean = false) => {
  const textarea = getGenericTextareaElement(type);
  if (!textarea) return;
  
  const popupMap = {
    rewriteSystem: rewriteSystemPopup,
    rewriteUser: rewriteUserPopup,
    fallback: fallbackPromptPopup,
  };
  const fieldKeyMap: Record<string, keyof typeof formData.value.config> = {
    rewriteSystem: 'rewrite_prompt_system',
    rewriteUser: 'rewrite_prompt_user',
    fallback: 'fallback_prompt',
  };
  
  const popup = popupMap[type];
  const fieldKey = fieldKeyMap[type];
  
  popup.value.show = false;
  popup.value.prefix = '';
  popup.value.selectedIndex = 0;
  
  nextTick(() => {
    const cursorPos = textarea.selectionStart;
    const currentValue = String(formData.value.config[fieldKey] || '');
    const textBeforeCursor = currentValue.substring(0, cursorPos);
    const textAfterCursor = currentValue.substring(cursorPos);
    
    // Only search for {{ and replace when selecting from dropdown list
    if (fromPopup) {
      let lastOpenPos = -1;
      for (let i = textBeforeCursor.length - 1; i >= 1; i--) {
        if (textBeforeCursor[i] === '{' && textBeforeCursor[i - 1] === '{') {
          lastOpenPos = i - 1;
          break;
        }
      }
      
      if (lastOpenPos !== -1) {
        const textBeforeOpen = currentValue.substring(0, lastOpenPos);
        const newValue = textBeforeOpen + `{{${placeholderName}}}` + textAfterCursor;
        (formData.value.config as any)[fieldKey] = newValue;
        
        nextTick(() => {
          const newCursorPos = textBeforeOpen.length + placeholderName.length + 4;
          textarea.setSelectionRange(newCursorPos, newCursorPos);
          textarea.focus();
        });
        return;
      }
    }
    
    // Directly insert complete placeholder at cursor position
    const newValue = textBeforeCursor + `{{${placeholderName}}}` + textAfterCursor;
    (formData.value.config as any)[fieldKey] = newValue;
    
    nextTick(() => {
      const newCursorPos = cursorPos + placeholderName.length + 4;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
      textarea.focus();
    });
  });
};

// Setup context template textarea event listeners
const setupContextTemplateEventListeners = () => {
  nextTick(() => {
    const textarea = getContextTemplateTextareaElement();
    if (textarea) {
      textarea.addEventListener('keydown', (e: KeyboardEvent) => {
        if (showContextPlaceholderPopup.value && filteredContextPlaceholders.value.length > 0) {
          if (e.key === 'ArrowDown') {
            e.preventDefault();
            e.stopPropagation();
            if (selectedContextPlaceholderIndex.value < filteredContextPlaceholders.value.length - 1) {
              selectedContextPlaceholderIndex.value++;
            } else {
              selectedContextPlaceholderIndex.value = 0;
            }
          } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            e.stopPropagation();
            if (selectedContextPlaceholderIndex.value > 0) {
              selectedContextPlaceholderIndex.value--;
            } else {
              selectedContextPlaceholderIndex.value = filteredContextPlaceholders.value.length - 1;
            }
          } else if (e.key === 'Enter' || e.key === 'Tab') {
            e.preventDefault();
            e.stopPropagation();
            const selected = filteredContextPlaceholders.value[selectedContextPlaceholderIndex.value];
            if (selected) {
              insertContextPlaceholder(selected.name, true);
            }
          } else if (e.key === 'Escape') {
            e.preventDefault();
            e.stopPropagation();
            showContextPlaceholderPopup.value = false;
            contextPlaceholderPrefix.value = '';
          }
        }
      }, true);
    }
  });
};

// Setup textarea event listeners
const setupTextareaEventListeners = () => {
  nextTick(() => {
    const textarea = getTextareaElement();
    if (textarea) {
      textarea.addEventListener('keydown', (e: KeyboardEvent) => {
        if (showPlaceholderPopup.value && filteredPlaceholders.value.length > 0) {
          if (e.key === 'ArrowDown') {
            e.preventDefault();
            e.stopPropagation();
            if (selectedPlaceholderIndex.value < filteredPlaceholders.value.length - 1) {
              selectedPlaceholderIndex.value++;
            } else {
              selectedPlaceholderIndex.value = 0;
            }
          } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            e.stopPropagation();
            if (selectedPlaceholderIndex.value > 0) {
              selectedPlaceholderIndex.value--;
            } else {
              selectedPlaceholderIndex.value = filteredPlaceholders.value.length - 1;
            }
          } else if (e.key === 'Enter' || e.key === 'Tab') {
            e.preventDefault();
            e.stopPropagation();
            const selected = filteredPlaceholders.value[selectedPlaceholderIndex.value];
            if (selected) {
              insertPlaceholder(selected.name, true);
            }
          } else if (e.key === 'Escape') {
            e.preventDefault();
            e.stopPropagation();
            showPlaceholderPopup.value = false;
            placeholderPrefix.value = '';
          }
        }
      }, true);
    }
  });
};

// Generic setup textarea event listeners
const setupGenericTextareaEventListeners = (
  type: 'rewriteSystem' | 'rewriteUser' | 'fallback',
  popup: typeof rewriteSystemPopup,
  filteredPlaceholders: () => PlaceholderDefinition[]
) => {
  nextTick(() => {
    const textarea = getGenericTextareaElement(type);
    if (textarea) {
      textarea.addEventListener('keydown', (e: KeyboardEvent) => {
        const filtered = filteredPlaceholders();
        if (popup.value.show && filtered.length > 0) {
          if (e.key === 'ArrowDown') {
            e.preventDefault();
            e.stopPropagation();
            if (popup.value.selectedIndex < filtered.length - 1) {
              popup.value.selectedIndex++;
            } else {
              popup.value.selectedIndex = 0;
            }
          } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            e.stopPropagation();
            if (popup.value.selectedIndex > 0) {
              popup.value.selectedIndex--;
            } else {
              popup.value.selectedIndex = filtered.length - 1;
            }
          } else if (e.key === 'Enter' || e.key === 'Tab') {
            e.preventDefault();
            e.stopPropagation();
            const selected = filtered[popup.value.selectedIndex];
            if (selected) {
              insertGenericPlaceholder(type, selected.name, true);
            }
          } else if (e.key === 'Escape') {
            e.preventDefault();
            e.stopPropagation();
            popup.value.show = false;
            popup.value.prefix = '';
          }
        }
      }, true);
    }
  });
};

// Handle placeholder tag click
const handlePlaceholderClick = (type: 'system' | 'context' | 'rewriteSystem' | 'rewriteUser' | 'fallback', placeholderName: string) => {
  if (type === 'system') {
    insertPlaceholder(placeholderName);
  } else if (type === 'context') {
    insertContextPlaceholder(placeholderName);
  } else {
    insertGenericPlaceholder(type, placeholderName);
  }
};

// Watch visible changes to setup event listeners
watch(() => props.visible, (val) => {
  if (val) {
    nextTick(() => {
      setupTextareaEventListeners();
      setupContextTemplateEventListeners();
      setupGenericTextareaEventListeners('rewriteSystem', rewriteSystemPopup, () => filteredRewriteSystemPlaceholders.value);
      setupGenericTextareaEventListeners('rewriteUser', rewriteUserPopup, () => filteredRewriteUserPlaceholders.value);
      setupGenericTextareaEventListeners('fallback', fallbackPromptPopup, () => filteredFallbackPlaceholders.value);
    });
  }
});

// Template selection handler functions
const handleSystemPromptTemplateSelect = (template: string) => {
  formData.value.config.system_prompt = template;
};

const handleContextTemplateSelect = (template: string) => {
  formData.value.config.context_template = template;
};

const handleRewriteSystemTemplateSelect = (template: string) => {
  formData.value.config.rewrite_prompt_system = template;
};

const handleRewriteUserTemplateSelect = (template: string) => {
  formData.value.config.rewrite_prompt_user = template;
};

const handleFallbackResponseTemplateSelect = (template: string) => {
  formData.value.config.fallback_response = template;
};

const handleFallbackPromptTemplateSelect = (template: string) => {
  formData.value.config.fallback_prompt = template;
};

// Helper function: check if prompt contains specified placeholder
const hasPlaceholder = (text: string | undefined, placeholder: string): boolean => {
  if (!text) return false;
  return text.includes(`{{${placeholder}}}`);
};

const handleSave = async () => {
  // Validate required fields (built-in agents don't validate name and system prompt)
  if (!isBuiltinAgent.value) {
    if (!formData.value.name || !formData.value.name.trim()) {
      MessagePlugin.error(t('agent.editor.nameRequired'));
      currentSection.value = 'basic';
      return;
    }

    // Custom agents must fill in system prompt
    if (!formData.value.config.system_prompt || !formData.value.config.system_prompt.trim()) {
      MessagePlugin.error(t('agent.editor.systemPromptRequired'));
      currentSection.value = 'basic';
      return;
    }

    // Custom agents in normal mode must fill in context template
    if (!isAgentMode.value && (!formData.value.config.context_template || !formData.value.config.context_template.trim())) {
      MessagePlugin.error(t('agent.editor.contextTemplateRequired'));
      currentSection.value = 'basic';
      return;
    }
  }

  // Validate placeholders (normal mode + knowledge base enabled)
  if (!isAgentMode.value && hasKnowledgeBase.value) {
    const contextTemplate = formData.value.config.context_template || '';
    if (!hasPlaceholder(contextTemplate, 'contexts')) {
      MessagePlugin.error(t('agent.editor.contextsMissing'));
      currentSection.value = 'basic';
      return;
    }
    if (!hasPlaceholder(contextTemplate, 'query')) {
      MessagePlugin.error(t('agent.editor.queryMissingInContext'));
      currentSection.value = 'basic';
      return;
    }
  }

  // Validate placeholders (Agent mode + knowledge base enabled)
  if (isAgentMode.value && hasKnowledgeBase.value) {
    const systemPrompt = formData.value.config.system_prompt || '';
    if (!hasPlaceholder(systemPrompt, 'knowledge_bases')) {
      MessagePlugin.warning(t('agent.editor.knowledgeBasesMissing'));
    }
  }

  // Validate placeholders (normal mode + multi-turn conversation rewrite enabled)
  if (!isAgentMode.value && formData.value.config.multi_turn_enabled && formData.value.config.enable_rewrite) {
    const rewritePrompt = formData.value.config.rewrite_prompt_user || '';
    // Only validate when user has customized rewrite prompt
    if (rewritePrompt.trim()) {
      if (!hasPlaceholder(rewritePrompt, 'query')) {
        MessagePlugin.error(t('agent.editor.queryMissingInRewrite'));
        currentSection.value = 'conversation';
        return;
      }
      if (!hasPlaceholder(rewritePrompt, 'conversation')) {
        MessagePlugin.error(t('agent.editor.conversationMissing'));
        currentSection.value = 'conversation';
        return;
      }
    }
  }

  // Validate placeholders (fallback strategy is model generation)
  if (!isAgentMode.value && formData.value.config.fallback_strategy === 'model') {
    const fallbackPrompt = formData.value.config.fallback_prompt || '';
    // Only validate when user has customized fallback prompt
    if (fallbackPrompt.trim() && !hasPlaceholder(fallbackPrompt, 'query')) {
      MessagePlugin.error(t('agent.editor.queryMissingInFallback'));
      currentSection.value = 'retrieval';
      return;
    }
  }

  if (!formData.value.config.model_id) {
    MessagePlugin.error(t('agent.editor.modelRequired'));
    currentSection.value = 'model';
    return;
  }

  // Validate ReRank model (required when needed)
  if (needsRerankModel.value && !formData.value.config.rerank_model_id) {
    MessagePlugin.error(t('agent.editor.rerankModelRequired'));
    currentSection.value = 'knowledge';
    return;
  }

  // Filter empty suggested prompts
  if (formData.value.config.suggested_prompts) {
    formData.value.config.suggested_prompts = formData.value.config.suggested_prompts.filter((p: string) => p.trim() !== '');
  }

  saving.value = true;
  try {
    if (props.mode === 'create') {
      await createAgent(formData.value);
      MessagePlugin.success(t('agent.messages.created'));
    } else {
      await updateAgent(formData.value.id, formData.value);
      MessagePlugin.success(t('agent.messages.updated'));
    }
    emit('success');
    handleClose();
  } catch (e: any) {
    MessagePlugin.error(e?.message || t('agent.messages.saveFailed'));
  } finally {
    saving.value = false;
  }
};
</script>

<style scoped lang="less">
// Reuse knowledge base creation styles
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.settings-modal {
  position: relative;
  width: 90vw;
  max-width: 1100px;
  height: 85vh;
  max-height: 750px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 32px;
  height: 32px;
  border: none;
  background: #f5f5f5;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: all 0.2s ease;
  z-index: 10;

  &:hover {
    background: #e5e5e5;
    color: #000;
  }
}

.settings-container {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.settings-sidebar {
  width: 200px;
  background: #fafafa;
  border-right: 1px solid #e5e5e5;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid #e5e5e5;
}

.sidebar-title {
  margin: 0;
  font-family: "PingFang SC";
  font-size: 18px;
  font-weight: 600;
  color: #000000e6;
}

.settings-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 4px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: "PingFang SC";
  font-size: 14px;
  color: #00000099;

  &:hover {
    background: #f0f0f0;
  }

  &.active {
    background: #07c05f1a;
    color: #07c05f;
    font-weight: 500;
  }
}

.nav-icon {
  margin-right: 8px;
  font-size: 18px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-label {
  flex: 1;
}

.settings-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 24px 32px;
}

.section {
  width: 100%;
}

// section-header style consistent with knowledge base settings
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

// settings-group style consistent with knowledge base settings
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

  &.setting-row-vertical {
    flex-direction: column;
    gap: 12px;
    
    .setting-info {
      max-width: 100%;
      padding-right: 0;
    }
  }
}

.setting-info {
  flex: 1;
  max-width: 55%;
  padding-right: 24px;

  &.full-width {
    max-width: 100%;
    padding-right: 0;
  }

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

    .required {
      color: #fa5151;
      margin-left: 2px;
    }
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
  min-width: 360px;
  display: flex;
  justify-content: flex-end;
  align-items: flex-start;

  &.setting-control-full {
    width: 100%;
    min-width: 100%;
    justify-content: flex-start;
  }

  // Make select and input fill the control area
  :deep(.t-select),
  :deep(.t-input),
  :deep(.t-textarea) {
    width: 100%;
  }

  :deep(.t-input-number) {
    width: 120px;
  }
}

// Name input with avatar preview
.name-input-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;

  .name-input {
    flex: 1;
  }
}

.settings-footer {
  padding: 16px 32px;
  border-top: 1px solid #e5e5e5;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  flex-shrink: 0;
}

// Mode hint style
.mode-hint {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  background: #f0faf5;
  border-radius: 6px;
  border: 1px solid #d4f0e2;
  color: #07c05f;
  font-size: 13px;
  line-height: 1.5;
}

// Transition animation
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .settings-modal {
    transform: scale(0.95);
  }
}

// Slider style
.slider-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;

  :deep(.t-slider) {
    flex: 1;
  }
}

.slider-value {
  width: 40px;
  text-align: right;
  font-family: monospace;
  font-size: 14px;
  color: #333;
}

// Suggested prompts list
.suggested-prompts-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.prompt-item {
  display: flex;
  align-items: center;
  gap: 8px;

  :deep(.t-input) {
    flex: 1;
  }
}

// Radio-group style optimization, consistent with project theme
:deep(.t-radio-group) {
  .t-radio-group--filled {
    background: #f5f5f5;
  }
  .t-radio-button {
    border-color: #d9d9d9;

    &:hover:not(.t-is-disabled) {
      border-color: #07c05f;
      color: #07c05f;
    }

    &.t-is-checked {
      background: #07c05f;
      border-color: #07c05f;
      color: #fff;

      &:hover:not(.t-is-disabled) {
        background: #05a04f;
        border-color: #05a04f;
        color: #fff;
      }
    }

    // Disabled state style
    &.t-is-disabled {
      background: #f5f5f5;
      border-color: #d9d9d9;
      color: #00000040;
      cursor: not-allowed;
      opacity: 0.6;

      &.t-is-checked {
        background: #f0f0f0;
        border-color: #d9d9d9;
        color: #00000066;
      }
    }
  }
}

// Tool selection style
.tools-checkbox-group {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  width: 100%;
}

.tool-checkbox-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;

  &:hover {
    border-color: #07c05f;
    background: #f0faf5;
  }

  :deep(.t-checkbox__input) {
    margin-top: 2px;
  }

  :deep(.t-checkbox__label) {
    flex: 1;
  }
}

.tool-item-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tool-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.tool-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.5;
}

.tool-disabled-hint {
  font-size: 11px;
  color: #f5a623;
  font-style: italic;
}

.tool-disabled {
  opacity: 0.6;
  
  .tool-name, .tool-desc {
    color: #999;
  }
}

// Checkbox selected style
:deep(.t-checkbox) {
  &.t-is-checked {
    .t-checkbox__input {
      border-color: #07c05f;
      background-color: #07c05f;
    }
  }
  
  &:hover:not(.t-is-disabled) {
    .t-checkbox__input {
      border-color: #07c05f;
    }
  }
}

// Switch style
:deep(.t-switch) {
  &.t-is-checked {
    background-color: #07c05f;
    
    &:hover:not(.t-is-disabled) {
      background-color: #05a04f;
    }
  }
}

// Slider style
:deep(.t-slider) {
  .t-slider__track {
    background-color: #07c05f;
  }
  
  .t-slider__button {
    border-color: #07c05f;
  }
}

// Button theme style
:deep(.t-button--theme-primary) {
  background-color: #07c05f;
  border-color: #07c05f;
  
  &:hover:not(.t-is-disabled) {
    background-color: #05a04f;
    border-color: #05a04f;
  }
}

// Input/Select focus style
:deep(.t-input),
:deep(.t-textarea),
:deep(.t-select) {
  &.t-is-focused,
  &:focus-within {
    border-color: #07c05f;
    box-shadow: 0 0 0 2px rgba(7, 192, 95, 0.1);
  }
}

// textarea and template selector container
.textarea-with-template {
  position: relative;
  width: 100%;
}

// System prompt input style
.system-prompt-textarea {
  width: 100%;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;

  :deep(textarea) {
    resize: vertical !important;
    min-height: 200px;
  }
}

// Placeholder tag group style
.placeholder-tags {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  line-height: 1.4;
  overflow-x: auto;
  white-space: nowrap;
  padding-bottom: 4px;
  
  // Hide scrollbar but keep scrollable
  scrollbar-width: thin;
  &::-webkit-scrollbar {
    height: 4px;
  }
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
  }

  .placeholder-label {
    color: var(--td-text-color-secondary, #666);
    flex-shrink: 0;
  }

  .placeholder-hint {
    color: var(--td-text-color-placeholder, #999);
    font-size: 11px;
    user-select: none;
    flex-shrink: 0;
  }

  .placeholder-tag {
    display: inline-flex;
    align-items: center;
    padding: 1px 5px;
    border-radius: 3px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 11px;
    color: var(--td-text-color-primary, #333);
    background-color: var(--td-bg-color-secondarycontainer, #f3f3f3);
    cursor: pointer;
    transition: all 0.2s;
    user-select: none;
    border: 1px solid transparent;
    flex-shrink: 0;

    &:hover {
      color: var(--td-brand-color, #0052d9);
      background-color: var(--td-brand-color-light, #ecf2fe);
      border-color: var(--td-brand-color-focus, #d0e0fd);
    }

    &:active {
      background-color: var(--td-brand-color-focus, #d0e0fd);
    }
  }
}

.placeholder-popup-wrapper {
  position: fixed;
  z-index: 10001;
  pointer-events: auto;
}

.placeholder-popup {
  background: var(--td-bg-color-container, #fff);
  border: 1px solid var(--td-component-stroke, #e5e7eb);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  max-width: 320px;
  max-height: 240px;
  overflow-y: auto;
  padding: 4px;
}

.placeholder-item {
  padding: 6px 10px;
  cursor: pointer;
  transition: background-color 0.15s;
  border-radius: 4px;

  &:hover,
  &.active {
    background-color: var(--td-bg-color-container-hover, #f5f7fa);
  }

  .placeholder-name {
    margin-bottom: 2px;

    code {
      background: var(--td-bg-color-container-hover, #f5f7fa);
      padding: 2px 5px;
      border-radius: 3px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 11px;
      color: var(--td-brand-color, #0052d9);
    }
  }

  .placeholder-desc {
    font-size: 11px;
    color: var(--td-text-color-secondary, #666);
  }
}

// Built-in agent notice
.builtin-agent-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 8px;
  margin-bottom: 16px;
  color: #d46b08;
  font-size: 14px;

  .t-icon {
    font-size: 16px;
    flex-shrink: 0;
  }
}

// Built-in agent avatar
.builtin-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
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

// Prompt toggle
.prompt-toggle {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;

  .prompt-toggle-label {
    font-size: 13px;
    color: #666;
  }
}

// Prompt disabled hint
.prompt-disabled-hint {
  color: #999;
  font-size: 13px;
  font-style: italic;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 6px;
}

// System prompt Tabs
.system-prompt-tabs {
  width: 100%;

  .prompt-variant-tabs {
    :deep(.t-tabs__nav) {
      margin-bottom: 12px;
    }
  }
}

// Knowledge base option style
.kb-option-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.kb-option-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
  
  // Document KB - Greenish
  &.doc-icon {
    color: #10b981;
  }
  
  // FAQ KB - Blueish
  &.faq-icon {
    color: #0052d9;
  }
}

.kb-option-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kb-option-count {
  flex-shrink: 0;
  font-size: 12px;
  color: #999;
}

// FAQ strategy section style
.faq-strategy-section {
  margin-top: 24px;
  padding: 16px;
  background: rgba(0, 82, 217, 0.04);
  border: 1px solid rgba(0, 82, 217, 0.15);
  border-radius: 8px;
}

.faq-strategy-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 600;
  color: #0052d9;
  
  .faq-icon {
    font-size: 18px;
  }
  
  .help-icon {
    font-size: 14px;
    color: #999;
    cursor: help;
  }
}

.faq-strategy-section .setting-row {
  padding: 12px 0;
  border-bottom: 1px solid rgba(0, 82, 217, 0.1);
  
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
  
  &:first-of-type {
    padding-top: 0;
  }
}
</style>
