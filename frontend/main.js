
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
});
