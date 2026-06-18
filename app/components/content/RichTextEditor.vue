<script setup lang="ts">
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import Placeholder from "@tiptap/extension-placeholder";
import TextAlign from "@tiptap/extension-text-align";

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

// Generate a unique ID for this editor instance
const editorId = `tiptap-${Math.random().toString(36).substring(2, 9)}`;

// Create editor on client only
const isClient = ref(false);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const editor = shallowRef<any>(null);

// Initialize editor after mount
onMounted(() => {
  isClient.value = true;
  editor.value = useEditor({
    content: props.modelValue || "",
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [1, 2, 3],
        },
      }),
      Underline,
      Placeholder.configure({
        placeholder: props.placeholder || "Write your content here...",
      }),
      TextAlign.configure({
        types: ["heading", "paragraph"],
      }),
    ],
    editorProps: {
      attributes: {
        class:
          "prose prose-sm sm:prose-base max-w-none focus:outline-none min-h-[150px] px-3 py-2",
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
      editor.value.commands.setContent(newValue || "", { emitUpdate: false });
    }
  }
);

onBeforeUnmount(() => {
  editor.value?.destroy();
});

// Toolbar button actions
function toggleBold() {
  editor.value?.chain().focus().toggleBold().run();
}
function toggleItalic() {
  editor.value?.chain().focus().toggleItalic().run();
}
function toggleUnderline() {
  editor.value?.chain().focus().toggleUnderline().run();
}
function toggleStrike() {
  editor.value?.chain().focus().toggleStrike().run();
}
function toggleHeading1() {
  editor.value?.chain().focus().toggleHeading({ level: 1 }).run();
}
function toggleHeading2() {
  editor.value?.chain().focus().toggleHeading({ level: 2 }).run();
}
function toggleHeading3() {
  editor.value?.chain().focus().toggleHeading({ level: 3 }).run();
}
function toggleBulletList() {
  editor.value?.chain().focus().toggleBulletList().run();
}
function toggleOrderedList() {
  editor.value?.chain().focus().toggleOrderedList().run();
}
function toggleBlockquote() {
  editor.value?.chain().focus().toggleBlockquote().run();
}
function toggleCodeBlock() {
  editor.value?.chain().focus().toggleCodeBlock().run();
}
function setAlignLeft() {
  editor.value?.chain().focus().setTextAlign("left").run();
}
function setAlignCenter() {
  editor.value?.chain().focus().setTextAlign("center").run();
}
function setAlignRight() {
  editor.value?.chain().focus().setTextAlign("right").run();
}
function undo() {
  editor.value?.chain().focus().undo().run();
}
function redo() {
  editor.value?.chain().focus().redo().run();
}
</script>

