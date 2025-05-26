import { useState } from "react";

export function Auth({ fetchState, setPlayerBalance }) {
    const [authFormType, setAuthFormType] = useState("login");

    async function auth(formdata) {
        const loader = document.getElementById("form-loader");
        if (loader) {
            loader.classList.remove("hidden");
        }

        try {
            const response = await fetch(`/api/${formdata.get("action")}`, {
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

            const data = await response.json();

            // player balance is returned from the login call
            setPlayerBalance(parseInt(data.data));
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
        <form className="flex column align-center" action={auth}>
            <div className="flex align-center space-between radio-wrap">
                <div className="backer" style={{ left: authFormType === "signup" ? ".3rem" : "50%" }}></div>
                <label style={{ marginRight: ".15rem" }} className={`radio-label ${authFormType === "signup" ? "active" : ""}`}>
                    Signup
                    <input
                        type="radio"
                        name="action"
                        value="signup"
                        checked={authFormType === "signup"}
                        onChange={() => setAuthFormType("signup")}
                    />
                </label>

                <label style={{ marginLeft: ".15rem" }} className={`radio-label ${authFormType === "login" ? "active" : ""}`}>
                    Login
                    <input
                        type="radio"
                        name="action"
                        value="login"
                        checked={authFormType === "login"}
                        onChange={() => setAuthFormType("login")}
                    />
                </label>
            </div>
            <label>
                Email<br></br>
                <input className="input-field" name="email" type="text" placeholder="-" />
            </label>
            <label>
                Password<br></br>
                <input className="input-field" name="password" type="password" placeholder="-" />
            </label>
            <button className="btn-solid" type="submit">Submit</button>
            <p id="login-error" className="error">something went wrong</p>
            <div id="form-loader"></div>
        </form>)
}