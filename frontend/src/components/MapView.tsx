import { useState } from "react";
import Map, { Marker, Popup } from "react-map-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { Post } from "@/types/post";
import { formatDistance } from "@/lib/utils";
import { MapPin, AlertTriangle, Navigation } from "lucide-react";
import Link from "next/link";
import { Button } from "./ui/button";

interface MapViewProps {
    lat: number;
    lng: number;
    posts: Post[];
    zoom?: number;
}

export function MapView({ lat, lng, posts, zoom = 13 }: MapViewProps) {
    const [selectedPost, setSelectedPost] = useState<Post | null>(null);

    const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_TOKEN;

    // Render a fallback if the token isn't provided (e.g., prototype environment)
    if (!mapboxToken) {
        return (
            <div className="w-full h-full min-h-[400px] bg-muted/30 rounded-xl border flex flex-col items-center justify-center p-6 text-center shadow-sm">
                <Navigation className="h-10 w-10 text-muted-foreground mb-4 opacity-50" />
                <h3 className="font-semibold text-lg mb-2">Map Unavailable</h3>
                <p className="text-sm text-muted-foreground max-w-sm mb-4">
                    A Mapbox access token is required to render the map visualizer. Please add NEXT_PUBLIC_MAPBOX_TOKEN to your environment variables.
                </p>
            </div>
        );
    }

    return (
        <div className="w-full h-[500px] md:h-[600px] rounded-xl overflow-hidden shadow-md border relative">
            <Map
                mapboxAccessToken={mapboxToken}
                initialViewState={{
                    longitude: lng,
                    latitude: lat,
                    zoom: zoom,
                }}
                mapStyle="mapbox://styles/mapbox/light-v11"
                style={{ width: "100%", height: "100%" }}
            >
                {/* User's Current Location Marker */}
                <Marker longitude={lng} latitude={lat} anchor="center">
                    <div className="relative flex items-center justify-center h-8 w-8">
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary/40 opacity-75"></span>
                        <div className="relative inline-flex rounded-full h-4 w-4 bg-primary border-2 border-white shadow-sm"></div>
                    </div>
                </Marker>

                {/* Nearby Posts Markers */}
                {posts.map((post) => {
                    // Validate post coordinates exist in production, fallback to 0 if missing.
                    const postLng = post.lng ?? lng;
                    const postLat = post.lat ?? lat;
                    const isHelp = post.type === "help";

                    return (
                        <Marker
                            key={post.id}
                            longitude={postLng}
                            latitude={postLat}
                            anchor="bottom"
                            onClick={(e: any) => {
                                e.originalEvent.stopPropagation();
                                setSelectedPost(post);
                            }}
                            style={{ cursor: "pointer" }}
                        >
                            <div
                                className={`flex items-center justify-center w-8 h-8 rounded-full shadow-lg border-2 border-white transition-transform hover:scale-110 ${isHelp ? "bg-yellow-500 text-white" : "bg-purple-500 text-white"
                                    }`}
                            >
                                {isHelp ? <AlertTriangle size={16} /> : <MapPin size={16} />}
                            </div>
                        </Marker>
                    );
                })}

                {/* Popup for Selected Marker */}
                {selectedPost && (
                    <Popup
                        longitude={selectedPost.lng ?? lng}
                        latitude={selectedPost.lat ?? lat}
                        anchor="top"
                        onClose={() => setSelectedPost(null)}
                        closeOnClick={false}
                        className="z-50 rounded-xl"
                        maxWidth="300px"
                    >
                        <div className="p-1 min-w-[200px]">
                            <div className="flex items-center justify-between mb-2">
                                <span
                                    className={`text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full ${selectedPost.type === "help"
                                        ? "bg-yellow-100 text-yellow-800"
                                        : "bg-purple-100 text-purple-800"
                                        }`}
                                >
                                    {selectedPost.type}
                                </span>
                                <span className="text-xs font-semibold text-muted-foreground">
                                    {formatDistance(selectedPost.distance)}
                                </span>
                            </div>

                            <h4 className="font-semibold text-sm line-clamp-2 leading-tight mb-2">
                                {selectedPost.title}
                            </h4>

                            <p className="text-xs text-muted-foreground line-clamp-2 mb-3">
                                {selectedPost.description}
                            </p>

                            <Button size="sm" className="w-full text-xs h-8" asChild>
                                <Link href={`/post/${selectedPost.id}`}>View Details</Link>
                            </Button>
                        </div>
                    </Popup>
                )}
            </Map>
        </div>
    );
}
