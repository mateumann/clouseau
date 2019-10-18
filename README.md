# Clouseau

Packet inspector and analyser. The name is obviously inspired by famous Inspector Clouseau and his adventures.

## Architecture

This packet inspecting tool consists of a few parts.  The first captures all packets going thorough given interface(s)
and is built upon `libpcap`.  The second part is a storage backend – it could be either standard SQL database (e.g.
conteinerised PostgreSQL), or maybe this database (?) that powers OpenLDAP (unfortunately I do not remember the exact
name, but it is something what evolved from Berleley DB).  The third component of this tool is a packet inspector and
analyser that focuses on those most interesting questions like where is the most of packets from my phone going to,
what kind of services is mostly used etc.  The fourth part is an API server (GraphQL maybe) that makes the service
available via HTTP (preferably v2 or even v3).  The last part is a Web UI that draws nice dashboards, graphs etc.
(actually maybe Grafana might be used for this purpose).
