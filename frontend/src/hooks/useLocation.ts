import { useState, useEffect } from "react";
import { useCurrentUser } from "@/hooks/useCurrentUser";
import { DEMO_EMAIL, DEMO_LOCATION } from "@/lib/demo";

interface LocationState {
    lat: number | null;
    lng: number | null;
    error: string | null;
    loading: boolean;
}

export function useLocation(): LocationState {
    const { user, isLoading: userLoading } = useCurrentUser();
    const isDemo = user?.email === DEMO_EMAIL;

    const [state, setState] = useState<LocationState>({
        lat: null,
        lng: null,
        error: null,
        loading: true,
    });

    useEffect(() => {
        // Wait until we know who is signed in, and skip the permission
        // prompt entirely for the demo account.
        if (userLoading || isDemo) return;

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
    }, [userLoading, isDemo]);

    if (isDemo) {
        return { ...DEMO_LOCATION, error: null, loading: false };
    }
    if (userLoading) {
        return { ...state, loading: true };
    }
    return state;
}
