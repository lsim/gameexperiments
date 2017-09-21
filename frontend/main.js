
$(() => {

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
    //console.debug("onmessage", gameState);
  };

  // let interval = setInterval(() => {
  //   socket.send(JSON.stringify({
  //     inMsg: "" + new Date().getSeconds()
  //   }))
  // }, 3000);

  function preload () {
    game.load.image('diamond', 'assets/diamond.png');
  }

  let diamond;
  let cursors;

  function create () {
    diamond = game.add.sprite(game.world.centerX, game.world.centerY, 'diamond');
    game.physics.arcade.enable(diamond);
    diamond.body.collideWorldBounds = true;
    cursors = game.input.keyboard.createCursorKeys();
    connectToServer()
  }

  function update () {
    //  Reset the players velocity (movement)
    diamond.body.velocity.x = 0;
    if (cursors.left.isDown) {
      //  Move to the left
      diamond.body.velocity.x = -150;
    } else if (cursors.right.isDown) {
      //  Move to the right
      diamond.body.velocity.x = 150;
    }

    if (cursors.up.isDown) {
      //  Move up
      diamond.body.velocity.y = -150;
    } else if (cursors.down.isDown) {
      //  Move down
      diamond.body.velocity.y = 150;
    }
  }

  let game = new Phaser.Game(800, 600, Phaser.AUTO, '', { preload: preload, create: create, update: update });

});
