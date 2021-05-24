<template>
  <div>
    <v-btn
        small
        @click.stop="downloadDialog = true"
        color="#8888ee"
        :disabled="disabled"
    ><v-icon>mdi-download</v-icon>
      Download
    </v-btn>

    <v-dialog
        v-model="downloadDialog" max-width="350"
    >
      <v-card>
        <v-card-title class="headline">
          Download the Picture
        </v-card-title>

        <v-card-text>
          <v-slider v-model="pixelSize" min="1" max="20" step="1" persistent-hint :hint="hint"></v-slider>
        </v-card-text>

        <v-text-field
            label="File Name"
            placeholder="untitled.png"
            v-model="fileName"
            :rules="rules.fileName"
            filled
        ></v-text-field>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
              color="green darken-1"
              text
              @click="downloadDialog = false"
          >
            <v-icon>mdi-cancel</v-icon>
            Cancel
          </v-btn>
          <v-btn
              color="green darken-1"
              text
              :disabled="!nameValid(fileName)"
              @click="downloadDialog = false; download(); fileName=''; pixelSize=3"
          >
            <v-icon>mdi-download</v-icon>
            Download
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import HatService from '../services'
export default {
  name: "DownloadButton",
  data() {
    return {
      downloadDialog: false,
      pixelSize: 3,
      fileName: '',
      rules: {
        fileName: [(fileName) => this.nameValid(fileName) || 'must be ends with ".png"']
      },
    }
  },
  props: [
    "disabled",
  ],
  methods: {
    download: function(){
      HatService.download({pixelSize: this.pixelSize, fileName: this.fileName})
    },
    nameValid: function (fileName) {
      return fileName.endsWith(".png")
    }
  },
  computed: {
    hint: function() {
      return `image pixel size: ${this.pixelSize} pixel${this.pixelSize === 1 ? '' : 's'}`
    },

  },
}
</script>

<style scoped>

</style>