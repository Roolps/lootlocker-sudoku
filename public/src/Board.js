export function Board({ fetchState, gameState, pencilState, selection, onCellClick, onNumberButtonClick, onActionButtonClick }) {
    let highlightValue = "row" in selection && "col" in selection && "value" in gameState[selection.row][selection.col] ? gameState[selection.row][selection.col].value : 0;

    let boardElements = [];

    const numberTotals = {};
    for (let i = 1; i <= 9; i++) {
        numberTotals[i] = 0;
    }

    // loop through each row of the game's state
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

            rowElements.push(<button className={classes} key={cellName} onClick={() => onCellClick(i, j)}>{value}</button>);
        });
        const rowName = `Row-${i}`
        boardElements.push(<div className="flex row" key={rowName}>{rowElements}</div>)
    });

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
        // set immutable to whether the current selection is mutable or not
        immutable = !!gameState[selection.row][selection.col]?.immutable;
    }

    async function exitGame() {
        try {
            const response = await fetch(`/api/state`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));

                throw new Error(`Start game failed with status ${response.status} : ${error.message}`);
            }

            // state was deleted - fetch a new one to reset
            fetchState();
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
                            className={`num-btn ${isActive ? 'active' : ''} ${isLocked ? 'locked' : ''} ${isDimmed ? 'dimmed' : ''}`}
                            onClick={() => onNumberButtonClick(i)}>
                            {i + 1}
                        </button>
                    )
                })}
            </div>
            <div id="action-btns" className="flex row">
                {["edit", "erase", "hint", "pencil", "undo"].map((key, i) => {
                    if (key == "pencil" && pencilState) {
                        return <button key={i} className={`action-btn ${key} active`} onClick={() => onActionButtonClick(key)}></button>
                    }
                    return <button key={i} className={`action-btn ${key}`} onClick={() => onActionButtonClick(key)}></button>
                })}
            </div>
            <div className="paused-overlay flex column align-center justify-center">
                <h3>Game Paused</h3>
                <div className="flex">
                    <button className="btn-solid">Resume</button>
                    <button className="btn-solid accent" onClick={exitGame}>Exit to Menu</button>
                </div>
                <p id="paused-overlay-error" className="error">something went wrong</p>
            </div>
        </div>
    );
}