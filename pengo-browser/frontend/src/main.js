import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";

// document.querySelector('#app').innerHTML = `<div class="hello">Hello World</div>`;
let currentLocation;
const appBody = document.querySelector(".content");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  currentLocation = document.querySelector(".search-bar-input").value;
  await appBodyFetch();
});

async function appBodyFetch() {
  const uri = currentLocation;
  try {
    if (uri.includes("http")) {
      appBody.innerHTML = wrongProtocol;
      return;
    }
    const response = await PengoFetch(uri);

    if (response.ContentType == ".html" || response.ContentType == "text") {
      const binary = atob(response.Body);
      const bytes = new Uint8Array(binary.length);
      for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
      }
      appBody.innerHTML = new TextDecoder().decode(bytes);
    }
    if (response.ContentType == ".png") {
      console.log(response.ContentType);
      //temp solution to test images
      const dataUrl = `data:image/png;base64,${response.Body}`;
      appBody.innerHTML = `<img src="${dataUrl}">`;
    }
  } catch (error) {
    console.log(error);
    appBody.innerHTML = connectionErrorPage;
  }
  document.querySelector(".search-bar-input").value = uri;
}

appBody.addEventListener("click", async () => {
  const e = event.target.closest("a");
  event.preventDefault();

  if (e) {
    const pengoedHref = e.getAttribute("href");
    const finalHref = new URL(pengoedHref, currentLocation);
    currentLocation = finalHref.href;
    await appBodyFetch();
  }
});
