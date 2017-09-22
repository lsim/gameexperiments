<template>
  <div id="app">
    <experiments class="experiments"></experiments>
  </div>
</template>

<script>
import Experiments from './components/experiments.vue'

export default {
  name: 'app',
  components: {
    Experiments
  }
}
</script>
<style lang="scss">
  .experiments {
    background-color: orange;
    color: deepskyblue;
  }
</style>
