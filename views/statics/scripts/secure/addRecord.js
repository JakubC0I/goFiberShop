const form = document.getElementById("form")
const image = form.photos

form.addEventListener("submit", async (e) => {
    const description = form.description.value
    const name = form.name.value
    const price = form.price.value
    const producer = form.producer.value
    e.preventDefault()
    const res = await fetch("/addRecord", {
        method: "POST",
        body: JSON.stringify({name, price: parseFloat(price), description, producer, "images": files}),
        headers: {
            "Content-Type": "application/json"
        }
    })
    res.json().then((result) => {
        console.log(result);
    })
})

let files = [];
image.addEventListener('change', async (e) => {
    const fileList = form.photos.files
    const fileNames = [];
    for (let index = 0; index < fileList.length; index++) {
        const element = fileList[index];
        const reader = new FileReader();
        reader.onload = (async () => {
            files.push(reader.result)
        })
        reader.readAsDataURL(element)
        fileNames.push(element.name)
    };
}
)