import type {APIRoute} from "astro";
import {userData, userPerms} from "../userStore.ts";

export const GET: APIRoute = async ({ cookies, redirect}) => {
    cookies.delete('token');

    userData.set({
        Username: "",
        UserID: 0
    })
    userPerms.set({
        DefaultRank: 0,
        SignUp: 0,
        CreateUsers: 0,
        DeleteUsers: 0,
        EditUsers: 0,
        ViewPosts: 0,
        CreatePosts: 0,
        DeleteOwnPosts: 0,
        DeleteOthersPosts: 0,
        EditOthersPosts: 0,
        CreateTags: 0,
        EditTags: 0,
        DeleteTags: 0
    })

    return redirect("/");
}
