const form = document.getElementById("form")

form.addEventListener("submit", async (e) => {
    const email = form.email.value
    const password = form.password.value
    e.preventDefault()
    const res = await fetch("/login", {
        method: "POST",
        body: JSON.stringify({email, password}),
        headers: {
            "Content-Type": "application/json"
        }
    })
    res.json().then((result) => {
        console.log(result.success);
        if (result.success) {
            window.location.replace("/")
        }
    })
})