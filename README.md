# Muxie



https://github.com/user-attachments/assets/83fc644d-7797-4439-9c06-d0f019d1115a



Hey there! Welcome to Muxie, your new best friend for managing tmux sessions with ease. Gone are the days of manually setting up your development environment every single time. With Muxie, you can define your sessions, windows, and panes in a simple YAML file and get everything up and running in a flash.

## What's Muxie all about?

Muxie is a terminal user interface (TUI) that allows you to select and start predefined tmux sessions. You can configure your sessions in a `config.yml` file, and Muxie will present them in a list for you to choose from. Select a session, hit enter, and Muxie will take care of the rest.

## Installation

Getting Muxie up and running is a breeze.

### Using Homebrew

If you're on macOS or Linux, you can install Muxie using [Homebrew](https://brew.sh/).

```bash
brew install phanorcoll/homebrew-muxie/muxie
```

### Manual Installation

1.  Head over to the [releases page](https://github.com/phanorcoll/muxie/releases).
2.  Download the appropriate asset for your operating system.
3.  Unzip the downloaded file.
4.  Place the `muxie` binary in a directory that's in your system's `PATH`.

And that's it! You're ready to start using Muxie.

## How to use it

Using Muxie is as simple as running a single command:

```bash
muxie
```

This will launch the Muxie TUI, where you'll see a list of all the sessions you've defined in your configuration file.

### Configuration

Muxie looks for a configuration file at `~/.config/muxie/config.yml`. Here's an example of what that file might look like:

```yaml
sessions:
  - name: "My Awesome Project"
    directory: "~/projects/my-awesome-project"
    windows:
      - name: "Code"
        layout: "vertical"
        #layout: "horizontal"
        panes:
          - command: "nvim"
          - command: "git status"
      - name: "Server"
        panes:
          - command: "npm run dev"
  - name: "Another Project"
    directory: "~/projects/another-project"
    windows:
      - name: "Editor"
        panes:
          - command: "vim"
```

In this example, we have two sessions defined: "My Awesome Project" and "Another Project". Each session has a name, a directory where it should be started, and a list of windows. Each window has a name, a layout, and a list of panes. Each pane has a command that will be executed when it's created.

### Keybindings

Muxie uses a simple set of keybindings to make it easy to navigate the TUI:

*   `↑` / `k`: Move up
*   `↓` / `j`: Move down
*   `enter`: Select a session
*   `q`: Quit
*   `a`: Add new session
*   `r`: Rename existing session
*   `s`: Start a session from config.yaml
*   `d`: Kill running session

## Tmux Integration

You can integrate Muxie with your `tmux.conf` to launch it with a key binding. This allows you to quickly bring up the Muxie interface without having to type the command in a shell.

Here's an example of how you can bind the `m` key to launch Muxie in a popup window:

```tmux
bind-key m display-popup \
  -w 100% -h 100% \
  -B \
  -E "~/<path>/muxie"
```

With this configuration, pressing `prefix + m` will open Muxie in a full-screen popup, allowing you to select and start a session. Make sure to replace `~/<path>/muxie` with the actual path to your Muxie binary if it's different.

## Contributing

We love contributions! If you have an idea for a new feature or have found a bug, please open an issue on our [GitHub repository](https://github.com/phanorcoll/muxie/issues).

## License

Muxie is open-source software licensed under the [MIT License](https://opensource.org/licenses/MIT).
