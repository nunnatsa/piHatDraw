module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    "proxy": "http://pihatdraw:8080/"
  }
}
