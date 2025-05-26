import { Board } from "./Board"
import { Auth } from "./Auth"
import { Menu } from "./Menu"

import { useState, useEffect } from 'react';

export default function Sudoku() {
  const [gameState, setGameState] = useState(null);
  const [pencilState, setPencilMarks] = useState(false);

  const [authenticated, setAuthenticated] = useState(false);

  const [playerState, setPlayerState] = useState({})
  const [initialStartTime, setInitialStartTime] = useState(0)

  useEffect(() => {
    fetchState();
  }, []);

  useEffect(() => {
    if (authenticated) {
      fetchPlayer();
    }
  }, [authenticated]);

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

      // moved this to only if false - this will stop calling the user's balance multiple times
      if (authenticated == false) {
        setAuthenticated(true);
      }
      setGameState(data.data.state);
      setInitialStartTime(data.data.start_time);

    } catch (err) {
      console.error(err.message);
    }
  }

  async function fetchPlayer() {
    try {
      const response = await fetch("/api/player", {
        method: "GET",
        headers: { "Accept": "application/json" },
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(`Get State failed with status ${response.status} : ${error.message}`);
      }

      const data = await response.json();

      setPlayerState(data["data"]);
    } catch (err) {
      console.error(err.message);
    }
  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">
        <h1>Lootlocker Sudoku</h1>

        {authenticated ?
          (<>
            {gameState && gameState.length > 0 ? (
              <Board
                setGameState={setGameState}
                setPencilMarks={setPencilMarks}
                fetchPlayer={fetchPlayer}

                gameState={gameState}
                pencilState={pencilState}

                playerBalance={playerState.balance}
                initialStartTime={initialStartTime}
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
            <Auth
              fetchState={fetchState}
            />
          )}

        <a id="powered-by-lootlocker" href="https://lootlocker.com/" target="_blank">Made using <span>Lootlocker</span></a>
      </div>
    </div>
  );
}