alpaca-mux
==========

A UDP proxy with multiple paths.

# Topology

```
                                       ┌────── path 1 ──────┐
app-client -> udp-server -> mux-client ┼────── path 2 ──────┼ mux-server -> udp-client -> app-server
                                       └────── path 3 ──────┘
```

A path can be implemented by adding iptable port forward rules on a middler server.

```sh
iptables -A PREROUTING  -t nat -p udp --dport 8001 -j DNAT --to-destination 192.168.1.200:8080
iptables -A POSTROUTING -t nat -p udp -d 192.168.1.200 --dport 8080 -j MASQUERADE
```
