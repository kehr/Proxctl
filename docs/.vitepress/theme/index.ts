import DefaultTheme from 'vitepress/theme'
import ArtifactNames from './components/ArtifactNames.vue'
import InstallExamples from './components/InstallExamples.vue'
import LatestVersion from './components/LatestVersion.vue'
import './style.css'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('ArtifactNames', ArtifactNames)
    app.component('InstallExamples', InstallExamples)
    app.component('LatestVersion', LatestVersion)
  }
}
