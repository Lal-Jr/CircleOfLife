"use client";

import { useState } from "react";
import { useLocation } from "@/hooks/useLocation";
import { useFeed } from "@/hooks/useFeed";
import { RadiusSelector } from "@/components/RadiusSelector";
import { PostCard } from "@/components/PostCard";
import { MapPinOff, Compass, PlusCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

function FeedSkeleton() {
    return (
        <div className="grid gap-4">
            {[1, 2, 3].map((i) => (
                <div key={i} className="animate-pulse bg-card border rounded-xl p-4 shadow-sm h-40 flex flex-col justify-between">
                    <div className="flex gap-3">
                        <div className="w-10 h-10 bg-muted rounded-full"></div>
                        <div className="flex-1 space-y-2 py-1">
                            <div className="h-4 bg-muted rounded w-2/4"></div>
                            <div className="h-3 bg-muted rounded w-1/4"></div>
                        </div>
                    </div>
                    <div className="space-y-2 mt-4">
                        <div className="h-4 bg-muted rounded w-3/4"></div>
                        <div className="h-3 bg-muted rounded w-full"></div>
                        <div className="h-3 bg-muted rounded w-5/6"></div>
                    </div>
                </div>
            ))}
        </div>
    );
}

export default function FeedPage() {
    const [radius, setRadius] = useState("5");
    const { lat, lng, error: locationError, loading: locationLoading } = useLocation();
    const { posts, isLoading: feedLoading, error: feedError } = useFeed(lat, lng, parseInt(radius));

    const isLoading = locationLoading || feedLoading;

    return (
        <div className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in duration-500">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Nearby Feed</h1>
                    <p className="text-muted-foreground mt-1">See what's happening around you</p>
                </div>

                <div className="flex items-center gap-4 bg-muted/40 p-2 rounded-xl border border-muted/60 backdrop-blur-sm self-start sm:self-auto shadow-sm">
                    <RadiusSelector value={radius} onChange={setRadius} />
                </div>
            </div>

            <div className="space-y-4">
                {locationError ? (
                    <div className="flex flex-col items-center justify-center py-20 text-center bg-card rounded-xl border-dashed border-2 shadow-sm border-destructive/20 animate-in zoom-in-95">
                        <div className="h-16 w-16 rounded-full bg-destructive/10 text-destructive flex items-center justify-center mb-4">
                            <MapPinOff className="h-8 w-8" />
                        </div>
                        <h3 className="text-xl font-semibold mb-2">Location Required</h3>
                        <p className="text-muted-foreground max-w-[280px] md:max-w-sm mb-6">
                            {locationError}
                        </p>
                        <Button variant="outline" className="gap-2" onClick={() => window.location.reload()}>
                            <Compass className="h-4 w-4" />
                            Try Again
                        </Button>
                    </div>
                ) : feedError ? (
                    <div className="flex flex-col items-center justify-center py-16 text-center bg-destructive/5 rounded-xl border border-destructive/20">
                        <h3 className="text-lg font-semibold text-destructive mb-2">Something went wrong</h3>
                        <p className="text-sm text-destructive/80 mb-4 text-balance">
                            We couldn't load the nearby posts. Please try again later.
                        </p>
                        <Button variant="outline" onClick={() => window.location.reload()} size="sm">
                            Refresh Feed
                        </Button>
                    </div>
                ) : isLoading ? (
                    <FeedSkeleton />
                ) : posts.length > 0 ? (
                    <div className="grid gap-4">
                        {posts.map((post) => (
                            <div key={post.id} className="animate-in slide-in-from-bottom-4 duration-500 fade-in" style={{ animationDelay: `${Math.random() * 200}ms`, animationFillMode: "both" }}>
                                <PostCard post={post} />
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 text-center bg-card rounded-xl border shadow-sm">
                        <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center mb-4 text-4xl">
                            🌍
                        </div>
                        <h3 className="text-xl font-semibold mb-2">No posts found</h3>
                        <p className="text-muted-foreground max-w-sm mb-6">
                            There aren't any posts within a {radius}km radius right now. Be the first to start a conversation!
                        </p>
                        <Button asChild className="gap-2 shadow-sm">
                            <Link href="/create">
                                <PlusCircle className="h-4 w-4" />
                                Create Post
                            </Link>
                        </Button>
                    </div>
                )}
            </div>
        </div>
    );
}
