import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Post } from "@/types/post";

export function usePost(id: string, lat?: number | null, lng?: number | null) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["post", id, lat, lng],
        queryFn: async () => {
            const response = await api.get(`/posts/${id}`, {
                params: lat != null && lng != null ? { lat, lng } : undefined,
            });
            return response.data.data as Post;
        },
        enabled: !!id,
    });

    return {
        post: data,
        isLoading,
        error,
    };
}
