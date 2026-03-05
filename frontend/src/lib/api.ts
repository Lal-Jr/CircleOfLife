import axios from "axios";

export const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
    headers: {
        "Content-Type": "application/json",
    },
});

// Intercept requests to add the JWT token
api.interceptors.request.use(
    (config) => {
        // Only run on the client side
        if (typeof window !== "undefined") {
            const token = localStorage.getItem("jwt_token");
            if (token && config.headers) {
                config.headers.Authorization = `Bearer ${token}`;
            }
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);
