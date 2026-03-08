"use client";

import { MapPinOff, AlertTriangle, RefreshCcw } from "lucide-react";
import { Button } from "./ui/button";

interface ErrorStateProps {
    title?: string;
    message?: string;
    onRetry?: () => void;
    type?: "location" | "general" | "network";
}

export function ErrorState({
    title = "Something went wrong",
    message = "An unexpected error occurred while loading this section.",
    onRetry,
    type = "general"
}: ErrorStateProps) {

    const renderIcon = () => {
        switch (type) {
            case "location":
                return <MapPinOff className="h-8 w-8" />;
            case "network":
                return <AlertTriangle className="h-8 w-8" />;
            default:
                return <AlertTriangle className="h-8 w-8" />;
        }
    };

    return (
        <div className="flex flex-col items-center justify-center p-8 bg-destructive/5 rounded-xl border border-destructive/20 text-center animate-in zoom-in-95 duration-300 w-full">
            <div className={`h-16 w-16 rounded-full flex items-center justify-center mb-4 ${type === 'location' ? 'bg-orange-100 text-orange-600' : 'bg-destructive/10 text-destructive'}`}>
                {renderIcon()}
            </div>
            <h3 className="text-xl font-semibold mb-2">{title}</h3>
            <p className="text-muted-foreground text-sm max-w-[280px] md:max-w-sm mb-6 text-balance">
                {message}
            </p>
            {onRetry && (
                <Button variant="outline" className="gap-2 shadow-sm" onClick={onRetry}>
                    <RefreshCcw className="h-4 w-4" />
                    Retry
                </Button>
            )}
        </div>
    );
}
