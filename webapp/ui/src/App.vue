<template>
  <v-app>
    <v-app-bar
      app
      color="#aaaaee"
      dark
    >
      <div class="d-flex align-center">
        <v-img
          alt="piHatDraw Logo"
          class="shrink mr-2"
          contain
          src="./assets/logo.png"
          transition="scale-transition"
          width="40"
        />
      </div>

      <span class="text-h4 main-title">Draw with RaspberryPi SenseHAT</span>
      <v-spacer></v-spacer>

    </v-app-bar>

    <v-main v-if="$store.state.canvas">
      <Board :gameOver="gameOver" />
    </v-main>
  </v-app>
</template>

<script>
import Board from './components/Board';
import store from './store'
import HatService from './services'

export default {

  name: 'App',
  metaInfo: {
    // if no subcomponents specify a metaInfo.title, this title will be used
    title: 'Pi HAT Draw',
  },
  components: {
    Board
  },

  mounted: function() {
    console.log(`Starting connection to WebSocket Server on ${window.location.host}`)
    const wsURL = `ws://${location.host}/api/canvas/register`;

    this.connection = new WebSocket(wsURL)

    this.connection.onmessage = function(event) {
      HatService.init()
      const data = JSON.parse(event.data)
      console.log("data: " + JSON.stringify(data))

      if (data) {
        store.commit('replace', data)
      }
    }

    this.connection.onclose = (event) => {
      console.log(event)
      this.gameOver = true
    }

    this.connection.onopen = function(event) {
      console.log(event)
      console.log("Successfully connected to the echo websocket server...")
    }

  },

  data: () => ({
    gameOver: false
  }),

};
</script>

<style>
  .v-main__wrap {
    background-color: #ccccff;
  }

  .main-title {
    color: #ddddff;
    text-shadow: 2px 2px #8888cc
  }
</style>
