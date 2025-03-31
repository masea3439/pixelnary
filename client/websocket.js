const currentUrl = window.location.href;
const roomKey = currentUrl.split('/').pop();
export const socket = new WebSocket(`ws://localhost:3333/ws?key=${roomKey}`);

export class PeriodicUpdateSocket {
    constructor(secondsBetweenRequests) {
        this.secondsBetweenRequests = secondsBetweenRequests;
        this.lastDataSent = null;
        this.newData = null;
        this.interval = setInterval(() => this.sendIfUpdate(), secondsBetweenRequests * 1000)
    }

    sendData(data) {
        this.newData = data;
    }

    sendIfUpdate() {
        if (this.newData != this.lastDataSent) {
            socket.send(this.newData);
            this.lastDataSent = this.newData;
        }
    }

    destroy() {
        clearInterval(this.interval)
    }
}