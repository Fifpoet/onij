<!-- 主导航栏组件 -->
<template>
  <MobileHeader class="lg:hidden" />
  <DesktopHeader class="hidden lg:block" />
</template>

<script setup lang="ts">
import MobileHeader from './MobileHeader.vue'
import DesktopHeader from './DesktopHeader.vue'
</script>