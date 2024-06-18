<template>
  <v-app>
    <v-main>
      <AppBar title="Tukan Navigator" />
      <Graph />
    </v-main>
  </v-app>
</template>

<script>
import AppBar from './components/AppBar.vue';
import Graph from './components/Graph.vue';

export default {
  name: 'App',
  components: {
    AppBar,
    Graph,
  },
};
</script>

<style>
html,
body {
  margin: 0;
  overscroll-behavior-y: contain;
  overscroll-behavior: contain;
  overflow: hidden;
  height: 100%;
}
</style>