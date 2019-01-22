# Changelog

All notable changes to this project will be documented in this file.

## [0.3.1] - 2019-01-22

### Added

* Documentation: updated instructions for how to fetch keys from SSM and KMS after vault init. (#31, [@innovia](https://github.com/innovia))
* Right now you can configure the `unseal-period` but it's not being used. This utilizes it in the sleep.  (#15, [@sheldonkwok](https://github.com/sheldonkwok))

### Changed

* Upgrade go to 1.11.4 (#40, [@JoshVanL](https://github.com/JoshVanL))
* Update dep and lock vault version to 0.9.6 (#39, [@JoshVanL](https://github.com/JoshVanL))
* Upgrade alpine image to 3.8 (#38, [@JoshVanL](https://github.com/JoshVanL))
* Use DOCKER_HOST: tcp://localhost:2375 for dind (#36, [@JoshVanL](https://github.com/JoshVanL))

## [0.3.0] - 2017-09-03

### Added

* Adds local mode for unseal keys in path (#21, [@JoshVanL](https://github.com/JoshVanL))

### Changed

* Use golang 1.10.4 (#20, [@simonswine](https://github.com/simonswine))

## [0.2.1] - 2018-02-12

### Changed

* Use logrus with lower case spelling #11
* Fix handling of kv.NotFound error #5 #7

## 0.2.0 - 2017-11-23

### Added

* Sign binaries using GPG key

### Changed

* Move the repository from jetstack-experimental to jetstack
* Updated to Golang 1.9.2

[Unreleased]: https://github.com/jetstack/tarmak/compare/0.3.1...HEAD
[0.3.1]: https://github.com/jetstack/tarmak/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/jetstack/tarmak/compare/0.2.1...0.3.0
[0.2.1]: https://github.com/jetstack/tarmak/compare/0.2.0...0.2.1
