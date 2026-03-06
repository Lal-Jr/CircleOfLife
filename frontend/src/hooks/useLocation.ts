import { useState, useEffect } from "react";

interface LocationState {
    lat: number | null;
    lng: number | null;
    error: string | null;
    loading: boolean;
}

export function useLocation() {
    const [state, setState] = useState<LocationState>({
        lat: null,
        lng: null,
        error: null,
        loading: true,
    });

    useEffect(() => {
        if (!("geolocation" in navigator)) {
            setState((s) => ({
                ...s,
                error: "Geolocation is not supported by your browser.",
                loading: false,
            }));
            return;
        }

        navigator.geolocation.getCurrentPosition(
            (position) => {
                setState({
                    lat: position.coords.latitude,
                    lng: position.coords.longitude,
                    error: null,
                    loading: false,
                });
            },
            (error) => {
                let errorMsg = "Unable to fetch location. Please enable location services.";
                if (error.code === error.PERMISSION_DENIED) {
                    errorMsg = "Location access denied. Please allow location to view nearby posts.";
                }
                setState((s) => ({
                    ...s,
                    error: errorMsg,
                    loading: false,
                }));
            },
            { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
        );
    }, []);

    return state;
}
