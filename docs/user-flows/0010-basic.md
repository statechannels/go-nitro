---
description: How channels are opened, executed and closed.
---

# Basic

This section describes how the Nitro framework can be used at increasing levels of scalability.

A very simple interaction involves just two parties, Alice and Bob, who transact off-chain. After prefund states are exchanged, at least one of them [deposits into the adjudicator](../protocol-tutorial/0060-funding-a-channel.md) in priority order. Then, after postfund states are exchanged, the channel may be [executed according to the rules of the channel](../protocol-tutorial/0020-execution-rules.md). Then, Alice and Bob may agree to finalize, conclude and liquidate the channel.

![Outcome priority](./basic-user-flow.png)

<!-- fontawesome f182 Alice
fontawesome f183 Bob
fontawesome f0e3 Adjudicator #red

group prefunding
Alice-#purple>Bob: create channel
Alice<#purple-Bob: join channel
end
Alice-#red>Adjudicator: deposit
Bob-#red>Adjudicator: deposit
group postfunding
Alice-#purple>Bob: confirm deposits
Alice<#purple-Bob: confirm deposits
end
group running
Alice-#purple>Bob: update
Bob-#purple>Alice: countersign
Alice-#purple>Bob: update
Bob-#purple>Alice: countersign
end
group finalizing
Alice-#purple>Bob: finalize channel
Alice<#purple-Bob: agree
end
Alice-#red>Adjudicator: concludeAndTransferAll -->
