import { eventEmitter } from "./event_emitter.js";
import { sendMessage } from "./websocket.js";
import { gameState } from "./game_state.js";

const incorrectIcon = `<svg height="1.5em" fill="#000000" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M9.172 16.242 12 13.414l2.828 2.828 1.414-1.414L13.414 12l2.828-2.828-1.414-1.414L12 10.586 9.172 7.758 7.758 9.172 10.586 12l-2.828 2.828z"></path><path d="M12 22c5.514 0 10-4.486 10-10S17.514 2 12 2 2 6.486 2 12s4.486 10 10 10zm0-18c4.411 0 8 3.589 8 8s-3.589 8-8 8-8-3.589-8-8 3.589-8 8-8z"></path></g></svg>`
const correctIcon = `<svg height="1.5em" fill="#000000" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M12 2C6.486 2 2 6.486 2 12s4.486 10 10 10 10-4.486 10-10S17.514 2 12 2zm0 18c-4.411 0-8-3.589-8-8s3.589-8 8-8 8 3.589 8 8-3.589 8-8 8z"></path><path d="M9.999 13.587 7.7 11.292l-1.412 1.416 3.713 3.705 6.706-6.706-1.414-1.414z"></path></g></svg>`
const guessInput = document.getElementById('guess-input');
const guessMessage = document.getElementById('guess-message');
const guessLog = document.getElementById("guess-log");

eventEmitter.on('game-state-updated', (data) => {
    if (gameState.matchState == "drawing" && gameState.playerId != gameState.drawRolePlayerId) {
        guessInput.style.visibility = 'visible';
        guessMessage.style.visibility = 'visible';
    } else {
        guessInput.style.visibility = 'hidden';
        guessMessage.style.visibility = 'hidden';
    }
});
eventEmitter.on('guess', addGuess);
eventEmitter.on('start-round', clearGuesses);

function addGuess(guessJson) {
    const guess = JSON.parse(guessJson);

    const newGuess = document.createElement('div');

    newGuess.classList.add('guess');
    newGuess.innerHTML = "";

    if (guess.isCorrect) {
        newGuess.classList.add('correct-guess');
        newGuess.innerHTML += correctIcon
    } else {
        newGuess.classList.add('incorrect-guess');
        newGuess.innerHTML += incorrectIcon
    }

    const newGuessText = document.createElement('div');
    newGuessText.innerHTML = guess.guess;

    newGuess.appendChild(newGuessText);
    guessLog.appendChild(newGuess);

    guessMessage.style.visibility = 'hidden';
}

function clearGuesses() {
    guessLog.replaceChildren();
}

guessInput.addEventListener('keydown', function(event) {
    if (event.key == 'Enter') {
        if (guessInput.value.trim() !== "") {
            sendMessage('guess', guessInput.value.trim());
        }
        guessInput.value = '';
    }
});