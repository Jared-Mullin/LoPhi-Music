<template>
  <v-app>
    <v-content>
      <v-app-bar
          color="#212121"
          dark
      >
        <v-toolbar-title id="home-link" class="nav-link">LoPhi</v-toolbar-title>
        <v-spacer></v-spacer>
        <router-link class="nav-link" to="/tracks">Tracks</router-link>
        <router-link class="nav-link" to="/artists">Artists</router-link>
        <router-link class="nav-link" to="/genres">Genres</router-link>
      </v-app-bar>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
export default {
  name: 'App',

  components: {
  },

  data: () => ({

  }),
};
</script>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Roboto+Condensed');
  @import url('https://fonts.googleapis.com/css2?family=Roboto');
  .v-content {
    font-family: 'Roboto Condensed', 'Roboto', Avenir, sans-serif !important;
  }
  .nav-link {
    padding-left: 1%;
    padding-right: 1%;
  }
  #home-link {
    font-size: 1.75em;
  }
  a {
    text-decoration: none;
  }
</style>