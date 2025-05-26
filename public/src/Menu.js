import { useState } from "react";

export function Menu({ fetchState }) {
    const [selectedDifficulty, setSelectedDifficulty] = useState("easy");

    // send a post request to start the game
    async function startgame() {
        try {
            const response = await fetch(`/api/game`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    difficulty: selectedDifficulty,
                })
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));

                throw new Error(`Start game failed with status ${response.status} : ${error.message}`);
            }

            // game was started - fetch the state!
            fetchState();
        } catch (err) {
            console.error(err.message);

            const errorMsg = document.getElementById("start-game-error");
            errorMsg.innerHTML = err.message;
            errorMsg.classList.add("active");

            setTimeout(() => {
                errorMsg.classList.remove("active");
            }, 5000);
        }
    }

    async function logout() {
        try {
            const response = await fetch(`/api/logout`, {
                method: "POST",
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));

                throw new Error(`Start game failed with status ${response.status} : ${error.message}`);
            }

            // game was started - fetch the state!
            fetchState();
        } catch (err) {
            console.error(err.message);

            alert(err.message);
        }
    }

    function getEmail() {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; email=`);
        if (parts.length === 2) return parts.pop().split(";").shift();
    }

    return (
        <div id="game-menu" className="flex column align-center">
            <div className="flex column align-start">
                <div className="new-game-btn-wrap flex column align-center">
                    {/* <h3>Select Difficulty</h3> */}
                    <button
                        className={`new-game-btn ${selectedDifficulty === "easy" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("easy")}>
                        <span>Easy<br></br><sup>+50 GridBits</sup></span>
                    </button>
                    <button
                        className={`new-game-btn coming-soon ${selectedDifficulty === "medium" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("medium")}>
                        <span>Medium<br></br><sup>+100 GridBits</sup></span>
                    </button>
                    <button
                        className={`new-game-btn coming-soon ${selectedDifficulty === "hard" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("hard")}>
                        <span>Hard<br></br><sup>+150 GridBits</sup></span>
                    </button>
                </div>

                <div className="divider">
                    <p>OR</p>
                </div>

                <button id="custom-upload-btn"
                    className={`new-game-btn coming-soon ${selectedDifficulty === "custom" ? "active" : ""}`}
                    onClick={() => setSelectedDifficulty("custom")}>
                    <span>Upload Custom<br></br><sup>Click to upload</sup></span>
                </button>

                <button className="btn-solid" onClick={startgame}>Start Game</button>
                <p id="start-game-error" className="error">something went wrong</p>
            </div>

            <div className="player-profile flex space-between align-center">
                <p>{getEmail()}</p>
                <button id="logout-btn" onClick={logout}>Logout</button>
            </div>
        </div>
    )
}