import type {APIRoute} from "astro";

export const GET: APIRoute = async ({ cookies, redirect}) => {
    cookies.delete('token');
    return redirect("/");
}
