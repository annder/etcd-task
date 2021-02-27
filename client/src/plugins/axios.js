import Axios from 'axios'

const instance = Axios.create({
  headers: {
    'content-type': 'application/x-www-form-urlencoded',
  },
})

instance.defaults.baseURL = 'http://localhost:3000'

instance.interceptors.request.use((config) => config)
instance.interceptors.response.use(
  (response) => {
    if (response.data.code === 200) {
      return response.data
    }
  },
  (err) => {
    console.log(err)
  }
)

export default instance
