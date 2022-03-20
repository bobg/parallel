# Parallel - generic functions for coordinating parallel workers

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/parallel.svg)](https://pkg.go.dev/github.com/bobg/parallel)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/parallel)](https://goreportcard.com/report/github.com/bobg/parallel)
[![Tests](https://github.com/bobg/parallel/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/parallel/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/parallel/badge.svg?branch=master)](https://coveralls.io/github/bobg/parallel?branch=master)

This is parallel,
a collection of generic functions for coordinating parallel workers.

This package includes these functions:

- `Consumers`, for managing a set of N workers consuming a stream of values produced by the caller
- `Producers`, for managing a set of N workers producing a stream of values consumed by the caller
- `Values`, for concurrently producing a set of N values
- `Pool`, for managing access to a pool of concurrent workers
