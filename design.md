# Design

## Goals

1. The library should be able to read credentials from Docker configuration file and external credentials stores such as the native keychain of the operating system.
2. The library should be able to save credentials in Docker configuration file and external credential stores. Besides, it should provide options to disable saving credentials in plain text.
3. The library should be able to remove the specified credentials from external credential stores or Docker configuration file.
4. The library should be easy to integrate with Go projects that uses `oras-go v2`.

## Non-Goals

1. The library will not support configuration formats other than the Docker configuration file.
2. The library will not handle encryption of credentials.

## Challenges

1. Depending on the version of the Docker CLI installed on the target machine, the format of the Docker configuration file may be different. The library should ensure that no config field loss when saving credentials to the Docker configuration file.
2. The credential helper configuration and the credential itself may be maintained in different configuration files.

## Proposals
