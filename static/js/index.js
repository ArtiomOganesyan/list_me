const newContainer = document.querySelector(".new-container");
const newContainerForm = document.forms["new"];
const container = document.getElementById("container");
const listContainer = document.getElementById("list-container");
const uuidSpan = document.getElementById("uuid");
const shareButton = document.getElementById("share");
const addItemForm = document.forms["add-item-form"];

// ==========================================
// CHECK FOR LIST ID
// ==========================================

let localData = JSON.parse(localStorage.getItem("listmeta") || "{}");
let { listID, secret } = localData;
const urlParams = new URLSearchParams(window.location.search);

if (!listID) {
  listID = urlParams.get("id");
}
if (!secret) {
  secret = urlParams.get("secret");
}

console.log(listID, secret);

if (listID && secret) {
  localStorage.setItem("listmeta", JSON.stringify({ listID, secret }));

  fetch(`/api/list?id=${listID}&secret=${secret}`)
    .then((res) => res.json())
    .then((data) => {
      console.log(data);
      data.rows.forEach((item) => {
        listContainer.innerHTML += createNewListItem(
          item.id,
          item.title,
          item.description,
          item.done
        );
      });
    })
    .catch(console.error);

  newContainer.classList.add("hidden");
  container.classList.remove("hidden");
}

// ==========================================
// MOBILE USER CHECK
// ==========================================
const isMobile = window.matchMedia("only screen and (max-width: 10cm)").matches;

// ==========================================
// CREATE NEW LIST
// ==========================================
newContainerForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const data = {
    ...Object.fromEntries(new FormData(newContainerForm)),
    listID,
  };
  localStorage.setItem(
    "listmeta",
    JSON.stringify({ listID, secret: data.secret })
  );

  history.pushState(null, "", `?id=${listID}&secret=${data.secret}`);

  fetch("/api/list", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((res) => {
      console.log(res);
      container.classList.remove("hidden");
      newContainer.classList.add("hidden");
    })
    .catch(console.error);
});

// ==========================================
// SHARE LIST
// ==========================================
shareButton.addEventListener("click", (e) => {
  e.preventDefault();
  console.log("clicked");
  let data = JSON.parse(localStorage.getItem("listmeta") || "{}");
  if (!data.listID || !data.secret) {
    console.error("Contact the developer.");
    return;
  }
  history.pushState(null, "", `?id=${data.listID}&secret=${data.secret}`);
  navigator.clipboard.writeText(window.location.href);
});

// ==========================================
// ADD ITEM
// ==========================================

const createNewListItem = (id, title, desc, done = false) => {
  return `<div class="list-item ${done && "checked"}" data-id="${id}">
    <div class="list-item__row" data.id="${id}">
      <div class="list-item__move">:</div>
      <div class="list-item__title">${title}</div>
      <div class="list-item__check">
      <img src="/static/assets/check.svg" alt="check mark" width="25px">
      </div>
      <div class="list-item__remove">
      <img src="/static/assets/remove.svg" alt="remove mark" width="25px">
      </div>
    </div>
    <div class="list-item__description">${desc}</div>
  </div>`;
};

addItemForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const data = Object.fromEntries(new FormData(addItemForm));

  if (data.title === "") {
    // shake the form
    addItemForm.classList.add("shake");
    setTimeout(() => {
      addItemForm.classList.remove("shake");
    }, 500);

    return;
  }

  fetch(`/api/list/${listID}/row`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((res) => res.json())
    .then((id) => {
      if (id) {
        const newItem = createNewListItem(id, data.title, data.description);
        listContainer.innerHTML += newItem;
        addItemForm.reset();
      }
    })
    .catch(console.error);
});

// ==========================================
// WORKING WITH LIST CONTAINER
// ==========================================

listContainer.addEventListener("click", (e) => {
  e.preventDefault();

  // ==========================================
  // CHECK ITEM
  // ==========================================
  if (
    e.target.classList.contains("list-item__check") ||
    e.target.closest(".list-item__check")
  ) {
    const id = e.target.closest(".list-item").dataset.id;
    const currentState = !e.target
      .closest(".list-item")
      .classList.contains("checked");

    fetch(`/api/row/${id}?state=${currentState}`, {
      method: "PATCH",
    })
      .then((res) => {
        if (res.status === 200) {
          e.target.closest(".list-item").classList.toggle("checked");
        }
      })
      .catch(console.error);
  }

  // ==========================================
  // REMOVE ITEM
  // ==========================================
  if (
    e.target.classList.contains("list-item__remove") ||
    e.target.closest(".list-item__remove")
  ) {
    const rowID = e.target.closest(".list-item").dataset.id;
    fetch(`/api/row/${rowID}`, {
      method: "DELETE",
    })
      .then((res) => {
        if (res.status === 204) {
          e.target.closest(".list-item").remove();
        }
      })
      .catch(console.error);

    e.target.closest(".list-item").remove();
  }
});
