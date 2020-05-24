<template>
    <div>
        <div v-for="artist in artists" :key="artist.name">
            <h3>{{artist.name}}</h3>
        </div>
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