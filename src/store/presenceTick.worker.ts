const HEARTBEAT_MS = 3000
setInterval(() => {
  postMessage(0)
}, HEARTBEAT_MS)
