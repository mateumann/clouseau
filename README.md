# Clouseau

[![<mateumann>](https://circleci.com/gh/mateumann/clouseau.svg?style=shield)](https://circleci.com/gh/mateumann/clouseau)
[![codecov](https://codecov.io/gh/mateumann/clouseau/branch/master/graph/badge.svg)](https://codecov.io/gh/mateumann/clouseau)
<!--[![Build](https://github.com/mateumann/clouseau/workflows/Build/badge.svg)](https://github.com/mateumann/clouseau/actions?query=workflow%3ABuild)
[![Linting](https://github.com/mateumann/clouseau/workflows/Linting/badge.svg)](https://github.com/mateumann/clouseau/actions?query=workflow%3ALinting) -->


Packet statistics collector (and analyser). The name is obviously inspired by famous Inspector Clouseau and his
adventures.

# Architecture

Project (will be) built from a few components each of them doing exactly one job.  The first one is collector which
listenes for incomming [IP Flow Information Export (IPFIX)](https://tools.ietf.org/html/rfc7011) packets.
