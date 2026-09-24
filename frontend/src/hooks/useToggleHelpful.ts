import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Post } from "@/types/post";

interface LikeResponse {
    liked: boolean;
    helpfulCount: number;
}

// Applies a liked/helpfulCount change to a single post, wherever it appears.
function applyToPost(post: Post, postId: string, liked: boolean, helpfulCount: number): Post {
    if (post.id !== postId) return post;
    return { ...post, likedByMe: liked, helpfulCount };
}

export function useToggleHelpful() {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async (postId: string) => {
            const response = await api.post(`/posts/${postId}/like`);
            return response.data.data as LikeResponse;
        },
        onMutate: async (postId: string) => {
            await queryClient.cancelQueries({ queryKey: ["feed"] });
            await queryClient.cancelQueries({ queryKey: ["post", postId] });

            const previousFeed = queryClient.getQueriesData({ queryKey: ["feed"] });
            const previousPost = queryClient.getQueriesData({ queryKey: ["post", postId] });

            // Optimistically flip the like state so the tap feels instant.
            let optimisticLiked = true;
            let optimisticCount = 0;

            queryClient.setQueriesData({ queryKey: ["feed"] }, (oldData: any) => {
                if (!oldData?.pages) return oldData;
                const newPages = oldData.pages.map((page: any) => ({
                    ...page,
                    posts: page.posts.map((p: Post) => {
                        if (p.id !== postId) return p;
                        optimisticLiked = !p.likedByMe;
                        optimisticCount = p.helpfulCount + (optimisticLiked ? 1 : -1);
                        return applyToPost(p, postId, optimisticLiked, optimisticCount);
                    }),
                }));
                return { ...oldData, pages: newPages };
            });

            queryClient.setQueriesData({ queryKey: ["post", postId] }, (oldPost: any) => {
                if (!oldPost) return oldPost;
                const liked = !oldPost.likedByMe;
                const count = oldPost.helpfulCount + (liked ? 1 : -1);
                return applyToPost(oldPost, postId, liked, count);
            });

            return { previousFeed, previousPost };
        },
        onError: (err, postId, context) => {
            context?.previousFeed?.forEach(([queryKey, data]) => {
                queryClient.setQueryData(queryKey, data);
            });
            context?.previousPost?.forEach(([queryKey, data]) => {
                queryClient.setQueryData(queryKey, data);
            });
        },
        onSuccess: (data, postId) => {
            // Reconcile with the server's real count in case of concurrent votes.
            queryClient.setQueriesData({ queryKey: ["feed"] }, (oldData: any) => {
                if (!oldData?.pages) return oldData;
                const newPages = oldData.pages.map((page: any) => ({
                    ...page,
                    posts: page.posts.map((p: Post) => applyToPost(p, postId, data.liked, data.helpfulCount)),
                }));
                return { ...oldData, pages: newPages };
            });
            queryClient.setQueriesData({ queryKey: ["post", postId] }, (oldPost: any) => {
                if (!oldPost) return oldPost;
                return applyToPost(oldPost, postId, data.liked, data.helpfulCount);
            });
        },
    });

    return {
        toggleHelpful: mutation.mutate,
        isToggling: mutation.isPending,
        error: mutation.error,
    };
}
