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
        console.log(JSON.stringify(grid))
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