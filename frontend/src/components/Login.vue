<template>
  <div>
    <input v-model="playerName" :disabled="playerRegistered"/>
    <button v-on:click="register" :disabled="registerDisabled">Register</button>
    <button v-on:click="unregister" :disabled="!playerRegistered">Disconnect</button>
  </div>
</template>

<script>
  import SocketService from '../services/socket'

  export default {
    name: 'experiments',
    data() {
      return {
        playerName: "",
        playerRegistered: false
      }
    },
    computed: {
      registerDisabled() {
        return !this.playerName || this.playerRegistered;
      }
    },
    methods: {
      register() {
        if (this.playerName) {
          SocketService.registerPlayer(this.playerName).then(
            (playerId) => this.playerRegistered = true,
            (error) => console.info("Error registering: ", error)
          );
          this.playerName = "";
        }
      },
      unregister() {
        SocketService.unregisterPlayer();
        this.playerRegistered = false;
      }
    },
    components: {

    }
  }
</script>
<style lang="sass">

</style>
