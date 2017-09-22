import Vue from 'vue'
import App from './App.vue'

// Nastiness made necessary because phaser 'was not built to be modular'
window.PIXI = require( '../../node_modules/phaser-ce/build/custom/pixi' );
window.p2 = require( '../../node_modules/phaser-ce/build/custom/p2' );
window.Phaser = require( '../../node_modules/phaser-ce/build/custom/phaser-split' );

new Vue({
  el: '#app',
  render: (h) => h(App)
});
