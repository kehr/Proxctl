<script setup lang="ts">
import { releaseAssets } from '../../generated/version'
</script>

<template>
  <ul>
    <li v-for="asset in releaseAssets" :key="asset">
      <code>{{ asset }}</code>
    </li>
  </ul>
</template>
