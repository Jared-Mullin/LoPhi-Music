<template>
<div>  
    <v-carousel 
    id="genre-carousel"
    hide-delimiters
    height="auto"
    width="auto">
        <v-carousel-item
        v-for="n in carouselLength"
        :key="n"
        >
            <v-sparkline
            class="genre-graph"
            :value="genreFreq[n - 1]"
            :labels="labels[n - 1]"
            color="#1db954"
            :smooth="2.5"
            :label-size="3"
            padding="16"
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
        const cfg = {
            headers: {
                'Authorization': `Bearer ${this.$cookies.get('token')}`
            }
        };
        this.getGenreData(cfg)
    },
    methods: {
        async getGenreData(cfg) {
            SpotifyService.getGenres(cfg).then(
                (genres => {
                    let freqSlices = [];
                    let labelSlices = [];
                    let vals = [];
                    let keys = [];
                    for(let [k, v] of Object.entries(genres)){
                        vals.push(v);
                        keys.push(k);
                        if (vals.length == 5)  {
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