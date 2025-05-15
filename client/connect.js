const roomCodeInput = document.getElementById('room-code-input');

roomCodeInput.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        connectToRoom();
    }
});

document.getElementById("room-connect-button").addEventListener("click", function() {
    connectToRoom();
});

function connectToRoom() {
    if (roomCodeInput.value) {
        window.location.href = `game/${roomCodeInput.value}`;
    }
}