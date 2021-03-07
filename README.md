# k8l

Lightweight Kubernetes native log aggregation.

> "Sometimes EFK is too much." cit. George Wahington

`k8l` (Kate-el for friends) is a scalable fault tolerant app that aggregates logs coming from Fluentd, using `dsqlite` as a backend, being written in `golang` its footprint in terms of resource is very low.

It is **self contained**, it doesn't need any external UI (but you can develop a new one if you don't like what's embedded) or any external backend to store the logs (you just need a Persistent Volume though).

When your cluster doesn't need all the features of the ELK stack, but you just need a centralized place to view and store logs, **k8l** is all you need.

## Install

TBD

## Similar software

- [Kubetail](https://github.com/johanhaleby/kubetail)
- [Kail](https://github.com/boz/kail)
- [Avologo](https://avologo.com/)

## Roadmap

- Retention control
- live tail with websocket
- configurable table clustering  `all`, `namespace`, `resource`
