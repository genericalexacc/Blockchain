# B (lol) kchain 🧱⛓

This is an implementation of a toy proof of work blockchain.

Every block has a payload and a nonce. The hash of the payload + nonce must have N zeros prepended where N is the current blockchain difficulty.
There is a work function that goes through the whole blockchain block by block and searches for the nonce by generating an array of random bits.

Sharing of blocks is done by keeping a list of peers with [memberlist](https://github.com/hashicorp/memberlist).
Currently a simple broadcast is done to merge remote states on join.
When a new hash has been found for a block it is also broacasted as a differing state.

### Proof of Work
<img width="508" alt="image" src="https://user-images.githubusercontent.com/23063635/160709062-2938ae18-058e-4615-a75d-959d3618d6c7.png">

### Gossip
<img width="618" alt="image" src="https://user-images.githubusercontent.com/23063635/160709284-ad438f94-261f-4d5b-9d35-06a0d46b661b.png">
