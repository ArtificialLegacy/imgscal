
# ![icon](./assets/favicon-32x32.png) ImgScal

A tool for automating image processing. Allows for creating of custom workflows
writing in lua.

* Run image processing tasks written using lua.
* CLI interface with portable exe.

## Build

### Windows

* Requires
  * Go
  * Makefile

```sh
make build-windows
```

**❗ Note: Removes the build directory, including any custom made workflows.**

## Documentation

Run `make doc` to generate the lua api documentation to `./docs/`.
