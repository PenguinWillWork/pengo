import { PengoFetch } from "../../wailsjs/go/main/App";
import { pengo } from "../../wailsjs/go/models";

/**
 * Fetches favicon.ico, if not found - generates a placeholder icon taking the first letter of the hostname
 * @param src src for which we are resolving the icon
 * @returns void
 */
export async function resolveIcon(src: string) {
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
    const iconUrl = new URL("/favicon.ico", src).href;
    const iconResponse = await PengoFetch(iconUrl);
    applySiteIcon(iconResponse, iconImg, iconContainer);
  } catch (error) {
    console.log(error);
    generatePlaceholderIcon(src);
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

function generatePlaceholderIcon(currentLocation: string) {
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
