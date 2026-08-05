import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import { pengo } from "../wailsjs/go/models";
import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";
import { resolveIcon } from "./services/icon.resolver";

// document.querySelector('#app').innerHTML = `<div class="hello">Hello World</div>`;
let currentLocation: string;
const appBody = document.querySelector(".content");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  currentLocation =
    document.querySelector<HTMLInputElement>(".search-bar-input").value;
  await appBodyFetch();
});

async function appBodyFetch() {
  const loadingSpinner = document.querySelector(".search-bar-spinner");
  const uri = currentLocation;
  resolveIcon(currentLocation);
  try {
    if (uri.includes("http")) {
      appBody.innerHTML = wrongProtocol;
      return;
    }
    if (loadingSpinner) loadingSpinner.removeAttribute("hidden");
    const response = await PengoFetch(uri);
    if (response.ContentType.includes("text")) {
      const binary = atob(response.Body as unknown as string);
      const bytes = new Uint8Array(binary.length);
      for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
      }
      appBody.innerHTML = new TextDecoder().decode(bytes);
      resolvePengoUrlImgs();
    }
    if (response.ContentType.includes("image")) {
      //temp solution to test images
      const dataUrl = `data:${response.ContentType};base64,${response.Body}`;
      appBody.innerHTML = `<img src="${dataUrl}">`;
    }
  } catch (error) {
    console.log(error);
    appBody.innerHTML = connectionErrorPage;
  } finally {
    loadingSpinner.setAttribute("hidden", "");
  }
  document.querySelector<HTMLInputElement>(".search-bar-input").value = uri;
}

//Temp pengo:// img fetcher that converts response to a base64 src.
//Will be replaced later since having all images on the page in base64 is not the most efficient thing
async function resolvePengoUrlImgs() {
  const imagesArr = document.querySelectorAll<HTMLImageElement>("img");

  for (const img of imagesArr) {
    if (!img) continue;
    try {
      const response = await PengoFetch(img.getAttribute("src"));
      const dataUrl = `data:${response.ContentType};base64,${response.Body}`;

      img.setAttribute("src", `${dataUrl}`);
    } catch (error) {
      console.error(error);
    }
  }
}

appBody.addEventListener("click", async () => {
  const e = (event.target as Element).closest("a");
  event.preventDefault();

  if (e) {
    const pengoedHref = e.getAttribute("href");
    const finalHref = new URL(pengoedHref, currentLocation);
    currentLocation = finalHref.href;
    await appBodyFetch();
  }
});
