import { useQuery } from "@tanstack/react-query";
import { Post } from "@/types/post";

// Mock data generator for initial UI validation
const mockPosts: Post[] = [
    {
        id: "1",
        title: "Need help moving furniture",
        description: "Looking for someone nearby who can help move a heavy table from the ground floor to the 2nd floor this weekend.",
        author: "Rahul M.",
        distance: "300m",
        type: "help",
        commentCount: 2,
        createdAt: new Date().toISOString(),
    },
    {
        id: "2",
        title: "Weekend cycling group",
        description: "Anyone interested in a 20km cycling loop around the neighborhood early Sunday morning?",
        author: "Sarah J.",
        distance: "1.2km",
        type: "meetup",
        commentCount: 5,
        createdAt: new Date(Date.now() - 3600000).toISOString(), // 1 hour ago
    },
    {
        id: "3",
        title: "Lost Golden Retriever",
        description: "Has anyone seen a very friendly golden retriever near main street? Wearing a red collar.",
        author: "Amit P.",
        distance: "500m",
        type: "help",
        commentCount: 12,
        createdAt: new Date(Date.now() - 7200000).toISOString(), // 2 hours ago
    },
    {
        id: "4",
        title: "Local coffee shop remote work session",
        description: "Going to be working from The Daily Grind today from 1pm to 5pm if any other remote workers want to join.",
        author: "Priya S.",
        distance: "800m",
        type: "meetup",
        commentCount: 0,
        createdAt: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
    }
];

export function useFeed(radiusMs: number = 5) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["posts", radiusMs],
        queryFn: async () => {
            // Future backend integration
            // const response = await api.get(`/posts?radius=${radiusMs}`);
            // return response.data;

            console.log(`Fetching mock posts for radius ${radiusMs}km`);
            // Simulate network request
            await new Promise(resolve => setTimeout(resolve, 800));
            return mockPosts;
        },
    });

    return {
        posts: data || [],
        isLoading,
        error,
    };
}
