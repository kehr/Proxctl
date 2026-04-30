<script setup lang="ts">
import { proxctlVersion, releaseUrl } from '../../generated/version'
</script>

<template>
  <div class="proxctl-version">
    <span class="proxctl-version__label">Latest version</span>
    <a :href="releaseUrl" class="proxctl-version__value">{{ proxctlVersion }}</a>
  </div>
</template>
