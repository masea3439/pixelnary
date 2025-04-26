document.getElementById("host-game").addEventListener("click", function() {
    fetch("/api/host", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => response.text())
    .then(roomKey => window.location.href = `game/${roomKey}`)
    .catch(error => console.error("Error:", error));
})

document.getElementById("join-game").addEventListener("click", function() {
    window.location.href = "connect"
})

const titleContainer = document.getElementById("title-container");
const svgData = [
`M2359 2735 c-3 -2 -126 -12 -275 -22 -282 -19 -373 -19 -1169 -8
l-490 6 -22 -23 c-19 -20 -24 -47 -39 -203 -16 -157 -19 -304 -21 -1135 -2
-606 -7 -958 -13 -964 -20 -20 -11 -71 17 -98 93 -87 223 -108 657 -108 182 0
392 5 466 11 299 23 469 42 545 60 110 25 500 89 608 98 77 7 92 12 115 35
l27 26 0 768 c0 819 2 773 -51 1112 -8 52 -18 165 -21 250 -7 208 1 200 -188
200 -78 0 -144 -2 -146 -5z`,

`M1615 2953 c-134 -18 -256 -26 -480 -33 -265 -7 -315 -11 -430 -35
-71 -14 -180 -29 -242 -32 -99 -5 -115 -8 -133 -28 -20 -21 -21 -31 -16 -106
3 -46 9 -131 12 -189 3 -58 11 -111 16 -118 12 -14 3 -208 -13 -282 -5 -25
-14 -144 -18 -265 -8 -192 -12 -230 -35 -303 -44 -145 -48 -215 -45 -767 l4
-514 28 -29 c28 -29 30 -29 165 -35 75 -4 166 -13 202 -22 49 -12 154 -18 425
-25 343 -8 504 -18 624 -40 31 -6 280 -10 581 -10 l526 0 40 43 c68 70 82 95
89 152 12 89 -13 556 -35 673 -19 97 -20 152 -20 926 l0 822 -27 26 c-31 27
-81 53 -128 63 -16 4 -39 13 -50 20 -11 6 -49 17 -85 24 -36 7 -118 28 -183
48 -225 66 -470 77 -772 36z`,

`M1985 2875 c-71 -8 -155 -18 -185 -24 -30 -5 -120 -13 -200 -16 -80
-4 -261 -22 -402 -41 -311 -41 -593 -54 -742 -34 -96 13 -97 12 -121 -10 -29
-27 -30 -49 -10 -145 32 -154 -2 -890 -56 -1180 -56 -303 -72 -503 -67 -845 3
-261 3 -266 28 -309 49 -87 41 -85 388 -92 169 -4 388 -14 487 -23 144 -12
334 -15 946 -16 l766 0 21 23 c20 21 21 32 19 118 -3 81 -7 103 -34 159 -36
77 -76 207 -93 305 -18 100 0 348 28 401 37 71 40 393 7 854 -8 118 -20 323
-25 455 -6 132 -13 256 -16 277 -8 52 -52 90 -122 105 -31 6 -88 20 -127 29
-91 22 -319 26 -490 9z`,

`M320 2722 c-38 -6 -70 -41 -70 -77 0 -34 31 -70 65 -77 13 -2 26 -18
35 -43 13 -35 15 -149 12 -830 -1 -434 -7 -936 -13 -1113 l-11 -324 22 -23
c21 -22 33 -24 149 -30 69 -3 571 -8 1115 -11 989 -4 990 -4 1011 16 22 21 22
23 28 503 3 344 9 521 20 617 53 453 63 863 28 1072 -11 63 -25 126 -30 142
-6 15 -11 51 -11 80 0 44 -4 57 -26 77 l-26 24 -1136 0 c-625 0 -1148 -2
-1162 -3z`,

`M142 2777 c-33 -35 -28 -75 13 -117 24 -25 38 -50 45 -85 5 -28 19
-93 31 -145 19 -86 22 -133 26 -480 4 -372 4 -392 -22 -590 -14 -113 -30 -248
-35 -300 -5 -52 -16 -131 -26 -175 -9 -44 -20 -100 -24 -125 -5 -25 -9 -122
-9 -216 l-1 -171 26 -24 26 -24 801 -5 c689 -4 806 -7 827 -20 20 -12 79 -16
320 -18 l295 -2 117 110 c173 164 281 332 342 532 23 74 26 102 26 229 0 118
-5 161 -24 239 -27 110 -100 309 -131 360 -12 19 -53 96 -92 170 -82 158 -90
189 -102 436 -8 169 -9 172 -39 216 -17 25 -47 60 -68 77 l-37 33 -131 -5
c-79 -2 -210 -17 -331 -37 -110 -18 -241 -35 -291 -37 -90 -5 -93 -5 -167 32
-76 38 -187 65 -267 65 -25 0 -94 8 -155 16 -60 9 -164 19 -230 24 -240 15
-531 40 -584 50 -81 15 -105 12 -129 -13z`
];

