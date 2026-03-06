import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Post } from "@/types/post";

export function useFeed(lat: number | null, lng: number | null, radiusMs: number = 5) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["posts", lat, lng, radiusMs],
        queryFn: async () => {
            if (lat === null || lng === null) return [];

            const response = await api.get(`/posts`, {
                params: { lat, lng, radius: radiusMs }
            });
            // The backend will return an array of posts. For now, if no backend, we just return empty or error.
            // But the integration plan states we should hit the real API.
            // So we assume the REST API returns something like `{ posts: Post[] }` or `Post[]`
            return response.data as Post[];
        },
        enabled: lat !== null && lng !== null,
    });

    return {
        posts: data || [],
        isLoading,
        error,
    };
}
