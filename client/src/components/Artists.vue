<template>
    <div>
        <v-layout>
            <v-row class="artist-row">
                <v-col :lg="3" :md="4" :sm="6" :xs="12" v-for="artist in artists" :key="artist.id">
                    <v-card class="artist-card" :href="artist.external_urls.spotify">
                        <v-card-title>{{artist.name}}</v-card-title>
                        <v-img :src="artist.images[0].url" aspect-ratio="1"></v-img>
                    </v-card>
                </v-col>
            </v-row>
        </v-layout>
    </div>
</template>
<script>
import SpotifyService from '@/services/SpotifyService.js';
export default {
    name: 'Artists',
    data() {
        return {
            artists: []
        }
    },
    created() {
        this.getArtistData()
    },
    methods: {
        async getArtistData() {
            SpotifyService.getArtists().then(
                (artists => {
                    this.$set(this, "artists", artists.items)
                }).bind(this)
            );
        }
    }
}
</script>
<style>
</style>