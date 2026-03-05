import { Post } from "@/types/post";
import { formatDistanceToNow } from "date-fns";
import { MessageSquare, MapPin, MoreHorizontal } from "lucide-react";

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

export function PostCard({ post, onClick }: PostCardProps) {
    const isHelp = post.type === "help";
    const formattedDate = formatDistanceToNow(new Date(post.createdAt), {
        addSuffix: true,
    });

    return (
        <Card
            className="overflow-hidden bg-card text-card-foreground shadow-sm hover:shadow-md transition-all cursor-pointer border-muted/50"
            onClick={() => onClick && onClick(post.id)}
        >
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
                            <span>{post.distance} away</span>
                            <span className="mx-1">•</span>
                            <span>{formattedDate}</span>
                        </div>
                    </div>
                </div>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground -mr-2">
                    <MoreHorizontal className="h-4 w-4" />
                </Button>
            </CardHeader>

            <CardContent className="p-4 pt-2">
                <div className="flex justify-between items-start mb-2 gap-2">
                    <h3 className="font-semibold text-base leading-tight">
                        {post.title}
                    </h3>
                    <Badge
                        variant="outline"
                        className={`capitalize whitespace-nowrap px-2 py-0.5 text-xs font-semibold
              ${isHelp
                                ? "bg-amber-100 text-amber-800 border-amber-200 dark:bg-amber-900/30 dark:text-amber-400 dark:border-amber-800"
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

            <CardFooter className="p-4 pt-0 border-t flex items-center justify-between mt-4">
                <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-foreground -ml-2 h-8 px-2 gap-1.5">
                    <MessageSquare className="h-4 w-4" />
                    <span className="text-xs font-medium">{post.commentCount}</span>
                </Button>
                <Button size="sm" variant={isHelp ? "secondary" : "default"} className="h-8 text-xs px-4 rounded-full">
                    View Details
                </Button>
            </CardFooter>
        </Card>
    );
}
