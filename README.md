# Axle Weight Planner

Truck axle weight calculator for steer, drive and trailer groups. It handles 3 and 4 axle tractors, jeep dollies, lowboys with boosters, step decks, vans and flatbeds, lift axles, sliding 5th wheel and tandem, and legal-limit checks (federal bridge formula or custom per-axle limits). Version 1.2.

## Ways to use it

- **Any browser:** open `index.html`.
- **Phone, works offline:** turn on GitHub Pages for this repo (Settings, Pages, deploy from branch `main`, folder `/ (root)`). Open the page address on the phone and choose "Add to Home Screen" (Safari) or "Install app" (Chrome). After the first visit it works without a signal.
- **Windows or Linux program:** `releases/AxleWeightPlanner.exe` (double-click) or `releases/axle-planner-linux` (run `sh releases/run-linux.sh`). Saved setups are kept in `setups.json` next to the program, so one thumb drive can serve both systems.

## Files

- `index.html` is the whole app in one file.
- `manifest.webmanifest`, `sw.js` and the `icon*` files make it installable and available offline.
- `program/` is the Go source for the Windows and Linux programs. `build.sh` rebuilds them from `index.html` (needs Go 1.24).

## Entering lengths

One number is inches (`144` becomes 12' 0"). Two numbers are feet then inches (`26 18` becomes 27' 6"). The metric switch at the top changes weights to kilograms and lengths to meters.
