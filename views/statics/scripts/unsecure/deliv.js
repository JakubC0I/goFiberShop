const prod = window.localStorage.getItem("productIds")
const products = JSON.parse(prod).products
const form = document.getElementById("form")

let cost = 0
products.forEach(element => {
    let e = element.quantity * element.price
    cost += e
    console.log(e);
});
document.getElementById("cost").innerText = "Total cost: " + cost.toFixed(2)


form.addEventListener("submit", async (e) => {
    e.preventDefault()
    const address = {
        "city": form.city.value,
        "postcode": form.postcode.value,
        "street": form.street.value,
        "hNum": form.hNum.value
    }
    const delivery = form.delivery.value
    const res = await fetch("/deliver", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({products, address, delivery})
    })

})