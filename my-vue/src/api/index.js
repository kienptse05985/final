import axios from 'axios'

const customHeaders = {
  'content-type': 'application/json',
};
export default {
  defaceBaseUrl() {
    return process.env.VUE_APP_DEFACE_BASE_URL
  },

  scanUrl(data) {
    return axios.post(`${this.defaceBaseUrl()}/api/v1/scan`, data, {headers: customHeaders}).then(response => response.data)
  },

  getUrlReport(data) {
    return axios({
      method: 'post',
      url: `${this.defaceBaseUrl()}/report`,
      data: data,
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(response => response.data)
  },
}
