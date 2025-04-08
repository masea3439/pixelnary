const roomCodeInput = document.getElementById('room-code-input');

roomCodeInput.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        window.location.href = `game/${roomCodeInput.value}`;
    }
});

document.getElementById("room-connect-button").addEventListener("click", function() {
    window.location.href = `game/${roomCodeInput.value}`;
});