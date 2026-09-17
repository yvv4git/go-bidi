(() => {
  "use strict";

  const button = document.getElementById("fetch-button");
  if (button === null) {
    return;
  }

  button.addEventListener("click", async () => {
    const response = await fetch("/api/echo?q=hello");
    const data = await response.json();
    button.textContent = "got " + data.q;
  });
})();