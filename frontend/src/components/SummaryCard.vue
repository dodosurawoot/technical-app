<template>
  <section class="summary-card">
    <div class="summary-icon" :class="tone">
      <component :is="icon" :size="22" />
    </div>
    <div>
      <small>{{ label }}</small>
      <strong>{{ value.toLocaleString('th-TH') }}</strong>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{ label: string; value: number; icon: unknown; tone?: string }>()
</script>

