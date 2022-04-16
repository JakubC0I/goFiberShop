let product = window.localStorage.getItem("productIds")
const items = document.getElementById("items")
let products = JSON.parse(product)
let summary = 0

products.products.forEach(element => {
    items.innerHTML += `
    <div id="${element.product}">
    <h4>${element.productName}</h4><br>
    <p>Quantity: ${element.quantity}</p>
    <p>Price:${element.price * element.quantity} (${element.price} each)</p><br>
    <img src="${element.img}">
    </div>
    `
    summary += element.price * element.quantity
});

document.getElementById("summary").innerText = `Summary: ${summary}`