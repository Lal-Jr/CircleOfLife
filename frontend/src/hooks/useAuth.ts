import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";

export function useAuth() {
    const router = useRouter();

    const loginMutation = useMutation({
        mutationFn: async (credentials: any) => {
            const response = await api.post("/auth/login", credentials);
            return response.data as { token: string };
        },
        onSuccess: (data) => {
            if (typeof window !== "undefined") {
                localStorage.setItem("jwt_token", data.token);
            }
            router.push("/feed");
        },
    });

    const signupMutation = useMutation({
        mutationFn: async (userData: any) => {
            const response = await api.post("/auth/signup", userData);
            return response.data as { token: string };
        },
        onSuccess: (data) => {
            if (typeof window !== "undefined") {
                localStorage.setItem("jwt_token", data.token);
            }
            router.push("/feed");
        },
    });

    const logout = () => {
        if (typeof window !== "undefined") {
            localStorage.removeItem("jwt_token");
        }
        router.push("/login");
    };

    return {
        login: loginMutation.mutate,
        isLoggingIn: loginMutation.isPending,
        signup: signupMutation.mutate,
        isSigningUp: signupMutation.isPending,
        logout,
    };
}
