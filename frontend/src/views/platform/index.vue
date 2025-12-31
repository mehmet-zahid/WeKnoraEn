<template>
    <div class="main" ref="dropzone">
        <Menu></Menu>
        <RouterView />
        <div class="upload-mask" v-show="ismask">
            <input type="file" style="display: none" ref="uploadInput" accept=".pdf,.docx,.doc,.txt,.md,.jpg,.jpeg,.png,.csv,.xls,.xlsx" />
            <UploadMask></UploadMask>
        </div>
        <!-- Global settings modal, used by all platform sub-routes -->
        <Settings />
    </div>
</template>
<script setup lang="ts">
import Menu from '@/components/menu.vue'
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router'
import useKnowledgeBase from '@/hooks/useKnowledgeBase'
import UploadMask from '@/components/upload-mask.vue'
import Settings from '@/views/settings/Settings.vue'
import { getKnowledgeBaseById } from '@/api/knowledge-base/index'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'

let { requestMethod } = useKnowledgeBase()
const route = useRoute();
let ismask = ref(false)
let uploadInput = ref();
const { t } = useI18n();

// Get current knowledge base ID
const getCurrentKbId = (): string | null => {
    return (route.params as any)?.kbId as string || null
}

// Check knowledge base initialization status
const checkKnowledgeBaseInitialization = async (): Promise<boolean> => {
    const currentKbId = getCurrentKbId();
    
    if (!currentKbId) {
        MessagePlugin.error(t('knowledgeBase.missingId'));
        return false;
    }
    
    try {
        const kbResponse = await getKnowledgeBaseById(currentKbId);
        const kb = kbResponse.data;
        
        if (!kb.embedding_model_id || !kb.summary_model_id) {
            MessagePlugin.warning(t('knowledgeBase.notInitialized'));
            return false;
        }
        return true;
    } catch (error) {
        MessagePlugin.error(t('knowledgeBase.getInfoFailed'));
        return false;
    }
}


// Global drag event handling
const handleGlobalDragEnter = (event: DragEvent) => {
    event.preventDefault();
    if (event.dataTransfer) {
        event.dataTransfer.effectAllowed = 'all';
    }
    ismask.value = true;
}

const handleGlobalDragOver = (event: DragEvent) => {
    event.preventDefault();
    if (event.dataTransfer) {
        event.dataTransfer.dropEffect = 'copy';
    }
    ismask.value = true;
}

const handleGlobalDrop = async (event: DragEvent) => {
    event.preventDefault();
    ismask.value = false;
    
    const DataTransferFiles = event.dataTransfer?.files ? Array.from(event.dataTransfer.files) : [];
    const DataTransferItemList = event.dataTransfer?.items ? Array.from(event.dataTransfer.items) : [];
    
    const isInitialized = await checkKnowledgeBaseInitialization();
    if (!isInitialized) {
        return;
    }
    
    if (DataTransferFiles.length > 0) {
        DataTransferFiles.forEach(file => requestMethod(file, uploadInput));
    } else if (DataTransferItemList.length > 0) {
        DataTransferItemList.forEach(dataTransferItem => {
            const fileEntry = dataTransferItem.webkitGetAsEntry() as FileSystemFileEntry | null;
            if (fileEntry) {
                fileEntry.file((file: File) => requestMethod(file, uploadInput));
            }
        });
    } else {
        MessagePlugin.warning(t('knowledgeBase.dragFileNotText'));
    }
}

// Add global event listeners when component mounts
onMounted(() => {
    document.addEventListener('dragenter', handleGlobalDragEnter, true);
    document.addEventListener('dragover', handleGlobalDragOver, true);
    document.addEventListener('drop', handleGlobalDrop, true);
});

// Remove global event listeners when component unmounts
onUnmounted(() => {
    document.removeEventListener('dragenter', handleGlobalDragEnter, true);
    document.removeEventListener('dragover', handleGlobalDragOver, true);
    document.removeEventListener('drop', handleGlobalDrop, true);
});
</script>
<style lang="less">
.main {
    display: flex;
    width: 100%;
    height: 100%;
    min-width: 600px;
}

.upload-mask {
    background-color: rgba(255, 255, 255, 0.8);
    position: fixed;
    width: 100%;
    height: 100%;
    z-index: 999;
    display: flex;
    justify-content: center;
    align-items: center;
}

img {
    -webkit-user-drag: none;
    -khtml-user-drag: none;
    -moz-user-drag: none;
    -o-user-drag: none;
    user-drag: none;
}
</style>