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

    setTimeout(() => {
      const loader = document.getElementById("form-loader");
      if (loader) {
        loader.classList.add("hidden");
      }
    }, 500);
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

  async function login(formdata) {
    const loader = document.getElementById("form-loader");
    if (loader) {
      loader.classList.remove("hidden");
    }

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: formdata.get("email"),
          password: formdata.get("password")
        })
      });

      if (!response.ok) {
        // hide the loader
        if (loader) loader.classList.add("hidden");

        const error = await response.json().catch(() => ({}));
        throw new Error(`Login failed with status ${response.status} : ${error.message}`);
      }

      // login was successful - fetch state :)
      fetchState();

    } catch (err) {
      // hide the loader
      if (loader) loader.classList.add("hidden");
      console.error(err.message);


      const errorMsg = document.getElementById("login-error")
      errorMsg.innerHTML = err.message
      errorMsg.classList.add("active")

      setTimeout(() => {
        errorMsg.classList.remove("active")
      }, 5000);
    }
  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">
        <h1>Lootlocker Sudoku</h1>

        {authenticated ?
          (<>
            <div id="timer-row" className="flex row space-between">
              <p id="player-tokens" className="flex align-center">0</p>
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
            <form className="flex column align-center" action={login}>
              <label>
                Email<br></br>
                <input className="input-field" name="email" type="text" placeholder="-" />
              </label>
              <label>
                Password<br></br>
                <input className="input-field" name="password" type="password" placeholder="-" />
              </label>
              <button className="btn-solid" type="submit">Submit</button>
              <p id="login-error">something went wrong</p>
              <div id="form-loader"></div>
            </form>
          )}

        <a id="powered-by-lootlocker" href="https://lootlocker.com/" target="_blank">Made using <span>Lootlocker</span></a>
      </div>
    </div>
  );
}