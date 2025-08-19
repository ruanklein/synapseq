# Contributing to SynapSeq

Thank you for your interest in contributing to SynapSeq! Please read the guidelines below before opening a pull request (PR).

---

## Contribution Policy

### Version 3.x (Go Rewrite)

- **New features will only be accepted for version 3.1 and above, and must be implemented in Go.**
- The current focus for version 3.0 (branch: `development-v3`) is the faithful porting of the original C code to idiomatic Go.
- **Only PRs that port code from C to Go, following the requirements and conventions described in [`ROADMAP.md`](./ROADMAP.md), will be accepted at this stage.**
- If you wish to help accelerate the Go rewrite, please consult the [`ROADMAP.md`](./ROADMAP.md) for detailed requirements and architectural guidelines.
- All PRs for the Go rewrite must target the `development-v3` branch.

### Version 2.x (C Codebase)

- The current C version (`synapseq.c`) is in maintenance mode and will be deprecated after the Go rewrite is released.
- **Only bug fixes will be accepted for the C codebase.**
- PRs for the C version must target the `development` branch.
- Contributions such as build scripts and documentation improvements for the C version are welcome if they add value to the project.

---

## Contributing Sequence Files (`.spsq`)

If you want to contribute new or improved sequence files (`.spsq`), please refer to the guidelines in [`contrib/README.md`](./contrib/README.md).  
The contribution process for sequence files is different from code contributions and is described in detail in that document.

---

## Pull Request Guidelines

- **For Go (v3) contributions:**

  - Open your PR against the `development-v3` branch.
  - Ensure your changes strictly follow the structure and naming conventions outlined in [`ROADMAP.md`](./ROADMAP.md).
  - Do not submit new features or refactors outside the porting scope until v3.1.

- **For C (current) contributions:**
  - Open your PR against the `development` branch.
  - Only bug fixes, build scripts, or documentation improvements will be considered.

---

## Not Sure? Open an Issue!

If you are unsure whether your contribution will be useful or fit the project's direction, **feel free to open an issue first**.  
We are happy to discuss ideas, answer questions, and help guide your contribution before you start working on a pull request.

---

## License

By contributing to SynapSeq, you agree that your contributions will be licensed under the same license as the project.  
Please make sure you have the right to submit your code or content under these terms.

---

Thank you for helping make SynapSeq better! We appreciate your contributions.
