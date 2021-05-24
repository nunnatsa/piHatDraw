<template>
  <v-card elevation="2" class="mx-auto my-12" color="#aaaaee"
          :v-if="!$store.state"
          :min-width="getWidth"
          :max-width="getWidth"
          :min-height="getHeight"
          :max-height="getHeight"
  >
    <v-row jestify="center" align="center">
      <v-col>
        <table id="canvas">
          <tr v-for="(line, y) in $store.state.canvas" v-bind:key="y">
            <Cell v-for="(cell, x) in line"
                  v-bind:key="x"
                  :bgColor="cell"
                  :tool="getToolChar(x, y)"
                  :borders="borders(x, y)"
            >
            </Cell>
          </tr>
        </table>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import Cell from './Cell'

export default {
  name: "Canvas",
  components: {
    Cell,
  },
  data: () => ({}),
  computed: {
    getHeight: function () {
      if (!this.$store.state.canvas) return 0;
      return this.$store.state.canvas.length * 21 + 16
    },
    getWidth: function () {
      if (!this.$store.state.canvas) return 0;
      return this.$store.state.canvas[0].length * 21 + 16
    },
  },
  methods: {
    getToolChar: function (x, y) {
      return x === this.$store.state.cursor.x && y === this.$store.state.cursor.y ? this.$store.state.toolChar : ''
    },
    borders: function (x, y) {
      const win = this.$store.state.window
      const b = {
        top: (y === win.top) && (x >= win.left && x <= win.right),
        left: (x === win.left) && (y >= win.top && y <= win.bottom),
        bottom: (y === win.bottom) && (x >= win.left && x <= win.right),
        right: (x === win.right) && (y >= win.top && y <= win.bottom),
      }
      return b
    },
  }
}
</script>

<style scoped>
table, tr {
  border: solid 1px white;
}

table {
  border-collapse: collapse;
  margin: auto;
}
</style>