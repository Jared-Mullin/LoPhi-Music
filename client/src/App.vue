<template>
  <v-app>
    <v-content>
      <Genres/>
      <Tracks/>
      <Artists/>
    </v-content>
  </v-app>
</template>

<script>
import Artists from './components/Artists';
import Tracks from './components/Tracks';
import Genres from './components/Genres';

export default {
  name: 'App',

  components: {
    Artists,
    Tracks,
    Genres,
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