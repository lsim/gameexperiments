<template>
  <div id="app">
    <login class="app-login"></login>
    <game class="app-game"></game>
  </div>
</template>

<script>
import Game from './components/Game.vue'
import Login from './components/Login.vue'

export default {
  name: 'app',
  data() {
    return {
    }
  },
  methods: {

  },
  components: {
    Login,
    Game
  }
}
</script>
<style lang="scss">
  #app {
    height: 100vh; // Set to 100% height of viewport
    display: grid; // Use new grid fanciness
    grid-template-columns: 1fr 10fr 1fr;
    grid-template-rows: 1fr 20fr 1fr;
    grid-template-areas:
      "lg hd rg"
      "lg main rg"
      "lg ft rg";

    .app-login {
      grid-area: hd;
    }
    .app-game {
      grid-area: main;
    }
  }
</style>
