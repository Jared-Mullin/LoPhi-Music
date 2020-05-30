<template>
    <div>
        <v-layout>
            <v-row class="artist-row">
                <v-col :lg="3" :md="4" :sm="6" :cols="12" v-for="artist in artists" :key="artist.id">
                    <v-card :href="artist.external_urls.spotify">
                        <v-card-title class="artist-card-title">{{artist.name}}</v-card-title>
                        <v-img :src="artist.images[0].url" aspect-ratio="1"></v-img>
                    </v-card>
                </v-col>
            </v-row>
        </v-layout>
    </div>
</template>
<script>
import axios from 'axios';

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
            let token = localStorage.getItem('token');
            let config = {
                headers: {
                    'Authorization': 'Bearer ' + token,
                },
            }
            axios.get('https://lophi.dev/spotify/artists', config).then(
                (artists => {
                    artists = artists.data;
                    this.$set(this, 'artists', artists.items)
                }).bind(this)
            );
        }
    }
}
</script>
<style>
    .artist-card-title {
        font-family: 'Roboto Condensed';
        font-size: 2vh;
    }
</style>