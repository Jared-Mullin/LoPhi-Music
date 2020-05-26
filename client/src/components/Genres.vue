<template>
    <v-container>
        <v-sparkline
        :value="genreFreq"
        :labels="labels"
        color="#1db954"
        :line-width="1"
        :smooth="5"
        :label-size="2"
        auto-draw
        ></v-sparkline>
    </v-container>
</template>
<script>
import SpotifyService from '@/services/SpotifyService.js';
export default {
    name: 'Genres',
    data() {
        return {
            genreFreq: [],
            labels: [],
        }
    },
    created() {
        this.getGenreData()
    },
    methods: {
        async getGenreData() {
            SpotifyService.getGenres().then(
                (genres => {
                    let vals = [];
                    let keys = [];
                    for(let [k, v] of Object.entries(genres)){
                        if (v > 1) {
                            vals.push(v);
                            keys.push(k);
                        }
                    }
                    this.$set(this, "genreFreq", vals);
                    this.$set(this, "labels", keys);
                }).bind(this)
            );
        }
    }
}
</script>
<style>
    .v-container {
        margin: 2.5%;
    }
</style>