import "./style.css";
import { Fetch } from "../wailsjs/go/main/App";

// document.querySelector('#app').innerHTML = `<div class="hello">Hello World</div>`;
const appBody = document.querySelector(".content");
document.querySelector(".search-bar-go").addEventListener("click", async () => {
  const uri = document.querySelector(".search-bar-input").value;
  const response = await Fetch(uri);
  appBody.innerHTML = response.Body;
});
