-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header

wrk.method = "POST"
wrk.body   = '{"sent_at":"2024-12-12T09:58:38.344Z","sdk":{"name":"sentry.javascript.vue","version":"8.26.0"},"dsn":"http://host.docker.internal:8081/0"}\n{"type":"session"}\n{"sid":"d959bbc52ffe475ca8ddf328e2a1570f","init":true,"started":"2024-12-12T09:58:38.343Z","timestamp":"2024-12-12T09:58:38.344Z","status":"ok","errors":0,"attrs":{"release":"nbcp-ncs-frontend@develop","environment":"development","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"}}'
wrk.headers["Content-Type"] = "text/plain;charset=UTF-8"
