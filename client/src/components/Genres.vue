<template>
<div>  
    <v-carousel hide-delimiters id="genre-carousel">
        <v-carousel-item
        v-for="n in carouselLength"
        :key="n"
        >
            <v-sparkline
            :value="genreFreq[n - 1]"
            :labels="labels[n - 1]"
            color="#1db954"
            :smooth="2.5"
            :label-size="2.25"
            auto-draw
            ></v-sparkline>
        </v-carousel-item>
    </v-carousel>
</div>
</template>
<script>
import SpotifyService from '@/services/SpotifyService.js';
export default {
    name: 'Genres',
    data() {
        return {
            genreFreq: [],
            labels: [],
            carouselLength: 0,
        }
    },
    created() {
        this.getGenreData()
    },
    methods: {
        async getGenreData() {
            SpotifyService.getGenres().then(
                (genres => {
                    let freqSlices = [];
                    let labelSlices = [];
                    let vals = [];
                    let keys = [];
                    for(let [k, v] of Object.entries(genres)){
                        vals.push(v);
                        keys.push(k);
                        if (vals.length == 10)  {
                            freqSlices.push(vals);
                            labelSlices.push(keys);
                            vals = [];
                            keys = [];
                        }
                    }
                    freqSlices.push(vals);
                    labelSlices.push(keys);
                    this.$set(this, "carouselLength", freqSlices.length)
                    this.$set(this, "genreFreq", freqSlices);
                    this.$set(this, "labels", labelSlices);
                }).bind(this)
            );
        }
    }
}
</script>
<style>
    #genre-carousel {
        font-family: 'Roboto Condensed'
    }
</style>