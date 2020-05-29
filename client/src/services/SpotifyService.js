import axios from 'axios'

export default {
    async getArtists() {
        let cookie = this.$cookies.get("token")
        const cfg = {
            headers: {
                'Authorization': `Bearer ${cookie}`
            }
        };
        let res = await axios.get("http://localhost:4200/spotify/artists", cfg);
        return res.data;
    },
    async getTracks(cfg) {
        let res = await axios.get("http://localhost:4200/spotify/tracks", cfg);
        return res.data;
    },
    async getGenres(cfg) {
        let res = await axios.get("http://localhost:4200/spotify/genres", cfg);
        return res.data;
    }
}