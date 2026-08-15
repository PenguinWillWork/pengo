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

EventsOn("pengo:navigated", updateNavBarUrl);

async function fetchPage() {
  const loadingSpinner = document.querySelector(".search-bar-spinner");
  resolveIcon(currentUrl);
  try {
    if (loadingSpinner) loadingSpinner.removeAttribute("hidden");
    const validatedUrl = validatePengoUrl(currentUrl);
    if (fallbackIfWrongUrl(validatedUrl)) {
      return;
    }
    await updateNavBarUrl(currentUrl);
    appBodyFrame.removeAttribute("srcdoc");
    appBodyFrame.src = convertUrlForPengoHandler(currentUrl);
  } catch (error) {
    console.log(error);
    appBodyFrame.srcdoc = connectionErrorPage;
  } finally {
    loadingSpinner.setAttribute("hidden", "");
  }
  document.querySelector<HTMLInputElement>(".search-bar-input").value =
    currentUrl;
}

function fallbackIfWrongUrl(validatedUrl: PengoUrlCheck): boolean {
  if (!validatedUrl.valid) {
    appBodyFrame.srcdoc =
      validatedUrl.error === "http" ? wrongProtocol : connectionErrorPage;
    return true;
  }
  return false;
}

function updateNavBarUrl(url: string) {
  currentUrl = url;
  document.querySelector<HTMLInputElement>(".search-bar-input").value = url;
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
