---
title: nitro-api
tags: nitro, api, rpc
author: Louis
source: 
---

# 0009 - Nitro as a service API definition

## Status

proposed

## Context

Objectives are a powerful abstraction for multi-party protocols, but when using nitro as a library it can be time consuming and cumbersome for new app developers.

Virtual application using voucher abstractions are the most lightweight and autonomous way of interacting with Nitro channels.

Nitro should abstract the virtual channel plumbing to the app developer and provide a simple interface to provision various channel construction protocols.

## Decision

Define the API features to be exposed using various encoders and transports such as:
- HTTP2
- Websocket
- NATS
- TCP wire

Encoding:
- JSON/BSON
- MsgPack
- Protobuf

### Notation

We will use the following notation to describe the API, but they can be implemented in different ways.

This notation remains the recommended serialization for JSON, BSON, MsgPack, Protobuf, it is more compact than json-rpc and compatible with Message Queue systems.

#### Request:
`[type, msgid, method, params]`

#### Response:
`[type, msgid, error, result]`

#### Event:
`[type, msgid, method, params]`

##### type
The message type, must be the integer zero (0) for "Request" messages. 
One (1) means that this message is the "Response" message.
Two (2) means that this message is the "Notification" message.

##### msgid
A 32-bit unsigned integer number. This number is used as a sequence number. The server's response to the "Request" will have the same msgid.

Events streams may also relate to a subscription "Request".

##### method
A string which represents the method name.

##### params
An array of the function arguments. The elements of this array are arbitrary objects.

#### error
If the method is executed correctly, this field is Nil. If the error occurred at the server-side, then this field is an arbitrary object which represents the error.

#### result
An arbitrary object, which represents the returned result of the function. If an error occurred, this field should be nil.

### create_objective

`[0, 42, "create_objective", [protocol, peers, app_definition]]`

Asynchronous provision of protocol, will return a valid state proof when protocol is ready or an error.

- protocol: **String**, name of protocol, "directfund", "directdefund", "virtualfund", "virtualdefund", "virtualchallenge"
- peers: **Array**, list of participants starting from the creator.
- app_definition: **String**, smart-contract address of the NitroApp.

#### Response:

`[1, 42,  nil, objective_id]`

Params:
- **objective_id**: **String**, unique identifier of the Objective request.

### objective_created Event

`[2, 42, "objective_created", [id, proof]]`

Params:
- id, **string**, Objective identifier
### cancel_objective

`[0, 43, "cancel_objective", id]`

#### Response

`[1, 43, error, result]`

## Consequences

Using nitro as a service as the following benefits:
- All application languages are supported
- Portability of legacy applications and protocols.
- Nitro service can simplify access to plumbing
- Enable applications to use high-performance property protocols
- App developer experience is greatly improved
- Enable a much wider adoption of Blockchain benefits
- Ability of having light web client 