function sortable() {
  const sortables = document.querySelectorAll(".sortable");
  for (let i = 0; i < sortables.length; i++) {
    const sortable = sortables[i];
    if (sortable.classList.contains("kanban")) {
      new Sortable(sortable, {
        animation: 150,
        onEnd: function () {
          const columns = document.querySelectorAll(".column");
          const order = Array.from(columns).map((el) =>
            Number(el.lastElementChild.value)
          );

          const xhr = new XMLHttpRequest();
          xhr.open("POST", "/columns");
          xhr.setRequestHeader("Content-type", "application/json");
          xhr.send(JSON.stringify({ columns: order }));
        },
      });
    } else {
      new Sortable(sortable, {
        group: "shared",
        animation: 150,
        onEnd: function (event) {
          const columnId = Number(
            Array.from(event.to.children).find((el) => el.localName === "input")
              .value
          );
          const items = Array.from(event.to.children)
            .filter((el) => el.localName !== "input")
            .map((el) => Number(el.firstElementChild.value));

          const xhr = new XMLHttpRequest();
          xhr.open("POST", "/items");
          xhr.setRequestHeader("Content-type", "application/json");
          xhr.send(JSON.stringify({ columnId: columnId, items: items }));
        },
      });
    }
  }
}

htmx.on("htmx:load", sortable);
htmx.on("htmx:afterRequest", sortable);
