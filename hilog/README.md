# hilog : highlight logs

Highlight log from various source with pre-defined parsing rule, LSP.

> Re-implemented with Go. Checkout branch, python for previous Python version.

![hilog_screenshot](https://lh6.googleusercontent.com/-JMq2QwuQh2g/UQYj9RFOvlI/AAAAAAAACHg/FFw3V7zg_9E/s912/hilog-go_dev-android.png)

## Install

[Install Go][1] and run;

    $ go get github.com/suapapa/hilog

Make sure you have `$GOPATH/bin` in your `$PATH`

[1]: http://golang.org/doc/install

## Choose LogSchemePack

You can use bundled LSP or can write your own in JSON.

### Bundles

    $ hilog -l
    red-stderr
    gotest
    android-logcat
    linux-kmsg

### External scheme written in JSON

Log `Source` can be selected:

- stdout (default)
- stderr

Search by `Ptn` in following search `Type`s:

- re
- startswith
- endwiths
- contains

Colorize by `FG`, `BG`, `Attrs`
Available formatting constants are:

    FG, BG  : BLACK, RED, GREEN, YELLOW, BLUE, MAGENTA, CYAN, WHITE
    Attrs   : BOLD, UNDERLINE, BLINK, INVERSE

You can dump bundled scheme to JSON with `-p` for reference.

    suapapa $ ./hilog -p linux-kmsg
    [
        {
            "Ptn": "^\u003c[012]\u003e",
            "Type": "re",
            "Source": "",
            "FG": "YELLOW",
            "BG": "RED",
            "Attrs": [
                "BOLD"
            ]
        },
    ...


## Usage

### Example

    $ hilog android-logcat adb logcat
    $ hilog linux-kmsg adb shell cat /proc/kmsg
    $ hilog red-stderr make -j6

### Make handy aliases

Add following lines to `~/.bash_aliases`

    alias hi-logcat='hilog android-logcat adb logcat '
    alias hi-kmsg='hilog linux-kmsg adb shell cat /proc/kmsg '

