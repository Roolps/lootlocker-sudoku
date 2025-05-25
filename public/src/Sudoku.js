import { Board } from "./Board"
import { Timer } from "./Timer"
import { Auth } from "./Auth"
import { Menu } from "./Menu"

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
      const response = await fetch("/api/state", {
        method: "GET",
        headers: { "Accept": "application/json" },
      });

      if (!response.ok) {
        // set user authentication status to not logged in
        if (response.status === 403) {
          setAuthenticated(false);

          setTimeout(() => {
            const loader = document.getElementById("form-loader");
            if (loader) {
              loader.classList.add("hidden");
            }
          }, 500);
        }

        const error = await response.json().catch(() => ({}));
        throw new Error(`Get State failed with status ${response.status} : ${error.message}`);
      }

      const data = await response.json();

      // user is logged in because it succeeded so
      // enforce that log in state is true
      setAuthenticated(true);
      setGameState(data["data"]);

    } catch (err) {
      console.error(err.message);
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
            {gameState && gameState.length > 0 ? (
              <>
                <div id="timer-row" className="flex row space-between">
                  <p id="player-tokens" className="flex align-center">0</p>
                  <Timer />
                </div>
                <Board
                  fetchState={fetchState}
                  gameState={gameState}
                  pencilState={pencilState}
                  selection={selection}
                  onCellClick={handleClick}
                  onNumberButtonClick={handleNumberClick}
                  onActionButtonClick={handleActionButton}
                />
              </>

            ) : gameState && gameState.length === 0 ? (
              <Menu
                fetchState={fetchState}
              />

            ) : (
              <div className="error-message">An error occurred loading the game board.</div>
            )}
          </>
          ) : (
            <Auth fetchState={fetchState} />
          )}

        <a id="powered-by-lootlocker" href="https://lootlocker.com/" target="_blank">Made using <span>Lootlocker</span></a>
      </div>
    </div>
  );
}