
$(() => {
  $.get('api/hello', function(data) {
    console.debug("received hello from server", JSON.parse(data));
  });

  let socket = new WebSocket(`ws://${location.host}/api/socket`);

  socket.onerror = (event) => {
    console.error(event);
  };

  socket.onopen = (event) => {
    console.debug("onopen", event);
  };

  socket.onmessage = (event) => {
    console.debug("onmessage", event);
  };

  // socket.send("foobar")

  function preload () {
    game.load.image('logo', 'phaser.png');
  }

  function create () {
    let logo = game.add.sprite(game.world.centerX, game.world.centerY, 'logo');
    logo.anchor.setTo(0.5, 0.5);
  }

  let game = new Phaser.Game(800, 600, Phaser.AUTO, '', { preload: preload, create: create });

});
