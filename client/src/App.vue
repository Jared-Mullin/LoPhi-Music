<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
    </v-app-bar>

    <v-content>
      <Tracks/>
      <Artists/>
    </v-content>
  </v-app>
</template>

<script>
import Artists from './components/Artists';
import Tracks from './components/Tracks';

export default {
  name: 'App',

  components: {
    Artists,
    Tracks,
  },

  data: () => ({
    //
  }),
};
</script>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Roboto+Condensed');
  @import url('https://fonts.googleapis.com/css2?family=Roboto');
  .v-application {
    font-family: 'Roboto Condensed', Avenir, sans-serif !important;
  }
</style>