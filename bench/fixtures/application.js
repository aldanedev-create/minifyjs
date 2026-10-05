// A realistic multi-module application, hand-written to exercise the
// shapes real code uses: imports, exports, async, classes, arrow
// functions, template literals, optional chaining, and destructuring.

import { EventEmitter } from "./events.js";
import { parseConfig } from "./config.js";

const DEFAULT_TIMEOUT = 30_000;
const MAX_RETRIES = 3;

class HttpClient extends EventEmitter {
    constructor(baseUrl, options = {}) {
        super();
        this.baseUrl = baseUrl;
        this.timeout = options.timeout ?? DEFAULT_TIMEOUT;
        this.headers = { "Content-Type": "application/json", ...options.headers };
    }

    async request(method, path, body) {
        const url = `${this.baseUrl}${path}`;
        const init = {
            method,
            headers: this.headers,
            signal: AbortSignal.timeout(this.timeout),
        };
        if (body !== undefined) {
            init.body = JSON.stringify(body);
        }

        for (let attempt = 0; attempt < MAX_RETRIES; attempt++) {
            try {
                const response = await fetch(url, init);
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`);
                }
                return await response.json();
            } catch (err) {
                this.emit("error", err);
                if (attempt === MAX_RETRIES - 1) {
                    throw err;
                }
                await new Promise((r) => setTimeout(r, 2 ** attempt * 100));
            }
        }
    }

    get(path) {
        return this.request("GET", path);
    }

    post(path, body) {
        return this.request("POST", path, body);
    }
}

function createApp(config) {
    const { baseUrl, timeout, retries = MAX_RETRIES } = parseConfig(config);
    const client = new HttpClient(baseUrl, { timeout });
    const state = { users: [], posts: [], loading: false };

    async function loadUsers() {
        state.loading = true;
        try {
            state.users = await client.get("/users");
        } finally {
            state.loading = false;
        }
    }

    async function loadPosts(userId) {
        state.posts = await client.get(`/users/${userId}/posts`);
    }

    return {
        client,
        state,
        loadUsers,
        loadPosts,
        retries,
    };
}

export { HttpClient, createApp, DEFAULT_TIMEOUT, MAX_RETRIES };