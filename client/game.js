import { eventEmitter } from "./event_emitter.js";
import { gameState } from "./game_state.js";
import { sendMessage } from "./websocket.js";

const pin = document.getElementById('pin');
const linkButton = document.getElementById('invite-link');
const copyLinkMessage = document.getElementById('copy-link-message');
const waitingScreen = document.getElementById('waiting-screen');
const outerGameContainer = document.getElementById('outer-game-container');
const word = document.getElementById('word');
const disconnectPopup = document.getElementById('disconnect-popup');
const gameOverPopup = document.getElementById('game-over-popup');
const rematchButton = document.getElementById('rematch');
const rematchDialogue = document.getElementById('rematch-dialogue');
const rematchLoading = document.getElementById('rematch-loading');

var copyLinkTimeout;

const currentUrl = window.location.href;
const roomKey = currentUrl.split('/').pop();
pin.innerHTML = `<strong>${roomKey}</strong>`;

eventEmitter.on('start-round', startRound);
eventEmitter.on('round-completed', roundCompleted);
eventEmitter.on('disconnect', disconnected);
eventEmitter.on('game-over', gameOver);

linkButton.addEventListener('click', function() {
    navigator.clipboard.writeText(`http://localhost:8080/game/${roomKey}`); //TODO replace
    clearTimeout(copyLinkTimeout)
    copyLinkMessage.style.visibility = 'visible';
    copyLinkTimeout = setTimeout(() => {
        copyLinkMessage.style.visibility = 'hidden';
    }, 3000);
});

Array.from(document.getElementsByClassName('menu')).forEach(function(menuButton) {
    menuButton.addEventListener('click', function() {
        window.location.href = "/";
    });
});

rematchButton.addEventListener('click', function() {
    rematchDialogue.style.display = 'none';
    rematchLoading.style.display = 'flex';
    sendMessage('rematch', '');
});

function updateWord() {
    if (gameState.word == null) {
        return;
    }
    if (gameState.word == "") {
        word.innerText = "";
        return;
    }
    let message;
    if (gameState.playerId == gameState.drawRolePlayerId) {
        message = "The word is: ";
    } else {
        message = "The word was: ";
    }
    word.innerText = message + gameState.word;
}

function startRound(roundJson) {
    gameOverPopup.style.display = 'none';

    const roundData = JSON.parse(roundJson);
    gameState.playerId = roundData.playerId;
    gameState.drawRolePlayerId = roundData.drawRolePlayerId;
    gameState.roundTimeLeft = roundData.roundTimeLeft;
    gameState.gridSize = roundData.gridSize;
    gameState.word = roundData.word;
    gameState.pixels = new Array(gameState.gridSize ** 2).fill('#bfbfbf');
    gameState.selectedColor = '#000000';

    gameState.matchState = "drawing";

    updateWord();

    outerGameContainer.style.display = "flex";
    eventEmitter.emit('game-state-updated', null);
    waitingScreen.style.display = "none";

}

function roundCompleted(completedMessageJson) {
    const completedMessage = JSON.parse(completedMessageJson);
    gameState.matchState = "completed";
    gameState.roundTimeLeft = completedMessage.timeUntilNextRound;
    eventEmitter.emit('game-state-updated', null);
}

function disconnected(data) {
    gameState.matchState = "disconnected";
    gameOverPopup.style.display = 'none';
    disconnectPopup.style.display = 'flex';
    eventEmitter.emit('game-state-updated', null);
}

function gameOver(gameOverMessageJson) {
    const gameOverMessage = JSON.parse(gameOverMessageJson);
    gameState.matchState = "game-over";
    gameState.word = gameOverMessage.word;
    gameState.roundTimeLeft = 0;

    updateWord();

    rematchDialogue.style.display = 'flex';
    rematchLoading.style.display = 'none';
    gameOverPopup.style.display = 'flex';
    eventEmitter.emit('game-state-updated', null);
}