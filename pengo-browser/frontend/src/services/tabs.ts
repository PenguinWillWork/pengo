import { validatePengoUrl } from "../main";
import { v4 as uuidv4 } from "uuid";

export interface ITabCreate {
  title: string | null;
  src: string | null;
}

export class Tab {
  constructor(tab: ITabCreate) {
    if (tab.title) {
      this._title = tab.title;
    } else {
      this._title = "New tab";
    }
    this._src = null;
    this.id = uuidv4();
  }

  readonly id: string;
  private _title: string | null;
  private _src: string | null;

  set title(title: string | null) {
    if (!title && this._src) {
      this._title = this._src;
      return;
    }
    if (!title && !this._src) {
      this._title = "New Tab";
    }
    this._title = title;
  }

  get getTitle() {
    let titleRes: string;
    // if (this.title) {
    titleRes = this._title;
    // }
    return titleRes;
  }

  set src(url: string) {
    const validatedUrl = validatePengoUrl(url);
    if (validatedUrl.error) {
      //...
    }
    this._src = url;
  }

  getSrc() {
    return this._src;
  }
}
