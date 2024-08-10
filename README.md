
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
