# Privateer Raid Wireframe for RDS

Use this a Raid for use with [Privateer](https://www.github.com/privateerproj/privateer).

This project was created from the base Privateer Project's [Raid Wireframe](https://www.github.com/privateerproj/raid-wireframe).

## Installation

Create a directory for the binary

```bash
mkdir -p ~/privateer/bin
```

Build the binaries and copy them to the bin directory

```bash
make release
```

## Primary Maintainer

![Krumware Logo](https://www.krum.io/assets/icons/logo-with-name.svg)

In accordance with the vision outlined by the 
[Compliant Financial Infrastructure](https://github.com/finos/compliant-financial-infrastructure)
project (CFI), this Raid has been adopted by [Krumware](https://www.krum.io/).

Contributions via pull requests and issue reports are welcome and encouraged from all parties or persons.

## Roadmap

### Now

This project is currently being developed by the CFI Runtime Validation Working Group and used as a 
proof of concept to demonstrate the value of the work being done by the
[Common Cloud Controls](https://github.com/finos/common-cloud-controls) project (CCC).

In collaboration with the CFI Reproducible Infrastructure Working Group, this Raid is being run as
part of a CI pipeline to validate a [CFI-managed RDS service](https://github.com/finos/cfi-ansible-aws-rds)
deployed by an ansible playbook.

### Near Future

As the CFI Reproducible Infrastructure WG continues to iterate on their infrastructure as code,
additional tests ("strikes") can be built to ensure holistic compliance with the CCC taxonomy.

### Optimistic Future

As the CCC project creates controls related to RDMS security, this Raid can support a new set of associated 
strikes ("tactic").

The new tactic will enable a user to select whether they would like to validate an RDS service's
compliance with CCC's service taxonomy or security controls.

New contributors to this group will allow that larger and more complex tactic to be developed more rapidly.

