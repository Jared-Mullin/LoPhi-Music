import axios from 'axios'

export default {
    async getArtists() {
        let res = await axios.get("http://localhost:4200/spotify/artists");
        return res.data;
    }
}