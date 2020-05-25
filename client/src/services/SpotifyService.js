import axios from 'axios'

export default {
    async getArtists() {
        let res = await axios.get("http://localhost:4200/spotify/artists");
        return res.data;
    },
    async getTracks() {
        let res = await axios.get("http://localhost:4200/spotify/tracks");
        return res.data;
    },
    async getGenres() {
        let res = await axios.get("http://localhost:4200/spotify/genres");
        return res.data;
    }
}