import {gameState} from "./game_state.js";

const colors = [
    '#ff0000', '#ff8000', '#ffff00', '#00ff00', '#00ffff',
    '#0080ff', '#0000ff', '#8000ff', '#ff00ff', '#ff0080', 
    '#ff8096', '#000000', '#bebebe', '#686868', '#ffffff',
    '#8e6900', '#c28000', '#008080', '#800080', '#009933'
];

const colorPicker = document.getElementById('color-picker');
colors.forEach(color => {
    const colorBox = document.createElement('div');
    colorBox.style.backgroundColor = color;
    colorBox.classList.add('color-box');
    colorBox.addEventListener('click', () => selectColor(color));
    colorPicker.appendChild(colorBox);
})

//colorPicker.style.visibility = 'hidden';

function selectColor(color) {
    gameState.selectedColor = color;
}