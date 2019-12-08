import axios from 'axios'

export default {
    defaceBaseUrl() {
        return process.env.VUE_APP_DEFACE_BASE_URL
    },

    scanUrl(data) {
        return axios({
            method: 'post',
            url: `${this.defaceBaseUrl()}/deface/scan`,
            data: data
        }).then(response => response.data)
    },

    getUrlReport(data) {
        return axios({
            method: 'post',
            url: `${this.defaceBaseUrl()}/deface/report`,
            data: data
        }).then(response => response.data)
    },
}