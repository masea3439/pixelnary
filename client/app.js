document.getElementById("host-game").addEventListener("click", function() {
    fetch("http://localhost:3333/api/host", {
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