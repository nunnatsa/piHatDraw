<template>
  <v-card width="380" class="mx-auto my-12" align="center" color="#aaaaee">
    <v-card-text>
      <v-spacer/>
      <v-row>
        <v-col>
          <ToolSelector v-if="!$store.state.initializing" :tool-name="$store.state.tool" :disabled="disabled"/>
        </v-col>
      </v-row>
      <v-spacer/>
      <v-row>
        <v-col>
          <v-card width="360" color="#8888ee">
            <v-card-text>
              <v-row>
                <v-col align="left">
                  <v-btn small @click="pickCursorColor" color="#8888ee" :disabled="disabled">
                    <v-icon>mdi-select-color</v-icon>
                    Pick Cursor Color
                  </v-btn>
                </v-col>
                <ColorButton v-if="!$store.state.initializing"  :color="$store.state.color" :disabled="disabled"/>
              </v-row>
              <v-spacer/>
              <v-row>
                <v-col align="left">
                  <DownloadButton :disabled="disabled"/>
                </v-col>
              </v-row>
              <v-spacer/>
              <v-row>
                <v-col align="left">
                  <v-btn small @click="undo" color="#8888ee" :disabled="disabled">
                    <v-icon>mdi-undo-variant</v-icon>
                    Undo
                  </v-btn>
                </v-col>
                <v-col align="right">
                  <ResetButton :disabled="disabled"/>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import ToolSelector from "./ToolSelector";
import {store} from '../store'
import HatService from '../services'
import ResetButton from "./ResetButton";
import DownloadButton from "./DownloadButton";
import ColorButton from "@/components/ColorButton";

export default {
  name: "Controls",
  components: {ColorButton, ResetButton, ToolSelector, DownloadButton},
  props: [
      "disabled",
  ],
  methods: {
    pickCursorColor: () => {
      const color = store.state.canvas[store.state.cursor.y][store.state.cursor.x]
      HatService.setColor(color)
    },
    undo: () => {
      HatService.undo()
    },
  },
}
</script>

<style scoped>

</style>
