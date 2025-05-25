export function ActionButtons({ setPencilMarks, pencilState }) {
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
        <div id="action-btns" className="flex row">
            {["edit", "erase", "hint", "pencil", "undo"].map((key, i) => {
                if (key == "pencil" && pencilState) {
                    return <button key={i} className={`action-btn ${key} active`} onClick={() => handleActionButton(key)}></button>
                }
                return <button key={i} className={`action-btn ${key}`} onClick={() => handleActionButton(key)}></button>
            })}
        </div>
    )
}