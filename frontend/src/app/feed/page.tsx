"use client";

import { useState } from "react";
import { useFeed } from "@/hooks/useFeed";
import { RadiusSelector } from "@/components/RadiusSelector";
import { PostCard } from "@/components/PostCard";
import { Loader2, PlusCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function FeedPage() {
    const [radius, setRadius] = useState("5");
    const { posts, isLoading } = useFeed(parseInt(radius));

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
                {isLoading ? (
                    <div className="flex flex-col items-center justify-center py-20 text-muted-foreground bg-card rounded-xl border shadow-sm">
                        <Loader2 className="h-10 w-10 animate-spin text-primary mb-4" />
                        <p className="font-medium">Discovering nearby posts...</p>
                    </div>
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
                        <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center mb-4">
                            <span className="text-2xl">🌍</span>
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
