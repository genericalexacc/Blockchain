# Blockchain

This is an implementation of a proof of work blockchain.
This project is more of an exploration of the idea for me, I didn't consult any material prior to starting it, and it is probably closest to a simple bitcoin type blockchain.

Every block has a payload and a hash. The hash must have N zeros prepended where N is the current blockchain difficulty.
There is a work function that goes through the whole blockchain block by block and calculates the hash by generating an array of random bits. There is definetly a better way of doing this.

Sharing of blocks is done by keeping a list of peers with [memberlist](https://github.com/hashicorp/memberlist).

Then using [controlled flooding](https://en.wikipedia.org/wiki/Flooding_(computer_networking)) we share work with all nodes.


<img width="1019" alt="image" src="https://user-images.githubusercontent.com/23063635/159838461-e1185cd6-9723-4c65-a8f7-280b369671ef.png">
