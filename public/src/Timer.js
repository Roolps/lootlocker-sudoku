import { useState, useEffect } from "react";

// GAME TIMER
export function Timer({ isPaused, initialStartTime }) {
    const [seconds, setSeconds] = useState(0);

    useEffect(() => {
        const now = Math.floor(Date.now() / 1000);
        setSeconds(now - initialStartTime);
    }, []);

    useEffect(() => {
        if (isPaused) return;

        const interval = setInterval(() => {
            setSeconds(prev => prev + 1);
        }, 1000);

        return () => clearInterval(interval);
    }, [isPaused]);

    function formatTime() {
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return <><span>{String(mins).padStart(2, "0")}</span>:<span>{String(secs).padStart(2, "0")}</span></>
    };

    return <p id="timer" className="flex">{formatTime()}</p>
}