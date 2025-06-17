import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'

const SERVER_URL = import.meta.env.VITE_SERVER_URL || "http://127.0.0.1:8889/";

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
        }
        if (config.method === 'post' || config.method === 'put') {
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
        return Promise.reject(error);
    }
);

// GET 请求
export const get = <T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.get(url, { params, ...config });
};

// POST 请求
export const post = <T = any, R = AxiosResponse<T>, _D = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<R> => {
    return apiClient.post(url, data, config);
};

// 上传文件专用 POST 请求
export const uploadFile = <T>(url: string, formData: FormData, config?: AxiosRequestConfig): Promise<T> => {
    const uploadConfig = {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
        ...config,
    };
    return apiClient.post(url, formData, uploadConfig).then(response => response.data);
};

export default apiClient;
