<script setup lang="ts">
import { Editor } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import { TextStyle } from "@tiptap/extension-text-style";
import Color from "@tiptap/extension-color";
import { onBeforeUnmount, watch } from "vue";

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const isClient = ref(false);
const editor = ref<Editor | null>(null);

// Initialize editor only on client side
onMounted(() => {
  isClient.value = true;
  editor.value = new Editor({
    content: props.modelValue || "",
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [1, 2, 3],
        },
      }),
      Underline,
      TextStyle,
      Color,
    ],
    editorProps: {
      attributes: {
        class: "prose dark:prose-invert max-w-none p-3 min-h-[150px] focus:outline-none",
      },
    },
    onUpdate: ({ editor }) => {
      emit("update:modelValue", editor.getHTML());
    },
  });
});

// Watch for external modelValue changes
watch(
  () => props.modelValue,
  (newValue) => {
    if (editor.value && newValue !== editor.value.getHTML()) {
      editor.value.commands.setContent(newValue || "", {});
    }
  }
);

// Cleanup
onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

<template>
  <div
    class="border border-gray-200 dark:border-gray-800 rounded-lg overflow-hidden bg-white dark:bg-gray-900"
  >
    <!-- Loading state during SSR -->
    <div
      v-if="!isClient"
      class="min-h-[150px] flex items-center justify-center bg-gray-50 dark:bg-gray-950"
    >
      <span class="text-sm text-gray-400">Loading editor...</span>
    </div>

    <template v-else>
      <!-- Toolbar -->
      <div
        v-if="editor"
        class="flex flex-wrap gap-1 border-b border-gray-200 dark:border-gray-800 p-2 bg-gray-50 dark:bg-gray-950"
      >
        <UButton
          size="xs"
          :variant="editor.isActive('bold') ? 'solid' : 'ghost'"
          class="text-black"
          icon="i-lucide-bold"
          @click="editor.chain().focus().toggleBold().run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('italic') ? 'solid' : 'ghost'"
          class="text-black"
          icon="i-lucide-italic"
          @click="editor.chain().focus().toggleItalic().run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('underline') ? 'solid' : 'ghost'"
          icon="i-lucide-underline"
          @click="editor.chain().focus().toggleUnderline().run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('bulletList') ? 'solid' : 'ghost'"
          class="text-black"
          icon="i-lucide-list"
          @click="editor.chain().focus().toggleBulletList().run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('orderedList') ? 'solid' : 'ghost'"
          class="text-black"
          icon="i-lucide-list-ordered"
          @click="editor.chain().focus().toggleOrderedList().run()"
        />
      </div>

      <!-- Editor Content -->
      <EditorContent :editor="editor" />
    </template>
  </div>
</template>

<style>
/* Global styles for Tiptap */
.tiptap:focus {
  outline: none;
}
.tiptap p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #9ca3af;
  pointer-events: none;
  height: 0;
}
</style>
