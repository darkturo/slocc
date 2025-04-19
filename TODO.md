# TODO List for slocc

## 🎯 Codebase Improvements

- [x] Refactor flag parsing into a dedicated `parseArgs()` function
- [x] Add `-v / --verbose` flag to print number of ignored files
- [x] Print binary name cleanly in help message using `filepath.Base(os.Args[0])`

## 🧹 Code Quality Enhancements

- [ ] Use `defer file.Close()` only when file opening succeeds (avoid closing nil)
- [ ] Replace `filepath.Walk` with `filepath.WalkDir` for improved traversal performance
- [ ] Sort language keys before printing output for deterministic language order
- [ ] Extract `.tmpl` output into a separate template file for easier customization (optional)

## 💡 Features Worth Exploring

- [ ] Add `--version` flag to print version info
- [ ] Add `--list-languages` to display supported languages
- [ ] Add support for excluding specific extensions via flags
- [ ] Add support for outputting JSON (for piping / scripts)
- [ ] Add ability to pass config file path via `--config=<file>`

## 🧪 Testing

- [ ] Add integration test for CLI execution with temp project folder
- [ ] Add tests for new flag behaviors
- [ ] Add golden file test for output formatting

---

Generated collaboratively with 🍵 and a dash of obsession for tidy CLIs.

