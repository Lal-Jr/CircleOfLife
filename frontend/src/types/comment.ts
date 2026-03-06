export interface Comment {
    id: string;
    postId: string;
    authorId: string;
    authorName: string;
    authorAvatar?: string;
    content: string;
    createdAt: string;
}
