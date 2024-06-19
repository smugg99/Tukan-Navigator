<template>
  <v-app-bar :elevation="2" prominent scroll-threshold="0">
    <template v-slot:prepend>
    </template>

    <v-app-bar-title>{{ title }}</v-app-bar-title>

    <template v-slot:append>
    </template>
    <v-btn class="mr-2" :icon="currentTheme === 'light' ? 'mdi-weather-sunny' : 'mdi-weather-night'"
      @click="toggleTheme">
    </v-btn>
    <v-btn class="mr-2" icon href="https://github.com/smugg99/Tukan-Navigator">
      <v-icon>mdi-github</v-icon>
    </v-btn>
    <slot></slot>
  </v-app-bar>
</template>

<script>
import { useTheme } from 'vuetify'

export default {
  name: 'AppBar',
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  setup() {
    const theme = useTheme()
    const currentTheme = theme.global.name
    const toggleTheme = () => {
      const newTheme = currentTheme.value === 'light' ? 'dark' : 'light'
      theme.global.name.value = newTheme
      localStorage.setItem('user-theme', newTheme)
    }
    const setInitialTheme = () => {
      const userTheme = localStorage.getItem('user-theme')
      if (userTheme) {
        theme.global.name.value = userTheme
      } else {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        theme.global.name.value = prefersDark ? 'dark' : 'light'
      }
    }

    setInitialTheme()

    return {
      currentTheme,
      toggleTheme,
    }
  },
};
</script>
