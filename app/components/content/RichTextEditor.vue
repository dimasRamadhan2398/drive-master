<script setup lang="ts">
import { Editor } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import { TextStyle } from "@tiptap/extension-text-style";
import Color from "@tiptap/extension-color";
import TextAlign from "@tiptap/extension-text-align";
import { onBeforeUnmount, watch } from "vue";

const { t } = useI18n()
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
      TextAlign.configure({
        types: ["heading", "paragraph"],
      }),
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
      <span class="text-sm text-gray-400">{{ t('common.loadingEditor') }}</span>
    </div>

    <template v-else>
      <!-- Toolbar -->
      <div
        v-if="editor"
        class="flex flex-wrap gap-1 border-b border-gray-200 dark:border-gray-800 p-2 bg-gray-50 dark:bg-gray-950"
      >
        <!-- Text Formatting -->
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

        <div class="w-px h-6 bg-gray-300 dark:bg-gray-600 mx-1 self-center" />

        <!-- Lists -->
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

        <div class="w-px h-6 bg-gray-300 dark:bg-gray-600 mx-1 self-center" />

        <!-- Quote -->
        <UButton
          size="xs"
          :variant="editor.isActive('blockquote') ? 'solid' : 'ghost'"
          icon="i-lucide-quote"
          @click="editor.chain().focus().toggleBlockquote().run()"
        />

        <div class="w-px h-6 bg-gray-300 dark:bg-gray-600 mx-1 self-center" />

        <!-- Alignment -->
        <UButton
          size="xs"
          :variant="editor.isActive({ textAlign: 'left' }) ? 'solid' : 'ghost'"
          icon="i-lucide-align-left"
          @click="editor.chain().focus().setTextAlign('left').run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive({ textAlign: 'center' }) ? 'solid' : 'ghost'"
          icon="i-lucide-align-center"
          @click="editor.chain().focus().setTextAlign('center').run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive({ textAlign: 'right' }) ? 'solid' : 'ghost'"
          icon="i-lucide-align-right"
          @click="editor.chain().focus().setTextAlign('right').run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive({ textAlign: 'justify' }) ? 'solid' : 'ghost'"
          icon="i-lucide-align-justify"
          @click="editor.chain().focus().setTextAlign('justify').run()"
        />

        <div class="w-px h-6 bg-gray-300 dark:bg-gray-600 mx-1 self-center" />

        <!-- Headings -->
        <UButton
          size="xs"
          :variant="editor.isActive('heading', { level: 1 }) ? 'solid' : 'ghost'"
          class="text-black font-bold"
          label="H1"
          @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('heading', { level: 2 }) ? 'solid' : 'ghost'"
          class="text-black font-bold"
          label="H2"
          @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
        />
        <UButton
          size="xs"
          :variant="editor.isActive('heading', { level: 3 }) ? 'solid' : 'ghost'"
          class="text-black font-bold"
          label="H3"
          @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
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

/* Quote styling */
.tiptap blockquote {
  border-left: 3px solid #6366f1;
  padding-left: 1rem;
  margin-left: 0;
  color: #6b7280;
  font-style: italic;
}

/* Quote styling */
.tiptap blockquote {
  border-left: 3px solid #6366f1;
  padding-left: 1rem;
  margin-left: 0;
  color: #6b7280;
  font-style: italic;
}
</style>
