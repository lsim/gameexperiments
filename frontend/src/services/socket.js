'use strict';

import { QueueingSubject } from 'queueing-subject'
import websocketConnect from 'rxjs-websockets'
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/take';
import 'rxjs/add/operator/retryWhen';
import { Subject } from 'rxjs'

let messageTypes = {
  Register: 0,
  UpdatePlayers: 1,
  Registered: 2,
  Unregister: 3,
  RotateClockWise: 4,
  RotateCounterClockWise: 5,
  AddThrust: 6,
  PlayerDied: 7,
  Shoot: 8,
  BulletDied: 9,
  WelcomeClient: 10,
};

const input = new QueueingSubject();
const { messages, connectionStatus } = websocketConnect(`ws://${location.host}/api/socket`, input);

const jsonReceivedSubject = new Subject();
messages.map((stringMessage) => JSON.parse(stringMessage)).subscribe(jsonReceivedSubject);

const playerUnregisteredSubject = new Subject();

const sendObject = (transmissionObj) => input.next(JSON.stringify(transmissionObj));

let playerId = null;

// If player dies, we notify 'unregister' listeners
jsonReceivedSubject
  .filter((m) => m.Type === messageTypes.PlayerDied && playerId !== null && m.Data === playerId)
  .subscribe(playerUnregisteredSubject);

playerUnregisteredSubject.subscribe(() => {
  playerId = null;
});

const getTypedSubject = (type) => jsonReceivedSubject.filter((m) => m.Type === type).map(m => m.Data);


const gameConstantsPromise = new Promise((resolve, reject) => {
  let welcomeSubscription = getTypedSubject(messageTypes.WelcomeClient).subscribe((gameConstants) => {
    resolve(gameConstants);
    welcomeSubscription.unsubscribe();
  });
});

export default {
  messageTypes: messageTypes,
  getMessageSubject: () => jsonReceivedSubject,
  getTypedMessageSubject: (type) => jsonReceivedSubject.filter((m) => m.Type === type).map(m => m.Data),
  getPlayerUnregisteredSubject: () => playerUnregisteredSubject,
  getConnectionStatusSubject: () => connectionStatus,
  send: sendObject,
  getGameConstants() {
    return gameConstantsPromise;
  },
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
    sendObject({Type: messageTypes.Unregister});
    playerUnregisteredSubject.next();
  },
  addThrust() {
    sendObject({Type: messageTypes.AddThrust})
  },
  rotate(direction) {
    if (direction > 0) {
      sendObject({Type: messageTypes.RotateClockWise});
    } else {
      sendObject({Type: messageTypes.RotateCounterClockWise});
    }
  },
  shoot() {
    sendObject({Type: messageTypes.Shoot});
  }

}
