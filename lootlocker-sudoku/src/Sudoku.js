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
  }

  return (
    <div className="container flex row align-center justify-center">
      <div className="board flex column align-center">        
        <h1>Roolps Sudoku</h1>
        <Board gameState={sampleGame} selection={selection} onCellClick={handleClick} />

        <div id="num-btns" className="flex row">
          {[...Array(9)].map((_, i) => (
            <button key={i} className={`num-btn ${i === 0 ? 'special' : ''}`}>
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
    </div>
  );
}