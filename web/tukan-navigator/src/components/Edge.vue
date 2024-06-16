<template>
  <g v-if="from && to" @click="selectEdge" :class="{ 'highlighted': highlighted }">
    <line
      :x1="from.x"
      :y1="from.y"
      :x2="to.x"
      :y2="to.y"
      stroke="black"
      stroke-dasharray="5, 5"
      cursor="pointer"
      :class="{ 'edge-highlighted': highlighted }"
    />
    <rect
      :x="(from.x + to.x) / 2 - 25"
      :y="(from.y + to.y) / 2 - 15"
      width="50"
      height="30"
      fill="transparent"
      cursor="pointer"
      class="rect-clickable"
    />
    <text
      :x="(from.x + to.x) / 2"
      :y="(from.y + to.y) / 2"
      text-anchor="middle"
      fill="black"
      cursor="pointer"
      :class="{ 'edge-highlighted': highlighted }"
      dominant-baseline="central"
      font-size="16px"
      font-weight="bold"
    >
      {{ edge.weight }}
    </text>
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
  stroke-width: 2px;
  transition: stroke 0.3s ease;
}

line.edge-highlighted {
  stroke: orange;
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
</style>
