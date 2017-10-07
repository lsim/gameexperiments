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
  .app-login {

  }
  .app-game {

  }
</style>
