<template>
  <g v-if="from && to" @click="selectEdge" :class="{ 'highlighted': highlighted }">
    <line
      :x1="from.x"
      :y1="from.y"
      :x2="to.x"
      :y2="to.y"
      stroke="black"
      cursor="pointer"
      :class="{ 'edge-highlighted': highlighted }"
    />
    <g class="clickable-text">
      <rect
        :x="(from.x + to.x) / 2 - 25"
        :y="(from.y + to.y) / 2 - 15"
        width="50"
        height="30"
        fill="transparent"
        cursor="pointer"
      />
      <text
        :x="(from.x + to.x) / 2"
        :y="(from.y + to.y) / 2"
        text-anchor="middle"
        fill="black"
        cursor="pointer"
        :class="{ 'edge-highlighted': highlighted }"
      >{{ edge.weight }}</text>
    </g>
  </g>
</template>

<script>
export default {
  props: {
    from: Object,
    to: Object,
    edge: Object,
    highlighted: Boolean,
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
  user-select: none;
  stroke-width: 6px;
  transition: stroke 0.3s ease;
}

line.edge-highlighted {
  stroke: orange;
  stroke-width: 6px;
}

text {
  cursor: pointer;
  user-select: none;
  transition: fill 0.3s ease;
  position: relative;
  z-index: 1;
  text-align: center; /* Add this line to center the text */
}

text.edge-highlighted {
  fill: orange;
}

.rect-clickable {
  cursor: pointer;
}

.highlighted {
  cursor: pointer;
}
</style>
