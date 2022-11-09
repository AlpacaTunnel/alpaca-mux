alpaca-mux
==========

A UDP proxy with multiple paths.

# Topology

```
dns-client ──┐                                                      ┌──> dns-server
             │                                                      │
app-client   │                                                      │    app-server
      ↓      ↓                                                      │      ↑
   port-1   port-2                                               conn-2   conn-1
      ↓      ↓                ┌────── path 1 ──────┐                ↑      ↑
     udp-server -> mux-client ┼────── path 2 ──────┼ mux-server -> udp-client
                              └────── path 3 ──────┘
```

Proxy multiple UDP ports over a single mux session, then transport a mux session over multiple paths. The second UDP port is useful when chaining shadowsocks -> kcptun -> alpaca-mux, because kcptun does not support UDP relay.

A path can be implemented by adding iptable port forward rules on a middler server.

```sh
iptables -A PREROUTING  -t nat -p udp --dport 8001 -j DNAT --to-destination 192.168.1.200:8080
iptables -A POSTROUTING -t nat -p udp -d 192.168.1.200 --dport 8080 -j MASQUERADE
```
