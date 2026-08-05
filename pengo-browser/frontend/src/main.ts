import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import { pengo } from "../wailsjs/go/models";
import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";
import { resolveIcon } from "./services/icon.resolver";

let currentUrl: string;
const appBodyFrame = document.querySelector("iframe");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  currentUrl =
    document.querySelector<HTMLInputElement>(".search-bar-input").value;
  await fetchPage();
});

function decodeHtml(response: pengo.Response): Uint8Array<ArrayBuffer> {
  const binary = atob(response.Body as unknown as string);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

function renderAsText(html: string) {
  appBodyFrame.srcdoc = html;
}

function renderAsImage(response: pengo.Response) {
  const dataUrl = `data:${response.ContentType};base64,${response.Body}`;
  appBodyFrame.srcdoc = `<img src="${dataUrl}">`;
}

async function resolveRequest(uri: string) {
  const response = await PengoFetch(uri);

  let parsedHtml = Document.parseHTMLUnsafe(
    atob(response.Body as unknown as string),
  );

  parsedHtml = await resolvePengoUrlImgs(parsedHtml);
  // console.log(parsedHtml);
  if (response.ContentType.includes("text")) {
    renderAsText(parsedHtml.documentElement.outerHTML);
  }
  if (response.ContentType.includes("image")) {
    renderAsImage(response);
  }
}

async function fetchPage() {
  const loadingSpinner = document.querySelector(".search-bar-spinner");
  resolveIcon(currentUrl);
  try {
    if (currentUrl.includes("http")) {
      appBodyFrame.srcdoc = wrongProtocol;
      return;
    }
    if (loadingSpinner) loadingSpinner.removeAttribute("hidden");
    await resolveRequest(currentUrl);
  } catch (error) {
    console.log(error);
    appBodyFrame.srcdoc = connectionErrorPage;
  } finally {
    loadingSpinner.setAttribute("hidden", "");
  }
  document.querySelector<HTMLInputElement>(".search-bar-input").value =
    currentUrl;
}

//Temp pengo:// img fetcher that converts response to a base64 src.
//Will be replaced later since having all images on the page in base64 is not the most efficient thing
async function resolvePengoUrlImgs(document: Document): Promise<Document> {
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
  return document;
}

appBodyFrame.addEventListener("click", async () => {
  const e = (event.target as Element).closest("a");
  event.preventDefault();

  if (e) {
    const pengoedHref = e.getAttribute("href");
    const finalHref = new URL(pengoedHref, currentUrl);
    currentUrl = finalHref.href;
    await fetchPage();
  }
});
