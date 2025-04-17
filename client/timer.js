import { eventEmitter } from "./event_emitter.js";
import { gameState } from "./game_state.js";

const timer = document.getElementById("time");
const clockIcon = document.getElementById("clock-icon");
const nextRoundTimerText = document.getElementById("next-round-timer-text");

let timerIntervalId;

eventEmitter.on('game-state-updated', (data) => {
    if (gameState.matchState == "drawing") {
        clockIcon.style.display = "block";
        nextRoundTimerText.style.display = "none";
    } else if (gameState.matchState == "completed") {
        clockIcon.style.display = "none";
        nextRoundTimerText.style.display = "block";
    } else {
        setTimer();
        stopTimer();
        return;
    }
    startTimer();
});

function setTimer() {
    let seconds = gameState.roundTimeLeft % 60;
    let minutes = Math.floor(gameState.roundTimeLeft/60);

    const strSeconds = String(seconds).padStart(2, "0");

    timer.textContent = `${minutes}:${strSeconds}`;
}

function tick() {
    if (gameState.roundTimeLeft == 0) {
        clearInterval(timerIntervalId);
        return;
    }
    gameState.roundTimeLeft--;
    setTimer();
}

//TODO Fix timer stopping when tab not opened?
function startTimer() {
    stopTimer();
    setTimer();
    timerIntervalId = setInterval(tick, 1000);
}

function stopTimer() {
    clearInterval(timerIntervalId);
}