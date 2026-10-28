# Qoin Shor AppImage

The desktop entry and AppImage recipe use `io.postquantumqoin.shor` as the
application ID. The icon is a simple packaging-specific Q monogram.

## Build

Build on Debian 12 x86_64 with the normal Shor native build dependencies and
`appimage-builder` installed:

```sh
make appimage GOCC=/usr/local/go/bin/go
```

The target builds `shor`, stages the desktop entry and icon in `AppDir`, and
passes the recipe explicitly to `appimage-builder`. The recipe bundles the
Debian 12 C/C++ runtime packages; additional native dependencies must be checked
against the actual release binary before distribution.

The default package version comes from `git describe --tags --always --dirty`.
For a release build, override it explicitly:

```sh
make appimage APPIMAGE_VERSION=1.0.0
```

The generated artifact is `QoinShor-<version>-x86_64.AppImage` in the project root.

## Run

Shor is a command-line client, not a graphical wallet. The desktop entry uses a
terminal and does not automatically start a node. With no arguments it displays
CLI help; the terminal may close afterwards depending on the desktop environment.
For persistent interaction, use an existing terminal:

```sh
./QoinShor-<version>-x86_64.AppImage --help
./QoinShor-<version>-x86_64.AppImage --version
```

Node data, keys and user configuration are not included in the package. This
packaging change does not migrate the proof-parameter cache or protocol names.
