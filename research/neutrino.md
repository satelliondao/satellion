#### Neutrino Client Protocol

BIP157 and BIP158 together define a privacy-focused, efficient protocol for lightweight Bitcoin clients, improving on legacy SPV (Simplified Payment Verification, BIP37).

##### Problem
- **Legacy SPV (BIP37):** Clients use Bloom filters to request only relevant transactions from full nodes.
- **Drawbacks:** Bloom filters can leak which addresses you’re interested in and are bandwidth-inefficient.

##### Solution
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

