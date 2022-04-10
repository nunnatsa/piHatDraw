module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    "proxy": "http://raspberrypi:8080/"
    // "proxy": "http://pihatdraw:8080/"
  }
}
