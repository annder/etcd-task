import axios from "axios"

const http = axios.create({
    timeout: 6000,
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
})
http.defaults.baseURL = "http://localhost:3000/"

http.defaults.transformRequest = [function (data) {
    let ret = ''
    for (let it in data) {
      ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&'
    }
    return ret
}]

http.interceptors.request.use(config => config)

http.interceptors.response.use(response => {
    if (response.data && response.data.code === 200) {
        return response.data.data
    }
})

export default http