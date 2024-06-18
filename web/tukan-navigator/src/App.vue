<template>
  <v-app>
    <AppBar title="Tukan Navigator" />
    <v-main>

      <Graph></Graph>
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
  overflow: hidden;
  height: 100%;
  margin: 0;
}
</style>