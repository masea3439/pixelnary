const canvas = document.getElementById('draw-canvas');
const gridSize = 5;
const squareMargin = 5;
let squareLength = null;
let pixels = new Array(gridSize ** 2).fill('#bfbfbf');

window.addEventListener('resize', handleResize);
canvas.addEventListener('mousemove', handleMouseMove);
canvas.addEventListener('mouseleave', handleMouseOut);

const ctx = canvas.getContext('2d');

function getMouseSquare(mouseX, mouseY) {
    let mouseSquareX = null;
    let mouseSquareY = null;

    if (mouseX && mouseY) {
        mouseSquareX = Math.floor((mouseX - squareMargin/2) / (squareLength + squareMargin));
        mouseSquareY = Math.floor((mouseY - squareMargin/2) / (squareLength + squareMargin));
    }
    if (mouseSquareX < 0 || mouseSquareX >= gridSize) {
        mouseSquareX = null
    }
    if (mouseSquareY < 0 || mouseSquareY >= gridSize) {
        mouseSquareY = null
    }
    return [mouseSquareX, mouseSquareY];
}

function drawGrid(ctx, gridSize, mouseX=null, mouseY=null) {
    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;
    
    const gridLength = Math.min(canvas.width, canvas.height);
    const squareX = centerX - gridLength / 2;
    const squareY = centerY - gridLength / 2;

    const [mouseSquareX, mouseSquareY] = getMouseSquare(mouseX, mouseY);
    
    for (let x = 0; x < gridSize; x++) {
        for (let y = 0; y < gridSize; y++) {
            if (mouseSquareX != null && mouseSquareY != null && mouseSquareX == x && mouseSquareY == y) {
                ctx.fillStyle = '#ffffff';
            } else {
                ctx.fillStyle = pixels[x*gridSize + y];
            }
            ctx.fillRect(
                squareMargin + squareX + (squareLength + squareMargin)*x,
                squareMargin + squareY + (squareLength + squareMargin)*y, 
                squareLength,
                squareLength
            );
        };
    };
}

function handleResize() {
    const gridLength = Math.min(window.innerWidth, window.innerHeight) * 0.7;
    squareLength = (gridLength - squareMargin*(gridSize+1)) / gridSize;
    canvas.width = gridLength;
    canvas.height = gridLength;
    drawGrid(ctx, gridSize);
}

function colorPixel(mouseX, mouseY, isDrawing) {
    const [mouseSquareX, mouseSquareY] = getMouseSquare(mouseX, mouseY);
    if (mouseSquareX != null && mouseSquareY != null && isDrawing) {
        pixels[mouseSquareX*gridSize + mouseSquareY] = '#188c41';
    }
}

function isPrimaryButtonPressed(buttons) {
    return (buttons % 2 !== 0)
}

function handleMouseMove(event) {
    const rect = canvas.getBoundingClientRect();
    isDrawing = isPrimaryButtonPressed(event.buttons)
    const mouseX = event.clientX - rect.left;
    const mouseY = event.clientY - rect.top;
    colorPixel(mouseX, mouseY, isDrawing)

  
    drawGrid(ctx, gridSize, mouseX, mouseY);
}
  
function handleMouseOut() {
    drawGrid(ctx, gridSize);
}

handleResize()