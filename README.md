# Blockchain

This is an implementation of a toy proof of work blockchain.

Every block has a payload and a nonce. The hash of the payload + nonce must have N zeros prepended where N is the current blockchain difficulty.
There is a work function that goes through the whole blockchain block by block and searches for the nonce by generating an array of random bits.

Sharing of blocks is done by keeping a list of peers with [memberlist](https://github.com/hashicorp/memberlist).
Currently a simple broadcast is done to merge remote states on join.
When a new hash has been found for a block it is also broacasted as a differing state.

### Proof of Work
<img width="508" alt="image" src="https://user-images.githubusercontent.com/23063635/160709062-2938ae18-058e-4615-a75d-959d3618d6c7.png">

### Gossip
![image](https://user-images.githubusercontent.com/23063635/160708409-9aa0d529-afae-4bb4-8d45-ae7e7f6adc39.png)
