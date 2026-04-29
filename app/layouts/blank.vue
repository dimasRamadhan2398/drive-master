<script setup lang="ts">
// Blank layout for authentication and specific pages
</script>

<template>
  <div class="min-h-screen bg-background flex flex-col relative">
    <!-- Top actions -->
    <div class="absolute top-4 right-4 z-50 flex items-center gap-2">
      <UButton
        to="/"
        icon="i-lucide-home"
        color="warning"
        variant="ghost"
        aria-label="Back to Home"
      />
      <UColorModeButton />
    </div>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col">
      <slot />
    </main>
  </div>
</template>
