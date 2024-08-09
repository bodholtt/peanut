# api

The api package handles the REST API routes for *peanut-server*. 

## Table of Contents
- [General](#General)
- [Authentication](#Authentication)
- [Posts](#Posts)
  - [Getting post count](#Getting-post-count)
  - [Getting multiple posts](#Getting-multiple-posts)
  - [Uploading a post](#Uploading-a-post)
  - [Getting a single post](#Getting-a-single-post)

## General

All routes should use the header `Content-Type: application/json`, 
with the exception of uploading files (the `/post POST` route),
which must use `Content-Type: multipart/form-data`.

All API responses will be of the following format, where `body` may be a 
string or JSON and `error` is always a string, which will be empty if a request
is successful.
```
{
  body: any
  error: ''
}
```


## Authentication

Authentication is done via JSON Web Tokens. Each JWT stores the username
and rank of a user along with the token's expiry date.
Requests should include the token in the Authorization header as such: 

`Authorization: Bearer your.token.here`

To provision a token, log in with a POST request to the `/login` API route.
The request body should be as follows:
```
{
  username: <your-username> 
  password:  <your-password>
}
```

If a login is successful, the response will be:
```
{
  body: {
    token: <token-here>,
    message: 'Login successful'
  },
  error: ''
}
```



## Posts

### Getting post count
`GET /postCount`

This route simply returns the number of posts in the instance as a raw number.
Requires no authorization.

### Getting multiple posts
`GET /posts`

By default, this route will get the newest posts up to a limit of 50, 
in newest to oldest order. 
If the `view_posts` permission is set to 0, this requires no authorization.

The response body will contain a list of post thumbnail objects, 
which each contain the post ID and the filepath of the post's thumbnail.
```
{ "Thumbs": [
    {"id": x, "image_path": y},
    {"id": x, "image_path": y},
    {"id": x, "image_path": y}
]}
```

Queries can be added to get a different output.

`GET /posts?limit=x&offset=y`

Get `x` posts (`x` limited to 50), newest to oldest, offset by `y` posts.

### Uploading a post
`POST /post`

A POST request to this route will initiate the uploading of a new post.
This must be done with a `multipart/form-data` request.

Example formdata:
```
image: file (required)
tags: string (space separated list, optional)
source: string (optional)
```

### Getting a single post
`GET /post/[id]`

This will get a single post's details by ID. Along with the details, it will also 
retrieve the IDs of the next and previous posts for pagination on the client.

[//]: # (Todo: Do I need to include this in the api route?)
If a search query was used to find the post, it will also return the query for pagination.
```
{
    "id": "1",
    "tags": [
        {
            "id": "1",
            "name": "name",
            "description": ""
        },
    ],
    "created_at": "2024-07-28 14:33:15.925 +0000 UTC",
    "image_path": "1.jpg",
    "author_id": "1",
    "source": "",
    "md5_hash": "",
    "previous": "0",
    "next": "2",
    "query": ""
}
```
