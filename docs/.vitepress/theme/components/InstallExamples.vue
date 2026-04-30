<script setup lang="ts">
import { latestInstallerUrl, proxctlVersion, releaseDownloadUrl } from '../../generated/version'

const defaultInstall = `curl -fsSL ${latestInstallerUrl} | sh`
const specificInstall = `curl -fsSL ${latestInstallerUrl} | PROXCTL_VERSION=${proxctlVersion} sh`
const manualInstall = `curl -LO ${releaseDownloadUrl}/proxctl-${proxctlVersion}-linux-amd64.tar.gz
tar -xzf proxctl-${proxctlVersion}-linux-amd64.tar.gz
sudo install -m 0755 proxctl-${proxctlVersion}-linux-amd64/proxctl /usr/local/bin/proxctl
proxctl version`
</script>

<template>
  <div class="proxctl-examples">
    <h3>One-line install</h3>
    <pre><code>{{ defaultInstall }}</code></pre>

    <h3>Install the documented version</h3>
    <pre><code>{{ specificInstall }}</code></pre>

    <h3>Manual Linux amd64 install</h3>
    <pre><code>{{ manualInstall }}</code></pre>
  </div>
</template>
