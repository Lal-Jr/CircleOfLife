import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";

export interface CreatePostPayload {
    title: string;
    description: string;
    type: "help" | "meetup";
    lat: number;
    lng: number;
    meetupTime?: string;
}

export function useCreatePost() {
    const queryClient = useQueryClient();
    const router = useRouter();

    const mutation = useMutation({
        mutationFn: async (payload: CreatePostPayload) => {
            const response = await api.post("/posts", payload);
            return response.data;
        },
        onMutate: async (newPost) => {
            await queryClient.cancelQueries({ queryKey: ["feed"] });
            const previousFeed = queryClient.getQueriesData({ queryKey: ["feed"] });

            queryClient.setQueriesData({ queryKey: ["feed"] }, (oldData: any) => {
                if (!oldData || !oldData.pages) return oldData;

                const optimisticPost = {
                    id: `temp-${Date.now()}`,
                    title: newPost.title,
                    description: newPost.description,
                    type: newPost.type,
                    author: "You",
                    authorAvatar: "",
                    distance: 0,
                    commentCount: 0,
                    createdAt: new Date().toISOString(),
                    meetupTime: newPost.meetupTime
                };

                const newPages = [...oldData.pages];
                if (newPages.length > 0) {
                    newPages[0] = {
                        ...newPages[0],
                        posts: [optimisticPost, ...newPages[0].posts],
                    };
                }

                return { ...oldData, pages: newPages };
            });

            // Redirect immediately for perceived zero-latency UX
            router.push("/feed");

            return { previousFeed };
        },
        onError: (err, newPost, context) => {
            if (context?.previousFeed) {
                context.previousFeed.forEach(([queryKey, data]) => {
                    queryClient.setQueryData(queryKey, data);
                });
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["feed"] });
        },
    });

    return {
        createPost: mutation.mutate,
        isCreating: mutation.isPending,
        error: mutation.error,
    };
}
