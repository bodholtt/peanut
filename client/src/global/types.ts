
export type Post = {
    id: number,
    tags: Tag[],
    image_path: string,
    created_at: string,
    author_id: number,
    source: string,
    previous: number,
    next: number,
}

export type PostThumb = {
    id: number,
    image_path: string
}

export type PostThumbs = {
    Thumbs: PostThumb[]
}

export type User = {
    id: number,
    username: string,
    rank: number,
    created_at: string,
}

export type Tag = {
    id: number,
    name: string,
    description: string,
}