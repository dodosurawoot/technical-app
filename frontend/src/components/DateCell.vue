<template>
  <span :class="{ empty: !value }">{{ formatDate(value) }}</span>
</template>

<script setup lang="ts">
import { formatDate } from '../utils'

defineProps<{ value?: string | null }>()
</script>

