import type {APIRoute} from "astro";
import {userData, userPerms} from "../userStore.ts";

export const POST: APIRoute = async ({ cookies, request}) => {

    const data = await request.json();
    const username = data.username;
    const password = data.password;

    const apiroute = `${import.meta.env.API_URL}/login`

    const response = await fetch(apiroute, {
        method: 'POST',
        body: JSON.stringify({
            username: username,
            password: password
        }),
        headers: {
            "Content-Type": "application/json"
        }
    });

    const tokendata = await response.json();
    if (!response.ok) {
        return new Response(JSON.stringify(tokendata.error), {
            status: 401
        });
    }

    const token = tokendata.body.token;

    const payload = JSON.parse(atob(token.split('.')[1])); // Token payload

    cookies.set('token', token, {
        path: '/',
        expires: new Date(payload.exp as number * 1000),
        httpOnly: true,
        secure: true,
        sameSite: "strict"
    });

    // should i store the token in the userdata store? something to think about

    userData.set({
        Username: payload.username,
        UserID: payload.user_id
    })

    const userPermsData = await (await fetch(`${import.meta.env.API_URL}/user/${userData.get().UserID}/permissions`)).json();
    userPerms.set(userPermsData.body);

    return new Response(JSON.stringify({
        message: "Token set"
    }), {
        status: 200
    });


}

