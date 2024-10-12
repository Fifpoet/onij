import axios from 'axios'

const SERVER_URL = "http://127.0.0.1:8080/";  // 添加 'http://'

const apiClient = axios.create({
    baseURL: SERVER_URL,
    timeout: 10000, // 请求超时时间
    headers: {
        'Content-Type': 'application/json', // 默认的请求头
        'Accept': 'application/json',
    },
});


apiClient.interceptors.request.use(
    (config) => {
        // 打印请求的 URL 和 body
        if (config.baseURL && config.url) {
            console.log('Request URL:', config.baseURL + config.url);
        }
        if (config.method === 'post' || config.method === 'put') {
            console.log('Request Body:', config.data);
        }

        // 在发送请求之前做点什么，比如添加 token
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    },
    (error) => {
        // 处理请求错误
        return Promise.reject(error);
    }
);


// 添加响应拦截器（可选）
apiClient.interceptors.response.use(
    (response) => {
        // 处理响应数据
        return response;
    },
    (error) => {
        // 处理响应错误
        console.error('API error:', error.response?.data || error.message);
        return Promise.reject(error);
    }
);

export default apiClient;
