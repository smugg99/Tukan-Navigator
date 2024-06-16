<template>
  <g @mousedown="select" :class="{'selectable': isSelectable, 'highlighted': highlighted}">
    <circle
      :cx="x"
      :cy="y"
      r="20"
      :class="{ 'node-selected': selected, 'node-highlighted': highlighted }"
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
    selected: Boolean,
    highlighted: Boolean,
    mode: String,
  },
  computed: {
    isSelectable() {
      return this.mode === 'remove' || this.mode === 'edit' || this.mode === 'addEdge' || this.mode === 'drag';
    }
  },
  methods: {
    select() {
      if(this.isSelectable) {
        this.$emit('select', this.id);
      }
    }
  },
};
</script>

<style scoped>
circle {
  z-index: 2;
  user-select: none;
  transition: fill 0.3s ease;
}

circle.node-selected {
  fill: rgb(75, 168, 255);
}

circle.node-highlighted {
  fill: rgb(75, 168, 255);
}

circle:not(.node-selected):not(.node-highlighted) {
  fill: rgb(0, 132, 255);
}

text {
  user-select: none;
}

.selectable {
  cursor: pointer;
}
</style>
