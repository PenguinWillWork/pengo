import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import { pengo } from "../wailsjs/go/models";
import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";

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
  resolveIcon();
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

async function resolveIcon() {
  const iconContainer =
    document.querySelector<HTMLDivElement>(".search-bar-icon");
  iconContainer.classList.remove("search-bar-icon--loaded");

  const iconImg = document.querySelector<HTMLImageElement>(
    ".search-bar-icon-image",
  );
  if (!iconImg || !iconContainer) {
    return;
  }
  if (iconImg) {
    iconImg.hidden = true;
  }

  const iconLetter = iconContainer.querySelector<HTMLSpanElement>(
    ".search-bar-icon-letter",
  );
  if (iconLetter) {
    iconLetter.hidden = true;
  }

  try {
    const iconUrl = new URL("/favicon.ico", currentLocation).href;
    const iconResponse = await PengoFetch(iconUrl);
    applySiteIcon(iconResponse, iconImg, iconContainer);
  } catch (error) {
    console.log(error);
    generatePlaceholderIcon();
  }
}

function applySiteIcon(
  iconResponse: pengo.Response,
  iconImg: HTMLImageElement,
  iconContainer: HTMLDivElement,
) {
  iconImg.src = `data:${iconResponse.ContentType};base64,${iconResponse.Body}`;
  iconImg.hidden = false;
  iconContainer.classList.add("search-bar-icon--loaded");
}

function generatePlaceholderIcon() {
  const firstLetter = currentLocation.replace("pengo://", "").trim().charAt(0);
  if (!firstLetter) return;

  const iconContainer =
    document.querySelector<HTMLDivElement>(".search-bar-icon");

  const iconImg = document.querySelector<HTMLImageElement>(
    ".search-bar-icon-image",
  );
  if (!iconImg || !iconContainer) {
    return;
  }

  iconImg.hidden = true;

  let iconLetter = iconContainer.querySelector<HTMLSpanElement>(
    ".search-bar-icon-letter",
  );
  if (!iconLetter) {
    iconLetter = document.createElement("span");
    iconLetter.className = "search-bar-icon-letter";
    iconContainer.append(iconLetter);
  }

  iconLetter.hidden = false;
  iconLetter.textContent = firstLetter.toUpperCase();
  iconContainer.classList.add("search-bar-icon--loaded");
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