const letterP = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 3, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
    {row: 3, col: 3},
];
const letterI = [
    {row: 1, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 2, col: 2},
    {row: 3, col: 2},
    {row: 4, col: 2},
    {row: 5, col: 2},
    {row: 1, col: 3},
    {row: 5, col: 3},
];
const letterX = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 3, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
    {row: 4, col: 3},
    {row: 5, col: 3},
];
const letterE = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 3, col: 2},
    {row: 5, col: 2},
    {row: 1, col: 3},
    {row: 3, col: 3},
    {row: 5, col: 3},
];
const letterL = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 5, col: 2},
    {row: 5, col: 3},
];
const letterN = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
    {row: 3, col: 3},
    {row: 4, col: 3},
    {row: 5, col: 3},
];
const letterA = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 3, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
    {row: 3, col: 3},
    {row: 4, col: 3},
    {row: 5, col: 3},
];
const letterR = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 1},
    {row: 4, col: 1},
    {row: 5, col: 1},
    {row: 1, col: 2},
    {row: 3, col: 2},
    {row: 4, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
    {row: 3, col: 3},
    {row: 5, col: 3},
];
const letterY = [
    {row: 1, col: 1},
    {row: 2, col: 1},
    {row: 3, col: 2},
    {row: 4, col: 2},
    {row: 5, col: 2},
    {row: 1, col: 3},
    {row: 2, col: 3},
];
const title = [letterP, letterI, letterX, letterE, letterL, letterN, letterA, letterR, letterY];
let squares = [];

function createSquareSVG(row, col, color, frame) {
    const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    svg.setAttribute('viewBox', '0 0 300 300');
    svg.setAttribute('width', '100%');
    svg.setAttribute('height', '100%');
    svg.setAttribute('preserveAspectRatio', 'xMidYMid meet');
    svg.style.gridRow = row;
    svg.style.gridColumn = col;

    const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    g.setAttribute('fill', color);
    g.setAttribute("transform", "translate(0.000000,300.000000) scale(0.100000,-0.100000)");
    g.setAttribute("stroke", "none");

    const path = document.createElementNS('http://www.w3.org/2000/svg', "path");
    path.setAttribute('d', svgData[frame]);

    g.appendChild(path);
    svg.appendChild(g);

    return [svg, path]
}

let letterNumber = 0;
title.forEach(letter => {
    letter.forEach(squareCoord => {
        const frame = Math.floor(Math.random() * 5);
        const [svg, path] = createSquareSVG(squareCoord.row, squareCoord.col + letterNumber * 4, '#686868', frame);
        titleContainer.appendChild(svg);
        squares.push({path: path, frame: frame});
    });
    letterNumber++;
});

const squareIntervalId = setInterval(() => {
    squares.forEach(square => {
        const frame = (square.frame + 1) % 5;
        square.path.setAttribute('d', svgData[frame]);
        square.frame = frame;
    })
}, 500);

const tree = {
    id: 'tree-drawing',
    rows: 5,
    cols: 5,
    squares: [
        {row: 2, col: 1, color: '#009933'},
        {row: 3, col: 1, color: '#009933'},
        {row: 1, col: 2, color: '#009933'},
        {row: 2, col: 2, color: '#009933'},
        {row: 3, col: 2, color: '#009933'},
        {row: 1, col: 3, color: '#009933'},
        {row: 2, col: 3, color: '#009933'},
        {row: 3, col: 3, color: '#009933'},
        {row: 4, col: 3, color: '#8e6900'},
        {row: 5, col: 3, color: '#8e6900'},
        {row: 1, col: 4, color: '#009933'},
        {row: 2, col: 4, color: '#009933'},
        {row: 3, col: 4, color: '#009933'},
        {row: 2, col: 5, color: '#009933'},
        {row: 3, col: 5, color: '#009933'},
    ]
};

const apple = {
    id: 'apple-drawing',
    rows: 5,
    cols: 4,
    squares: [
        {row: 3, col: 1, color: '#ff0000'},
        {row: 4, col: 1, color: '#ff0000'},
        {row: 2, col: 2, color: '#ff0000'},
        {row: 3, col: 2, color: '#ff0000'},
        {row: 4, col: 2, color: '#ff0000'},
        {row: 5, col: 2, color: '#ff0000'},
        {row: 1, col: 3, color: '#8e6900'},
        {row: 2, col: 3, color: '#ff0000'},
        {row: 3, col: 3, color: '#ff0000'},
        {row: 4, col: 3, color: '#ff0000'},
        {row: 5, col: 3, color: '#ff0000'},
        {row: 1, col: 4, color: '#009933'},
        {row: 3, col: 4, color: '#ff0000'},
        {row: 4, col: 4, color: '#ff0000'},
    ]
};

