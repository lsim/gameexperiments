'use strict';

import { QueueingSubject } from 'queueing-subject'
import websocketConnect from 'rxjs-websockets'
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/take';
import { Subject } from 'rxjs'

let messageTypes = {
  Register: 0,
  UpdatePlayers: 1,
  Registered: 2,
  Unregister: 3,
};

const input = new QueueingSubject();
const { messages, connectionStatus } = websocketConnect(`ws://${location.host}/api/socket`, input);
// const jsonMessages = messages.map((stringMessage) => JSON.parse(stringMessage));

const jsonReceivedSubject = new Subject();
messages.map((stringMessage) => JSON.parse(stringMessage)).subscribe(jsonReceivedSubject);

const playerUnregisteredSubject = new Subject();

const sendObject = (transmissionObj) => input.next(JSON.stringify(transmissionObj));

let playerId = null;

export default {
  getMessageSubject: () => jsonReceivedSubject,
  getPlayerUpdateSubject: () => jsonReceivedSubject.filter((m) => m.Type === messageTypes.UpdatePlayers).map(m => m.Data),
  getPlayerRegisteredSubject: () => jsonReceivedSubject.filter((m) => m.Type === messageTypes.Registered).map(m => m.Data),
  getPlayerUnregisteredSubject: () => playerUnregisteredSubject,
  getConnectionStatusSubject: () => connectionStatus,
  send: sendObject,
  registerPlayer(playerName) {
    return new Promise((resolve, reject) => {
      if (playerId !== null) {
        reject("Player already registered with id " + playerId);
        return;
      }
      sendObject({Type: messageTypes.Register, Data: playerName });
      jsonReceivedSubject.take(1).subscribe((message) => {
        if (message.Type !== messageTypes.Registered) {
          reject("Received unexpected reply upon registering: " + message.Type);
        } else {
          playerId = message.Data;
          resolve(playerId);
        }
      });
    });
  },
  unregisterPlayer() {
    if (playerId === null) {
      console.debug("unregisterPlayer - no player registration exists! Doing nothing");
    }
    playerId = null;
    sendObject({Type: messageTypes.Unregister});
    playerUnregisteredSubject.next();
  }
}
