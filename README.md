# Auto Staging Tower

[![Maintainability](https://api.codeclimate.com/v1/badges/6a32b4e0cf97b4cf635d/maintainability)](https://codeclimate.com/github/auto-staging/tower/maintainability)
[![GoDoc](https://godoc.org/github.com/auto-staging/tower?status.svg)](https://godoc.org/github.com/auto-staging/tower)
[![Go Report Card](https://goreportcard.com/badge/github.com/auto-staging/tower)](https://goreportcard.com/report/github.com/auto-staging/tower)
[![Build Status](https://travis-ci.com/auto-staging/tower.svg?branch=master)](https://travis-ci.com/auto-staging/tower)

> Tower is the central management Lambda function for auto-staging, it is called through the AWS API Gateway.

## Tower API Documentation

[OpenAPI Specification](https://app.swaggerhub.com/apis-docs/auto-staging/auto-staging-tower/1.0.0)

## Requirements

- Golang
- Go dep

## Usage

### Install dependencies

Go dep must be installed

```bash
make prepare
```

### Build binary

```bash
make build
```

compiles to bin/auto-staging-tower

## License and Author

Author: Jan Ritter

License: MIT