# SynapSeq Roadmap - Version 3.5 and Beyond

This roadmap outlines the next phase of evolution for SynapSeq following the completion of the V3 migration. The focus shifts from core engine stability toward ecosystem growth, interoperability, and accessible creation tools.

---

## 1. Licensing Upgrade to GPLv3

**Status: Under Consideration**

SynapSeq will migrate from GPL v2 to GPL v3 to enable broader compatibility with modern open-source ecosystems, including permissive licenses such as Apache 2.0.

This change expands the range of audio and encoding libraries that can be legally integrated, allowing significant improvements in export formats and tooling.

**Updated:**  
A potential migration to GPLv3 is under evaluation. This transition depends on the following condition:

A full rewrite of the remaining engine components that still conceptually derive from the SynapSeq v2 / GPL v2 (approximately 20% of the current audio engine).

Because this rewrite is substantial and not currently a priority, this migration may ultimately never happen.

---

## 2. Native Export to Compressed Audio Formats

**Status: In Research**

SynapSeq will support direct export to:

- MP3
- OGG Vorbis
- (Maybe) OPUS

Compressed formats drastically reduce file sizes, making it easier for users to store, distribute, and share their generated sessions. WAV will continue as the high-fidelity default, with optional compressed export.

**Updated:**
Initial MP3 export experiments using libmp3lame revealed that native support for compressed formats introduces significant technical complexity into the build process. The dependency on external C libraries, platform-specific toolchains, and custom compilation flags breaks the simplicity and portability that SynapSeq is designed to maintain.

This issue is especially severe on Windows, where reliable integration of LAME is difficult and harms reproducibility of builds.

---

## 3. SynapSeq Hub - Ecosystem Expansion

**Status: Ongoing**

The SynapSeq Hub will become a more central component of the ecosystem, evolving from a simple sequence repository to a platform with broader goals:

- community-driven sharing
- curated collections
- discoverability of sessions
- versioned metadata
- potential future marketplace for creators

The Hub may be split into its own standalone project as the ecosystem grows.

---

## 4. MP3 Distribution Through the Hub

**Status: Planned**

Every published session in the Hub will automatically offer ready-to-download MP3 versions, allowing users to listen without installing SynapSeq locally.

This complements the core tool while making SynapSeq accessible to non-technical audiences.

---

## 5. SynapSeq Playground - Web IDE

**Status: Concept**

A lightweight web-based IDE will be developed under the name SynapSeq Playground.

The goal is to provide:

- in-browser SPSQ editing
- live syntax validation
- optional real-time audio previews
- integration with the Hub for saving and publishing

This enables users to explore SynapSeq without installing anything.

---

## 6. Cross-Format Import and Conversion

**Status: Exploring Feasibility**

Planned support for importing and converting sessions from other brainwave generators:

- SBaGen (.sbg)
- Gnaural (.gnaural / XML)

This feature would unify the ecosystem, allowing legacy users to migrate their sessions to the modern SPSQ format.

Complexity varies per format and will require deeper analysis before implementation.

**Updated:**
Early investigation showed that translating formats like SBG into SPSQ is substantially more complex than anticipated.

SBaGen relies on implicit behaviors (relative timelines, auto-fade rules, automatic transitions, inherited track states, and NOW-dependent timing) that do not map cleanly to SynapSeq’s fully explicit and deterministic model.

Producing an automated converter would require interpreting intent that the source files simply do not encode, and the resulting SPSQ would still need manual correction.

Maintaining such a tool would add disproportionate overhead without meaningful benefit.

In practice, direct human translation is far more predictable, produces clearer results, and preserves the explicit nature of SPSQ without introducing hidden heuristics.
