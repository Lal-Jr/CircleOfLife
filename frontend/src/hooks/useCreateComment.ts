import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useCreateComment() {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async ({ postId, content }: { postId: string; content: string }) => {
            const response = await api.post(`/posts/${postId}/comments`, { content });
            return response.data.data;
        },
        onMutate: async (newComment) => {
            await queryClient.cancelQueries({ queryKey: ["comments", newComment.postId] });
            const previousComments = queryClient.getQueryData(["comments", newComment.postId]);

            queryClient.setQueryData(["comments", newComment.postId], (old: any) => {
                const optimisticComment = {
                    id: `temp-${Date.now()}`,
                    postId: newComment.postId,
                    author: "You",
                    authorAvatar: "",
                    content: newComment.content,
                    createdAt: new Date().toISOString(),
                };
                return old ? [...old, optimisticComment] : [optimisticComment];
            });

            return { previousComments };
        },
        onError: (err, newComment, context) => {
            if (context?.previousComments) {
                queryClient.setQueryData(["comments", newComment.postId], context.previousComments);
            }
        },
        onSettled: (data, error, variables) => {
            queryClient.invalidateQueries({ queryKey: ["comments", variables.postId] });
            queryClient.invalidateQueries({ queryKey: ["post", variables.postId] });
        },
    });

    return {
        createComment: mutation.mutate,
        isCreating: mutation.isPending,
        error: mutation.error,
    };
}
