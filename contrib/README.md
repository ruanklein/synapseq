# Contributing Guide

This directory contains community contributions to the SynapSeq project. We're grateful for everyone who shares their brainwave entrainment sequences with our community!

**Thank you for contributing!**

Your creativity and expertise help others discover new ways to achieve relaxation, focus, and meditation through sound. Every contribution, no matter how small, makes the SynapSeq ecosystem richer and more diverse.

Whether you're sharing your first sequence or you're an experienced creator, your work inspires others and helps build a supportive community around brainwave entrainment.

Keep creating and sharing!

---

## Table of Contents

- [Contributing Projects](#contributing-projects)
- [How to Contribute](#how-to-contribute)
- [Rules for Project Contributions](#rules-for-project-contributions)
- [Example Structure](#example-structure)
- [Example CREDITS.md](#example-creditsmd)
- [Steps to Contribute](#steps-to-contribute)

---

## Contributing Projects

This directory is specifically for sharing `.spsq` sequence files created by the community.

### How to Contribute

You can contribute projects in two ways:

1. **Fork and Pull Request**: Follow the steps below
2. **Discussion**: Create a topic in the [Discussions/Sequences](https://github.com/ruanklein/synapseq/discussions/categories/sequences) section

### Rules for Project Contributions

If you fork and create a pull request, please follow these rules:

- The contribution must be placed within the `contrib` directory
- Create a folder named after your project
- Project and sequence names must be lowercase with spaces replaced by hyphens (e.g., "Deep Relax" becomes "deep-relax")
- Inside the project folder, create a `README.md` file with a description of the project
- You may create additional `.md` files within that folder for additional details
- Create a `CREDITS.md` file with your GitHub username or other information that identifies you as the creator of the project
- You can include multiple `.spsq` files in the same project folder
- All documentation files (`.md` files) must be written in English

### Example Structure

```
contrib/
├── my-awesome-project/
│   ├── README.md
│   ├── CREDITS.md
│   ├── additional-details.md (optional)
│   ├── sequence-1.spsq
│   └── sequence-2.spsq
└── another-project/
    ├── README.md
    ├── CREDITS.md
    └── main-sequence.spsq
```

### Example `CREDITS.md`

```
# Credits

Creator: [@your-github-username](https://github.com/your-github-username)

Contact: your-email@example.com (optional)
Website: https://your-website.com (optional)
Description: Brief description about yourself or your work (optional)
```

### Steps to Contribute

1. Fork the repository
2. Create a new branch with the naming convention: `contrib/project-name`
3. Create your project folder in the `contrib` directory
4. Add your `.spsq` file(s) and required documentation
5. Submit a pull request
