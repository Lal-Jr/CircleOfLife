import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { User } from "@/types/user";

export function useCurrentUser() {
    const { data, isLoading, error } = useQuery({
        queryKey: ["me"],
        queryFn: async () => {
            const response = await api.get("/users/me");
            return response.data.data as User;
        },
    });

    return {
        user: data,
        isLoading,
        error,
    };
}
