# Axle Weight Planner

Truck axle weight calculator for steer, drive and trailer groups, with tandem and 5th wheel slides, lift axles, jeep dollies, lowboys with boosters, step decks and legal-limit checks (federal bridge formula or custom per-axle limits).

- **index.html** is the whole app in one file. Open it in any browser.
- **releases/** holds the Windows program (`AxleWeightPlanner.exe`) and the Linux program (`axle-planner-linux`, started with `sh run-linux.sh`). Saved setups are kept in `setups.json` next to the program, so one thumb drive can serve both systems.
- **program/** is the Go source for those programs. `build.sh` rebuilds them from `index.html` (needs Go 1.24).

Lengths are typed as feet and inches: one number is inches (`144` becomes 12' 0"), two numbers are feet then inches (`26 18` becomes 27' 6").
