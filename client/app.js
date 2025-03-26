const socket = new WebSocket("ws://localhost:3333/ws");

export class RateLimitedSocket {
    constructor(minSecondsBetweenRequests) {
        this.minSecondsBetweenRequests = minSecondsBetweenRequests;
        this.lastSendTime = 0;
    }

    send(data) {
        const now = Date.now();
        const elapsedTime = now - this.lastSendTime;
        if (elapsedTime > this.minSecondsBetweenRequests * 1000) {
            this.lastSendTime = now;
            socket.send(data);
        }
    }
}