const fish = {
    id: 'fish-drawing',
    rows: 3,
    cols: 5,
    squares: [
        {row: 1, col: 1, color: '#ff8000'},
        {row: 3, col: 1, color: '#ff8000'},
        {row: 2, col: 2, color: '#ff8000'},
        {row: 1, col: 3, color: '#ff8000'},
        {row: 2, col: 3, color: '#ff8000'},
        {row: 3, col: 3, color: '#ff8000'},
        {row: 1, col: 4, color: '#ff8000'},
        {row: 2, col: 4, color: '#ff8000'},
        {row: 3, col: 4, color: '#ff8000'},
        {row: 1, col: 5, color: '#ff8000'},
        {row: 2, col: 5, color: '#000000'},
        {row: 3, col: 5, color: '#ff8000'},
    ]
};

const pencil = {
    id: 'pencil-drawing',
    rows: 4,
    cols: 1,
    squares: [
        {row: 1, col: 1, color: '#ff8096'},
        {row: 2, col: 1, color: '#ffff00'},
        {row: 3, col: 1, color: '#ffff00'},
        {row: 4, col: 1, color: '#000000'},
    ]
};

const earth = {
    id: 'earth-drawing',
    rows: 3,
    cols: 5,
    squares: [
        {row: 1, col: 1, color: '#009933'},
        {row: 2, col: 1, color: '#009933'},
        {row: 3, col: 1, color: '#009933'},
        {row: 4, col: 1, color: '#0080ff'},
        {row: 1, col: 2, color: '#0080ff'},
        {row: 2, col: 2, color: '#009933'},
        {row: 3, col: 2, color: '#0080ff'},
        {row: 4, col: 2, color: '#0080ff'},
        {row: 1, col: 3, color: '#0080ff'},
        {row: 2, col: 3, color: '#0080ff'},
        {row: 3, col: 3, color: '#0080ff'},
        {row: 4, col: 3, color: '#009933'},
        {row: 1, col: 4, color: '#009933'},
        {row: 2, col: 4, color: '#0080ff'},
        {row: 3, col: 4, color: '#009933'},
        {row: 4, col: 4, color: '#009933'},
    ]
};

const lightning = {
    id: 'lightning-drawing',
    rows: 4,
    cols: 4,
    squares: [
        {row: 2, col: 1, color: '#bebebe'},
        {row: 1, col: 2, color: '#bebebe'},
        {row: 2, col: 2, color: '#bebebe'},
        {row: 3, col: 2, color: '#ffff00'},
        {row: 1, col: 3, color: '#bebebe'},
        {row: 2, col: 3, color: '#bebebe'},
        {row: 4, col: 3, color: '#ffff00'},
        {row: 2, col: 4, color: '#bebebe'},
    ]
};

const fire = {
    id: 'fire-drawing',
    rows: 5,
    cols: 5,
    squares: [
        {row: 3, col: 1, color: '#ff0000'},
        {row: 4, col: 1, color: '#ff0000'},
        {row: 2, col: 2, color: '#ff0000'},
        {row: 3, col: 2, color: '#ff8000'},
        {row: 4, col: 2, color: '#ffff00'},
        {row: 5, col: 2, color: '#8e6900'},
        {row: 1, col: 3, color: '#ff0000'},
        {row: 2, col: 3, color: '#ff8000'},
        {row: 3, col: 3, color: '#ffff00'},
        {row: 4, col: 3, color: '#ffff00'},
        {row: 5, col: 3, color: '#8e6900'},
        {row: 2, col: 4, color: '#ff0000'},
        {row: 3, col: 4, color: '#ff8000'},
        {row: 4, col: 4, color: '#ffff00'},
        {row: 5, col: 4, color: '#8e6900'},
        {row: 3, col: 5, color: '#ff0000'},
        {row: 4, col: 5, color: '#ff0000'},
    ]
};


const backgroundDrawings = [tree, apple, fish, pencil, earth, lightning, fire];

backgroundDrawings.forEach(drawing => {
    const container = document.createElement('div');
    container.id = drawing.id;

    container.style.gridTemplateRows = `repeat(${drawing.rows}, 1fr);`;
    container.style.gridTemplateColumns = `repeat(${drawing.cols}, 1fr);`;
    const width = 80 / 35 * drawing.cols;
    container.style.width = `${width}%`;

    drawing.squares.forEach((squareCoord) => {
        const frame = Math.floor(Math.random() * 5);
        const [svg, path] = createSquareSVG(squareCoord.row, squareCoord.col, squareCoord.color, frame);

        container.appendChild(svg);
        squares.push({path: path, frame: frame});
    })

    document.body.appendChild(container);
})