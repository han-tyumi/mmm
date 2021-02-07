## v0.1.0 (2021-02-07)

### Perf

- :zap: cache mods by slug initially
- cache slug mods
- :zap: add more concurrency
- :zap: introduce concurrency

### Fix

- :bug: remove second mod slug check
- :bug: use correct dep key
- :bug: correct ModsBySlug

### Feat

- add mod file permissions (755)
- :fire: remove search for get and add
- :sparkles: add table formatting for search
- improve search results
- :sparkles: add initial remove command
- :sparkles: add update command
- add initial slug support
- :sparkles: initial add command
- add init command
- add cwd flag
- add get version flag
- :sparkles: add initial get command
- :zap: update to latest mcf
- :sparkles: add initial search command
- :sparkles: add cobra

### Refactor

- :truck: move non-command related packages out of cmd
- :recycle: extract config logic
- extract latest files for args
- :recycle: extract download and latest file logic
- :recycle: move get functions into package
- update version usage
- :recycle: use utils.Error
