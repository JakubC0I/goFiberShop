let product = window.localStorage.getItem("productIds")
const items = document.getElementById("items")
let products = JSON.parse(product)
let summary = 0
let host = "http://localhost:3000/product"

products.products.forEach(element => {
    items.innerHTML += `
    <div id="${element.product}">
    <b><a href="${host}/${element.product}">${element.productName}</a></b><br>
    <form id="form">
    <input type="button" name="minus" value="-">
    <input type="number" name="number" value="${element.quantity}" max="${element.qMax}">
    <input type="button" name="plus" value="+"><br>
    </form>
    <p id="${element.product}_price">Price:${(element.price * element.quantity).toFixed(2)} (${element.price} each)</p><br>
    <img src="${element.img}">
    </div>
    `
    // summary += element.price * element.quantity
});

let forms = document.querySelectorAll("form")
let nums = []
for (let index = 0; index < forms.length; index++) {
    const ele = forms[index];
    const form = ele
    const element = products.products[index]
    element.quantity = form.number.value
    nums.push(parseFloat((element.price * element.quantity).toFixed(2)))
    form.addEventListener("submit", (e) => {
        if (parseInt(form.number.value) > parseInt(form.number.max)) {
        } else {
            e.preventDefault()
            element.quantity = parseInt(form.number.value)
            //SPRAWDZIC NA WIEKSZJE ILOSCI ELEMENTOW
            nums[index] = parseFloat((element.price * element.quantity).toFixed(2))
            document.getElementById(`${element.product}_price`).innerText = `Price:${(element.price * element.quantity).toFixed(2)} (${element.price} each)`
            products.products.forEach((elem) => {
                if (elem.product === element.product) {
                    elem.quantity = parseInt(element.quantity)
                    products.products.splice(index, 1, elem)
                    // products.products.push(elem)
                    window.localStorage.setItem("productIds", JSON.stringify(products))
                }
            })
            sum()
        }
    })
    //Dodawanie i odejmowanie itemów nie działa
    form.minus.addEventListener("click", (e) => { })
    form.plus.addEventListener("click", (e) => { })
    
}
function sum() {
    summary = 0
    nums.forEach((v) => {
        summary += v
    })
    document.getElementById("summary").innerText = `Summary: ${summary.toFixed(2)}`
}
sum()

const deliv = document.getElementById("deliv")
deliv.addEventListener("click", (e) => {
    e.preventDefault()
    window.location.assign("http://localhost:3000/deliv")
})