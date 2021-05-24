<template>
  <v-card width="360" color="#8888ee">
    <v-card-title class="text-body-1 color-title">Chose Color</v-card-title>
    <v-card-text>
      <v-color-picker
          :value="color"
          mode="hexa"
          hide-inputs
          @input.self="onColorChange"
          elevation="4"
          flat
          :disabled="disabled"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import HatService from "../services";
import _ from 'lodash'

export default {
  name: "ColorPickerWrapper",
  props: [
    'color',
    'disabled',
  ],
  methods: {
    onColorChange: _.throttle( (input) => {
      if (input) {
        HatService.setColor(input)
      }
    }, 200)
  },
}
</script>

<style scoped>
  .color-title {
    color: #ccccff;
    text-shadow: 1px 1px #666688;
  }
</style>
