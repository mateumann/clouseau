# Clouseau

[![<mateumann>](https://circleci.com/gh/mateumann/clouseau.svg?style=shield)](https://circleci.com/gh/mateumann/clouseau)
[![Codacy Badge](https://img.shields.io/codacy/grade/f194807bd2fe48eabedef320abc91f0b)](https://www.codacy.com/manual/mateumann/clouseau?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mateumann/clouseau&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/mateumann/clouseau/branch/master/graph/badge.svg)](https://codecov.io/gh/mateumann/clouseau)
[![license](https://img.shields.io/github/license/mateumann/clouseau)](https://github.com/mateumann/clouseau/blob/master/LICENSE)

Packet statistics collector (and analyser). The name is inspired by the famous
Inspector Clouseau and his adventures.

## Architecture

Project (will be) built from components, each of them doing exactly one job.
The first one is collector which listenes for incomming
[IP Flow Information Export (IPFIX)](https://tools.ietf.org/html/rfc7011)
packets.
