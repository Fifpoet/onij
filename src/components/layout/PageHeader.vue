<!-- 主导航栏组件 -->
<template>
  <MobileHeader v-if="isMobile" />
  <DesktopHeader v-else />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import MobileHeader from './MobileHeader.vue'
import DesktopHeader from './DesktopHeader.vue'

const isMobile = useMediaQuery('(max-width: 1024px)')
</script>