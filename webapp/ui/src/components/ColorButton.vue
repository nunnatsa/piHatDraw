<template>
  <div>
    <v-tooltip bottom :color="color">
      <template v-slot:activator="{ on, attrs }">
        <v-btn
            class="mx-2"
            fab
            small
            :color="color"
            :disabled="disabled"
            @click.stop="onClick"
            v-bind="attrs"
            v-on="on"
        >
          <v-icon :color="textColor">
            mdi-palette
          </v-icon>
        </v-btn>
      </template>
      <span :style="{'color': textColor}">Select Color</span>
    </v-tooltip>
    <v-dialog v-model="showDialog" max-width="350">
      <v-card>
        <v-card-title class="headline">
          Choose a color
        </v-card-title>
        <v-card-text>
          <v-color-picker
              v-model="selectedColor"
              mode="hex"
              hide-inputs="true"
              elevation="4"
              flat
              :disabled="disabled"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
              color="green darken-1"
              text
              @click="cancel"
          >
            <v-icon>mdi-cancel</v-icon>
            Cancel
          </v-btn>
          <v-btn
              color="green darken-1"
              text
              @click="chooseColor"
          >
            <v-icon>mdi-eyedropper</v-icon>
            Choose Color
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import HatService from "../services";

export default {
  name: "ColorButton",
  data: function () {
    return {
      showDialog: false,
      selectedColor: "",
    }
  },
  props: [
    "disabled",
    "color",
  ],
  methods: {
    chooseColor: function () {
      if (this.selectedColor) {
        HatService.setColor(this.selectedColor)
      }
      this.showDialog = false
    },
    cancel: function () {
      this.showDialog = false
      this.selectedColor = this.color
    },
    onClick: function () {
      this.selectedColor = this.color
      this.showDialog = true
      if (document.activeElement instanceof HTMLElement) {
        document.activeElement.blur()
      }
    },
  },
  // mounted() {
  //   this.selectedColor = this.color
  // },
  computed: {
    textColor: function() {
      let clr = parseInt(Number("0x" + this.color.substring(1)))
      if ( clr && ((clr & 0xff0000) < 0x800000) || ((clr & 0xff00) < 0x8000) || ((clr & 0xff) < 0x80)) {
        return "white"
      }
      return "black"
    }
  },
}
</script>

<style scoped>
button {
  margin-top: 1em;
}
</style>
