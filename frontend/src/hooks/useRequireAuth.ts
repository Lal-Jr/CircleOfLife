"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

// Guards a client page behind an existing JWT. Redirects to /login when
// no token is present instead of letting the page render and fail every
// request with a 401 (the axios interceptor still handles tokens that
// go bad after this check runs).
export function useRequireAuth() {
    const router = useRouter();
    const [checked, setChecked] = useState(false);

    useEffect(() => {
        const token = typeof window !== "undefined" ? localStorage.getItem("jwt_token") : null;
        if (!token) {
            router.replace("/login");
            return;
        }
        setChecked(true);
    }, [router]);

    return { isAuthChecked: checked };
}
