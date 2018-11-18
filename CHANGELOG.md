# Changelog

All notable changes to this project will be documented in this file.

## [0.3.4] - 2018-10-11

### Added

* Adds new mode aws-sec-ssm which is using ssm parameters type SEcureString which can be automatically decrypted by AWS KMS key
* Update Dockerfile to multistage build format

## [0.3.0] - 2017-09-03

### Added

* Adds local mode for unseal keys in path (#21, [@JoshVanL](https://github.com/JoshVanL))

### Changed

* Use golang 1.10.4 (#20, [@simonswine](https://github.com/simonswine))

## [0.2.1] - 2018-02-12

### Changed

* Use logrus with lower case spelling #11
* Fix handling of kv.NotFound error #5 #7

## [0.2.0] - 2017-11-23

### Added

* Sign binaries using GPG key

### Changed

* Move the repository from jetstack-experimental to jetstack
* Updated to Golang 1.9.2
