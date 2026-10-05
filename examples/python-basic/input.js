// A small file with constants and dead code, so the optimizer has
// something to do.

const DEBUG = false;
const API_URL = "https://api.example.com/v1";

function fetchUser(id) {
    if (DEBUG) {
        console.log("fetching user " + id);
    }
    return fetch(API_URL + "/users/" + id).then(function (r) {
        return r.json();
    });
}

function unusedHelper() {
    return "this function is never called";
}

const user = { name: "Alice", age: 30 };

console.log("loaded user " + user.name);