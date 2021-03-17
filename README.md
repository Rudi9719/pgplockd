# pgplockd
A logind locker using PGP

### Proposal
This program should be runnable as a systemd user service to periodically query the user's PGP key and lock the screen if the key is locked.

## Current Status
* journald is being used for logging
* connection to logind is being made
* user's ~/.pgplockd file can be read for PGP Fingerprint
* current session is being found
* timeout loop can lock session using logind\

## To Do
* Implement an actual check for unlock status
* Random message generator?

## Workflow
* User logs in
* pgplockd starts
* Random message is generated and requested to be signed, timeout countdown starts
* User unlocks PGP key and signs message OR countdown locks screen
* Timer starts, then another message is generated


If the user does/can not unlock the PGP key (be it a removed smartcard or otherwise) logind should lock the screen.