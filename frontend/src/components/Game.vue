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
        worldWidth: 8000,
        worldHeight: 6000,
        planet: null,
        cursors: null,
        clientPlayers: {}
      }
    },
    methods: {
      listenForPlayerUpdates() {
        if (!this.playerUpdateSubscription) {
          this.playerUpdateSubscription = SocketService.getPlayerUpdateSubject().subscribe((players) => {
            for(let player of players) {
              let clientPlayer = this.clientPlayers[player.Id];
              if (!clientPlayer) {
                clientPlayer = this.clientPlayers[player.Id] = this.createPlayerSprite(player.Name);
                clientPlayer.centerX = player.Pos[0];
                clientPlayer.centerY = player.Pos[1];
                clientPlayer.rotation = player.Angle;
                clientPlayer.text.centerX = player.Pos[0];
                clientPlayer.text.centerY = player.Pos[1] - 20;
              } else {
                // Use tweens to get smooth movement out of sparse server updates
                this.phaserGame.add.tween(clientPlayer).to(
                  {
                    centerX: player.Pos[0],
                    centerY: player.Pos[1],
                    rotation: player.Angle
                  },
                  1000.0/10.0,// TODO: we need to know the server fps here - they should be part of the initial handshake
                  Phaser.Easing.Linear.None,
                  true
                );
                this.phaserGame.add.tween(clientPlayer.text).to(
                  {
                    centerX: player.Pos[0],
                    centerY: player.Pos[1] - 20,
                  },
                  1000.0/10.0,
                  Phaser.Easing.Linear.None,
                  true
                );
              }
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
        if (this.cursors.left.isDown) {
//          this.phaserGame.camera.x -= 5;
          SocketService.rotate(-1);
        } else if (this.cursors.right.isDown) {
//          this.phaserGame.camera.x += 5;
          SocketService.rotate(1);
        }

        if (this.cursors.up.isDown) {
//          this.phaserGame.camera.y -= 5;
          SocketService.addThrust();
        } else if (this.cursors.down.isDown) {
//          this.phaserGame.camera.y += 5;
        }
      },
      phaserRender() {
        // this.phaserGame.debug.cameraInfo(this.[phaserGame.camera, 32, 32);
//        this.phaserGame.debug.spriteInfo(this.planet, 370, 32);
      },
      createPlanetSprite() {
        let planet = this.phaserGame.add.graphics(0, 0);

//        graphics.lineStyle(2, 0xffd900, 1);
        planet.beginFill(0xFF0000, 1);
        planet.drawCircle(0, 0, 70);

        return planet;
      },
      createPlayerSprite(name) {
        let player = this.phaserGame.add.graphics(0, 0);

        player.beginFill(0x00FF00, 1);
        player.moveTo(-5, 5);
        player.lineTo(5, 0);
        player.lineTo(-5, -5);
        player.lineTo(-3, 0);
        player.lineTo(-5, 5);
        player.endFill();
        player.anchor.setTo(0.5, 0.5);
        player.width = player.height = 25;
        player.text = this.phaserGame.add.text(100, 100, name, {
          font: "9px Arial",
          fill: "#fff",
          align: "center"
        });
        return player;
      }
    },
    mounted() {
      this.playerRegisteredSubscription = SocketService.getPlayerRegisteredSubject().subscribe((playerId) => {
        this.playerId = playerId;
        this.listenForPlayerUpdates();
      });
      this.playerUnregisteredSubscription = SocketService.getPlayerUnregisteredSubject().subscribe(() => {
        if (this.playerId !== null && this.clientPlayers[this.playerId] !== null) {
          this.clientPlayers[this.playerId].text.destroy()
          this.clientPlayers[this.playerId].destroy();
          delete this.clientPlayers[this.playerId];
        }
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
      this.phaserGame.destroy();
      let container = this.$refs.canvasContainer;
      while (container.firstChild) {
        container.removeChild(container.firstChild);
      }
    },
    components: {

    }
  }
</script>
<style lang="sass">

</style>
