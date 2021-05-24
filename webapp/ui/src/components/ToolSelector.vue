<template>
  <v-card elevation="1" width="360" color="#8888ee">
    <v-card-title class="text-body-1 tool-title">Select Tool</v-card-title>
    <v-card-text>
      <v-btn-toggle background-color="green" tile
                    :value="toolName"
                    color="#444488"
                    group
                    mandatory
                    @change.self="selectTool"
      >
        <v-btn small elevation="2" value="pen" :disabled="disabled"><v-icon>mdi-pen</v-icon>Pen</v-btn>
        <v-btn small elevation="2" value="eraser" :disabled="disabled"><v-icon>mdi-eraser-variant</v-icon>Eraser</v-btn>
        <v-btn small elevation="2" value="bucket" :disabled="disabled"><v-icon>mdi-format-color-fill</v-icon>Bucket</v-btn>
      </v-btn-toggle>
    </v-card-text>
  </v-card>
</template>

<script>
import HatService from '../services'

export default {
  name: "ToolSelector",
  methods: {
    selectTool: function (value) {
      if (value) {
        HatService.setTool(value)
      }
    }
  },
  props: [
    'toolName',
    'disabled',
  ],
}
</script>

<style scoped>
  .tool-title {
    color: #ccccff;
    text-shadow: 1px 1px #666688;
  }
</style>