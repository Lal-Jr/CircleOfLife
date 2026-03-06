"use client";

import { useState } from "react";
import { useComments } from "@/hooks/useComments";
import { useCreateComment } from "@/hooks/useCreateComment";
import { formatTimeAgo } from "@/lib/utils";
import { Send, Loader2 } from "lucide-react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function CommentList({ postId }: { postId: string }) {
    const { comments, isLoading, error } = useComments(postId);
    const { createComment, isCreating } = useCreateComment();
    const [content, setContent] = useState("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!content.trim()) return;

        createComment(
            { postId, content },
            {
                onSuccess: () => {
                    setContent("");
                },
            }
        );
    };

    return (
        <div className="pt-4 border-t mt-4 mb-2">
            <h4 className="text-sm font-semibold mb-3">Comments ({comments.length})</h4>

            <div className="space-y-4 mb-4 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
                {isLoading ? (
                    <div className="flex flex-col space-y-3">
                        {[1, 2].map((i) => (
                            <div key={i} className="flex gap-3 animate-pulse">
                                <div className="h-8 w-8 rounded-full bg-muted"></div>
                                <div className="flex-1 space-y-2 py-1">
                                    <div className="h-4 bg-muted rounded w-1/4"></div>
                                    <div className="h-3 bg-muted rounded w-3/4"></div>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : error ? (
                    <div className="text-sm text-destructive text-center py-4 bg-destructive/10 rounded-lg">
                        Could not load comments.
                    </div>
                ) : comments.length === 0 ? (
                    <div className="text-sm text-muted-foreground text-center py-4 bg-muted/30 rounded-lg border border-dashed border-muted">
                        No comments yet. Be the first to reply!
                    </div>
                ) : (
                    <div className="space-y-4">
                        {comments.map((comment) => (
                            <div key={comment.id} className="flex gap-3 text-sm animate-in fade-in slide-in-from-bottom-2">
                                <Avatar className="h-8 w-8 border border-muted ring-1 ring-background ring-offset-background">
                                    <AvatarImage src={comment.authorAvatar} alt={comment.authorName} />
                                    <AvatarFallback className="bg-primary/10 text-primary font-medium text-xs">
                                        {comment.authorName.substring(0, 2).toUpperCase()}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="flex max-w-[85%] flex-col">
                                    <div className="bg-muted/50 rounded-2xl rounded-tl-none px-4 py-2 border border-muted">
                                        <span className="font-semibold block text-xs mb-0.5">{comment.authorName}</span>
                                        <span className="text-foreground leading-relaxed">{comment.content}</span>
                                    </div>
                                    <span className="text-[10px] text-muted-foreground mt-1 ml-1">
                                        {formatTimeAgo(comment.createdAt)}
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <form onSubmit={handleSubmit} className="flex gap-2">
                <Input
                    placeholder="Write a comment..."
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    className="flex-1 h-10 rounded-full bg-muted/30 focus-visible:ring-primary/50 px-4"
                    disabled={isCreating}
                />
                <Button
                    type="submit"
                    size="icon"
                    className="h-10 w-10 shrink-0 rounded-full transition-all"
                    disabled={!content.trim() || isCreating}
                >
                    {isCreating ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
                </Button>
            </form>
        </div>
    );
}
