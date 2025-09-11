# Contributing to SynapSeq

First off, thank you for considering contributing to **SynapSeq**!
This project grows stronger with community support, whether through code, docs, bug reports, build scripts, or new ideas.

## TL;DR (Quick Summary)

🚀 The `main` branch contains the latest Go (V3) codebase.  
🔀 All pull requests (PRs) should target the `development` branch.  
🎵 Sequence files (`.spsq`) → see [`contrib/README.md`](./contrib/README.md)

🗂️ The legacy C code (V2, inherited from SBaGen) is available in the `v2` branch. If you wish to view, modify, or fork the old version, use that branch.

---

## Contribution Policy

### Version 3.x (Go)

- The `main` branch contains the latest Go (V3) codebase.
- All new features, bug fixes, and improvements should be submitted as PRs to the `development` branch.
- Follow Go best practices and maintain clean, readable code.
- Ensure backward compatibility when possible.

### Version 2.x (C, legacy)

- The legacy C codebase (V2, inherited from SBaGen) is available in the `v2` branch.
- No new features will be accepted for V2. Only maintenance or forks should use this branch.

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

- [ ] Your PR targets the correct branch (`development` for Go V3, `v2` for legacy C code).
- [ ] You wrote clear, descriptive commit messages (see above).
- [ ] Tests (if applicable) run successfully.
- [ ] Your changes are limited to the scope of the PR (no unrelated edits).
- [ ] Code follows Go best practices and conventions (for V3 contributions).
- [ ] Documentation is updated if your changes affect user-facing features.

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
