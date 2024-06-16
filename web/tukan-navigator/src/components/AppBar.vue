<template>
  <v-app-bar :elevation="2" prominent>
    <template v-slot:prepend>
    </template>

    <v-app-bar-title>{{ title }}</v-app-bar-title>

    <template v-slot:append>
    </template>

    <slot></slot>
  </v-app-bar>
</template>

<script>
export default {
  name: 'AppBar',
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
