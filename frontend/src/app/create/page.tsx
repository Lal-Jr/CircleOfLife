"use client";

import { useState } from "react";
import { useLocation } from "@/hooks/useLocation";
import { useCreatePost } from "@/hooks/useCreatePost";
import { ArrowLeft, Compass, Send, Calendar, Clock } from "lucide-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";

export default function CreatePostPage() {
    const { lat, lng, error: locationError, loading: locationLoading } = useLocation();
    const { createPost, isCreating } = useCreatePost();

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [type, setType] = useState<"help" | "meetup">("help");
    const [date, setDate] = useState("");
    const [time, setTime] = useState("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!lat || !lng) return;

        let meetupTime: string | undefined = undefined;
        if (type === "meetup" && date && time) {
            meetupTime = new Date(`${date}T${time}`).toISOString();
        }

        createPost({
            title,
            description,
            type,
            lat,
            lng,
            meetupTime,
        });
    };

    const isLocationReady = lat !== null && lng !== null;
    const isFormValid = title.trim() && description.trim() && (type === "help" || (type === "meetup" && date && time));

    return (
        <div className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in zoom-in-95 duration-500">
            <Button variant="ghost" size="sm" asChild className="mb-6 -ml-3 text-muted-foreground hover:text-foreground">
                <Link href="/feed" className="flex items-center gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Back to Feed
                </Link>
            </Button>

            <div className="mb-8">
                <h1 className="text-3xl font-bold tracking-tight mb-2">Create Post</h1>
                <p className="text-muted-foreground">Share something with your local community</p>
            </div>

            <Card className="shadow-lg shadow-primary/5 border-muted">
                <CardHeader className="bg-muted/30 border-b border-muted/50 pb-6 rounded-t-xl mb-4">
                    <CardTitle className="text-xl">Post Details</CardTitle>
                    <CardDescription>Fill out the information below. Your current location will be attached securely.</CardDescription>
                </CardHeader>
                <CardContent>
                    {locationLoading ? (
                        <div className="flex items-center justify-center p-8 bg-muted/20 rounded-xl">
                            <Compass className="h-6 w-6 animate-spin text-primary mr-3" />
                            <span className="font-medium text-muted-foreground">Acquiring your location...</span>
                        </div>
                    ) : locationError ? (
                        <div className="p-4 bg-destructive/10 text-destructive rounded-xl border border-destructive/20 flex flex-col items-start gap-2 mb-6 text-sm">
                            <div className="flex items-center font-semibold">
                                <Compass className="h-4 w-4 mr-2" />
                                Location Access Denied
                            </div>
                            <span>We need your location to attach this post correctly to the map. Please enable location permissions in your browser.</span>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit} className="space-y-6">
                            <div className="space-y-2">
                                <Label htmlFor="type" className="font-semibold">Post Type</Label>
                                <Select value={type} onValueChange={(val: "help" | "meetup") => setType(val)}>
                                    <SelectTrigger id="type" className="bg-background focus:ring-primary/50 h-11 transition-all">
                                        <SelectValue placeholder="Request Help or Organize a Meetup" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="help">Request Help</SelectItem>
                                        <SelectItem value="meetup">Organize Meetup</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="title" className="font-semibold">Title</Label>
                                <Input
                                    id="title"
                                    placeholder="E.g. Need help moving desk"
                                    value={title}
                                    onChange={(e) => setTitle(e.target.value)}
                                    className="bg-background focus-visible:ring-primary/50 h-11 transition-all"
                                    maxLength={60}
                                    required
                                />
                                <p className="text-xs text-muted-foreground text-right">{title.length}/60</p>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="description" className="font-semibold">Description</Label>
                                <Textarea
                                    id="description"
                                    placeholder="Add some details about what you need or what the meetup is about..."
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                    className="bg-background focus-visible:ring-primary/50 min-h-[120px] resize-none transition-all p-3"
                                    maxLength={500}
                                    required
                                />
                                <p className="text-xs text-muted-foreground text-right">{description.length}/500</p>
                            </div>

                            {type === "meetup" && (
                                <div className="grid grid-cols-2 gap-4 animate-in slide-in-from-top-4 fade-in duration-300 bg-muted/20 p-4 rounded-xl border border-muted/50">
                                    <div className="space-y-2">
                                        <Label htmlFor="date" className="font-semibold text-sm flex items-center gap-1.5"><Calendar className="w-3.5 h-3.5" /> Date</Label>
                                        <Input
                                            id="date"
                                            type="date"
                                            value={date}
                                            onChange={(e) => setDate(e.target.value)}
                                            className="bg-background h-10 transition-all focus-visible:ring-primary/50"
                                            required={type === "meetup"}
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="time" className="font-semibold text-sm flex items-center gap-1.5"><Clock className="w-3.5 h-3.5" /> Time</Label>
                                        <Input
                                            id="time"
                                            type="time"
                                            value={time}
                                            onChange={(e) => setTime(e.target.value)}
                                            className="bg-background h-10 transition-all focus-visible:ring-primary/50"
                                            required={type === "meetup"}
                                        />
                                    </div>
                                </div>
                            )}

                            <Button
                                type="submit"
                                className="w-full h-12 text-base font-semibold mt-4 shadow-md transition-all gap-2"
                                disabled={!isLocationReady || isCreating || !isFormValid}
                            >
                                {isCreating ? "Posting..." : "Create Post"}
                                {!isCreating && <Send className="h-4 w-4" />}
                            </Button>
                        </form>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
