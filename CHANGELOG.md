<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given the same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [v7.2.0](https://github.com/cosmos/ibc-go/releases/tag/v7.2.0) - 2023-06-22

### Dependencies

* [\#3810](https://github.com/cosmos/ibc-go/pull/3810) Update Cosmos SDK to v0.47.3.
* [\#3862](https://github.com/cosmos/ibc-go/pull/3862) Update CometBFT to v0.37.2.

### State Machine Breaking

* [\#3907](https://github.com/cosmos/ibc-go/pull/3907) Re-implemented missing functions of `LegacyMsg` interface to fix transaction signing with ledger.

## [v7.1.0](https://github.com/cosmos/ibc-go/releases/tag/v7.1.0) - 2023-06-09 

### Dependencies

* [\#3542](https://github.com/cosmos/ibc-go/pull/3542) Update Cosmos SDK to v0.47.2 and CometBFT to v0.37.1.
* [\#3457](https://github.com/cosmos/ibc-go/pull/3457) Update to ics23 v0.10.0.

### Improvements

* (apps/transfer) [\#3454](https://github.com/cosmos/ibc-go/pull/3454) Support transfer authorization unlimited spending when the max `uint256` value is provided as limit.

### Features

* (light-clients/09-localhost) [\#3229](https://github.com/cosmos/ibc-go/pull/3229) Implementation of v2 of localhost loopback client.
* (apps/transfer) [\#3019](https://github.com/cosmos/ibc-go/pull/3019) Add state entry to keep track of total amount of tokens in escrow.

### Bug Fixes

* (core/04-channel) [\#3346](https://github.com/cosmos/ibc-go/pull/3346) Properly handle ordered channels in `UnreceivedPackets` query.
* (core/04-channel) [\#3593](https://github.com/cosmos/ibc-go/pull/3593) `SendPacket` now correctly returns `ErrClientNotFound` in favour of `ErrConsensusStateNotFound`.


