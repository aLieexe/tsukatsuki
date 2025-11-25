## Description
The `init` command initializes Tsukatsuki in the current project directory. It walks you through an interactive prompt to collect required information (runtime, server host, user preferences) and generates the base configuration file `tsukatsuki.yaml`.
Additionally, it prepares the initial directory structure under `.deploy/` unless specified otherwise

## Usage
```bash
tsukatsuki init [options]
```

## Options
```
--output         Specify output directory
```
## Examples
```bash
tsukatsuki init
tsukatsuki init --output=newdir
```