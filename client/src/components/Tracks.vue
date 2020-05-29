<template>
    <div>
        <v-layout>
            <v-row class="track-row">
                <v-col :lg="3" :md="4" :sm="6" :cols="12" v-for="track in tracks" :key="track.id">
                    <v-card :href="track.external_urls.spotify">
                        <v-card-title class="tracks-card-title">{{track.name}}</v-card-title>
                        <v-img :src="track.album.images[0].url" aspect-ratio="1"></v-img>
                    </v-card>
                </v-col>
            </v-row>
        </v-layout>
    </div>
</template>
<script>
import SpotifyService from '@/services/SpotifyService.js';
export default {
    name: 'Tracks',
    data() {
        return {
            tracks: []
        }
    },
    created() {
        const cfg = {
            headers: {
                'Authorization': `Bearer ${this.$cookies.get('token')}`
            }
        };
        this.getTrackData(cfg)
    },
    methods: {
        async getTrackData(cfg) {
            SpotifyService.getTracks(cfg).then(
                (tracks => {
                    this.$set(this, "tracks", tracks.items)
                }).bind(this)
            );
        }
    }
}
</script>
<style>
    .tracks-card-title {
        font-family: 'Roboto Condensed';
        font-size: 2vh;
    }
</style>