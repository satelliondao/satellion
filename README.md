# Satellion - open-source community driven bitcoin wallet

Annotation

We will develop a lightweight, secure, and transparent command-line Bitcoin wallet featuring a friendly, step-by-step ask-and-confirm interaction model. The project is non-profit and aimed primarily at experienced users who want straightforward, trustworthy self-custody without trading functions or unnecessary complexity.

The initial release will support Bitcoin only; consideration of other assets will be deferred until much later, and only after careful evaluation.

Focus: macOS and Linux systems, with hierarchical deterministic (HD) wallet generation and support for modern Bitcoin address types (SegWit and Taproot) exclusively.

## Principles

* **Simplicity**: minimal surface area; clear prompts; safe defaults.
* **Minimal dependencies**: prefer Go stdlib + a few vetted Bitcoin packages.
* **Privacy & security by default**: no telemetry, rotating addresses, encrypted secrets.
* **Open governance & transparency**: public roadmap, reproducible builds, signed releases.
* **Modern Bitcoin**: HD wallet generation with SegWit and Taproot address support only.

#### BIP157/158 (“Neutrino”): Modern LightweightClient Protocol

BIP157 and BIP158 together define a privacy-focused, efficient protocol for lightweight Bitcoin clients, improving on legacy SPV (Simplified Payment Verification, BIP37).
Neutrino is ONLY cool to wallet developers who do not want to deal with different APIs for querying the state of their lightning node. [Why I dont celebrate neutrino - nicolasdorier ](https://medium.com/@nicolasdorier/why-i-dont-celebrate-neutrino-206bafa5fda0)

**Problem**:

- **Legacy SPV (BIP37):** Clients use Bloom filters to request only relevant transactions from full nodes.
- **Drawbacks:** Bloom filters can leak which addresses you’re interested in and are bandwidth-inefficient.

**Neutrino Protocol**

- **Neutrino** is the name commonly used for this protocol.
- **Key change:** Instead of Bloom filters, it uses compact block filters.

**BIP158 – Compact Block Filters**

- Full nodes build a Golomb-Rice filter for each block.
- Each filter summarizes the addresses and scripts in the block.
- Filters are much smaller than full blocks, making them efficient to download.

**BIP157 – Client/Server Interaction**

- Defines how clients request and receive filters from full nodes.
- Clients download block headers and filters, then check locally for relevant transactions.
- If a filter matches, the client fetches the full block or just the relevant transactions.

**Benefits**

- **Privacy:** Clients don’t reveal their addresses to the server; filtering happens locally.
- **Efficiency:** Less data transferred compared to Bloom filters.
- **Scalability:** Full nodes don’t need to process custom filters for each client.
- **Security:** Clients can query multiple peers and only need one honest node for correct results.

**Adoption**

- Used by wallets such as Lightning Labs’ Neutrino, BTCPay Server mobile apps, and several Lightning Network wallets.
