export function Generate({ difficultySelected }) {

    const difficulty = {
        "easy": 62,
        "medium": 53,
        "hard": 44,
        "very-hard": 35,
        "insane": 26,
        "inhuman": 17
    };


    function newBoard() {
        // set up empty board
        const grid = Array.from({ length: 9 }, () =>
            Array.from({ length: 9 }, () => ({}))
        );

        for (let box = 0; box < 3; box++) {
            const rowStart = box * 3;
            const colStart = box * 3;

            const numbers = shuffle()
            let idx = 0;
            for (let r = 0; r < 3; r++) {
                for (let c = 0; c < 3; c++) {
                    grid[rowStart + r][colStart + c] = { value: numbers[idx++] };
                }
            }
        }
        populateGrid(grid);
    }

    function populateGrid(grid) {
        const nextEmptyCell = findNextEmptyCell(grid);
        // if find next empty cell returns null then the grid is full
        if (!nextEmptyCell) return true;
        
        // set row and column and get a new shuffled array
        const { r, c } = nextEmptyCell;
        var numbers = shuffle();

        numbers.forEach((val) => {
            console.log(validCellValue(grid, r, c, val))
        })
    }

    function findNextEmptyCell(grid) {
        // loop through the grid and find the first square with no value in
        for (var r = 0; r < 9; r++) {
            for (var c = 0; c < 9; c++) {
                if (!grid[r][c].value) return { r, c };
            }
        }
        return null;
    }

    function validCellValue(grid, row, col, val) {
        // search down the column to check if the value is valid
        for (let c = 0; c < 9; c++) {
            if (grid[row][c].value === val) return false;
        }

        // search down the current row to check if the value is valid
        for (let r = 0; r < 9; r++) {
            if (grid[r][col].value === val) return false;
        }

        // search the current box to find out if the value is valid
        const boxRow = row - (row % 3);
        const boxCol = col - (col % 3);
        for (let r = 0; r < 3; r++) {
            for (let c = 0; c < 3; c++) {
                if (grid[boxRow + r][boxCol + c].value === val) return false;
            }
        }

        return true;
    }

    function shuffle() {
        const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9];

        for (let i = numbers.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [numbers[i], numbers[j]] = [numbers[j], numbers[i]]; // swap
        }

        return numbers;
    }

    return <>{newBoard()}</>
}