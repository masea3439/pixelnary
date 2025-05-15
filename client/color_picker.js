import {gameState} from "./game_state.js";
import { eventEmitter } from "./event_emitter.js";

const colors = [
    '#ff0000', '#ff8000', '#ffff00', '#00ff00', '#00ffff',
    '#0080ff', '#0000ff', '#8000ff', '#ff00ff', '#ff0080', 
    '#ff8096', '#008080', '#800080', '#009933', '#c28000',
    '#8e6900', '#ffffff', '#bebebe', '#686868', '#000000'
];

const colorPickerWrapper = document.getElementById('color-picker-wrapper');
const colorPicker = document.getElementById('color-picker');
const drawContainer = document.getElementById('draw-container');

colors.forEach(color => {
    const colorBox = document.createElement('div');
    colorBox.style.backgroundColor = color;
    colorBox.classList.add('color-box');
    colorBox.addEventListener('pointerdown', () => selectColor(color));
    colorPicker.appendChild(colorBox);
});

eventEmitter.on('game-state-updated', (data) => {
    if (gameState.matchState == "drawing" && gameState.playerId == gameState.drawRolePlayerId) {
        colorPickerWrapper.style.display = 'flex';
        drawContainer.classList.add('draw-container-two-rows');
    } else {
        colorPickerWrapper.style.display = 'none';
        drawContainer.classList.remove('draw-container-two-rows');
    }
});

function selectColor(color) {
    gameState.selectedColor = color;
}