<script setup lang="ts">
import EddyterEditor from "./EddyterEditor.vue";

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

function handleUpdate(value: string) {
  emit("update:modelValue", value);
}
</script>

<template>
  <EddyterEditor
    :model-value="modelValue"
    :placeholder="placeholder"
    @update:model-value="handleUpdate"
  />
</template>
