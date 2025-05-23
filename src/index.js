import React from 'react';
import ReactDOM from 'react-dom/client';
import Board from './Sudoku';
import './assets/style.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Board />
  </React.StrictMode>
);