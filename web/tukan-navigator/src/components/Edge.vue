<template>
  <g v-if="from && to" @click="selectEdge">
    <line
      :x1="from.x"
      :y1="from.y"
      :x2="to.x"
      :y2="to.y"
      stroke="black"
    />
    <text
      :x="(from.x + to.x) / 2"
      :y="(from.y + to.y) / 2"
      text-anchor="middle"
      fill="black"
    >{{ edge.weight }}</text>
  </g>
</template>

<script>
export default {
  props: {
    from: Object,
    to: Object,
    edge: Object
  },
  methods: {
    selectEdge() {
      this.$emit('select', this.edge);
    }
  },
};
</script>

<style scoped>
line {
  cursor: pointer;
  /* user-select: none; */
}

text {
  cursor: pointer;
  user-select: none;
}
</style>