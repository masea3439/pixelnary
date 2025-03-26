const colors = ['#ff0000', '#ff8000', '#ffff00', '#00ffff', '#0000ff', '#8000ff', '#ff00ff', '#ffffff', '#8c8c8c', '#000000'];

const colorPicker = document.getElementById('color-picker');
export let selectedColor = '#000000';
colors.forEach(color => {
    const colorBox = document.createElement('div');
    colorBox.style.backgroundColor = color;
    colorBox.classList.add('color-box');
    colorBox.addEventListener('click', () => selectColor(color));
    colorPicker.appendChild(colorBox);
})

function selectColor(color) {
    selectedColor = color;
}