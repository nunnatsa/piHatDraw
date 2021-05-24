<template>
  <div>
    <v-btn
        small
        @click.stop="confirmDialog = true"
        color="#8888ee"
        :disabled="disabled"
    >
      <v-icon>mdi-restart</v-icon>
      Reset
    </v-btn>

    <v-dialog
        v-model="confirmDialog" max-width="350"
    >
      <v-card>
        <v-card-title class="headline">
          Reset the Whole Picture?
        </v-card-title>

        <v-card-text>
          This will delete the picture and will set the color and the tool back to default; Do you want to continue?
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
              color="green darken-1"
              text
              @click="confirmDialog = false"
          >
            <v-icon>mdi-cancel</v-icon>
            No
          </v-btn>

          <v-btn
              color="green darken-1"
              text
              @click="confirmDialog = false; resetCanvas()"
          >
            <v-icon>mdi-restart</v-icon>
            Yes
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import HatService from '../services'
export default {
  name: "ResetButton",
  data() {
    return {
      confirmDialog: false
    }
  },
  props: [
      "disabled",
  ],
  methods: {
    resetCanvas: () => {
      HatService.reset()
    }
  }
}
</script>

<style scoped>

</style>