import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
import { pengo } from "../wailsjs/go/models";
import { resolveIcon } from "./services/icon.resolver";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { ITabCreate, Tab } from "./services/tabs";

import connectionErrorPage from "./pages/connection-error.html?raw";
import wrongProtocol from "./pages/wrong-protocol.html?raw";
import tabMarkup from "./components/tab.html?raw";
import newTabPage from "./pages/new-tab.html?raw";

let currentUrl: string;
let activeTab: Tab;
const tabs: Tab[] = [];

const appBodyFrame = document.querySelector("iframe");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  currentUrl =
    document.querySelector<HTMLInputElement>(".search-bar-input").value;
  await fetchPage();
});

document.querySelector(".tab-new").addEventListener("click", () => {
  createTab();
});

const tabTemplate = document.createElement("template");
tabTemplate.innerHTML = tabMarkup;
interface PengoUrlCheck {
  valid: boolean;
  error: "http" | "malformed";
}
EventsOn("pengo:navigated", onNavigated);

function onNavigated(src: string, title: string) {
  updateTabData({ title, src });
  updateActiveTab();
  updateNavBarUrl(src);
}

async function fetchPage() {
  if (!currentUrl) {
    appBodyFrame.srcdoc = newTabPage;
    return;
  }

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
    updateTabData({ title: "Oops", src: "" });
    updateActiveTab();
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

function createTab() {
  const newTab = new Tab({ title: null, src: null });
  tabs.push(newTab);
  activeTab = newTab;

  const tabStrip = document.querySelector(".tab-strip");
  tabStrip.append(buildTabElement());

  assignNewTabId(tabStrip);
  updateTabData({ title: newTab.getTitle, src: null });
  updateActiveTab();
  displayActiveTab();
}

function attachEventListenersToTab(tab: HTMLElement) {
  tab.addEventListener("click", () => {
    const tabId = tab.getAttribute("tab-id");
    activeTab = tabs.find((t) => t.id === tabId);
    updateActiveTab();
    displayActiveTab();
  });
  tab.querySelector(".tab-close").addEventListener("click", () => {
    event.stopPropagation();
    destroyTab(tab.getAttribute("tab-id"));
  });
}

function buildTabElement(): HTMLElement {
  const element = tabTemplate.content.firstElementChild.cloneNode(
    true,
  ) as HTMLElement;
  attachEventListenersToTab(element);
  return element;
}

function assignNewTabId(tabStrip: Element) {
  const tabElements = tabStrip.querySelectorAll(".tab");
  for (const tabElement of tabElements) {
    if (!tabElement.getAttribute("tab-id")) {
      tabElement.setAttribute("tab-id", activeTab.id);
    }
  }
}

function displayActiveTab() {
  currentUrl = activeTab.getSrc();
  fetchPage();
}

function destroyTab(tabId: string) {
  const tabToDestroy = tabs.find((t) => t.id == tabId);
  const tabElements = document.querySelectorAll(".tab");
  for (const tabElement of tabElements) {
    if (tabElement.getAttribute("tab-id") == tabToDestroy.id) {
      tabElement.remove();
    }
  }
  const index = tabs.indexOf(tabToDestroy);
  if (index > -1) {
    tabs.splice(index, 1);
  }

  if (tabs.length > 0) {
    if (activeTab.id === tabToDestroy.id) {
      activeTab = tabs[index + 1 < tabs.length ? index : index - 1];
    }
    updateActiveTab();
    displayActiveTab();
  } else {
    createTab();
    return;
  }
}
function updateTabData(updateTabData?: ITabCreate) {
  activeTab.title = updateTabData.title;
  activeTab.src = updateTabData.src;
  document.querySelector(".active-tab-title").textContent = updateTabData.title;
  updateActiveTab();
}

async function updateActiveTab() {
  const tabElements = document.querySelectorAll(".tab");
  for (const tabElement of tabElements) {
    tabElement.classList.remove("tab--active");
    if (tabElement.getAttribute("tab-id") == activeTab.id) {
      tabElement.classList.add("tab--active");
      tabElement.querySelector(".tab-title").innerHTML = activeTab.getTitle;
      document.querySelector(".active-tab-title").textContent =
        activeTab.getTitle;
    }
  }
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
