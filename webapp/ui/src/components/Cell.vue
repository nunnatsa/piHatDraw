<template>
  <td :style="cssVars">{{tool}}</td>
</template>

<script>
export default {
  name: "Cell",
  props: [
    'bgColor', 'tool', 'borders',
  ],
  computed: {
    cssVars() {
      return {
        /* variables you want to pass to css */
        '--bgColor': this.bgColor,
        '--color': this.reverseBgColor,
        '--topBorder': `${(this.borders.top) ? '3px' : '1px'}`,
        '--leftBorder': `${(this.borders.left) ? '3px' : '1px'}`,
        '--bottomBorder': `${(this.borders.bottom) ? '3px' : '1px'}`,
        '--rightBorder': `${(this.borders.right) ? '3px' : '1px'}`,
      }
    },
    reverseBgColor() {
      const rx = /#(..)(..)(..)/
      let colors = this.bgColor.match(rx)
      if (!colors) {
        return '#ffffff'
      }

      const r = (255 - parseInt(colors[1], 16)).toString(16).padStart(2, '0')
      const g = (255 - parseInt(colors[2], 16)).toString(16).padStart(2, '0')
      const b = (255 - parseInt(colors[3], 16)).toString(16).padStart(2, '0')

      return `#${r}${g}${b}`
    },
  }
}

</script>

<style scoped>
  td {
    height: 20px;
    width: 20px;
    color: var(--color);
    background-color: var(--bgColor);
    text-align: center;
    vertical-align: middle;
    border-style: solid;
    border-color: white;
    margin: 0;
    font-size: 0.7em;
    font-weight: bolder;
    border-top-width: var(--topBorder);
    border-left-width: var(--leftBorder);
    border-bottom-width: var(--bottomBorder);
    border-right-width: var(--rightBorder);
    padding: 0;
  }
</style>
