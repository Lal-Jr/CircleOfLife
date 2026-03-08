import { useEffect, useState } from "react";
import { useAuth } from "./useAuth";

export function useRealtimeFeed() {
    const [hasNewPosts, setHasNewPosts] = useState(false);
    const { token } = useAuth();
    const baseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

    useEffect(() => {
        if (!token) return;

        // We pass the token via query param because EventSource does not inherently support Custom Headers out of the box in browsers.
        // The backend middleware would technically accept a query `?token=` if we want to be strict, but we'll assume it handles standard cookies or we adjust the Auth Middleware.
        // For this prototype, I am utilizing an open EventSource mapping to the endpoint directly since Next handles the proxying or the Go CORS block permits it.
        const eventSource = new EventSource(`${baseURL}/events?token=${token}`);

        eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === "post_created") {
                    setHasNewPosts(true);
                }
            } catch (err) {
                console.error("Error parsing SSE event", err);
            }
        };

        eventSource.onerror = (err) => {
            console.error("EventSource encountered an error", err);
            eventSource.close();
        };

        return () => {
            eventSource.close();
        };
    }, [baseURL, token]);

    const clearNewPosts = () => {
        setHasNewPosts(false);
    };

    return { hasNewPosts, clearNewPosts };
}
