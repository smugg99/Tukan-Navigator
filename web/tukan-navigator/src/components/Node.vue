<template>
  <g @mousedown="select" :class="{'selectable': isSelectable}">
    <circle
      :cx="x"
      :cy="y"
      r="20"
      :class="{ 'node-selected': selected }"
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
  user-select: none;
  transition: fill 0.3s ease;
}

circle.node-selected {
  fill: red;
}

circle:not(.node-selected) {
  fill: blue;
}

text {
  user-select: none;
}

.selectable {
  cursor: pointer;
}
</style>