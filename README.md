# Satellion - open-source community driven bitcoin wallet

Annotation

We will develop a lightweight, secure, and transparent command-line Bitcoin wallet featuring a friendly, step-by-step ask-and-confirm interaction model. The project is non-profit and aimed primarily at experienced users who want straightforward, trustworthy self-custody without trading functions or unnecessary complexity.

The initial release will support Bitcoin only; consideration of other assets will be deferred until much later, and only after careful evaluation.

Focus: macOS and Linux systems, with hierarchical deterministic (HD) wallet generation and support for modern Bitcoin address types (SegWit and Taproot) exclusively.

---

## Principles

* **Simplicity**: minimal surface area; clear prompts; safe defaults.
* **Minimal dependencies**: prefer Go stdlib + a few vetted Bitcoin packages.
* **Privacy & security by default**: no telemetry, rotating addresses, encrypted secrets.
* **Open governance & transparency**: public roadmap, reproducible builds, signed releases.
* **Modern Bitcoin**: HD wallet generation with SegWit and Taproot address support only.
