# AGENTS.md

Phigros score lookup service: Go (gin + gorm) backend, Vue 3 + Element Plus frontend. Repo dir is `phi-rank-query` but the Go module is `github.com/xeonds/phi-plug-go` — imports use the latter.

## Run / Build

- Runtime cwd is `build/` (config + dist live there). `make run` does `cd build && ./phi-rank-query-linux-amd64`.
- `make` default target (`all`) = `linux-amd64 web`. `make -j` builds both in parallel.
- Frontend: `make web` = `cd web && pnpm i && pnpm run build --outDir=../build/dist --emptyOutDir`.
- Dev server: `cd web && pnpm dev` (Vite proxies `/api/v1` → `http://localhost:8542`).
- Backend serves `./dist` via `NoRoute`; `config.Data.Difficulty`/`Info` must point at `./dist/*.tsv`.

## Game data extraction (Go, replaces the old Python scripts)

- `make fetch` → `go run ./cmd/unpack fetch` downloads the latest Phigros APK from TapTap to `build/base.apk`.
- `make unpack` → `go run ./cmd/unpack run build/base.apk web/public` regenerates `web/public/{difficulty,info}.tsv`, `web/public/config.json` (game version + update date) and `web/public/assets/illustrations/*.png`.
- `go run ./cmd/unpack version build/base.apk` prints the game version (from the binary `AndroidManifest.xml`).
- The tool lives in `cmd/unpack`; there is no Python/UnityPy dependency anymore.
- `build/config.yaml` holds real config; `lib.LoadConfig`-equivalent is `config.Load`, which auto-creates a template and `log.Fatal`s if missing.

## Backend

- `main.go`: routes. `GET /api/v1/rank_table` → `[{id,title,EZ,HD,IN,AT,Legacy}]`; `GET /api/v1/leaderboard` → `[{id,username,rks}]`; `POST /api/v1/player` `{session}` → `{player,rks,challenge,date,records,phi}`.
- `service.Client` (LeanCloud + AES save decryption + rks). `Query` returns all records (sorted desc by rks) and `phi` = top 3 charts with 100% acc. LeanCloud endpoints/keys are hardcoded; save files are AES-CBC encrypted zips read via `lib.ByteReader`.
- `service.UpdateUser`/`Leaderboard` keep a GORM `User` leaderboard. `model.Record` carries the chart result (json lowerCamel).
- `service.RankTable` merges `difficulty.tsv` + `info.tsv` at startup into `[{id,title,...}]`.
- `lib`: `LoadTSV`, `Logger` middleware, `NoRoute`, `NewDB`, `ByteReader`. `config`: viper `Load`. Note the config key is `databaseconfig` (see `build/config.yaml`).
- DB: sqlite `data.db` by default; `build/*.db`, `build/log.txt`, `build/users.csv` are runtime artifacts (gitignored `build/`).

## Frontend

- Vue 3 `<script setup>` + TypeScript, Element Plus, hash router (`src/router.ts`).
- `src/api.ts` (typed API client), `src/types.ts`, `src/utils.ts` (rating / challenge decode / rks formula), `src/sessions.ts` (localStorage session + per-session query history, with migration from the old `sessions`/`aliases` format).
- `components/SongCard.vue` (one chart card), `components/PlayerPanel.vue` (query + header + grid; `all` prop switches B27 vs all-records). Views: `HomeView`, `RankView`, `B19View`, `BNView`, `CalcView`, `HistoryView`, `LeaderboardView`, `SessionView`.
- `vue-tsc` runs with `strict` + `noUnusedLocals` + `noUnusedParameters` — unused vars break `pnpm run build`. No lazy-load plugin; images use native `loading="lazy"`.

## Testing

- `go test ./...` compiles all packages (currently no test files).
- Always run `go build ./...` + `go vet ./...` and `cd web && pnpm run build` after changes.

## Current Phigros data layout (v4+)

- Game info is in `assets/bin/Data/data.unity3d` (a UnityFS bundle) under `level0`; illustrations are ~2691 Addressables bundles in `assets/aa/Android/`.
- Catalog keys look like `<songId>.0/IllustrationLowRes.jpg`, bundle value `<hash>_<name>.bundle`.
- Illustration Texture2D is uncompressed RGB24 (`m_TextureFormat:3`) streamed in the bundle's `.resS` file — no ETC2/ASTC decoder needed.
- `internal/unity` deps: `github.com/pierrec/lz4/v4` (block decompress), `github.com/ulikunitz/xz` (raw LZMA1); `internal/phigros` uses `github.com/shogo82148/androidbinary` for the manifest.

## Deploy

- `make deploy` = `docker-compose up -d`; mounts `./build` to `/app`, external `services` network, port 8542.