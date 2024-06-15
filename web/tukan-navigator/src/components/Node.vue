<template>
  <g @mousedown="select">
    <circle
      :cx="x"
      :cy="y"
      r="20"
      :fill="selected ? 'red' : 'blue'"
    />
    <text
      :x="x"
      :y="y"
      text-anchor="middle"
      fill="white"
      dy=".3em"
    >{{ id }}</text>
  </g>
</template>

<script>
export default {
  props: {
    id: String,
    x: Number,
    y: Number,
    selected: Boolean
  },
  methods: {
    select() {
      this.$emit('select', this.id);
    }
  },
};
</script>

<style scoped>
circle {
  cursor: pointer;
  user-select: none;
}

text {
  user-select: none;
}
</style>