<template>
  <div ref="canvasContainer">

  </div>
</template>

<script>
  import SocketService from '../services/socket'
  import _ from 'lodash'

  export default {
    name: 'game',
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
        clientPlayers: {},
        explosionPool: null,
      }
    },
    methods: {
      listenForPlayerUpdates() {
        if (!this.playerUpdateSubscription) {
          this.playerUpdateSubscription = SocketService.getTypedMessageSubject(SocketService.messageTypes.UpdatePlayers).subscribe((players) => {
            if(players) {
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
                    1000.0/10.0,// TODO: we need to know the server update frequency here - it should be part of the initial handshake
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
            }
          });
          SocketService.getTypedMessageSubject(SocketService.messageTypes.PlayerDied).subscribe((playerId) => {
            this.killPlayer(playerId);
          });
        }
      },
      killPlayer(playerId) {
        let clientPlayer = this.clientPlayers[playerId];
        if (clientPlayer) {
          let explosion = this.explosions.getFirstExists(false);
          // TODO: sometimes we get null back - perhaps when all sprites have run their animation and have been destroyed
          explosion.reset(clientPlayer.x, clientPlayer.y);
          explosion.play('explosion', 30, false, true);
          clientPlayer.kill();
          clientPlayer.text.kill();
          delete this.clientPlayers[playerId];
        }
      },
      phaserPreload() {
        this.phaserGame.load.image('planet', 'assets/planet.png');
        this.phaserGame.load.image('player', 'assets/fighter-01.png');
        this.phaserGame.load.spritesheet('explosion', 'assets/explosion.png', 157, 229, 19);
        // Keep game running when it doesn't have focus
        this.phaserGame.stage.disableVisibilityChange = true;
        // Make game canvas resize with window size changes
        this.phaserGame.scale.scaleMode = Phaser.ScaleManager.RESIZE;
        this.phaserGame.scale.setResizeCallback(() => {
          // Re-center the camera on resize
          this.phaserGame.camera.x = -this.phaserGame.camera.view.width / 2;
          this.phaserGame.camera.y = -this.phaserGame.camera.view.height / 2;
        });
      },
      phaserCreate() {
        this.phaserGame.world.setBounds(-this.worldWidth/2, -this.worldHeight/2, this.worldWidth, this.worldHeight);
        this.cursors = this.phaserGame.input.keyboard.createCursorKeys();
        this.planet = this.createPlanetSprite();
        this.phaserGame.camera.x = -this.phaserGame.camera.view.width / 2;
        this.phaserGame.camera.y = -this.phaserGame.camera.view.height / 2;
        // Set up explosions
        this.explosions = this.phaserGame.add.group();
        this.explosions.createMultiple(30, 'explosion');
        this.explosions.forEach((explosionSprite) => {
          explosionSprite.anchor.x = explosionSprite.anchor.y = 0.5;
          explosionSprite.animations.add('explosion');
        });
        this.listenForPlayerUpdates();
      },
      rotateLeft: _.throttle(() => SocketService.rotate(-1), 50),
      rotateRight: _.throttle(() => SocketService.rotate(1), 50),
      thrust: _.throttle(() => SocketService.addThrust(), 50),
      phaserUpdate() {
        if (this.playerId) {
          if (this.cursors.left.isDown) {
            this.rotateLeft();
          } else if (this.cursors.right.isDown) {
            this.rotateRight();
          }
          if (this.cursors.up.isDown) {
            this.thrust();
          } else if (this.cursors.down.isDown) {
          }
        }
      },
      phaserRender() {
        // this.phaserGame.debug.cameraInfo(this.[phaserGame.camera, 32, 32);
//        this.phaserGame.debug.spriteInfo(this.planet, 370, 32);
      },
      createPlanetSprite() {
        let planet = this.phaserGame.add.sprite(0, 0, 'planet');
        planet.height = planet.width = 140; // TODO: should come from the server
        planet.anchor.setTo(0.5, 0.5);
        return planet;
      },
      createPlayerSprite(name) {
        let player = this.phaserGame.add.sprite(0, 0, 'player');
        player.height = 20; // TODO: should come from the server
        player.width = 20;

        player.anchor.setTo(0.5, 0.5);
        player.text = this.phaserGame.add.text(100, 100, name, {
          font: "9px Arial",
          fill: "#fff",
          align: "center"
        });
        return player;
      }
    },
    mounted() {
      this.playerRegisteredSubscription = SocketService.getTypedMessageSubject(SocketService.messageTypes.Registered).subscribe((playerId) => {
        this.playerId = playerId;
      });
      this.playerUnregisteredSubscription = SocketService.getPlayerUnregisteredSubject().subscribe(() => {
        if (this.playerId !== null && this.clientPlayers[this.playerId] !== null) {
          this.killPlayer(this.playerId);
        }
        this.playerId = null;
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
