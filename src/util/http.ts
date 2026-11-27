import axios, { AxiosRequestConfig } from 'axios'

const PROD_API = 'https://onij.fun/go'

function resolveServerUrl(): string {
  const fromEnv = (import.meta.env.VITE_SERVER_URL as string | undefined)?.replace(/\/$/, '') || ''
  if (!import.meta.env.PROD) {
    // 本机 Go 走 Vite 同源反代，避免浏览器把 127.0.0.1:8889 当局域网请求拦掉
    if (!fromEnv || /127\.0\.0\.1|localhost/i.test(fromEnv)) return ''
    return fromEnv
  }
  if (fromEnv && !/127\.0\.0\.1|localhost/i.test(fromEnv)) return fromEnv
  return PROD_API
}

const SERVER_URL = resolveServerUrl()

export function getServerBaseUrl(): string {
  return SERVER_URL
}
const NETEASE_URL = import.meta.env.VITE_NETEASE_URL;

export const apiClient = axios.create({
    baseURL: SERVER_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    },
});


apiClient.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

apiClient.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export const get = async <T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.get(url, { params, ...config }).then(response => response.data);
};

/**
 * POST 并直接返回响应体（已解包 response.data）。
 * 调用方不要再取一次 `.data`，否则拿到的是 undefined。
 */
export const post = async <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return apiClient.post<T>(url, data, config).then(response => response.data);
};

export const uploadFile =  async<T>(url: string, formData: FormData, config?: AxiosRequestConfig): Promise<T> => {
    const uploadConfig = {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
        ...config,
    };
    return apiClient.post(url, formData, uploadConfig).then(response => response.data);
};



const neteaseApiClient = axios.create({
    baseURL: NETEASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    },
});

export const getNetease =  async <T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
    return neteaseApiClient.get(url, { params, ...config }).then(response => response.data);
};

export const postNetease = async <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return neteaseApiClient.post<T>(url, data, config).then(response => response.data);
};