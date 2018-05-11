# Neversink filter updater

This is a simple tool for updating of [Neversink's loot filter](https://github.com/NeverSinkDev/NeverSink-Filter).


## Installation

* Download the latest version from [release section](https://github.com/ondrowan/neversink-filter-updater/releases).
* In case you prefer normal style filters, just run the file.
* In case you wish to use different style of filters:
    * Create a shortcut of exe file.
    * Right click it and select "Properties".
    * In "Shortcut" section, append `--style`, followed by desired style name behind "Target".
      Currently, valid choices are: `blue`, `purple`, `slick` and `streamsound`. The result
      should look, for example, like: `C:\neversink-filter-updater.exe --style blue`

If you don't mind using command line, here is the usage:

```bash
Usage:
  -help
        Prints help.
  -style string
        Style of filters. Can be one of: blue, purple, slick, streamsound.
  -version
        Prints version.
```

This program was tested on 64-bit Windows 10.