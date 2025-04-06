import { eventEmitter } from "./event_emitter.js";

const timer = document.getElementById("time")
let timeLeft = 60;
let timerIntervalId;

function tick() {
    if (timeLeft == 0) {
        clearInterval(timerIntervalId);
        return;
    }
    timeLeft--;
    let seconds = timeLeft % 60;
    let minutes = Math.floor(timeLeft/60);

    const strSeconds = String(seconds).padStart(2, "0");

    timer.textContent = `${minutes}:${strSeconds}`;
}

function startTimer() {
    timerIntervalId = setInterval(tick, 1000);
}

function stopTimer() {
    clearInterval(timerIntervalId);
}

startTimer();
