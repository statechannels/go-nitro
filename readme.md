#  ⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️ TODO ⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️⚠️

## What was I doing?
Trying to switch the message format to allow custom payloads defined by objectives. `VirtualDefund`  on a message payload is defined as `interface{}` to allow the virtual defund objective to decide the exact shape of the message. However  I'm running into issues around casting and types trying to go from a message to either a `StartMessage` or`UpdateSigMessage`.

Am I just trying to reimplement intersection types in a language that doesn't want that?

I also started switching some payloads to be `byte[]` but I think that just mostly broke things. Currently some tests pass but most are broken by my payload typing shenanigans.

Probably should pause and have a think. Maybe take a look at existing libraries.