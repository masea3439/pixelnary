export const eventEmitter = {
    listeners: {},
    on(eventName, callback) {
        if (this.listeners[eventName]) {
            this.listeners[eventName].push(callback)
        } else {
            this.listeners[eventName] = [callback]
        }
    },
    emit(eventName, data) {
        this.listeners[eventName].forEach(callback => {
            callback(data) 
        });
    }
}