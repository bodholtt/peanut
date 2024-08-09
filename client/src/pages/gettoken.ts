import type {APIRoute} from "astro";

export const GET: APIRoute = async ({ cookies}) => {

    // This is probably not secure at all, but this was originally for uploading a new post.
    // Having the logic back here did not want to work.

    const token = cookies.get("token")
    const tokenValue = token? token.value : ""

    return new Response(JSON.stringify(tokenValue))
}