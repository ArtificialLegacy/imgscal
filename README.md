
# <div style="display: flex; align-items: center;"><img style="margin-right: 10px; width: 48px;" src="./assets/favicon-32x32.png" alt="icon"/> ImgScal </div>

A tool for automating image processing. Allows for creating of custom workflows
writing in lua.

* Run image processing tasks including AI upscaling on images or entire directories of images.
* Automatically download and manage tools such as Real-ESRGAN.
* CLI interface with portable exe.

## Build

### Windows

* Requires
  * Go
  * Makefile

```sh
make build-windows
```

## Lua Workflow Documentation

* A workflow lua file should include 2 top-level functions:

```lua
config(options) -- Sets options for the workflow when initially loaded.
```

```lua
main(func) -- Sets the function to be called when the workflow is run on a file.
```

### `config(options)`

* `options` : `{table}` - A table including the workflow config values.
* `options.name` : `{string}` - The name of the workflow to be displayed in the selection menu.
* `options.version` : `{string}` - The current version of the workflow.
* `options.requires` : `{Array<string>}` - An array of libraries that are used by the workflow.

```lua
config({
    name= "ESRGAN_X4",
    version= "1.0.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})
```

### `main(func)`

* `func` : `{function}` - The function to be called when the workflow is called.
* The function includes 1 argument as a `{string}` for the initial file name,
    this is only the file name not the path, the file is stored in the `/temp/` directory.

```lua
main(function (file)
    job("esrgan.x4", file)
    file = job("imgscal.rename", file, {prefix= "up_"})
    job("imgscal.output", file)
end)
```

### `job(job, file, options?) fileout`

* `job` : `{string}` - A string to determine what job to run in the form of `lib.job`.
* `file` : `{string}` - The file to run the job on, generally should be the file passed into the job.
* `options?` : `{table}` - An optional table for configuring how a job should be run.
* `fileout`: `{string}` - The name of the file after the job has been run, only changes on certain jobs.

* This should only be run within the main function.

```lua
file = job("imgscal.rename", file, {prefix= "up_"})
```

## Libraries

### imgscal

* This is the core library that is always available.
* Should be required for most workflows.

#### Jobs - imgscal

* `rename` - Renames a file and returns the new name.
* Available options:
  * `prefix` - Appends a string to the beginning of the filename.
  * `suffix` - Appends a string to the end of the filename (not including extension).
  * `name` - Renames the entire filename (not including extension), and happens before prefix and suffix are applied.

```lua
-- before: name.png
file = job("imgscal.rename", file, {prefix= "prefix_", suffix= "_suffix", name= "newname"})
-- after: prefix_newname_suffix.png
```

* `output` - Copies the file into the `/outputs/` directory, the outputted file should be treated as final.

```lua
job("imgscal.output", file)
```

### esrgan

* This library requires Real-ESRGAN to be installed.

#### Jobs - esrgan

* `x4` - Upscales an image using the x4 model.
* Available options:
  * `scale` - Sets the scale to upscale by, supports the following values: `[2, 3, 4]`. Defaults to `4`.

```lua
job("esrgan.x4", file)
-- or
job("esrgan.x4", file, {scale= 3})
```

* `animex4` - Upscales an image using the animex4 model.
* Available options:
  * `scale` - Sets the scale to upscale by, supports the following values: `[2, 3, 4]`. Defaults to `4`.

```lua
job("esrgan.animex4", file)
-- or
job("esrgan.animex4", file, {scale= 3})
```
