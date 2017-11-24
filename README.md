# Rabbit IO - Work in Progress  
This is a tool to support backup and restoring of RabbitMQ messages, currently work in progress and might not be functional

## Requirements

You will need following to build `rabbitio` locally:

- [Golang](https://golang.org/dl/)
- [dep](https://github.com/golang/dep)

## Getting started

If you plan to work on `rabbitio` you will need to:

1. Create directories
```
mkdir -p $GOPATH/src/github.com/meltwater
```

2. Clone `rabbitio`:
```
cd $GOPATH/src/github.com/meltwater
git clone git@github.com:meltwater/rabbitio.git
```

3. Make:
```
cd rabbitio
make && make build
```

## Maintainers

For any bug reports or change requests, please create a Github issue or submit a PR.

Also feel free to drop a line to the maintainers:

- Joel ([@vorce](https://github.com/vorce), [joel.carlbark@meltwater.com](mailto:joel.carlbark@meltwater.com))
- Stian ([@stiangrindvoll](https://github.com/stiangrindvoll), [stian.grindvoll@meltwater.com](mailto:stian.grindvoll@meltwater.com))
- Team Blacksmiths ([all.blacksmiths@meltwater.com](mailto:all.blacksmiths@meltwater.com))
