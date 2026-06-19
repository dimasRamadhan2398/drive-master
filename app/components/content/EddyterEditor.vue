<script setup lang="tsx">
import { useRef, onMounted, onBeforeUnmount, watch } from "vue";
import { createRoot } from "react-dom/client";
import React from "react";
import ConfigurableEditor, {
  editorConfig as defaultEditorConfig,
  EditorConfigTypes,
} from "eddyter";

interface Props {
  modelValue: string;
  placeholder?: string;
  apiKey?: string;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  placeholder: "Write your content here...",
  apiKey: "",
});

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const editorContainer = ref<HTMLElement | null>(null);
const reactRootRef = ref<any>(null);
const isClient = ref(false);

// Create a simple config based on placeholder
const simpleConfig: EditorConfigTypes = {
  ...defaultEditorConfig,
  enableToolbar: true,
  toolbarOptions: {
    ...defaultEditorConfig.toolbarOptions,
    enableAIChat: false, // Disable AI features without API key
  },
};

// Initialize React root
onMounted(() => {
  isClient.value = true;

  if (editorContainer.value) {
    const root = createRoot(editorContainer.value);
    reactRootRef.value = root;

    // Create the React element
    const editorElement = React.createElement(ConfigurableEditor, {
      config: simpleConfig,
      defaultFontFamilies: defaultEditorConfig.defaultFontFamilies,
      mentionUserList: [],
      initialContent: props.modelValue,
      onChange: (html: string) => {
        emit("update:modelValue", html);
      },
      apiKey: props.apiKey,
      editorStyle: {
        minHeight: "150px",
        padding: "8px 12px",
      },
      contentClassName: "focus:outline-none",
    });

    root.render(editorElement);
  }
});

// Watch for external modelValue changes
watch(
  () => props.modelValue,
  (newValue) => {
    // The eddyter editor handles its own state, so we don't need to do anything here
    // unless we want to force update from outside
  }
);

// Cleanup
onBeforeUnmount(() => {
  if (reactRootRef.value) {
    reactRootRef.value.unmount();
  }
});
</script>

<template>
  <div class="border border-default rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary/50 focus-within:border-primary transition-all">
    <ClientOnly>
      <div ref="editorContainer" class="bg-background" />
      <template #fallback>
        <div class="min-h-[200px] flex items-center justify-center bg-muted/30">
          <span class="text-sm text-muted">Loading editor...</span>
        </div>
      </template>
    </ClientOnly>
  </div>
</template>
