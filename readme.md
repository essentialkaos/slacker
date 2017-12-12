<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/slacker.v6"><img src="https://godoc.org/pkg.re/essentialkaos/slacker.v6?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/slacker"><img src="https://goreportcard.com/badge/github.com/essentialkaos/slacker"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-slacker-master"><img src="https://codebeat.co/badges/849c74bd-e041-44e6-9d9a-f2d46408b286"></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.io/ekol.svg"></a>
</p>

<p align="center"><a href="#readme"><img src="https://gh.kaos.io/slacker.svg"/></a></p>

`slacker` is simple go package for bootstraping Slack bots.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.7+ workspace ([instructions](https://golang.org/doc/install)), then:

````
go get pkg.re/essentialkaos/slacker.v5
````

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/slacker.v5
```

### License

[EKOL](https://essentialkaos.com/ekol)
