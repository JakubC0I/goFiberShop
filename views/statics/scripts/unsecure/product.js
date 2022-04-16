const btn = document.getElementById("commentButton")
const vCbtn = document.getElementById("viewComments")
const addCart = document.getElementById("cart")
const id = location.pathname.substr(8)
const field = document.getElementById("comments")
const img = document.getElementById("images").firstElementChild.src
const price = document.getElementById("price").innerText
const productName = document.getElementById("productName").innerText
btn.addEventListener("click", (e) => {
    e.preventDefault()
    document.getElementById("commentForm").innerHTML = `<form action="/addComment${id}" method="POST" id="form"><input type="text" name="body"><input type="submit"></form>`
    const form = document.getElementById("form")
    form.addEventListener("submit", async (e) => {
        const body = form.body.value
        e.preventDefault()
        const res = await fetch(`/addComment${id}`, {
            method: "POST",
            body: JSON.stringify({ body }),
            headers: {
                "Content-Type": "application/json"
            }
        })
        res.json().then((result) => {
            console.log(result.success);
        })
    })
})

vCbtn.addEventListener("click", async (e) => {
    e.preventDefault()
    //Jak prawidłowo umieścić resultaty? 
    field.innerHTML = ""
    const res = await fetch(`/viewComments${id}`, {
        method: "GET"
    })
    res.json().then((result) => {
        result.comments.forEach(element => {
            field.innerHTML += `<h4>${element.username}</h4><br><p>${element.body}</p><br><p>${element.created_at}</p>`
        });
        console.log(result);
    })
})

function create(prod, index) {
    prod.products.splice(index, 1)
    prod.products.push({ "product": id, "quantity": parseInt(quantity), img, "price": parseFloat(price), productName })
    console.log(JSON.stringify(prod));
    window.localStorage.setItem("productIds", JSON.stringify(prod))
}

addCart.addEventListener("click", (e) => {
    let quantity = parseInt(document.getElementById("quantity").value)
    const q = document.getElementById("quantity")
    if (quantity > parseInt(q.max) || quantity < parseInt(q.min)) {
    } else {
        e.preventDefault()
        const id = document.getElementById("oID").innerText
        const ids = []
        if (window.localStorage.getItem("productIds") == null) {
            ids.push(JSON.stringify({ "products": [{ "product": id, "quantity": parseInt(quantity), img, "price": parseFloat(price), productName }] }))
            window.localStorage.setItem("productIds", ids)
        } else {
            let products = window.localStorage.getItem("productIds")
            let prod = JSON.parse(products)
            let x = true
            for (let index = 0; index < prod.products.length; index++) {
                const element = prod.products[index];
                if (element.product == id) {
                    quantity = element.quantity + quantity
                    prod.products.splice(index, 1)
                    if (quantity > parseInt(q.max)) {
                        quantity = parseInt(q.max)
                        prod.products.push({ "product": id, "quantity": parseInt(quantity), img, "price": parseFloat(price), productName })
                        console.log(JSON.stringify(prod));
                        window.localStorage.setItem("productIds", JSON.stringify(prod))
                        x = false
                    } else {
                        prod.products.push({ "product": id, "quantity": parseInt(quantity), img, "price": parseFloat(price), productName })
                        console.log(JSON.stringify(prod));
                        window.localStorage.setItem("productIds", JSON.stringify(prod))
                        x = false
                    }
                }
            }
            if (x) {
                prod.products.push({ "product": id, "quantity": parseInt(quantity), img, "price": parseFloat(price), productName })
                console.log(JSON.stringify(prod));
                window.localStorage.setItem("productIds", JSON.stringify(prod))
            }

        }
    }


})
