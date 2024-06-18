<template>
  <g v-if="from && to" @click="selectEdge" :class="{ 'highlighted': highlighted, 'traversed': traversed }">
    <line :x1="from.x" :y1="from.y" :x2="to.x" :y2="to.y" :stroke="lineColor" stroke-dasharray="5, 5" cursor="pointer"
      :class="{
        'edge-highlighted': highlighted,
        'edge-traversed': traversed
      }" />
    <rect :x="middleX - 25" :y="middleY - 25" width="50" height="30" fill="transparent" cursor="pointer"
      class="rect-clickable" />
    <text :x="middleX" :y="middleY" text-anchor="middle" :fill="textColor" cursor="pointer" :class="{
      'edge-highlighted': highlighted,
      'edge-traversed': traversed
    }" dominant-baseline="central" font-size="16px" font-weight="bold"
      :transform="`rotate(${adjustedAngle}, ${middleX}, ${middleY}) translate(0, -12)`">
      {{ edge.weight }}
    </text>
  </g>
</template>

<script>
import { useTheme } from 'vuetify'

export default {
  props: {
    from: Object,
    to: Object,
    edge: Object,
    highlighted: Boolean,
    traversed: Boolean,
  },
  setup() {
    const theme = useTheme()
    const currentTheme = theme.global.name

    return {
      theme,
      currentTheme
    };
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
    },
    lineColor() {
      if (this.currentTheme === 'dark') {
        return this.highlighted ? 'gray' : 'white';
      } else {
        return this.highlighted ? 'white' : 'black'; // Change the color to black for light theme
      }
    },
    textColor() {
      if (this.currentTheme === 'dark') {
        return this.highlighted ? 'gray' : 'white';
      } else {
        return this.highlighted ? 'white' : 'black'; // Change the color to black for light theme
      }
    },
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

line.edge-traversed {
  stroke: yellow;
  /* Change stroke color for traversed edges */
}

text {
  cursor: pointer;
  user-select: none;
  transition: fill 0.3s ease;
}

text.edge-highlighted {
  fill: gray;
}

text.edge-traversed {
  fill: yellow;
  /* Change text color for traversed edges */
}

.rect-clickable {
  cursor: pointer;
}
</style>
