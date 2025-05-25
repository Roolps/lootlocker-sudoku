import { Board } from "./Board"
import { Auth } from "./Auth"
import { Menu } from "./Menu"

import { useState, useEffect } from 'react';

export default function Sudoku() {
  const [gameState, setGameState] = useState(null);
  const [pencilState, setPencilMarks] = useState(false);

  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    fetchState();
  }, []);

  useEffect(() => {
    if (Array.isArray(gameState)) {
      verifyResult();
    }
  }, [gameState]);

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
    console.log(allRowsValid && allColsValid);
  }

  function isValidGroup(group) {
    const sorted = [...group].sort((a, b) => a - b);
    const expected = [1, 2, 3, 4, 5, 6, 7, 8, 9];
    return JSON.stringify(sorted) === JSON.stringify(expected);
  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">
        <h1>Lootlocker Sudoku</h1>

        {authenticated ?
          (<>
            {gameState && gameState.length > 0 ? (
              <Board
                fetchState={fetchState}
                setGameState={setGameState}
                setPencilMarks={setPencilMarks}

                gameState={gameState}
                pencilState={pencilState}
              />

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