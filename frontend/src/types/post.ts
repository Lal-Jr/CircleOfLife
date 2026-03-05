export interface Post {
    id: string;
    title: string;
    description: string;
    author: string;
    distance: string; // e.g., "300m", "1.2km"
    type: "help" | "meetup";
    commentCount: number;
    createdAt: string;
    authorAvatar?: string;
}
