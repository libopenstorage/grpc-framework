# gRPC Framework
[![Build Status](https://app.travis-ci.com/libopenstorage/grpc-framework.svg?branch=master)](https://app.travis-ci.com/libopenstorage/grpc-framework)

This framework makes it simple for developers to add gRPC and automated REST
interfaces for their Golang applications.

## Usage

## Usage

Please see our [Documentation](https://libopenstorage.github.io/grpc-framework) for
more information.

* FYI, we are in the process of creating tutorials and adding more documentation
to the framework

## grpc-framework Development

### Documentation for the website

* Setup the environment to write documentation

```
$ make doc-env
```

* Bring up the webserver

```
$ make doc-serve
```

Now edit the files in the direcory `website`.

* When done, build the website and add the `docs` dir to git:

```
$ make doc-build
```

For more information on mkdocs go to:

* https://www.mkdocs.org/
* https://squidfunk.github.io/mkdocs-material/getting-started/ 
