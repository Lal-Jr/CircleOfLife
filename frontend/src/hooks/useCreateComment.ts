import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Comment } from "@/types/comment";

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

            queryClient.setQueryData(["comments", newComment.postId], (old: Comment[] | undefined) => {
                const optimisticComment: Comment = {
                    id: `temp-${Date.now()}`,
                    postId: newComment.postId,
                    authorId: "me",
                    authorName: "You",
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
            // The feed's post cards also show commentCount, so they'd otherwise
            // go stale until the query's own staleTime elapses.
            queryClient.invalidateQueries({ queryKey: ["feed"] });
        },
    });

    return {
        createComment: mutation.mutate,
        isCreating: mutation.isPending,
        error: mutation.error,
    };
}
