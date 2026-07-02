<template>
  <div v-if="open" class="modal-backdrop">
    <div class="modal">
      <h3>{{ title }}</h3>
      <p>{{ message }}</p>
      <div class="modal-actions">
        <button class="ghost" @click="$emit('cancel')">ยกเลิก</button>
        <button class="danger" @click="$emit('confirm')">ยืนยัน</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ open: boolean; title: string; message: string }>()
defineEmits<{ confirm: []; cancel: [] }>()
</script>

