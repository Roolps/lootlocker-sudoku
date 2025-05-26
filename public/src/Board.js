import { useState, useEffect } from "react";
import { Timer } from "./Timer";
import { ActionButtons } from "./ActionButtons";

export function Board({ setGameState, setPencilMarks, fetchPlayer, gameState, pencilState, exitGame, finishGame, playerBalance, initialStartTime }) {
    const [isPaused, setIsPaused] = useState(false);
    const [selection, setSelection] = useState({ row: 0, col: 0 });

    // check if the current cell is immutable
    var immutable = false
    if (
        Array.isArray(gameState) &&
        Number.isInteger(selection.row) &&
        Number.isInteger(selection.col) &&
        selection.row >= 0 &&
        selection.col >= 0 &&
        selection.row < gameState.length &&
        selection.col < (gameState[selection.row]?.length || 0)
    ) {
        immutable = !!gameState[selection.row][selection.col]?.immutable;
    }

    let highlightValue = "row" in selection && "col" in selection && "value" in gameState[selection.row][selection.col] ? gameState[selection.row][selection.col].value : 0;

    let boardElements = [];

    const numberTotals = {};
    for (let i = 1; i <= 9; i++) {
        numberTotals[i] = 0;
    }

    // generate the board from the game state
    gameState.forEach((row, i) => {
        var rowElements = []
        row.forEach((cell, j) => {
            let classes = "cell";
            let value = null;

            if (cell["immutable"]) {
                classes += " initial"
            }

            // add border-bottom to every 3rd row except row 9
            if ((i + 1) % 3 === 0 && i + 1 !== 9) {
                classes += " border-bottom"
            }
            // border-top if the previous row index was a multiple of 3
            if (i % 3 === 0 && i !== 0) {
                classes += " border-top";
            }
            // add border-right to every 3rd column except column 9
            if ((j + 1) % 3 === 0 && j + 1 !== 9) {
                classes += " border-right"
            }
            // border-left if the previous column index was a multiple of 3
            if (j % 3 === 0 && j !== 0) {
                classes += " border-left"
            }

            if (
                selection.row === i ||
                selection.col === j ||
                (Math.floor(i / 3) === Math.floor(selection.row / 3) &&
                    Math.floor(j / 3) === Math.floor(selection.col / 3))
            ) {
                classes += " cross";
            }

            if ("value" in cell) {
                value = cell.value;
                if (value === highlightValue) {
                    classes += " active";
                }
            } else {
                if (selection.row === i && selection.col === j) {
                    classes += " active";
                }
            }

            if ("pencil" in cell) {
                classes += " pencil-cell flex space-between";

                value = [...Array(9)].map((_, k) => {
                    const key = `Pencil-${i}-${j}-${k}`
                    const mark = cell.pencil[k]

                    return <span className="pencil" key={key} position={k + 1}>{mark ? (k + 1) : ""}</span>
                });
            }

            let cellName = `Cell-${i}-${j}`;

            if (Number.isInteger(value)) {
                numberTotals[value] = (Number(numberTotals[value]) || 0) + 1;
            }

            rowElements.push(<button className={classes} key={cellName} onClick={() => setSelection({ row: i, col: j })}>{value}</button>);
        });
        const rowName = `Row-${i}`
        boardElements.push(<div className="flex row" key={rowName}>{rowElements}</div>)
    });

    // handle numberbutton click
    function handleNumberButton(number) {
        if (isPaused) return;

        let newGameState = gameState.map(row => row.slice());
        let cell = gameState[selection.row][selection.col];

        if ("value" in cell && cell.immutable) {
            return;
        }

        if (pencilState) {
            delete cell.value

            // set the pencil state
            if (!cell.pencil || !Array.isArray(cell.pencil)) {
                cell.pencil = Array(9).fill(false);
            }
            if (cell.pencil[number]) {
                cell.pencil[number] = false
            } else {
                cell.pencil[number] = true
            }
        } else {
            if (cell.value === number + 1) {
                newGameState[selection.row][selection.col] = {};
            } else {
                newGameState[selection.row][selection.col] = { value: number + 1, immutable: false };
            }
        }

        // post new state to backend
        fetch("/api/state", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(newGameState),
        }).catch((error) => {
            // to do: make an error popup
            console.log(error)
        });

        setGameState(newGameState);
    }


    function togglePause() {
        setIsPaused(prev => !prev)
    }

    useEffect(() => {
        if (Array.isArray(gameState)) {
            verifyResult();
        }
    }, [gameState]);

    // check to see if the puzzle is complete
    function verifyResult() {
        var rows = Array.from({ length: 9 }, () => []);
        var columns = Array.from({ length: 9 }, () => []);

        gameState.forEach((row, i) => {
            row.forEach((cell, j) => {
                const val = cell.value ?? 0;
                rows[i].push(val);
                columns[j].push(val);
            })
        })

        const allRowsValid = rows.every(isValidGroup);
        const allColsValid = columns.every(isValidGroup);

        // if these are both valid then the solution is complete + valid!
        if (allRowsValid && allColsValid) {
            // run a complete game request and reset the game state
            document.getElementById("game-finished-popup").classList.add("active");
            setIsPaused(true);
        }
    }

    function isValidGroup(group) {
        const sorted = [...group].sort((a, b) => a - b);
        const expected = [1, 2, 3, 4, 5, 6, 7, 8, 9];
        return JSON.stringify(sorted) === JSON.stringify(expected);
    }

    // call to finish the game
    async function finishGame() {
        try {
            const response = await fetch(`/api/game`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));

                throw new Error(`Finish game failed with status ${response.status} : ${error.message}`);
            }

            setGameState([]);
            fetchPlayer();
        } catch (err) {
            console.error(err.message);

            const errorMsg = document.getElementById("finished-overlay-error");
            errorMsg.innerHTML = err.message;
            errorMsg.classList.add("active");

            setTimeout(() => {
                errorMsg.classList.remove("active");
            }, 5000);
        }
    }

    // clear the board and exit game
    async function exitGame() {
        try {
            const response = await fetch(`/api/state`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));

                throw new Error(`Exit game failed with status ${response.status} : ${error.message}`);
            }

            setGameState([]);
            fetchPlayer();
        } catch (err) {
            console.error(err.message);

            const errorMsg = document.getElementById("paused-overlay-error");
            errorMsg.innerHTML = err.message;
            errorMsg.classList.add("active");

            setTimeout(() => {
                errorMsg.classList.remove("active");
            }, 5000);
        }
    }

    return (
        <>
            <div id="timer-row" className="flex row space-between align-center">
                <div className="flex align-center">
                    <div id="player-tokens-icon"></div>
                    <p id="player-tokens" className="flex align-center">{playerBalance}</p>
                </div>
                <div className="flex align-center">
                    <Timer
                        isPaused={isPaused}
                        initialStartTime={initialStartTime}
                    />
                    <button id="pause-btn" className={isPaused ? "paused" : ""} onClick={togglePause}></button>
                </div>
            </div>
            <div className="flex column align-center" style={{ position: "relative" }}>
                <div className="cells flex column">{boardElements}</div>
                <div id="num-btns" className="flex row">
                    {[...Array(9)].map((_, i) => {
                        const isActive = (i + 1) === highlightValue;
                        const isLocked = isActive && immutable;
                        const isDimmed = !isActive && (immutable || numberTotals[i + 1] === 9);

                        return (
                            <button
                                key={i}
                                className={`num-btn ${isActive ? "active" : ""} ${isLocked ? "locked" : ""} ${isDimmed ? "dimmed" : ""}`}
                                onClick={() => handleNumberButton(i)}>
                                {i + 1}
                            </button>
                        )
                    })}
                </div>
                <ActionButtons
                    setPencilMarks={setPencilMarks}
                    pencilState={pencilState}
                ></ActionButtons>
                <div className={`popup flex column align-center justify-center ${isPaused ? "active" : ""}`}>
                    <h3>Game Paused</h3>
                    <p>Returning to menu will reset your progress.</p>
                    <div className="flex">
                        <button className="btn-solid" onClick={togglePause}>Resume</button>
                        <button className="btn-solid accent" onClick={exitGame}>Exit to Menu</button>
                    </div>
                    <p id="paused-overlay-error" className="error">something went wrong</p>
                </div>
                <div id="game-finished-popup" className="popup flex column align-center justify-center">
                    <h3>Woohoo! Good Job</h3>
                    <p>You have completed this level.</p>
                    <p className="score">+100 GridBits</p>
                    <button className="btn-solid accent" onClick={finishGame}>Exit to Menu</button>
                    <p id="finished-overlay-error" className="error">something went wrong</p>
                </div>
            </div>
        </>
    );
}