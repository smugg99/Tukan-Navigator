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
      :x="middleX - 25"
      :y="middleY - 25"
      width="50"
      height="30"
      fill="transparent"
      cursor="pointer"
      class="rect-clickable"
    />
    <text
      :x="middleX"
      :y="middleY"
      text-anchor="middle"
      fill="black"
      cursor="pointer"
      :class="{ 'edge-highlighted': highlighted }"
      dominant-baseline="central"
      font-size="16px"
      font-weight="bold"
      :transform="`rotate(${adjustedAngle}, ${middleX}, ${middleY}) translate(0, -12)`"
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
  computed: {
    angle() {
      const dx = this.to.x - this.from.x;
      const dy = this.to.y - this.from.y;
      return (Math.atan2(dy, dx) * 180) / Math.PI;
    },
    adjustedAngle() {
      let angle = this.angle;
      if (angle > 90) {
        angle -= 180;
      } else if (angle < -90) {
        angle += 180;
      }
      return angle;
    },
    middleX() {
      return (this.from.x + this.to.x) / 2;
    },
    middleY() {
      return (this.from.y + this.to.y) / 2;
    }
  },
  methods: {
    selectEdge() {
      this.$emit('select', this.edge);
    },
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
  stroke: slategray;
}

text {
  cursor: pointer;
  user-select: none;
  transition: fill 0.3s ease;
}

text.edge-highlighted {
  fill: gray;
}

.rect-clickable {
  cursor: pointer;
}
</style>
