package router

// key => room:uuid value => { ...user tokens }
// key => user:token => { username, activeRoom, ... }

// Start websocket connection steps
// - Send token in query parameters to authorize
// - if token is valid we upgrade connection, if not refuse
// - Change token property in redis to connected but no room so far

// Connect to room
// At this point we could re-validate token, but I don't think it's necessary
// room:id:user:username
// Update room list (watch out for race conditions in this resource) <-- this could be a Set to prevent duplication
// Update user active room property

// Leave room
// Update room list
// Update user active room property

// Send message
// - Validate token
// - Retrieve all connected users in room and broadcast message
// - update ui
