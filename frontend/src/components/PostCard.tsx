import { Post } from "@/types/post";
import { formatDistance, formatTimeAgo } from "@/lib/utils";
import { MessageSquare, MapPin, Share2, ThumbsUp } from "lucide-react";
import Link from "next/link";
import React from "react";

import {
    Card,
    CardContent,
    CardFooter,
    CardHeader,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface PostCardProps {
    post: Post;
    onClick?: (id: string) => void;
}

export const PostCard = React.memo(({ post, onClick }: PostCardProps) => {
    const isHelp = post.type === "help";
    const formattedDate = formatTimeAgo(post.createdAt);

    return (
        <Card
            className="overflow-hidden bg-card text-card-foreground shadow-sm hover:shadow-md transition-all border-muted/50 group"
        >
            <Link href={`/post/${post.id}`} onClick={(e) => {
                if (onClick) {
                    e.preventDefault();
                    onClick(post.id);
                }
            }} className="flex flex-col h-full cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-transparent rounded-lg">
                <CardHeader className="p-4 pb-2 flex flex-row items-center justify-between space-y-0">
                    <div className="flex items-center gap-3">
                        <Avatar className="h-10 w-10 border border-muted ring-1 ring-background ring-offset-1 ring-offset-background">
                            <AvatarImage src={post.authorAvatar} alt={post.author} />
                            <AvatarFallback className="bg-primary/10 text-primary font-medium">
                                {post.author.substring(0, 2).toUpperCase()}
                            </AvatarFallback>
                        </Avatar>
                        <div className="flex flex-col">
                            <span className="font-semibold text-sm leading-none">{post.author}</span>
                            <div className="flex items-center text-xs text-muted-foreground mt-1 gap-1">
                                <MapPin className="h-3 w-3" />
                                <span>{formatDistance(post.distance)} away</span>
                                <span className="mx-1">•</span>
                                <span>{formattedDate}</span>
                            </div>
                        </div>
                    </div>
                </CardHeader>

                <CardContent className="p-4 pt-2">
                    <div className="flex justify-between items-start mb-2 gap-2">
                        <h3 className="font-semibold text-base leading-tight group-hover:text-primary transition-colors">
                            {post.title}
                        </h3>
                        <Badge
                            variant="outline"
                            className={`capitalize whitespace-nowrap px-2 py-0.5 text-xs font-semibold
                ${isHelp
                                    ? "bg-yellow-100 text-yellow-800 border-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400 dark:border-yellow-800"
                                    : "bg-purple-100 text-purple-800 border-purple-200 dark:bg-purple-900/30 dark:text-purple-400 dark:border-purple-800"
                                }`}
                        >
                            {post.type}
                        </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground line-clamp-3 leading-relaxed">
                        {post.description}
                    </p>
                </CardContent>

                <CardFooter className="p-2 pt-2 border-t flex items-center justify-between mt-auto bg-muted/5 gap-1">
                    <Button variant="ghost" size="sm" className="flex-1 text-muted-foreground hover:text-primary hover:bg-primary/10 h-9 gap-1.5 rounded-md" onClick={(e) => { e.preventDefault(); e.stopPropagation(); }}>
                        <ThumbsUp className="h-4 w-4" />
                        <span className="text-xs font-medium">Helpful</span>
                    </Button>
                    <Button variant="ghost" size="sm" className="flex-1 text-muted-foreground hover:text-primary hover:bg-primary/10 h-9 gap-1.5 rounded-md" asChild>
                        <div className="flex items-center pointer-events-none">
                            <MessageSquare className="h-4 w-4" />
                            <span className="text-xs font-medium">{post.commentCount} <span className="hidden sm:inline">Comments</span></span>
                        </div>
                    </Button>
                    <Button variant="ghost" size="sm" className="flex-1 text-muted-foreground hover:text-primary hover:bg-primary/10 h-9 gap-1.5 rounded-md" onClick={(e) => { e.preventDefault(); e.stopPropagation(); }}>
                        <Share2 className="h-4 w-4" />
                        <span className="text-xs font-medium">Share</span>
                    </Button>
                </CardFooter>
            </Link>
        </Card>
    );
});
PostCard.displayName = "PostCard";
