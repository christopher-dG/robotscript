# RobotScript

[![Build Status](https://travis-ci.com/christopher-dG/robotscript.svg?branch=master)](https://travis-ci.com/christopher-dG/robotscript)

**RobotScript is a scriptable, cross-platform desktop automation tool.**

## Installation

Statically linked binaries will be available for download [here](https://github.com/christopher-dG/robotscript/releases) once the project is functional.
In the meantime, feel free to build it yourself.

## Usage

```sh
robotscript script.yml
```

## Script Reference

Script files are written in [YAML format](http://yaml.org/).
The commands are fairly self-explanatory/self-documenting.
Here is a comprehensive example:

```yaml
commands:
  - mouse:                 # command: moves the mouse cursor
      x: 200               # required: integer: x pixel
      y: 200               # required: integer: y pixel
      relative: false      # optional: boolean: moves cursor relative to original position
  - click:                 # command: clicks a mouse button
      button: left         # required: string: mouse button (left, right, or center)
  - keypress:              # command: presses a key
      key: c               # required: string: key to be pressed (case insensitive)
      mods:                # optional: list<string>: modifier keys
        - shift            # note: this must be a list, even if there's only one
  - sleep:                 # command: does nothing for some time
      seconds: 1           # required: integer: seconds to wait
  - type:                  # command: types some text
      text: hello, world   # required: string: text to be typed
  - exec:                  # command: executes a program
      program: stat        # required: string: program to run
      args: /tmp/fakefile  # optional: string: command line arguments (not a list)
```

For a key naming reference, see [here](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md).


### Acknowledgements

The following open source libraries have made my life much easier:

* [robotgo](https://github.com/go-vgo/robotgo)
* [mapstructure](https://github.com/mitchellh/mapstructure)
