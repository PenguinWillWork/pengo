## What Pengo is for

### Before

![before](assets/before.png)

### After

![after](assets/after.png)

## Structure

Web for penguins, written from scratch:

| Web for humans | Pengo for PENGUINS         |
| -------------- | -------------------------- |
| HTTP           | the Pengo protocol         |
| DNS            | Pengo DNS server (`:7007`) |
| Chrome         | pengo-browser              |

| Folder           | What                                                                   |
| ---------------- | ---------------------------------------------------------------------- |
| `protocol/`      | the protocol itself — the wire format and the client side of it        |
| `resolver/`      | DNS client — turns a name into an ip, expects a `:7007` dns server     |
| `cmd/server/`    | serves a folder (`-root`) on a port (`-port`, `:2719` by default)      |
| `cmd/dns/`       | DNS server that resolves names via `registry.json` on `:7007`          |
| `cmd/client/`    | CLI client — owns its own command-line syntax, not the protocol's      |
| `pengo-browser/` | the browser (Wails: Go backend + JS frontend)                          |
| `sites/`         | example sites — plain folders of files, nothing pengo-specific in them |

Everything under `cmd/` is runnable. The rest is a library.

**Important:** this project changes its core ideas and structure constantly. Code that looks junky now will be fixed, replaced, or quietly deleted like it never existed. There's no roadmap (currently), the author figures things out as he goes. Penguins don't mind as soon as pengo lets them connect

## How it works

`pengo-browser` is a Go app on Wails (similar to Tauri) that opens a window with a web page in it. The JS frontend is the address bar, the tab strip and the page area; the Go backend does the networking using solely Pengo protocol.

Tabs work now — each one keeps its own address, title and favicon. The title can't be read out of the page from JS (the iframe is sandboxed), so Go digs `<title>` out of the response and hands it to the frontend as an event.

Typing `pengo://welcome`:

```
frontend  → hands the address to the Go backend
backend   → splits into host (welcome) + path (/ if none given)
          → asks DNS :7007 where welcome is  → 127.0.0.1:2719
          → sends  PENGO/1.0 fetch welcome /
          ← response comes back
frontend  ← renders it
```

![pengo-browser](assets/pengo-browser-demo.gif)

Request format - request line, headers, blank line:

```
PENGO/1.0 fetch welcome /about
PebbleNonce: 123456

```

The request line is version, method, host, path — all four, always. Headers are optional and the blank line is what ends the request. Nothing reads those Pebble headers yet — that's [Pebble](#pebble).

Response format - headers, blank line, body:

```
PENGO/0.1
200 OKY
Content-Length:412
Content-Type:text/html

<div>...the page...</div>
```

Nothing new but when you realize it's the only internet for penguins you start to think it's really cool

## Methods

| Method   | What                                                   |
| -------- | ------------------------------------------------------ |
| `fetch`  | give me what's at this path                            |
| `submit` | here's some data — parsed and routed, does nothing yet |

Two layers decide separately, and it's worth keeping them apart. `protocol/` only knows whether a method is a word Pengo has at all — anything else is a malformed request. What a method _does_ is the server's call: `cmd/server` serves files out of a folder, so it answers `fetch` and has nothing to submit to.

`submit` also has nowhere to put a body yet — requests don't carry a `Content-Length`, only responses do. That's the next thing.

## Pebble

Identity, built into the protocol instead of into every site. WIP — the headers are reserved, nothing reads them yet.

### Why

Every site on the human web solves the same problem twice: a signup form nobody wants to fill in, and a password nobody wants to remember. The escape hatch is "sign in with Google", which works and costs you a company sitting between you and every site you visit.

Pengo doesn't have that problem to solve, because the request already says who sent it.

### How

You have an ed25519 keypair. The seed phrase _is_ the key — deriving one from the other is a pure function, so there's no stored secret to lose or steal.

A signed request carries four headers:

| Header            | What                 |
| ----------------- | -------------------- |
| `PebblePublicKey` | who's asking         |
| `PebbleTimestamp` | when                 |
| `PebbleNonce`     | issued by the server |
| `PebbleSignature` | the proof            |

### Why it's better

|                       | Human web                 | Pengo                  |
| --------------------- | ------------------------- | ---------------------- |
| Signing up            | a form, per site          | nothing                |
| A site's auth code    | sessions, hashing, resets | a pubkey column        |
| What a breach leaks   | the password database     | nothing worth stealing |
| Who's in the middle   | Google, Facebook, Apple   | nobody                 |
| Going to a second app | log in again, or OAuth    | already logged in      |

The last row is the one you actually feel. Two pengo apps written by strangers can store their own user data, but can both verify user based on the user's Pebble instantly.

Note: currently a plan

## CLI client

```bash
go run ./cmd/client
pengo> fetch pengo://welcome/about headers:[]
```

`headers:[]` is required even when empty, for now. The client also lowercases everything you type, so header values don't survive intact — which Pebble signatures will care about.

That syntax belongs to the client, not to Pengo — `protocol/` never sees the string you type. Swap it for curl-style flags tomorrow and nothing in `protocol/` moves.

## Serving files

A site is just a folder. You point the server at it with `-root` and it reads whatever gets asked for.

Paths are resolved by trying things in order:

```
/cat.png  → cat.png exists  → serve it
/about    → no about        → about.html exists → serve that
/         → /index          → index.html
nothing?  → 404.html, or a built-in one if the site has none
```

So pages don't need `.html` in the address, and files keep theirs. The client sends the path exactly as typed, the server does the guessing.

Bodies are bytes now instead of text, so images work.

The browser reads it and picks: html gets rendered, images get turned into a `data:` url. There's base64 in the middle because Wails talks to the frontend in JSON, and JSON can't hold raw bytes.

## vs HTTP

The whole "browser -> pengo protocol -> server" flow currently doesn't differ from http much. But.. Pengo doesn't really have TLDs:
`registry.json` currently is one flat file, name -> ip.

## What pengo doesn't replace (spoiler: JS is staying)

The web stack is the one piece Pengo deliberately doesn't replace.

For two reasons:

- **HTML/CSS/JS is a web standart** and every platform already ships an engine that renders them. Pengo is a separate **net**, not a different way to write a page. Anyone who can build a website can build a pengo site. Pengo is not trying to be different for the sake of being different.

But most importantly:

- **A rendering engine isn't part of a "net".** Layout, text shaping, fonts, the CSS, a JS engine. The hardest part of a browser and the part with nothing to
  do with the net. Writing it would eat the project and pengo would never get past the address bar.

P.S browser was migrated to TypeScript

## Running

A terminal each:

```bash
cp cmd/dns/registry.example.json cmd/dns/registry.json
cd cmd/dns && go run .                                  # DNS — reads registry.json from its own folder
go run ./cmd/server -root ./sites/welcome               # a site, on :2719
go run ./cmd/server -root ./sites/fishwrap -port 2721   # another one
cd pengo-browser && wails dev                           # browser
```

The server serves whatever `-root` points at, on `-port` (`.` and `2719` if you don't say). Any folder of files works — no config, nothing to register.

Then type `pengo://127.0.0.1`. For names instead of IPs, put them in `registry.json`:

```json
{ "welcome": "127.0.0.1:2719", "fishwrap": "127.0.0.1:2721" }
```
