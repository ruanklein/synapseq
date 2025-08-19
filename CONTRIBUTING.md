# Contributing to SynapSeq

First off, thank you for considering contributing to **SynapSeq**!
This project grows stronger with community support, whether through code, docs, bug reports, build scripts, or new ideas.

## TL;DR (Quick Summary)

- 🚀 PRs for the Go rewrite (v3) → branch: `development-v3`
- 🛠️ PRs for the C version (v2) → branch: `development` (bug fixes, docs, build scripts only)
- 🎵 Sequence files (`.spsq`) → see [`contrib/README.md`](./contrib/README.md)

---

## Contribution Policy

### Version 3.x (Go Rewrite)

- The current focus for version 3.0 (branch: `development-v3`) is the faithful porting of the original C code to idiomatic Go.
- **Only PRs that port code from C to Go, following the requirements and conventions described in [`ROADMAP.md`](./ROADMAP.md), will be accepted at this stage.**
- New features will only be considered starting in version **3.1+**, after the port is complete.
- Please open your PRs against `development-v3`.

### Version 2.x (C Codebase)

- The current C version (`synapseq.c`) is in **maintenance mode**.
- We still welcome contributions that add value for users:
  - **Bug fixes**
  - **Documentation improvements**
  - **Build scripts or integration tweaks** to support more systems and tools
- All PRs for the C codebase must target the `development` branch.

---

## Contributing Sequence Files (`.spsq`)

If you’d like to contribute new or improved sequence files (`.spsq`), please see [`contrib/README.md`](./contrib/README.md).  
This process is separate from code contributions.

---

## Commit Convention

We use the **Conventional Commits** format.  
Examples:

- `feat: add new waveform option`
- `fix: correct parsing bug for noise sequences`
- `docs: update README with usage examples`
- `build: add Makefile for macOS`
- `chore: clean up unused code in parser`

Following this format keeps the commit history clear and enables automated changelog generation in the future.

---

## Pull Request Guidelines

✅ Before opening a PR, please make sure:

- [ ] Your PR targets the correct branch (`development-v3` for Go, `development` for C).
- [ ] You followed the conventions in [`ROADMAP.md`](./ROADMAP.md) (v3).
- [ ] You wrote clear, descriptive commit messages (see above).
- [ ] Tests (if applicable) run successfully.
- [ ] Your changes are limited to the scope of the PR (no unrelated edits).

---

## Not Sure? Open an Issue First!

If you’re not sure whether your contribution fits, **open an issue**.  
We’ll be happy to discuss your idea before you start coding — saving you time and aligning with the project’s roadmap.

---

## License

By contributing to SynapSeq, you agree that your contributions will be licensed under the same license as the project.  
Please ensure you have the right to submit your code or content under these terms.

---

💡 Thank you for helping make SynapSeq better! Even small contributions - fixing typos, improving docs, or sharing ideas - help this project grow.
