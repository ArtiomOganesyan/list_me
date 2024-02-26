const newContainer = document.querySelector(".new-container");
const newContainerForm = document.forms["new"];
const container = document.getElementById("container");
const listContainer = document.getElementById("list-container");
const uuidSpan = document.getElementById("uuid");
const shareButton = document.getElementById("share");
const addItemForm = document.forms["add-item-form"];
const listLabel = document.getElementById("label-show-input-container");

// ==========================================
// CHECK FOR LIST ID
// ==========================================
const urlParams = new URLSearchParams(window.location.search);
id = urlParams.get("id");
secret = urlParams.get("secret");

let localData = JSON.parse(localStorage.getItem("listmeta") || "{}");

if (!id) {
  id = localData.id;
}
if (!secret) {
  secret = localData.secret;
}

if (id && secret) {
  localStorage.setItem("listmeta", JSON.stringify({ id, secret }));
  history.pushState(null, "", `?id=${id}&secret=${secret}`);

  fetch(`/api/list?id=${id}&secret=${secret}`)
    .then((res) => res.json())
    .then((data) => {
      if (data.error) {
        throw new Error(data.error);
      }
      data.rows?.forEach((item) => {
        listContainer.innerHTML += createNewListItem(
          item.id,
          item.title,
          item.description,
          item.done
        );
      });
      newContainer.classList.add("hidden");
      container.classList.remove("hidden");
      listLabel.classList.remove("hidden");
    })
    .catch((err) => {
      console.error(err);
      id = null;
      secret = null;
    })
    .then(() => {
      if (!id) {
        id = crypto.randomUUID();
        uuidSpan.innerText = id;
      }
    });
}

if (!id) {
  console.log("No id found");
  id = crypto.randomUUID();
  uuidSpan.innerText = id;
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
    id,
  };

  if (!data.id || !data.secret) {
    newContainerForm.classList.add("shake");
    setTimeout(() => {
      newContainerForm.classList.remove("shake");
    }, 500);
    return;
  }

  localStorage.setItem("listmeta", JSON.stringify({ id, secret: data.secret }));

  history.pushState(null, "", `?id=${id}&secret=${data.secret}`);

  fetch("/api/list", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((res) => {
      container.classList.remove("hidden");
      newContainer.classList.add("hidden");
      listLabel.classList.remove("hidden");
    })
    .catch(console.error);
});

// ==========================================
// SHARE LIST
// ==========================================
shareButton.addEventListener("click", (e) => {
  e.preventDefault();
  let data = JSON.parse(localStorage.getItem("listmeta") || "{}");
  if (!data.id || !data.secret) {
    console.error("Contact the developer.");
    return;
  }
  history.pushState(null, "", `?id=${data.id}&secret=${data.secret}`);
  navigator.clipboard.writeText(window.location.href);
});

// ==========================================
// ADD ITEM
// ==========================================

const createNewListItem = (id, title, desc, done = false) => {
  return `<div class="list-item ${done && "checked"}" data-rowid="${id}">
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

  fetch(`/api/list/${id}/row`, {
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
    const rowId = e.target.closest(".list-item").dataset.rowid;
    const currentState = !e.target
      .closest(".list-item")
      .classList.contains("checked");

    fetch(`/api/row/${rowId}?state=${currentState}`, {
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
    const rowID = e.target.closest(".list-item").dataset.rowid;
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
