# k8l

Lightweight Kubernetes native log aggregation.

> Sometimes EFK is too much Enter FL

`k8l` (Katelogs for friends) is a scalable fault tolerant app that aggregates logs coming from Fluentd, using `dsqlite` as a backend, being written in `golang` its footprint in terms of resource is very low.

It is self contained, doesn't need an external UI or external backend to store the logs.

When your cluster doesn't need all the features of the ELK stack, but you just need a centralized place to view container logs, k8l is all that you need.

## competitors

- [https://avologo.com/]
- not very much really

## MVP

- ingestion endpoint  
  - `/_bulk` act liek ES
- yaml resource e/o helm chart
- REST api
- UI
  - filter indexed search
  - fulltext search

## Roadmap

- retention control
- live tail of logs with websocket
