import { eventEmitter } from "./event_emitter.js";
import { gameState } from "./game_state.js";

const pin = document.getElementById('pin');
const linkButton = document.getElementById('invite-link');
const copyLinkMessage = document.getElementById('copy-link-message');
const waitingScreen = document.getElementById('waiting-screen');
const outerGameContainer = document.getElementById('outer-game-container');

var copyLinkTimeout;

const currentUrl = window.location.href;
const roomKey = currentUrl.split('/').pop();
pin.innerHTML = `<strong>${roomKey}</strong>`;

linkButton.addEventListener('click', function() {
    navigator.clipboard.writeText(`http://localhost:8080/game/${roomKey}`); //TODO replace
    clearTimeout(copyLinkTimeout)
    copyLinkMessage.style.visibility = 'visible';
    copyLinkTimeout = setTimeout(() => {
        copyLinkMessage.style.visibility = 'hidden';
    }, 3000)
});

function startRound(roundJson) {
    const roundData = JSON.parse(roundJson);
    gameState.playerId = roundData.playerId;
    gameState.drawRolePlayerId = roundData.drawRolePlayerId;
    gameState.roundTimeLeft = roundData.roundTimeLeft;
    gameState.gridSize = roundData.gridSize;
    gameState.pixels = new Array(gameState.gridSize ** 2).fill('#bfbfbf');
    gameState.selectedColor = '#000000';

    gameState.matchState = "drawing";

    eventEmitter.emit('game-state-updated', null);
    waitingScreen.style.display = "none";
    outerGameContainer.style.display = "visible";
}

eventEmitter.on('start-round', startRound);