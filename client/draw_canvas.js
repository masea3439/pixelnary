import { PeriodicUpdateSocket } from "./websocket.js";
import { eventEmitter } from "./event_emitter.js";
import { gameState } from "./game_state.js";

const drawSocket = new PeriodicUpdateSocket("canvas", 0.25);
const drawMessage = document.getElementById('draw-message');
const canvas = document.getElementById('draw-canvas');
let drawPixels = [];

eventEmitter.on('canvas', (data) => {
    console.log('canvas')
    gameState.pixels = data.split(',');
    updatePixels();
});

eventEmitter.on('game-state-updated', (data) => {
    resetGrid();
    showMessage();
    handleResize();
});

window.addEventListener('resize', handleResize);

function handleResize() {
    const parentElement = canvas.parentElement;
    const canvasLength = Math.min(parentElement.offsetWidth, parentElement.offsetHeight);
    canvas.style.width = `${canvasLength}px`;
    canvas.style.height = `${canvasLength}px`;
}

function showMessage() {
    if (gameState.matchState == "drawing" && gameState.playerId == gameState.drawRolePlayerId) {
        drawMessage.style.visibility = 'visible';
    } else {
        drawMessage.style.visibility = 'hidden';
    }
}

function hideMessage() {
    drawMessage.style.visibility = 'hidden';
}

function resetGrid() {
    canvas.innerHTML = "";
    drawPixels = [];
    canvas.removeEventListener('pointerdown', colorPixelIfDrawing);
    canvas.removeEventListener('pointermove', colorPixelIfDrawing);

    canvas.style.gridTemplateColumns = `repeat(${gameState.gridSize}, 1fr)`;
    canvas.style.gridTemplateRows = `repeat(${gameState.gridSize}, 1fr)`;
    for (let i = 0; i < gameState.gridSize**2; i++) {
        const pixel = document.createElement('div');
        pixel.classList.add('draw-pixel');
        pixel.style.backgroundColor = gameState.pixels[i];
        drawPixels.push(pixel);
        canvas.appendChild(pixel);
    }
    if (gameState.matchState == "drawing" && gameState.playerId == gameState.drawRolePlayerId) {
        canvas.addEventListener('pointerdown', colorPixelIfDrawing);
        canvas.addEventListener('pointermove', colorPixelIfDrawing);
    };
}

function updatePixels() {
    for (let i = 0; i < drawPixels.length; i++) {
        drawPixels[i].style.backgroundColor = gameState.pixels[i];
    }
}

function colorPixelIfDrawing(event) {
    if (event.buttons % 2 !== 0) {
        colorPixel(event);
    }
}

function colorPixel(event) {
    let pixel = document.elementFromPoint(event.clientX, event.clientY);
    if (pixel && pixel.classList.contains('draw-pixel')) {
        const index = Array.prototype.indexOf.call(drawPixels, pixel);
        gameState.pixels[index] = gameState.selectedColor;
        drawSocket.sendData(gameState.pixels.toString());
        hideMessage();
        pixel.style.backgroundColor = gameState.selectedColor;
    }
}

handleResize();