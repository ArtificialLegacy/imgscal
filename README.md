
# ![icon](./pkg/assets/icons/favicon-32x32.png) ImgScal

Automate image processing programmatically.

* Built around concurrency.
* Workflows writteng in lua.
* Builtin handling for image encodings and color models.
* Included ImGui wrapper for building custom GUI tools.
* Spritesheet support.
* Command-line support, e.g. `imgscal resize ./image.png 100 100`.

## Documentation - [Live](https://artificiallegacy.github.io/imgscal/)

Run `make doc` to generate the lua api documentation to `./docs/`.

## Run

```sh
make start
# or
make dev
```

* `make dev` runs `make examples` then `make start`.
* Use `make start` for a first time launch, as `make examples` requires a valid config.

## Examples

> Examples can be installed by running `make examples`.
> This requires ImgScal to have been run at least **once**,
> as the config file and workflow directory must exist.

### GUI/Imger - [Source](/examples/gui/imger.lua)

Demo workflow that creates an interface with controls to apply filters from the Imger library onto images.

![imger example](assets/demos/example_imger.png)

### GUI/Noise - [Source](/examples/gui/noise.lua)

Demo workflow that creates an interface with controls to generate noise maps.

![noise example](assets/demos/example_noise.png)

### GUI/Filter - [Source](/examples/gui/filter.lua)

Demo workflow that creates an interface to apply filters to an image.

![filter example](assets/demos/example_filters.png)

## Build

* Requires
  * Go >= 1.22.6
  * Makefile
  * A C compiler (mingw, TDM-GCC or g++)
  * OpenCL

```sh
make build-windows
# or
make build-linux
```

## Install

### From Local

* Requires
  * Go >= 1.22.6
  * Makefile
  * A C compiler (mingw, TDM-GCC or g++)
  * OpenCL

```sh
make install-slim
# or
make install
```

Tools installed by each command:

* **`make install-slim`**
  * `imgscal`
* **`make install`**
  * `imgscal`
  * `imgscal-new`
  * `imgscal-entrypoint`
  * `imgscal-log`

To best make use of these, `$GOPATH/bin` should be added to your system path.

## Additional Tools

> These tools will require ImgScal to have been run at least once before using.

### `imgscal-new`

A terminal form for quickly creating a new `workflow.json`. This also validates that the created workflow is unique.

Install with `go install ./cmd/imgscal-new`.

### `imgscal-entrypoint`

Command for adding new entry points to the current workflow.

```sh
imgscal-entrypoint <name> <path> [-c]
```

* `<name>`: The name of the entry point, use `\*` to bind to workflow name.
* `<path>`: Path including the `.lua` file to create as the entry point. Any subdirectories included will be created if needed.
* `[-c]`: Optional flag to create a cli entry point.
* `[-w]`: Optional flag to to set the relative path to search for the `workflow.json` file.

Install with `go install ./cmd/imgscal-entrypoint`.

### `imgscal-workspace`

Command for quickly creating an empty `.luarc.json` file. This allows the Lua LSP to find the root of a workflow or plugin.

```sh
imgscal-workspace
```

Install with `go install ./cmd/imgscal-workspace`.

### `imgscal-log`

Prints the log file `@latest.txt` if it exists. `make log` is a shortcut for calling this.

> This can also be used to pipe the output.

```sh
make log | grep '! ERROR'
make log | kate -i
make log > latest.txt
# or if installed
imgscal-log | grep '! ERROR'
```

Install with `go install ./cmd/imgscal-log`.

### `imgscal-examples`

Copies all workflows from the `./examples/` directory to your workflow directory.

Example workflows are copied into a directory named `examples`, and this directory is cleared each time. Any additional workflows added here will be deleted.

### `imgscal-doc`

Generates documentation for the lua api to `./docs/`.
When using an up to date version of the tool from the main branch,
refer to the [live site](https://artificiallegacy.github.io/imgscal/) instead.

### `imgscal-types`

Generates type definition files for the lua api to `./types/`.
Copy these to the library directory set in your editor config.

## Editor Configs - [README](/editor_configs/README.md)

Included in this repo is `./editor_configs/`. These can be used as a reference for setting up your dev environment for workflows and plugins.

## Known Issues

> These have been tentatively fixed.

* There is an upstream issue related to `mainthread` when running on Windows.
  * Sometimes causes the application to crash within the first frame of the master window.
  * Sometimes prevents `gui.update()` from updating the UI when called from a scheduled function.
  * Recommended to run through WSL when using the gui library on Windows.
* Some parts of the `gui` library may not work properly.
  * `gui.css_parse()` is added but the underlying `g.ParseCSSStyleSheet()` is currently broken.
  * `gui.wg_css_tag()` can still be used, but will not have any affect until the upstream issue is fixed.
* When a lua panic occurs outside of the `gui.window_run()` loop, it may cause the window to not close until the process is closed.
  * Looking for a solution as calling `.Close()` on an already closed window causes GLFW to break until the process is restarted, and there is no publically exported field to check if a window is active.
