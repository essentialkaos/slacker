<p align="center"><a href="#readme"><img src="https://gh.kaos.st/slacker.svg"/></a></p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/essentialkaos/slacker"><img src="https://pkg.go.dev/badge/github.com/essentialkaos/slacker" /></a>
  <a href="https://travis-ci.org/essentialkaos/slacker"><img src="https://travis-ci.org/essentialkaos/slacker.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/slacker"><img src="https://goreportcard.com/badge/github.com/essentialkaos/slacker"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-slacker-master"><img src="https://codebeat.co/badges/849c74bd-e041-44e6-9d9a-f2d46408b286"></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

`slacker` is simple go package for bootstraping Slack bots.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.12+ workspace ([instructions](https://golang.org/doc/install)), then:

````
go get pkg.re/essentialkaos/slacker.v9
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/slacker.v9
```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/slacker.svg?branch=master)](https://travis-ci.org/essentialkaos/slacker) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/slacker.svg?branch=develop)](https://travis-ci.org/essentialkaos/slacker) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
