import { Board } from "./Board"
import { useState } from "react";

export default function Sudoku() {
  const sampleGame = [[{}, { "value": 1, "immutable": true }, { "value": 2, "immutable": true }, {}, { "value": 5, "immutable": true }, {}, {}, {}, {}],
  [{ "value": 9, "immutable": true }, {}, {}, {}, {}, { "value": 3, "immutable": true }, {}, { "value": 1, "immutable": true }, { "value": 4, "immutable": true }],
  [{}, {}, {}, {}, {}, { "value": 8, "immutable": true }, { "value": 2, "immutable": true }, { "value": 3, "immutable": true }, {}],
  [{}, {}, { "value": 1, "immutable": true }, {}, {}, {}, {}, { "value": 2, "immutable": true }, { "value": 6, "immutable": true }],
  [{}, {}, {}, {}, {}, {}, {}, {}, {}],
  [{ "value": 6, "immutable": true }, { "value": 5, "immutable": true }, {}, {}, {}, {}, { "value": 7, "immutable": true }, {}, {}],
  [{}, { "value": 9, "immutable": true }, { "value": 7, "immutable": true }, { "value": 6, "immutable": true }, {}, {}, {}, {}, {}],
  [{ "value": 5, "immutable": true }, { "value": 4, "immutable": true }, {}, { "value": 2, "immutable": true }, {}, {}, {}, {}, { "value": 1, "immutable": true }],
  [{}, {}, {}, {}, { "value": 7, "immutable": true }, {}, { "value": 9, "immutable": true }, { "value": 5, "immutable": true }, {}]];

  const [gameState, setGameState] = useState(sampleGame);
  const [selection, setSelection] = useState({ row: 0, col: 0 });

  const [pencilMarks, setPencilMarks] = useState(false);


  function handleClick(i, j) {
    setSelection({ row: i, col: j });
  }

  function handleNumberClick(number) {
    let newGameState = gameState.map(row => row.slice());
    let cell = gameState[selection.row][selection.col];

    if ("value" in cell && cell.immutable) {
      return;
    }

    if (cell.value === number + 1) {
      newGameState[selection.row][selection.col] = {};
    } else {
      newGameState[selection.row][selection.col] = { value: number + 1, immutable: false };
    }
  
    setGameState(newGameState);
  }

  function handleActionButton(action) {
    switch (action) {
      case "edit":
        return handleEditAction()
      case "erase":
        return handleEraseAction()
      case "hint":
        return handleHintAction()
      case "pencil":
        return handlePencilAction()
      case "undo":
        return handleUndoAction()
    }
  }

  function handleEditAction() {

  }

  function handleEraseAction() {

  }

  function handleHintAction() {

  }

  function handlePencilAction() {
    // toggle the value of the 
    setPencilMarks(prev => !prev);
  }

  function handleUndoAction() {

  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">
        <h1>Lootlocker Sudoku</h1>

        <Board gameState={gameState} selection={selection} onCellClick={handleClick} onNumberButtonClick={handleNumberClick} onActionButtonClick={handleActionButton} />

        <a id="powered-by-lootlocker" href="https://lootlocker.com/" target="_blank">Made using <span>Lootlocker</span></a>
      </div>
    </div>
  );
}