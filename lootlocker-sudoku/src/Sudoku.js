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
  const [selection, setSelection] = useState({});
  
  
  function handleClick(i, j) {
      setSelection({ row: i, col: j });
  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">        
        <h1>Lootlocker Sudoku</h1>
        
        <Board gameState={sampleGame} selection={selection} onCellClick={handleClick} />

        <p id="powered-by-lootlocker">Powered by Lootlocker's free plan</p>
      </div>
    </div>
  );
}