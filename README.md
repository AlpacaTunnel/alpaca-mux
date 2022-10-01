alpaca-mux
==========

A UDP proxy with multiple paths.

# Topology

                         ┌────── path 1 ──────┐
udp-client -> mux-client ┼────── path 2 ──────┼ mux-server -> udp-server
                         └────── path 3 ──────┘

A path can be implemented by adding iptable port forward rules on a middler server.
