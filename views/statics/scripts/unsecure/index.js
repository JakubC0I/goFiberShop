const form = document.getElementById("form")

form.addEventListener("submit", async (e) => {
    const searchbar = form.searchbar.value
    e.preventDefault()
    const res = await fetch("/", {
        method: "POST",
        body: JSON.stringify({searchbar}),
        headers: {
            "Content-Type": "application/json"
        }
    })
    res.json().then((result) => {
        console.log(result);
        document.getElementById("searchResults").innerHTML = ""
        document.getElementById("searchResults").innerHTML = `<p>${result}</p>`
    })
})