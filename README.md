
# ![icon](./assets/favicon-32x32.png) ImgScal

A tool for automating image processing. Allows for creating of custom workflows
writing in lua.

* Run image processing tasks written using lua.
* CLI interface with portable exe.

## Examples

### GUI/Noise

Demo workflow that creates an interface with controls to generate noise maps.

> [Source File](/workflows/example_noise.lua)

![noise example](assets/demos/example_noise.png)

### GUI/Filter

Demo workflow that creates an interface to apply filters to an image.

> [Source File](/workflows//example_filter.lua)

![filter example](assets/demos/example_filters.png)

## Known Issues

* There is an upstream issue related to `mainthread` when running on Windows.
  * Sometimes causes the application to crash within the first frame of the master window.
  * Sometimes prevents `gui.update()` from updating the UI when called from a scheduled function.
  * Recommended to run through WSL when using the gui library on Windows.
* Some parts of the `gui` library may not work properly.
  * `gui.css_parse()` is added but the underlying `g.ParseCSSStyleSheet()` is currently broken.
    * `wg_css_tag()` can still be used, but will not have any affect until the upstream issue is fixed.
  * Bindings for the text editor widget have not been added, as it is currently disabled upstream.
  * Bindings for the markdown widget have also not been added for the same reason.
  * No lua bindings for manually pushing and popping `wg_style` as widgets only exist as tables within lua.
* When a lua panic occurs outside of the `gui.window_run()` loop, it may cause the window to now close until the process is closed.
  * Looking for a solution as calling `.Close()` on an already closed window causes GLFW to break until the process is restarted, and there is no publically exported field to check if a window is active.
* It is currently possible to deadlock in certain circumstances.
  * Passing the same collection item twice into a function that schedules on it. e.g. `image.draw()`
  * Calling a function that schedules on a collection item within a function already running for that
    collection item. e.g. Calling `image.size()` within the callback of `image.map()` for the same images.

## Build

### Windows

* Requires
  * Go
  * Makefile
  * Either mingw or TDM-GCC to use the gui library

```sh
make build-windows
```

**❗ Note: Removes the build directory, including any custom made workflows.**

## Documentation

Run `make doc` to generate the lua api documentation to `./docs/`.

## Logs

The make file includes a few shortcuts for log files:

* make log - prints latest log to terminal using cat.
* make logview - opens latest log in notepad.
* make logclear - rm all log files.
