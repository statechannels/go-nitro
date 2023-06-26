# Multi-node tests

This folder contains the `node_test` package, including several tests which orchestrate interactions between several `go-nitro` nodes.

The nodes run in the same process as the tests, and communicate (via `go chans`) with a local messaging system. The messaging system has the capability to add random delays to message dispatch, causing message reordering.

The tests check for:

- protocols succesfully completing (as indicated by the `Node.CompletedObjectives()` API)
- Node `stores` containing expected information
- The duration of test runs

The tests output:

- basic logs to any chosen file (often `../artifacts/name_of_test.log`)
- vector clock logs suitable for visualization with https://bestchai.bitbucket.io/shiviz/
