document.body.onload = function() {
    console.log("hello");
fetch("../matherial/9", { mode: "no-cors"})
  .then((response) => {
    return response.text();
  })
  .then((data) => {
    console.log(data);
  });

};

const myListL = document.getElementById("left");
const myListR = document.getElementById("right");

function make_list(url, node) {
  fetch(url, { mode: "no-cors"})
    .then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error, status = ${response.status}`);
        }
        return response.json();
    })
    .then((data) => {
        for (const product of data.value) {
            const listItem = document.createElement("li");
            const nameElement = document.createElement("strong");
            nameElement.textContent = product.name;

            const priceElement = document.createElement("strong");
            priceElement.textContent = `${product.cost} грн.`;
            priceElement.classList.add("cost");

            listItem.append(
                nameElement,
                ` ${product.matherial_group} `,
                priceElement
            );
            node.appendChild(listItem);
        }
    })
    .catch(
        (error) => {
            const p = document.createElement("p");
            p.appendChild(document.createTextNode(`Error: ${error.message}`));
            document.body.insertBefore(p, node);
        }
    );
}

make_list("../w_matherial_get_all", myListL);
make_list("../w_product_get_all", myListR);
