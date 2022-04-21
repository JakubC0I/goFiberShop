const form = document.getElementById("form")

form.addEventListener("submit", async (e) => {
    const username = form.username.value
    const email = form.email.value
    const password = form.password.value
    e.preventDefault()
    const res = await fetch("/register", {
        method: "POST",
        body: JSON.stringify({email, username, password}),
        headers: {
            "Content-Type": "application/json"
        }
    })
    res.json().then((result) => {
        console.log(result);//adwad
    })
})