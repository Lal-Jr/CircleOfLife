import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useCreateComment() {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async ({ postId, content }: { postId: string; content: string }) => {
            const response = await api.post(`/posts/${postId}/comments`, { content });
            return response.data;
        },
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ["comments", variables.postId] });
            queryClient.invalidateQueries({ queryKey: ["post", variables.postId] }); // updates commentCount
        },
    });

    return {
        createComment: mutation.mutate,
        isCreating: mutation.isPending,
        error: mutation.error,
    };
}
