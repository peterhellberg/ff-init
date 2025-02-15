# ff-init :sparkles:

This is a command line tool that acts as a companion to my
[ff](https://github.com/peterhellberg/ff) module
for [Zig](https://ziglang.org/) :zap:

`ff-init` is used to create a directory containing code that
allows you to promptly get started coding on an app (or game) for the
lovely little WebAssembly console [Firefly Zero](https://fireflyzero.com/).

The Zig build `.target` is declared as `.{ .cpu_arch = .wasm32, .os_tag = .freestanding }`
and `.optimize` is set to `.ReleaseSmall`

> [!Important]
> No need to specify `-Doptimize=ReleaseSmall`

## Installation

(Requires you to have [Go](https://go.dev/) installed)

```sh
go install github.com/peterhellberg/ff-init@latest
```

## Usage

(Requires you to have an up to date (_nightly_) version of
[Zig](https://ziglang.org/download/#release-master) installed.

```sh
ff-init myapp
cd myapp
zig build run
```

> [!Note]
> There is also a `zig build spy` command.

:seedling:

## Screenshot

![Firefly Zero app showing the Zig logo](https://assets.c7.se/imgur/fBWmZgU.png)
