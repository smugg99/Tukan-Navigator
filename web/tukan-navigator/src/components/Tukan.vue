<template>
  <circle
    :cx="currentPosition.x"
    :cy="currentPosition.y"
    :r="10"
    fill="blue"
  />
</template>

<script>
import { ref, watch, onMounted } from 'vue';

export default {
  props: ['path'],
  setup(props) {
    const currentPosition = ref({ x: 0, y: 0 });
    const pathIndex = ref(0);

    const moveToNextNode = () => {
      if (pathIndex.value < props.path.length) {
        const nextNode = props.path[pathIndex.value];
        currentPosition.value = { x: nextNode.x, y: nextNode.y };
        pathIndex.value++;
        setTimeout(moveToNextNode, 1000); // Adjust the interval as needed
      }
    };

    const start = () => {
      pathIndex.value = 0;
      moveToNextNode();
    };

    watch(() => props.path, (newPath) => {
      if (newPath.length > 0) {
        start();
      }
    });

    onMounted(() => {
      if (props.path.length > 0) {
        start();
      }
    });

    return {
      currentPosition,
      start,
    };
  },
};
</script>
