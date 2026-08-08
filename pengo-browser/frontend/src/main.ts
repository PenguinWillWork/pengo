import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import { pengo } from "../wailsjs/go/models";
import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";
import { resolveIcon } from "./services/icon.resolver";
import { EventsOn } from "../wailsjs/runtime/runtime";

let currentUrl: string;
const appBodyFrame = document.querySelector("iframe");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  currentUrl =
    document.querySelector<HTMLInputElement>(".search-bar-input").value;
  await fetchPage();
});

interface PengoUrlCheck {
  valid: boolean;
  error: "http" | "malformed";
}

EventsOn("pengo:navigated", showNavigatedUrl);

async function fetchPage() {
  const loadingSpinner = document.querySelector(".search-bar-spinner");
  resolveIcon(currentUrl);
  try {
    if (loadingSpinner) loadingSpinner.removeAttribute("hidden");
    await showNavigatedUrl(currentUrl);
  } catch (error) {
    console.log(error);
    appBodyFrame.srcdoc = connectionErrorPage;
  } finally {
    loadingSpinner.setAttribute("hidden", "");
  }
  document.querySelector<HTMLInputElement>(".search-bar-input").value =
    currentUrl;
}

function showNavigatedUrl(url: string) {
  const validatedUrl = validatePengoUrl(url);
  if (!validatedUrl.valid) {
    appBodyFrame.srcdoc =
      validatedUrl.error === "http" ? wrongProtocol : connectionErrorPage;
    throw new Error("malformed");
  }

  currentUrl = url;
  appBodyFrame.removeAttribute("srcdoc");
  appBodyFrame.src = convertUrlForPengoHandler(url);
  document.querySelector<HTMLInputElement>(".search-bar-input").value = url;
  resolveIcon(url);
}

function convertUrlForPengoHandler(url: string) {
  const pengoProtocolConverted = url.replace("pengo://", "/pengo/");
  return pengoProtocolConverted;
}

export function validatePengoUrl(input: string): PengoUrlCheck {
  let parsed: URL;
  try {
    parsed = new URL(input.trim());
  } catch {
    return { valid: false, error: "malformed" };
  }

  if (parsed.protocol === "http:" || parsed.protocol === "https:") {
    return { valid: false, error: "http" };
  }
  if (parsed.protocol !== "pengo:" || !parsed.host) {
    return { valid: false, error: "malformed" };
  }

  return {
    valid: true,
    error: null,
  };
}
