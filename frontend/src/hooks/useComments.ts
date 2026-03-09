import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Comment } from "@/types/comment";

export function useComments(postId: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["comments", postId],
        queryFn: async () => {
            const response = await api.get(`/posts/${postId}/comments`);
            return response.data.data as Comment[];
        },
        enabled: !!postId,
    });

    return {
        comments: data || [],
        isLoading,
        error,
    };
}
