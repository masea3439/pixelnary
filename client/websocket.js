import { eventEmitter } from "./event_emitter.js";

const currentUrl = window.location.href;
const roomKey = currentUrl.split('/').pop();
const socket = new WebSocket(`/ws?key=${roomKey}`);

export function sendMessage(messageType, data) {
    const message = {
        messageType: messageType,
        data: data
    };
    socket.send(JSON.stringify(message));
}

export class PeriodicUpdateSocket {
    constructor(messageType, secondsBetweenRequests) {
        this.messageType = messageType;
        this.secondsBetweenRequests = secondsBetweenRequests;
        this.lastDataSent = null;
        this.newData = null;
        this.interval = setInterval(() => this.sendIfUpdate(), secondsBetweenRequests * 1000);
    }

    sendData(data) {
        this.newData = data;
    }

    sendIfUpdate() {
        if (this.newData != this.lastDataSent) {
            sendMessage(this.messageType, this.newData);
            this.lastDataSent = this.newData;
        }
    }

    destroy() {
        clearInterval(this.interval);
    }
}

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    eventEmitter.emit(message.messageType, message.data);
};

socket.onerror = (event) => {
    eventEmitter.emit('does-not-exist', '');
};