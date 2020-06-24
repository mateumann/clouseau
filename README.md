# Clouseau

[![<mateumann>](https://circleci.com/gh/mateumann/clouseau.svg?style=shield)](https://circleci.com/gh/mateumann/clouseau)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b0582c09e1144e8bbf5c8c25fca89285)](https://www.codacy.com/manual/mateumann/clouseau?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mateumann/clouseau&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/mateumann/clouseau/branch/master/graph/badge.svg)](https://codecov.io/gh/mateumann/clouseau)
[![license](https://img.shields.io/github/license/mateumann/clouseau)](https://github.com/mateumann/clouseau/blob/master/LICENSE)
<!--[![Build](https://github.com/mateumann/clouseau/workflows/Build/badge.svg)](https://github.com/mateumann/clouseau/actions?query=workflow%3ABuild)
[![Linting](https://github.com/mateumann/clouseau/workflows/Linting/badge.svg)](https://github.com/mateumann/clouseau/actions?query=workflow%3ALinting) -->



Packet statistics collector (and analyser). The name is obviously inspired by the famous Inspector
Clouseau and his adventures.

## Architecture

Project (will be) built from a few components each of them doing exactly one job.  The first one is collector which
listenes for incomming [IP Flow Information Export (IPFIX)](https://tools.ietf.org/html/rfc7011) packets.
