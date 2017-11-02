# Slacker [![GoDoc](https://godoc.org/pkg.re/essentialkaos/slacker.v4?status.svg)](https://godoc.org/pkg.re/essentialkaos/slacker.v4) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/slacker)](https://goreportcard.com/report/github.com/essentialkaos/slacker) [![codebeat badge](https://codebeat.co/badges/849c74bd-e041-44e6-9d9a-f2d46408b286)](https://codebeat.co/projects/github-com-essentialkaos-slacker-master) [![License](https://gh.kaos.io/ekol.svg)](https://essentialkaos.com/ekol)

`slacker` is simple go package for bootstraping Slack bots.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.5+ workspace ([instructions](https://golang.org/doc/install)), then:

````
go get pkg.re/essentialkaos/slacker.v4
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/slacker.v4
```

### License

[EKOL](https://essentialkaos.com/ekol)
