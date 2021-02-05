# Clouseau

[![github](https://img.shields.io/github/checks-status/mateumann/clouseau/master)](https://github.com/mateumann/clouseau)
[![<circleci>](https://circleci.com/gh/mateumann/clouseau.svg?style=shield)](https://circleci.com/gh/mateumann/clouseau)
[![DepShield Badge](https://depshield.sonatype.org/badges/mateumann/clouseau/depshield.svg)](https://depshield.github.io)
[![Codacy Grade](https://img.shields.io/codacy/grade/f194807bd2fe48eabedef320abc91f0b)](https://www.codacy.com/gh/mateumann/clouseau/dashboard)
[![codecov](https://img.shields.io/codecov/c/gh/mateumann/clouseau?token=d4658c59e19d48d9b937bf8fd825714c)](https://codecov.io/gh/mateumann/clouseau)
[![license](https://img.shields.io/github/license/mateumann/clouseau)](https://github.com/mateumann/clouseau/blob/master/LICENSE)

Packet statistics collector (and analyser). The name comes from the famous
Inspector Clouseau and his adventures.

## Architecture

Project (will be) built from components, each of them performing one task.
The first one is collector which listenes for incomming
[IP Flow Information Export (IPFIX)](https://tools.ietf.org/html/rfc7011)
packets.
