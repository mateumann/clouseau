# Clouseau

[![<circleci>](https://circleci.com/gh/mateumann/clouseau.svg?style=shield)](https://circleci.com/gh/mateumann/clouseau)
[![Codacy Grade](https://img.shields.io/codacy/grade/f194807bd2fe48eabedef320abc91f0b)](https://www.codacy.com/gh/mateumann/clouseau/dashboard)
[![Codeac](https://static.codeac.io/badges/2-216085254.svg "Codeac")](https://app.codeac.io/github/mateumann/clouseau)
[![codecov](https://codecov.io/gh/mateumann/clouseau/branch/master/graph/badge.svg)](https://codecov.io/gh/mateumann/clouseau)
[![license](https://img.shields.io/github/license/mateumann/clouseau)](https://github.com/mateumann/clouseau/blob/master/LICENSE)

Packet statistics collector (and analyser). The name comes from the famous
Inspector Clouseau and his adventures.

## Architecture

Project (will be) built from components, each of them performing one task.
The first one is collector which listenes for incomming
[IP Flow Information Export (IPFIX)](https://tools.ietf.org/html/rfc7011)
packets.
