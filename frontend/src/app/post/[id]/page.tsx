"use client";

import { usePost } from "@/hooks/usePost";
import { useLocation } from "@/hooks/useLocation";
import { useParams, useRouter } from "next/navigation";
import { formatDistance, formatTimeAgo } from "@/lib/utils";
import { CommentList } from "@/components/CommentList";
import { ArrowLeft, MapPin, Compass, AlertCircle, CalendarClock } from "lucide-react";
import Link from "next/link";
import { format } from "date-fns";

import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

export default function PostDetailPage() {
    const params = useParams();
    const id = params.id as string;
    const router = useRouter();

    // Technically we already have lat/lng but usePost doesn't necessarily need it to fetch down if the backend handles distance or if we compute it. 
    // Let's assume the API already returns the distance from coordinate headers/params on global state or similar, per specs.
    const { post, isLoading, error } = usePost(id);

    if (isLoading) {
        return (
            <div className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in">
                <Button variant="ghost" size="sm" asChild className="mb-6 -ml-3 text-muted-foreground hover:text-foreground">
                    <Link href="/feed" className="flex items-center gap-2">
                        <ArrowLeft className="h-4 w-4" />
                        Back to Feed
                    </Link>
                </Button>
                <div className="animate-pulse space-y-6">
                    <div className="h-8 bg-muted rounded w-3/4 mb-4"></div>
                    <div className="flex gap-4 items-center">
                        <div className="w-12 h-12 bg-muted rounded-full"></div>
                        <div className="space-y-2 flex-1">
                            <div className="h-4 bg-muted rounded w-1/4"></div>
                            <div className="h-3 bg-muted rounded w-1/3"></div>
                        </div>
                    </div>
                    <div className="h-24 bg-muted rounded w-full"></div>
                </div>
            </div>
        );
    }

    if (error || !post) {
        return (
            <div className="container max-w-2xl mx-auto px-4 py-8">
                <Button variant="ghost" size="sm" asChild className="mb-6 -ml-3 text-muted-foreground hover:text-foreground">
                    <Link href="/feed" className="flex items-center gap-2">
                        <ArrowLeft className="h-4 w-4" />
                        Back to Feed
                    </Link>
                </Button>
                <div className="flex flex-col items-center justify-center py-20 bg-destructive/5 rounded-xl border border-destructive/20 text-center">
                    <AlertCircle className="h-10 w-10 text-destructive mb-4" />
                    <h3 className="text-xl font-semibold mb-2">Post not found</h3>
                    <p className="text-muted-foreground mb-6">This post may have been removed or doesn't exist.</p>
                    <Button onClick={() => router.push("/feed")}>Return to Feed</Button>
                </div>
            </div>
        );
    }

    const isHelp = post.type === "help";

    return (
        <div className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <Button variant="ghost" size="sm" asChild className="mb-6 -ml-3 text-muted-foreground hover:text-foreground transition-colors hover:bg-muted/50">
                <Link href="/feed" className="flex items-center gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Back to Feed
                </Link>
            </Button>

            <Card className="border-muted shadow-md overflow-hidden bg-card">
                <div className="bg-muted/20 p-6 border-b border-muted">
                    <div className="flex justify-between items-start mb-4 gap-4">
                        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight leading-tight">{post.title}</h1>
                        <Badge
                            variant="outline"
                            className={`capitalize px-3 py-1 text-xs sm:text-sm font-semibold shadow-sm shrink-0
                ${isHelp
                                    ? "bg-yellow-100 text-yellow-800 border-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400 dark:border-yellow-800"
                                    : "bg-purple-100 text-purple-800 border-purple-200 dark:bg-purple-900/30 dark:text-purple-400 dark:border-purple-800"
                                }`}
                        >
                            {post.type}
                        </Badge>
                    </div>

                    <div className="flex items-center gap-3">
                        <Avatar className="h-12 w-12 border border-muted ring-2 ring-background ring-offset-1 ring-offset-background shadow-sm">
                            <AvatarImage src={post.authorAvatar} alt={post.author} />
                            <AvatarFallback className="bg-primary/10 text-primary font-bold text-lg">
                                {post.author.substring(0, 2).toUpperCase()}
                            </AvatarFallback>
                        </Avatar>
                        <div className="flex flex-col">
                            <span className="font-semibold text-base">{post.author}</span>
                            <div className="flex flex-wrap items-center text-sm text-muted-foreground mt-0.5 gap-x-2 gap-y-1">
                                <div className="flex items-center text-primary/80 font-medium">
                                    <Compass className="h-3.5 w-3.5 mr-1" />
                                    <span>{formatDistance(post.distance)}</span>
                                </div>
                                <span className="hidden sm:inline">•</span>
                                <span>{formatTimeAgo(post.createdAt)}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <CardContent className="p-6">
                    <div className="prose prose-sm sm:prose-base dark:prose-invert max-w-none">
                        <p className="whitespace-pre-wrap leading-relaxed text-card-foreground/90">{post.description}</p>
                    </div>

                    {post.meetupTime && (
                        <div className="mt-6 flex items-center gap-3 bg-primary/5 border border-primary/20 text-primary p-4 rounded-xl shadow-sm">
                            <div className="bg-primary/10 p-2 rounded-lg">
                                <CalendarClock className="h-5 w-5" />
                            </div>
                            <div className="flex flex-col">
                                <span className="text-xs font-semibold uppercase tracking-wider text-primary/80">Meetup Scheduled</span>
                                <span className="font-medium text-sm sm:text-base">
                                    {format(new Date(post.meetupTime), "EEEE, MMMM do 'at' h:mm a")}
                                </span>
                            </div>
                        </div>
                    )}

                    <div className="mt-8">
                        <CommentList postId={post.id} />
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
