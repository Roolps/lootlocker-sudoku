import { Board } from "./Board"
import { Timer } from "./Timer"

import { useState, useEffect } from 'react';

export default function Sudoku() {
  const [gameState, setGameState] = useState(null);
  const [pencilState, setPencilMarks] = useState(false);

  const [authenticated, setAuthenticated] = useState(false);

  const [selection, setSelection] = useState({ row: 0, col: 0 });

  useEffect(() => {
    fetchState();
  }, []);

  async function fetchState() {
    try {
      const response = await fetch("/api/state");

      if (!response.ok) {
        if (response.status === 403) {
          setAuthenticated(false);
        }
        throw new Error(`Response status: ${response.status}`);
      }

      const data = await response.json();
      setAuthenticated(true);
      setGameState(data);

    } catch (error) {
      console.error(error.message);
    }
  }

  function handleClick(i, j) {
    setSelection({ row: i, col: j });
  }

  function handleNumberClick(number) {
    let newGameState = gameState.map(row => row.slice());
    let cell = gameState[selection.row][selection.col];

    if ("value" in cell && cell.immutable) {
      return;
    }

    if (pencilState) {
      delete cell.value

      // set the pencil state
      if (!cell.pencil || !Array.isArray(cell.pencil)) {
        cell.pencil = Array(9).fill(0);
      }
      if (cell.pencil[number]) {
        cell.pencil[number] = 0
      } else {
        cell.pencil[number] = 1
      }
    } else {
      if (cell.value === number + 1) {
        newGameState[selection.row][selection.col] = {};
      } else {
        newGameState[selection.row][selection.col] = { value: number + 1, immutable: false };
      }
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
    // toggle the value of the pencil mark state
    setPencilMarks(prev => !prev);
  }

  function handleUndoAction() {

  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">
        <h1>Lootlocker Sudoku</h1>

        {authenticated ?
          (<>
            <div id="timer-row" className="flex row space-between">
              <p id="player-tokens" className="flex align-center">1042</p>
              <Timer />
            </div>

            {gameState && (
              <Board
                gameState={gameState}
                pencilState={pencilState}
                selection={selection}
                onCellClick={handleClick}
                onNumberButtonClick={handleNumberClick}
                onActionButtonClick={handleActionButton}
              />
            )}
            {!gameState && <div className="error-message">An error occurred loading the game board.</div>}
          </>
          ) : (
            <form className="flex column align-center">
              <label>
                Email<br></br>
                <input className="input-field" type="text" placeholder="-" />
              </label>
              <label>
                Password<br></br>
                <input className="input-field" type="password" placeholder="-" />
              </label>
              <button className="btn-solid" type="submit">Submit</button>
            </form>
          )}

          <a id="powered-by-lootlocker" href="https://lootlocker.com/" target="_blank">Made using <span>Lootlocker</span></a>
      </div>
    </div>
  );
}