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
        onSuccess: () => {
            // Invalidate posts list so feed updates
            queryClient.invalidateQueries({ queryKey: ["posts"] });
            router.push("/feed");
        },
    });

    return {
        createPost: mutation.mutate,
        isCreating: mutation.isPending,
        error: mutation.error,
    };
}
