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

// If the token is missing/expired/invalid, clear it and send the user back to login
// instead of leaving every page stuck on a generic error state.
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (
            typeof window !== "undefined" &&
            error.response?.status === 401 &&
            !window.location.pathname.startsWith("/login") &&
            !window.location.pathname.startsWith("/signup")
        ) {
            localStorage.removeItem("jwt_token");
            window.location.href = "/login";
        }
        return Promise.reject(error);
    }
);
