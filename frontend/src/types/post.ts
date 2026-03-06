export interface Post {
    id: string;
    title: string;
    description: string;
    author: string;
    distance: number; // in meters returned from backend
    type: "help" | "meetup";
    commentCount: number;
    createdAt: string;
    authorAvatar?: string;
    lat?: number;
    lng?: number;
    meetupTime?: string;
}
