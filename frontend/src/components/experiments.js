
export default () => {
// Alright .. so we more or less have the vue.js browserify setup
// The assets are loaded from the root of the project
// The index.html is loaded from the root of the project. It would be nice if it could be served without exposing the entire project ..

  let messageTypes = {
    Register: 0,
  };
  let socket;

  let connectToServer = () => {
    socket = new WebSocket(`ws://${location.host}/api/socket`);

    socket.onerror = (event) => {
      console.error(event);
      setTimeout(connectToServer, 1000)
    };

    socket.onopen = (event) => {
      console.debug("onopen", event);
      // TODO: create vue UI for registering
      // TODO: get hold of our player's Id. Maybe the registration should be an http post?
      sendMessage(messageTypes.Register, "Foobar")
    };

    socket.onmessage = (event) => {
      handleMessage(JSON.parse(event.data))
    };
  };

  let sendMessage = (messageType, data) => {
    socket.send(JSON.stringify({
      Type: messageType,
      Data: data
    }));
  };

  let handleMessage = (gameState) => {
    // let planetPos = gameState.PlanetPos;
    // let planetRadius = gameState.PlanetRadius;
    let players = gameState.Players;
    for(let player of players) {
      let clientPlayer = clientPlayers[player.Id];
      if (!clientPlayer) {
        clientPlayer = clientPlayers[player.Id] = createPlayerSprite(player);
      }
      clientPlayer.centerX = player.Pos[0];
      clientPlayer.centerY = player.Pos[1];
    }
  };

  let clientPlayers = {};
  let viewPortWidth = 800;
  let viewPortHeight = 600;
  let worldWidth = 800;
  let worldHeight = 600;

  function preload () {
    game.load.image('diamond', 'assets/diamond.png');
  }

  let planet;
  let cursors;

  function createPlanetSprite() {
    let planet = game.add.sprite(0, 0, 'diamond');
    planet.scale.setTo(5, 5);
    planet.centerX = 0;
    planet.centerY = 0;
    return planet;
  }

  function createPlayerSprite(player) {
    return game.add.sprite(0, 0, 'diamond');
  }

  function create () {
    game.world.setBounds(-worldWidth/2, -worldHeight/2, worldWidth, worldHeight);
    cursors = game.input.keyboard.createCursorKeys();
    planet = createPlanetSprite();
    game.camera.x = -game.camera.view.width / 2;
    game.camera.y = -game.camera.view.height / 2;
    connectToServer();
  }

  function update () {
    //TODO: send the change in velocity to the server?
    //  Reset the players velocity (movement)
    // diamond.body.velocity.x = 0;
    // diamond.body.velocity.y = 0;
    if (cursors.left.isDown) {
      //  Move to the left
      // diamond.body.velocity.x -= 15;
      game.camera.x -= 5;
    } else if (cursors.right.isDown) {
      //  Move to the right
      // diamond.body.velocity.x += 15;
      game.camera.x += 5;
    }

    if (cursors.up.isDown) {
      //  Move up
      // diamond.body.velocity.y -= 15;
      game.camera.y -= 5;
    } else if (cursors.down.isDown) {
      //  Move down
      // diamond.body.velocity.y += 15;
      game.camera.y += 5;
    }
  }

  function render() {
    // game.debug.cameraInfo(this.game.camera, 32, 32);
    // game.debug.spriteInfo(planet, 370, 32);
  }
  console.debug("Creating Phaser instance..");
  let game = new Phaser.Game(viewPortWidth, viewPortHeight, Phaser.AUTO, '', { preload: preload, create: create, update: update, render: render });

}