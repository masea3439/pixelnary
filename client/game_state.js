export let gameState = {
    playerId: null,
    drawRolePlayerId: null,
    matchState: "drawing",
    roundTimeLeft: 60,
    gridSize: 6,
    pixels: new Array(6 ** 2).fill('#bfbfbf'),
    selectedColor: '#000000'
};

//let pixels = new Array(gridSize ** 2).fill('#bfbfbf');