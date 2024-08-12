import { atom } from 'nanostores';

export const userData = atom({
    Username: "",
    UserID: 0
});

export const userPerms = atom({
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