<template>
  <circle
    v-if="animationRunning && animatedPath.length > 0 && pathIndex < animatedPath.length"
    ref="toucanSvg"
    :cx="currentPosition.x + panX"
    :cy="currentPosition.y + panY"
    r="10"
    fill="green"
  />
</template>

<script>
export default {
  props: ['animatedPath', 'pathIndex', 'animationRunning', 'panX', 'panY'],
  data() {
    return {
      currentPosition: { x: 0, y: 0 }
    };
  },
  methods: {
    moveToucanTo(x, y) {
      this.$refs.toucanSvg.setAttribute('transform', `translate(${x},${y})`);
    }
  },
  watch: {
    pathIndex(newIndex) {
      if (this.animationRunning && this.animatedPath.length > 0 && newIndex < this.animatedPath.length) {
        const currentNode = this.animatedPath[newIndex];
        this.currentPosition = { x: currentNode.x, y: currentNode.y };
      }
    }
  }
};
</script>
