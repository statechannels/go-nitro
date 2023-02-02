---
description: ForceMove, Nitro and SATP -- Evolution of the framework.
---

# Research papers

Nitro is based on a culmination of research spanning the following papers:

## (2018) [:fontawesome-solid-file-pdf:](https://magmo.com/force-move-games.pdf) ForceMove

An n-party state channel protocol. ForceMove contains early descriptions of a _generalized_ state channel protocol, i.e. one which allows for [user-defined applications](./protocol-tutorial/0010-states-channels.md#appdefinition), and of a turn-taking paradigm that allows individual participants to force changes to channel state (as opposed to all participants co-signing a state change). This turn-taking paradigm is also covered in the blog post [Putting the 'state' in state channels](https://blog.statechannels.org/putting-the-state-in-state-channels/).

Note: much of the terminology from this paper is out of date with respect to current implementation.

## (2019) [:fontawesome-solid-file-pdf:](https://magmo.com/nitro-protocol.pdf) Nitro

A protocol for state channel networks. Describes protocols for the secure adjustment of channel outcomes to:

- fund subchannels
- perform top ups of running channels, and
- to construct [virtual channels](https://blog.statechannels.org/virtual-channels/) between users who have no common channel, via a shared intermediary

Note: much of the terminology from this paper is out of date with respect to current implementation.

## (2022) [:fontawesome-solid-file-pdf:](https://statechannels.github.io/satp_paper/satp.pdf) Stateful Asset Transfer Protocol

This paper describes refinements to Nitro protocol. Highlights are a protocol for flat, multi-hop virtual channel construction, and the replacement of guarantee channels with a richer [Outcome format](./protocol-tutorial/0030-outcomes.md). In essence, the virtual funding protocol in this paper achieves the same functionality as the Nitro paper, but with

- fewer rounds of network communication (reducing latency for construction of channels), and
- fewer channels required in the construction (reducing gas cost of the challenge path for channels in case of disputes)

Terminology and protocol descriptions from this paper are up to date, and should match the current implementation reasonably well.
