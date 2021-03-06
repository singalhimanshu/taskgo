# Taskgo

[![codecov](https://codecov.io/gh/singalhimanshu/taskgo/branch/main/graph/badge.svg?token=1KX3Y961FG)](https://codecov.io/gh/singalhimanshu/taskgo)

Fast and simple terminal based Kanban board manager

## Demo

![Taskgo demo](./demo/taskgo.gif)

## Features

- Vim Style keybinds
- Markdown file to store data
- Undo/Redo operations

## Install

The easiest way to get the application is to download the precompiled binaries from the [release](https://github.com/singalhimanshu/taskgo/releases) section.

### Install latest version

Make sure that you have [setup](https://golang.org/doc/install) go properly on your system and you have `$GOPATH/bin` in your `$PATH` variable (for linux/macOS) or environment variable for Windows.

```sh
$ go get -u github.com/singalhimanshu/taskgo/cmd/taskgo
```

This will create a taskgo binary under `$GOPATH/bin` directory.

## Usage

Simply run `taskgo`. This will create a taskgo.md file in your current directory.

There is a `-f` flag to provide custom file name. 

Example: `taskgo -f file_name.md`.

### Keybinds

You can press `?` in the application itself to see the keybinds. But for reference they are here as well -

| Key          | Description                    |
| ------------ | ------------------------------ |
| j/down arrow | Move down                      |
| k/up arrow   | Move up                        |
| l/h          | Move left/right                |
| J/K          | Move task down/up the list     |
| L/H          | Move task left/right the lists |
| a            | add task under the cursor      |
| A            | add task at the end of list    |
| D            | Delete a task                  |
| d            | Mark a task as done            |
| e            | Edit a task                    |
| Enter        | View task information          |
| g            | focus first item of list       |
| G            | focus last item of list        |
| u            | undo                           |
| \<C-r>      | redo                            |
| ?            | To view all these keybinds     |
| q            | Quit application               |
