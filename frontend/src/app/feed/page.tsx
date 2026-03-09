"use client";

import { PostCard } from "@/components/PostCard";
import { RadiusSelector } from "@/components/RadiusSelector";
import { Button } from "@/components/ui/button";
import { useFeed } from "@/hooks/useFeed";
import { useLocation } from "@/hooks/useLocation";
import { useRealtimeFeed } from "@/hooks/useRealtimeFeed";
import { useQueryClient } from "@tanstack/react-query";
import { Compass, MapPinOff, PlusCircle, RefreshCw } from "lucide-react";
import Link from "next/link";
import { useState, useEffect } from "react";
import { useInView } from "react-intersection-observer";

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

    const {
        posts,
        isLoading: feedLoading,
        error: feedError,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage
    } = useFeed(lat, lng, parseInt(radius));

    const queryClient = useQueryClient();
    const { hasNewPosts, clearNewPosts } = useRealtimeFeed();
    const { ref: observerRef, inView } = useInView();

    // Pull to refresh UX states
    const [pullDistance, setPullDistance] = useState(0);
    const [startY, setStartY] = useState(0);
    const [isRefreshing, setIsRefreshing] = useState(false);

    useEffect(() => {
        if (inView && hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
        }
    }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

    const handleRefresh = async () => {
        setIsRefreshing(true);
        await queryClient.invalidateQueries({ queryKey: ["feed"] });
        clearNewPosts();
        setIsRefreshing(false);
    };

    const handleTouchStart = (e: React.TouchEvent) => {
        if (window.scrollY === 0) setStartY(e.touches[0].clientY);
    };

    const handleTouchMove = (e: React.TouchEvent) => {
        if (startY > 0) {
            const distance = e.touches[0].clientY - startY;
            if (distance > 0 && window.scrollY === 0) {
                setPullDistance(Math.min(distance * 0.4, 80));
            }
        }
    };

    const handleTouchEnd = async () => {
        if (pullDistance > 60 && !isRefreshing) {
            await handleRefresh();
        }
        setStartY(0);
        setPullDistance(0);
    };

    const isLoading = locationLoading || feedLoading;

    return (
        <div
            className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in duration-500 relative transition-transform"
            style={{ transform: `translateY(${isRefreshing ? 60 : pullDistance}px)` }}
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
        >

            {/* Pull to refresh visual UI mounted behind absolute bounds */}
            {(pullDistance > 0 || isRefreshing) && (
                <div className="absolute top-0 left-0 right-0 h-16 -mt-16 flex items-center justify-center w-full">
                    <RefreshCw className={`h-6 w-6 text-primary transition-all duration-300 ${isRefreshing ? 'animate-spin opacity-100' : 'opacity-70'}`} style={{ transform: `rotate(${pullDistance * 3}deg)` }} />
                </div>
            )}

            {/* Real-time Refresh Banner */}
            {hasNewPosts && (
                <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50 animate-in slide-in-from-top-4 fade-in duration-300">
                    <Button
                        onClick={handleRefresh}
                        className="rounded-full shadow-lg shadow-primary/20 gap-2 border border-primary/20 backdrop-blur-md"
                    >
                        <RefreshCw className="h-4 w-4" />
                        New posts nearby! Refresh
                    </Button>
                </div>
            )}

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
                    <>
                        <div className="grid gap-4">
                            {posts.map((post) => (
                                <div key={post.id} className="animate-in slide-in-from-bottom-4 duration-500 fade-in">
                                    <PostCard post={post} />
                                </div>
                            ))}
                        </div>

                        {/* Infinite Scroll trigger anchor */}
                        <div ref={observerRef} className="h-14 w-full flex items-center justify-center mt-4">
                            {isFetchingNextPage && <RefreshCw className="h-6 w-6 animate-spin text-muted-foreground" />}
                        </div>
                    </>
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
