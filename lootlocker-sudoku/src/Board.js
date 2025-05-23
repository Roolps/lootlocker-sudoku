export function Board({ gameState, selection, onCellClick }) {

    let highlightValue = "row" in selection && "col" in selection && "value" in gameState[selection.row][selection.col] ? gameState[selection.row][selection.col].value : 0;

    let boardElements = [];

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

            if (selection.row === i) {
                classes += " cross";
            }

            if (selection.col === j) {
                classes += " cross";
            }

            if ('value' in cell) {
                value = cell.value;
                if (value === highlightValue) {
                    classes += " active";
                }
            } else {
                if (selection.row === i && selection.col === j) {
                    classes += " active";
                }
            }

            // if ('pencil' in cell) {
            //     classes += " pencil";
            //     value = <PencilMark marks={cell.pencil}></PencilMark>;
            // }

            let cellName = `Cell-${i}-${j}`;
            rowElements.push(<button className={classes} key={cellName} onClick={() => onCellClick(i, j)}>{value}</button>);
        });
        const rowName = `Row-${i}`
        boardElements.push(<div className="flex row" key={rowName}>{rowElements}</div>)
    });

    console.log(highlightValue)

    return (
        <div className="flex column align-center">
            <div className="cells flex column">{boardElements}</div>
            <div id="num-btns" className="flex row">
                {[...Array(9)].map((_, i) => (
                    <button key={i} className={`num-btn ${i + 1 === highlightValue ? 'active' : ''}`}>
                        {i + 1}
                    </button>
                ))}
            </div>
            <div id="action-btns" className="flex row">
                {["edit", "erase", "hint", "pencil", "undo"].map((key, i) => (
                    <button key={i} className={`action-btn ${key}`}></button>
                ))}
            </div>
        </div>
    );
}