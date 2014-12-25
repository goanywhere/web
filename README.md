Rex
======

Rex is a powerful starter kit for modular web applications/services in Golang.

## Getting Started

Install the package (**go 1.4** and greater is required):

```shell
$ go get -v github.com/goanywhere/rex
```

Command line tool (Optional but highly recommended)

```shell
$ go get -v github.com/goanywhere/rex/cmd/rex
```


## Features
* Flexible Env-based configurations.
* Non-intrusive/Modular design, extremely easy to use.
* Awesome routing system provided by [Gorilla/Mux](http://www.gorillatoolkit.org/pkg/mux).
* Flexible middleware system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Works nicely with other Golang packages.
* Command line tools (incl. Live reload supports).
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first server, we named it `server.go` here.

``` go
package main

import (
    "github.com/goanywhere/rex"
    "github.com/goanywhere/rex/web"
)

func main() {
    server := rex.New()
    server.Get("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := web.NewContext(w, r)
        ctx.String("Hello World")
    })
    server.Get("/hello", func(ctx *web.Context) {
        ctx.String("Hello Again")
    })
    server.Run()
}
```

Then start your server:
``` shell
rex
```

You will now have a HTTP server running on `localhost:5000`.


## Template

The standard template (html/template) package implements data-driven templates for generating HTML output safe against code injection, sounds nice? But once you step into the real world, you will soon find your code to be spaghetti. To parse multiple files with pieces of "define", say you have a "index.html", and header source defined in "header.html", footer source in "footer.html", you will need this:

```go
template.Must(template.ParseFiles("index.html", "header.html", "footer.html"))
```

What if another page say "contact.html" will share the same header & footer? Oops & yes, you'll need to do this again,

```go
template.Must(template.ParseFiles("contact.html", "header.html", "footer.html"))
```

Inheritance? They are pretty much the same, yes, you'll have to do this over & over again like this:

```go
template.Must(template.ParseFiles("layout.html", "index.html", "header.html", "footer.html"))

template.Must(template.ParseFiles("layout.html", "contact.html", "header.html", "footer.html"))
```

Rex's solution? Simple, in addition to the standard tags, we introduce two "new" (not really if you have ever used Django/Tornado/Jinja/Liquid) tags, "extends" & "include". You simply add the these two into the html pages as previous, the code will then will be like:

```go
import "github.com/goanywhere/rex/template"

loader := template.NewLoader("templates")
template := loader.Parse("index.html")
```

There you Go now, simple as that.


## Context

Context is a very useful helper shipped with Rex. It allows you to access incoming requests & responsed data, there are also shortcuts for rendering HTML/JSON/XML.


``` go
package main

import (
    "github.com/goanywhere/rex"
    "github.com/goanywhere/rex/web"
)

func index (ctx *web.Context) {
    ctx.HTML("index.html")  // Context.HTML has the extends/include tag supports by default.
}

func json (ctx *web.Context) {
    ctx.JSON(rex.H{"data": "Hello Rex", "success": true})
}

func main() {
    server := rex.New()
    server.Get("/", index)
    server.Get("/api", json)
    server.Run()
}
```


## Settings

All settings on Rex can be accessed via `rex.Settings`, which essentially stored in `os.Environ`. By using this approach you can compile your own settings files into the binary package for deployment without exposing the sensitive settings, it also makes configuration extremly easy & flexible via both command line & application.

``` go
package main

import (
    "github.com/goanywhere/env"
    "github.com/goanywhere/rex"
    "github.com/goanywhere/rex/web"
)

func index (ctx *web.Context) {
    ctx.HTML("index.html")
}

func main() {
    // Override default 5000 port here.
    env.Set("Port", "9394")

    server := rex.New()
    server.Get("/", index)
    server.Run()
}
```

You will now have the HTTP server running on `0.0.0.0:9394`.

Hey, dude, why not just use those popular approaches, like file-based config? We know you'll be asking & we have the answer as well, [here](//12factor.net/config).


## Middleware

Middlewares work between http requests and the router, they are no different than the standard http.Handler. Existing middlewares from other frameworks like logging, authorization, session, gzipping are very easy to integrate into Rex. As long as the middleware comply the `rex.Middleware` interface (shorcut to standard `func(http.Handler) http.Handler`), you can simply add one like this:

``` go
app.Use(middleware.XSRF)
```


Since the middleware is just the standard http.Handler, writing a custom middleware is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        rex.Debug("Custom Middleware Started")
        next.ServeHTTP(writer, request)
        rex.Debug("Custom Middleware Ended")
    })
})
```

## Frameworks comes & dies, will this be supported?

Positive! Rex is an internal/fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0


- [X] Sharding Supports
- [X] Env-Based Configurations
- [X] Project Home page
- [X] Test Suite
- [X] New Project Template
- [X] CLI Apps Integrations 
- [X] Improved Template Rendering
- [X] Performance Boost
- [X] Hot-Compile Runner
- [X] Live Reload Integration
- [ ] Better Logging
- [ ] Template Functions
- [ ] i18n Supports
- [ ] More Middlewares
- [ ] Form Validations
- [ ] Project Wiki
- [ ] Continuous Integration
- [ ] Stable API
