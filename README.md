# Taskgo

Fast and simple terminal based Kanban board manager

## Demo

![Taskgo demo](./demo/taskgo.gif)

## Features
  * Vim Style keybinds
  * Uses markdown file to store data (which is very readable)

## Install 
**Note**: Currently the only way to install taskgo is to build from source.

Make sure that you have [setup](https://golang.org/doc/install) go properly on your system and you have `$GOPATH/bin` in your `$PATH` variable (for linux/macOS) or environment variable for Windows.

```sh
$ go get github.com/singalhimanshu/taskgo
```

This will create a taskgo binary under `$GOPATH/bin` directory.

## Usage

Simply run `taskgo`. This will create a taskgo.md file in your current directory.

### Keybinds

You can press `?` in the application itself to see the keybinds. But for reference they are here as well - 

| Key              | Description                    |
| ---              | ---                            |
| j/down arrow     | Move down                      |
| k/up arrow       | Move up                        |
| l/h              | Move left/right                |
| J/K              | Move task down/up the list     |
| L/H              | Move task left/right the lists |
| a                | Add a new task                 |
| D                | Delete a task                  |
| d                | Mark a task as done            |
| e                | Edit a task                    |
| Enter            | View task information          |
| ?                | To view all these keybinds     |
| q                | Quit application               |
