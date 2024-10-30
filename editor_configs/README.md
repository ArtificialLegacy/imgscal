
# What are these?

These configs are examples used for setting up your editor for writing lua workflows and plugins.
This is not for development of ImgScal itself.

* LSP configs for lua should be placed in the top level directory you have open for development.
  * This can be inside `%home%/imgscal`, or in each workflow/plugin directory.
* `.ignore` is for `%home%/imgscal`. This prevents unneeded files from showing up in your fuzzyfinder.
* In `./workflow_configs` there is a `.gitignore`, this is used for excluding workflow secrets from any git repo capturing the main imgscal directory.
  * While included here, it is created automatically by ImgScal in the `%home%/imgscal/config/` directory.

## Links

### LSP Configs

* [Neovim (neoconf)](/editor_configs/lsp/.neoconf.json)

### Misc

* [.ignore](/editor_configs/.ignore)
* [workflow_configs/.gitignore](/editor_configs/workflow_configs/.gitignore)
