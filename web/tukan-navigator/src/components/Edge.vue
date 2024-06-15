<template>
  <g v-if="from && to" @click="selectEdge" :class="{ 'highlighted': highlighted }">
    <line
      :x1="from.x"
      :y1="from.y"
      :x2="to.x"
      :y2="to.y"
      stroke="black"
      :class="{ 'edge-highlighted': highlighted }"
    />
    <g class="clickable-text">
      <rect
        :x="(from.x + to.x) / 2 - 20"
        :y="(from.y + to.y) / 2 - 10"
        width="40"
        height="20"
        fill="transparent"
        cursor="pointer"
      />
      <text
        :x="(from.x + to.x) / 2"
        :y="(from.y + to.y) / 2"
        text-anchor="middle"
        fill="black"
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
  transition: stroke 0.3s ease;
}

line.edge-highlighted {
  stroke: orange;
  stroke-width: 3px;
}

text {
  cursor: pointer;
  user-select: none;
  transition: fill 0.3s ease;
}

text.edge-highlighted {
  fill: orange;
}

.rect-clickable {
  cursor: pointer;
}

.highlighted {
  pointer-events: none;
}
</style>