<template>
  <div
    :id="editorId"
    class="border border-default rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary/50 focus-within:border-primary transition-all"
  >
    <!-- Loading state during SSR or before initialization -->
    <div
      v-if="!isClient || !editor"
      class="min-h-[200px] flex items-center justify-center bg-muted/30"
    >
      <span class="text-sm text-muted">Loading editor...</span>
    </div>

    <template v-else>
      <!-- Toolbar -->
      <div
        class="flex flex-wrap items-center gap-0.5 px-2 py-1.5 border-b border-default bg-muted/30"
      >
        <!-- Text Formatting -->
        <UButton
          icon="i-lucide-bold"
          variant="ghost"
          size="xs"
          :color="editor.isActive('bold') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('bold') }"
          @click="toggleBold"
        />
        <UButton
          icon="i-lucide-italic"
          variant="ghost"
          size="xs"
          :color="editor.isActive('italic') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('italic') }"
          @click="toggleItalic"
        />
        <UButton
          icon="i-lucide-underline"
          variant="ghost"
          size="xs"
          :color="editor.isActive('underline') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('underline') }"
          @click="toggleUnderline"
        />
        <UButton
          icon="i-lucide-strikethrough"
          variant="ghost"
          size="xs"
          :color="editor.isActive('strike') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('strike') }"
          @click="toggleStrike"
        />

        <div class="w-px h-5 bg-default mx-1" />

        <!-- Headings -->
        <UButton
          icon="i-lucide-heading-1"
          variant="ghost"
          size="xs"
          :color="editor.isActive('heading', { level: 1 }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('heading', { level: 1 }) }"
          @click="toggleHeading1"
        />
        <UButton
          icon="i-lucide-heading-2"
          variant="ghost"
          size="xs"
          :color="editor.isActive('heading', { level: 2 }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('heading', { level: 2 }) }"
          @click="toggleHeading2"
        />
        <UButton
          icon="i-lucide-heading-3"
          variant="ghost"
          size="xs"
          :color="editor.isActive('heading', { level: 3 }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('heading', { level: 3 }) }"
          @click="toggleHeading3"
        />

        <div class="w-px h-5 bg-default mx-1" />

        <!-- Lists & Blocks -->
        <UButton
          icon="i-lucide-list"
          variant="ghost"
          size="xs"
          :color="editor.isActive('bulletList') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('bulletList') }"
          @click="toggleBulletList"
        />
        <UButton
          icon="i-lucide-list-ordered"
          variant="ghost"
          size="xs"
          :color="editor.isActive('orderedList') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('orderedList') }"
          @click="toggleOrderedList"
        />
        <UButton
          icon="i-lucide-quote"
          variant="ghost"
          size="xs"
          :color="editor.isActive('blockquote') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('blockquote') }"
          @click="toggleBlockquote"
        />
        <UButton
          icon="i-lucide-code"
          variant="ghost"
          size="xs"
          :color="editor.isActive('codeBlock') ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive('codeBlock') }"
          @click="toggleCodeBlock"
        />

        <div class="w-px h-5 bg-default mx-1" />

        <!-- Alignment -->
        <UButton
          icon="i-lucide-align-left"
          variant="ghost"
          size="xs"
          :color="editor.isActive({ textAlign: 'left' }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive({ textAlign: 'left' }) }"
          @click="setAlignLeft"
        />
        <UButton
          icon="i-lucide-align-center"
          variant="ghost"
          size="xs"
          :color="editor.isActive({ textAlign: 'center' }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive({ textAlign: 'center' }) }"
          @click="setAlignCenter"
        />
        <UButton
          icon="i-lucide-align-right"
          variant="ghost"
          size="xs"
          :color="editor.isActive({ textAlign: 'right' }) ? 'primary' : 'neutral'"
          :class="{ 'bg-primary/10': editor.isActive({ textAlign: 'right' }) }"
          @click="setAlignRight"
        />

        <div class="w-px h-5 bg-default mx-1" />

        <!-- Undo/Redo -->
        <UButton
          icon="i-lucide-undo"
          variant="ghost"
          size="xs"
          @click="undo"
        />
        <UButton
          icon="i-lucide-redo"
          variant="ghost"
          size="xs"
          @click="redo"
        />
      </div>

      <!-- Editor Content -->
      <EditorContent :editor="editor" class="bg-background" />
    </template>
  </div>
</template>

<style>
/* Prose styles for the editor content */
.tiptap p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #a1a1aa;
  pointer-events: none;
  height: 0;
}

.tiptap h1 {
  font-size: 1.875rem;
  font-weight: 700;
  margin-top: 1rem;
  margin-bottom: 0.5rem;
}

.tiptap h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-top: 0.875rem;
  margin-bottom: 0.5rem;
}

.tiptap h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 0.75rem;
  margin-bottom: 0.5rem;
}

.tiptap p {
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.tiptap ul,
.tiptap ol {
  padding-left: 1.5rem;
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.tiptap ul {
  list-style-type: disc;
}

.tiptap ol {
  list-style-type: decimal;
}

.tiptap li {
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}

.tiptap blockquote {
  border-left: 3px solid #6366f1;
  padding-left: 1rem;
  margin: 0.5rem 0;
  font-style: italic;
  color: #71717a;
}

.tiptap pre {
  background-color: #27272a;
  color: #fafafa;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
    Liberation Mono, Courier New, monospace;
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  margin: 0.5rem 0;
}

.tiptap code {
  background-color: #f4f4f5;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
}

.tiptap pre code {
  background-color: transparent;
  padding: 0;
}

.tiptap a {
  color: #6366f1;
  text-decoration: underline;
}
</style>