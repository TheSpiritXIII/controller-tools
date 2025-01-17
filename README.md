> NOTE: This is a fork of controller-tools with https://github.com/kubernetes-sigs/controller-tools/pull/878

[![Go Reference](https://pkg.go.dev/badge/sigs.k8s.io/controller-tools.svg)](https://pkg.go.dev/sigs.k8s.io/controller-tools)
[![Build Status](https://travis-ci.org/kubernetes-sigs/controller-tools.svg?branch=master)](https://travis-ci.org/kubernetes-sigs/controller-tools "Travis")
[![Go Report Card](https://goreportcard.com/badge/sigs.k8s.io/controller-tools)](https://goreportcard.com/report/sigs.k8s.io/controller-tools)

# Kubernetes controller-tools Project

The Kubernetes controller-tools Project is a set of go libraries for building Controllers.

## Development

Clone this project, and iterate on changes by running `./test.sh`.

This project uses Go modules to manage its dependencies, so feel free to work from outside
of your `GOPATH`. However, if you'd like to continue to work from within your `GOPATH`, please
export `GO111MODULE=on`.

## Releasing and Versioning

See [VERSIONING.md](VERSIONING.md).

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

controller-tools is a subproject of the [kubebuilder](https://sigs.k8s.io/kubebuilder) project
in sig apimachinery.

You can reach the maintainers of this project at:

- Slack channel: [#kubebuilder](http://slack.k8s.io/#kubebuilder)
- Google Group: [kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
