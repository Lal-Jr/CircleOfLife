import { useInfiniteQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Post } from "@/types/post";

export function useFeed(lat: number | null, lng: number | null, radiusKm: number = 5) {
    const fetchPosts = async ({ pageParam = 1 }) => {
        if (lat === null || lng === null) return { posts: [], nextPage: undefined };

        const response = await api.get(`/posts`, {
            params: { lat, lng, radius: radiusKm, page: pageParam, limit: 10 }
        });

        const payloadData = response.data.data as Post[];
        const hasNext = response.data.meta?.hasNext || false;

        return {
            posts: payloadData,
            nextPage: hasNext ? pageParam + 1 : undefined,
        };
    };

    const query = useInfiniteQuery({
        queryKey: ["feed", lat, lng, radiusKm],
        queryFn: fetchPosts,
        initialPageParam: 1,
        getNextPageParam: (lastPage) => lastPage.nextPage,
        enabled: lat !== null && lng !== null,
    });

    // Flatten pages for easy iteration inside the UI
    const posts = query.data?.pages.flatMap((page) => page.posts) || [];

    return {
        ...query,
        posts,
    };
}
