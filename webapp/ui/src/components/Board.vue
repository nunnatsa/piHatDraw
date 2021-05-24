<template>
  <v-container fluid>
    <v-row>
      <v-col align="center">
        <Controls :disabled="gameOver" />
      </v-col>
      <v-col align="center">
        <Canvas />
      </v-col>
    </v-row>
    <v-row>
      <v-col align="center">
        <v-alert
            :value="gameOver"
            color="#8888ee"
            transition="scale-transition"
            elevation="2"
            width="25em"
        >
          <p class="alert-text">The piHatDraw application was closed on the raspberry pi or we just lost connection with it.</p>
          <p class="bye-text">
            Thanks for using piHatDraw. See you next time!</p>
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Canvas from "./Canvas";
import Controls from "./Controls";
export default {
  name: "Board",
  components: {Controls, Canvas},
  props: [
      "gameOver"
  ],
}
</script>

<style scoped>
.alert-text {
  color: #222277;
  text-shadow: 1px 1px #666688;
  text-align: left;
}
.bye-text {
  color: #ccccff;
  text-shadow: 1px 1px #666688;
  text-align: left;
}
</style>
