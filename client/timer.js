import { eventEmitter } from "./event_emitter.js";
import { gameState } from "./game_state.js";

const timer = document.getElementById("time")
let timerIntervalId;

eventEmitter.on('game-state-updated', startTimer);

function tick() {
    if (gameState.roundTimeLeft == 0) {
        clearInterval(timerIntervalId);
        return;
    }
    gameState.roundTimeLeft--;
    let seconds = gameState.roundTimeLeft % 60;
    let minutes = Math.floor(gameState.roundTimeLeft/60);

    const strSeconds = String(seconds).padStart(2, "0");

    timer.textContent = `${minutes}:${strSeconds}`;
}

function startTimer() {
    timerIntervalId = setInterval(tick, 1000);
}

function stopTimer() {
    clearInterval(timerIntervalId);
}