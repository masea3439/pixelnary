const currentUrl = window.location.href;
const roomKey = currentUrl.split('/').pop();
export const socket = new WebSocket(`ws://localhost:3333/ws?key=${roomKey}`);

export class RateLimitedSocket {
    constructor(minSecondsBetweenRequests) {
        this.minSecondsBetweenRequests = minSecondsBetweenRequests;
        this.lastSendTime = 0;
    }

    // TODO cache and send most recent request after time elapsed
    send(data) {
        const now = Date.now();
        const elapsedTime = now - this.lastSendTime;
        if (elapsedTime > this.minSecondsBetweenRequests * 1000) {
            this.lastSendTime = now;
            socket.send(data);
        }
    }
}