async function handleRequest(request) {
    const clientIP = request.headers.get('CF-Connecting-IP');
    return new Response(clientIP + "\n");
}

addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request));
});
