<template>
  <div ref="canvasContainer">

  </div>
</template>

<script>
  import SocketService from '../services/socket'
  import _ from 'lodash'

  let vect = (x, y) => new Phaser.Point(x, y);

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
        clientBullets: {},
        explosionPool: null,
        bulletGroup: null,
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
                  if (player.Id === this.playerId) {
                    this.player = clientPlayer;
                  }
                }
                clientPlayer.centerX = player.Pos[0];
                clientPlayer.centerY = player.Pos[1];
                clientPlayer.rotation = player.Angle;
                clientPlayer.text.centerX = player.Pos[0];
                clientPlayer.text.centerY = player.Pos[1] - 20;
                clientPlayer.body.velocity.setTo(player.Velocity[0], player.Velocity[1])
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
      killBullet(bulletId) {
        let bullet = this.clientBullets[bulletId];
        if (bullet) {
          bullet.kill();
          delete this.clientBullets[bulletId]
        }
      },
      phaserPreload() {
        this.phaserGame.load.image('planet', 'assets/planet.png');
        this.phaserGame.load.image('player', 'assets/player-ship.png');
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
        this.phaserGame.physics.startSystem(Phaser.Physics.ARCADE);
        this.phaserGame.world.setBounds(-this.worldWidth/2, -this.worldHeight/2, this.worldWidth, this.worldHeight);
        this.cursors = this.phaserGame.input.keyboard.createCursorKeys();
        this.planet = this.createPlanetSprite();
        this.phaserGame.camera.x = -this.phaserGame.camera.view.width / 2;
        this.phaserGame.camera.y = -this.phaserGame.camera.view.height / 2;
        // Set up explosions
        this.explosions = this.phaserGame.add.group();
        this.explosions.enableBody = false;
        this.explosions.createMultiple(30, 'explosion');
        this.explosions.forEach((explosionSprite) => {
          explosionSprite.anchor.x = explosionSprite.anchor.y = 0.5;
          explosionSprite.animations.add('explosion');
        });
        this.bulletGroup = this.phaserGame.add.group();
        this.bulletGroup.enableBody = true;
        this.bulletGroup.physicsBodyType = Phaser.Physics.ARCADE;
        this.listenForPlayerUpdates();
      },
      rotateLeft: _.throttle((player) => {
        player.rotation -= 0.1;
        SocketService.rotate(-1);
      }, 50, { trailing: false }),
      rotateRight: _.throttle((player) => {
        player.rotation += 0.1;
        SocketService.rotate(1)
      }, 50, { trailing: false }),
      thrust: _.throttle((player) => {
        let thrustVector = new Phaser.Point(Math.cos(player.rotation), Math.sin(player.rotation));
        thrustVector = thrustVector.multiply(5, 5);
        player.body.velocity.add(thrustVector.x, thrustVector.y);
        SocketService.addThrust();
      }, 50, { trailing: false }),
      shoot: _.throttle(() => SocketService.shoot(), 100, { trailing: false }),
      phaserUpdate() {
        _.values(this.clientPlayers).forEach((player) => {
          this.applyPlanetaryGravitation(player);
        });
        if (this.player !== null) {
          if (this.cursors.left.isDown) {
            this.rotateLeft(this.player);
          } else if (this.cursors.right.isDown) {
            this.rotateRight(this.player);
          }
          if (this.cursors.up.isDown) {
            this.thrust(this.player);
          } else if (this.cursors.down.isDown) {
            this.shoot();
          }
        }
      },
      applyPlanetaryGravitation(player) {
        let gravityStrength = 1e6; // TODO: should come from server
        let deltaTime = this.phaserGame.time.physicsElapsedMS / 1000;
        let playerPos = player.position;
        let sqDist = playerPos.getMagnitudeSq();
        let multiplier = -gravityStrength/(sqDist*Math.sqrt(sqDist));
        let gravityVect = playerPos.clone().multiply(multiplier, multiplier);
        this.updatePlayerVelocity(player, gravityVect, deltaTime);
      },
      updatePlayerVelocity(player, gravityVect, deltaTime) {
        let timeAdjustedGravityContribution = gravityVect.multiply(deltaTime, deltaTime);
        player.body.velocity = player.body.velocity.add(timeAdjustedGravityContribution.x, timeAdjustedGravityContribution.y);
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
        this.phaserGame.physics.enable(player);
        return player;
      },
      createBulletSprite(bulletInfo) {
        let bulletGraphics = this.phaserGame.add.graphics(bulletInfo.Pos[0], bulletInfo.Pos[1])
        bulletGraphics.beginFill(0x00FFFF, 1);
        bulletGraphics.drawRect(0, 0, 3, 1);
        bulletGraphics.rotation = bulletInfo.Angle;
        bulletGraphics.anchor.setTo(0.5, 0.5);
        this.clientBullets[bulletInfo.Id] = bulletGraphics;
        this.bulletGroup.add(bulletGraphics);
        bulletGraphics.body.velocity.setTo(bulletInfo.Velocity[0], bulletInfo.Velocity[1]);
        return bulletGraphics;
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
        this.player = null;
      });
      this.playerShootSubscription = SocketService.getTypedMessageSubject(SocketService.messageTypes.Shoot).subscribe((bulletInfo) => {
        this.createBulletSprite(bulletInfo);
      });
      this.bulletDiedSubscription = SocketService.getTypedMessageSubject(SocketService.messageTypes.BulletDied).subscribe((bulletId) => {
        this.killBullet(bulletId);
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
      this.playerShootSubscription.unsubscribe();
      this.bulletDiedSubscription.unsubscribe();
      this.phaserGame.destroy();
      let container = this.$refs.canvasContainer;
      while (container.firstChild) {
        container.removeChild(container.firstChild);
      }
    },
    components: {}
  }
</script>
<style lang="sass">

</style>
