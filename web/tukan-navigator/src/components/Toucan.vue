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
  methods: {
    moveTo(x, y) {
      this.x = x;
      this.y = y;
    }
  }
};
</script>