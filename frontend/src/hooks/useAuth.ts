import { useMutation } from "@tanstack/react-query";
// import { api } from "@/lib/api"; // For future backend integration
import { useRouter } from "next/navigation";

export function useAuth() {
    const router = useRouter();

    const loginMutation = useMutation({
        mutationFn: async (credentials: any) => {
            // Mock API call
            // const response = await api.post("/auth/login", credentials);
            // return response.data;

            console.log("Mock login for:", credentials.email);
            // Simulate network delay
            await new Promise((resolve) => setTimeout(resolve, 500));
            return { token: "mock_jwt_token_12345" };
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
            // Mock API call
            // const response = await api.post("/auth/signup", userData);
            // return response.data;

            console.log("Mock signup for:", userData.email);
            await new Promise((resolve) => setTimeout(resolve, 500));
            return { token: "mock_jwt_token_12345" };
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
