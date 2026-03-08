"use client";

import { useState } from "react";
import { useLocation } from "@/hooks/useLocation";
import { useFeed } from "@/hooks/useFeed";
import { RadiusSelector } from "@/components/RadiusSelector";
import { MapView } from "@/components/MapView";
import { MapPinOff, Compass, List } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function MapPage() {
    const [radius, setRadius] = useState("5");
    const { lat, lng, error: locationError, loading: locationLoading } = useLocation();
    const { posts, isLoading: feedLoading, error: feedError } = useFeed(lat, lng, parseInt(radius));

    return (
        <div className="container max-w-5xl mx-auto px-4 py-8 animate-in fade-in duration-500">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Activity Map</h1>
                    <p className="text-muted-foreground mt-1">Explore posts and clusters around you visually</p>
                </div>

                <div className="flex flex-wrap items-center gap-3">
                    <div className="bg-muted/40 p-2 rounded-xl border border-muted/60 backdrop-blur-sm self-start sm:self-auto shadow-sm">
                        <RadiusSelector value={radius} onChange={setRadius} />
                    </div>
                    <Button variant="outline" asChild className="gap-2 shadow-sm">
                        <Link href="/feed">
                            <List className="h-4 w-4" />
                            <span className="hidden sm:inline">List View</span>
                        </Link>
                    </Button>
                </div>
            </div>

            <div className="rounded-xl overflow-hidden border shadow-sm bg-card">
                {locationError ? (
                    <div className="flex flex-col items-center justify-center py-20 text-center bg-card">
                        <div className="h-16 w-16 rounded-full bg-destructive/10 text-destructive flex items-center justify-center mb-4">
                            <MapPinOff className="h-8 w-8" />
                        </div>
                        <h3 className="text-xl font-semibold mb-2">Location Required</h3>
                        <p className="text-muted-foreground max-w-sm mb-6">
                            {locationError}
                        </p>
                        <Button variant="outline" className="gap-2" onClick={() => window.location.reload()}>
                            <Compass className="h-4 w-4" />
                            Try Again
                        </Button>
                    </div>
                ) : feedError ? (
                    <div className="flex flex-col items-center justify-center py-20 text-center bg-destructive/5">
                        <h3 className="text-lg font-semibold text-destructive mb-2">Something went wrong</h3>
                        <p className="text-sm text-destructive/80 mb-4 text-balance">
                            We couldn't load the map data. Please try again.
                        </p>
                        <Button variant="outline" onClick={() => window.location.reload()} size="sm">
                            Refresh Map
                        </Button>
                    </div>
                ) : locationLoading || feedLoading ? (
                    <div className="w-full h-[500px] md:h-[600px] bg-muted/20 animate-pulse flex items-center justify-center">
                        <div className="flex flex-col items-center opacity-50">
                            <Compass className="h-8 w-8 animate-spin mb-3 text-primary" />
                            <span className="font-medium">Loading map tiles...</span>
                        </div>
                    </div>
                ) : lat && lng ? (
                    <div className="animate-in zoom-in-95 duration-700">
                        <MapView lat={lat} lng={lng} posts={posts} zoom={radius === "1" ? 14 : radius === "3" ? 13 : radius === "5" ? 12 : 11} />
                    </div>
                ) : null}
            </div>
        </div>
    );
}
