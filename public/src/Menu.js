import { useState } from 'react';

export function Menu() {
    const [selectedDifficulty, setSelectedDifficulty] = useState("easy");

    // send a post request to start the game
    function startgame(){

    }

    return (
        <div id="game-menu" className="flex column align-center">
            <div className="flex column align-start">
                <div className="new-game-btn-wrap flex column align-center">
                    {/* <h2>Select Difficulty</h2> */}
                    <button
                        className={`new-game-btn ${selectedDifficulty === "easy" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("easy")}
                    ><span>Easy</span></button>
                    <button
                        className={`new-game-btn ${selectedDifficulty === "medium" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("medium")}
                    ><span>Medium</span></button>
                    <button
                        className={`new-game-btn ${selectedDifficulty === "hard" ? "active" : ""}`}
                        onClick={() => setSelectedDifficulty("hard")}
                    ><span>Hard</span></button>
                </div>

                <div className="divider">
                    <p>OR</p>
                </div>

                <button id="custom-upload-btn"
                    className={`new-game-btn ${selectedDifficulty === "custom" ? "active" : ""}`}
                    onClick={() => setSelectedDifficulty("custom")}
                ><span>Upload Custom</span>
                </button>

                <button className="btn-solid">Start Game</button>
            </div>

            <div className="player-profile flex space-between align-center">
                <p>yara@clubnode.com</p>
                <button id="logout-btn">Logout</button>
            </div>
        </div>
    )
}