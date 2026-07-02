<template>
  <span class="badge" :class="`badge-${statusTone(status)}`">{{ statusLabel(status) }}</span>
</template>

<script setup lang="ts">
import { statusLabel, statusTone } from '../utils'

defineProps<{ status?: string }>()
</script>

