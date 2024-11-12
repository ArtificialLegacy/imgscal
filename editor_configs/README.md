
# What are these?

These configs are examples used for setting up your editor for writing lua workflows and plugins.
This is not for development of ImgScal itself.

* LSP configs for lua should be placed in the top level directory you have open for development.
  * This can be inside `%home%/imgscal`, or in each workflow/plugin directory.
  * The `.luarc.json` allows for the lsp config to be at a higher directory than the root directory, place this in the root of a workflow or plugin.
* `.ignore` is for `%home%/imgscal`. This prevents unneeded files from showing up in your fuzzyfinder.
* In `./workflow_configs` there is a `.gitignore`, this is used for excluding workflow secrets from any git repo capturing the main imgscal directory.
  * While included here, it is created automatically by ImgScal in the `%home%/imgscal/config/` directory.
* In `./schema` there are files that are used for provided JSON schema to different `.json` files used by the tool.

## Links

### LSP Configs

* [Neovim (neoconf)](/editor_configs/lsp/.neoconf.json)

### JSON Schemas

* [workflow.json](/editor_configs/schema/workflow.json)
* [config.json](/editor_configs/schema/config.json)

#### Gist Links

* [workflow.json](https://gist.githubusercontent.com/ArtificialLegacy/9711f20511e76b519aedb729a6762b9f/raw/de77e999654060a38d7a4e7eea8aeb4f5ee1273e/imgscal_workflow.json)
* [config.json](https://gist.githubusercontent.com/ArtificialLegacy/bf37b79d4fc943006f333cc35467266c/raw/933fdffd6d871d3bf5a281a7815b7d408fcd51b2/imgscal_config.json)

### Misc

* [.ignore](/editor_configs/.ignore)
* [workflow_configs/.gitignore](/editor_configs/workflow_configs/.gitignore)
