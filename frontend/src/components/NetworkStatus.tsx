"use client";

import { useEffect, useState } from "react";
import { WifiOff } from "lucide-react";

export function NetworkStatus() {
    const [isOnline, setIsOnline] = useState(true);

    useEffect(() => {
        // Safe check for SSR
        if (typeof window !== "undefined") {
            setIsOnline(navigator.onLine);

            const handleOnline = () => setIsOnline(true);
            const handleOffline = () => setIsOnline(false);

            window.addEventListener("online", handleOnline);
            window.addEventListener("offline", handleOffline);

            return () => {
                window.removeEventListener("online", handleOnline);
                window.removeEventListener("offline", handleOffline);
            };
        }
    }, []);

    if (isOnline) return null;

    return (
        <div className="fixed top-0 left-0 right-0 z-[100] bg-destructive text-destructive-foreground py-2 px-4 shadow-md text-sm font-medium flex items-center justify-center gap-2 animate-in slide-in-from-top-full duration-300">
            <WifiOff className="h-4 w-4" />
            <div className="flex flex-col sm:flex-row sm:items-center sm:gap-2 text-center">
                <span>Network connection lost.</span>
                <span className="opacity-80 text-xs sm:text-sm">Trying to reconnect...</span>
            </div>
        </div>
    );
}
