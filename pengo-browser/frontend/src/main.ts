import "./style.css";
import { PengoFetch } from "../wailsjs/go/main/App";
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
  updateSeachbarUrl(src);
}

// Address-bar and tab navigation: points the iframe at /pengo/…, which the
// middleware intercepts and translates into a Pengo protocol request.
// Links clicked inside the page bypass this and hit the handler directly.
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
    await updateSeachbarUrl(currentUrl);
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

//Falls back to built-in error pages in case something is obviosly not right, e.g http, wrong format of the url etc.
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

//Updates search bar url, for example on opening a hyperlink
function updateSeachbarUrl(url: string) {
  currentUrl = url;
  document.querySelector<HTMLInputElement>(".search-bar-input").value = url;
}

//iframe eats /pengo/[something], user enters a pengo:// url, this function converts it so pengo handler can work with it
function convertUrlForPengoHandler(url: string) {
  const pengoProtocolConverted = url.replace("pengo://", "/pengo/");
  return pengoProtocolConverted;
}

//Creates tab, sets it as active, displays it ("new tab" page by default)
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

//Attach event listeners for tab click and tab "x" button click
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

//Builds a tab element from a template in the tab nav area
function buildTabElement(): HTMLElement {
  const element = tabTemplate.content.firstElementChild.cloneNode(
    true,
  ) as HTMLElement;
  attachEventListenersToTab(element);
  return element;
}

//On create appends the id to the tab-id attribute of the new tab. Id is generated in the tab class on init
function assignNewTabId(tabStrip: Element) {
  const tabElements = tabStrip.querySelectorAll(".tab");
  for (const tabElement of tabElements) {
    if (!tabElement.getAttribute("tab-id")) {
      tabElement.setAttribute("tab-id", activeTab.id);
    }
  }
}

//Sets current url to activetab src and fetches it
//In the future: pages that were already fetched have to be cached
function displayActiveTab() {
  currentUrl = activeTab.getSrc();
  fetchPage();
}

//Destroys a tab and handles the switch afterwards. Triggers on "x"
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

//Updates tab data: title, src, icon in the future, then updates active tabs
function updateTabData(updateTabData?: ITabCreate) {
  activeTab.title = updateTabData.title;
  activeTab.src = updateTabData.src;
  document.querySelector(".active-tab-title").textContent = updateTabData.title;
  updateActiveTab();
}

//Updates active tab in case there was a switch, a tab close or other changes requiring to set and re-render active
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
