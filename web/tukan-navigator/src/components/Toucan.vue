<template>
  <circle
    ref="toucanSvg"
    :cx="x"
    :cy="y"
    :r="25"
    fill="red"
  />
</template>

<script>
export default {
  props: {
    x: Number,
    y: Number,
  },
};
</script>

<style scoped>
circle {
  position: absolute;
  z-index: 6;
  user-select: none;
}
</style>
