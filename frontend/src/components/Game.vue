<template>
  <div ref="canvasContainer">

  </div>
</template>

<script>
//  import FooBarBaz from './experiments'
  import SocketService from '../services/socket'

  export default {
    name: 'experiments',
    data() {
      return {
        playerId: null,
        playerUpdateSubscription: null,
        phaserGame: null,
        viewPortWidth: 800,
        viewPortHeight: 600,
        worldWidth: 800,
        worldHeight: 600,
        planet: null,
        cursors: null,
        clientPlayers: {}
      }
    },
    methods: {
      listenForPlayerUpdates() {
        if (!this.playerUpdateSubscription) {
          this.playerUpdateSubscription = SocketService.getPlayerUpdateSubject().subscribe((players) => {
            console.debug("Game.vue got player update", players);
            for(let player of players) {
              let clientPlayer = this.clientPlayers[player.Id];
              if (!clientPlayer) {
                clientPlayer = this.clientPlayers[player.Id] = this.createPlayerSprite(player);
              }
              clientPlayer.centerX = player.Pos[0];
              clientPlayer.centerY = player.Pos[1];
            }
          });
        }
      },
      stopListeningForPlayerUpdates() {
        if (this.playerUpdateSubscription) {
          this.playerUpdateSubscription.unsubscribe();
          this.playerUpdateSubscription = null;
        }
      },
      phaserPreload() {
        this.phaserGame.load.image('diamond', 'assets/diamond.png');
      },
      phaserCreate() {
        this.phaserGame.world.setBounds(-this.worldWidth/2, -this.worldHeight/2, this.worldWidth, this.worldHeight);
        this.cursors = this.phaserGame.input.keyboard.createCursorKeys();
        this.planet = this.createPlanetSprite();
        this.phaserGame.camera.x = -this.phaserGame.camera.view.width / 2;
        this.phaserGame.camera.y = -this.phaserGame.camera.view.height / 2;
      },
      phaserUpdate() {

      },
      phaserRender() {

      },
      createPlanetSprite() {
        let planet = this.phaserGame.add.sprite(0, 0, 'diamond');
        planet.scale.setTo(5, 5);
        planet.centerX = 0;
        planet.centerY = 0;
        return planet;
      },
      createPlayerSprite(player) {
        return this.phaserGame.add.sprite(0, 0, 'diamond');
      }
    },
    mounted() {
      this.playerRegisteredSubscription = SocketService.getPlayerRegisteredSubject().subscribe((playerId) => {
        this.playerId = playerId;
        this.listenForPlayerUpdates();
      });
      this.playerUnregisteredSubscription = SocketService.getPlayerUnregisteredSubject().subscribe(() => {
        this.playerId = null;
        this.stopListeningForPlayerUpdates();
      });
      this.phaserGame = new Phaser.Game(
        this.viewPortWidth,
        this.viewPortHeight,
        Phaser.AUTO,
        this.$refs.canvasContainer, {
          preload: this.phaserPreload,
          create: this.phaserCreate,
          update: this.phaserUpdate,
          render: this.phaserRender
        });
    },
    beforeDestroy() {
      this.playerRegisteredSubscription.unsubscribe();
      this.playerUnregisteredSubscription.unsubscribe();
    },
    components: {

    }
  }
</script>
<style lang="sass">

</style>
