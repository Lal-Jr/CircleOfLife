import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Post } from "@/types/post";

export function usePost(id: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["post", id],
        queryFn: async () => {
            const response = await api.get(`/posts/${id}`);
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
