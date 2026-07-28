# Background

The `dotFiles` tui app takes the existing modular (stowed) dotfile installation and turns it into an interactive TUI app. The app will allow users to easily manage their dotfiles, including installing, updating, and removing them.

# Features

## Interactive TUI built with https://github.com/charmbracelet/lipgloss library

The app will replace the `bootstrap.sh` script with an interactive TUI built using the `lipgloss` library. This will provide a more user-friendly interface for managing dotfiles.

- If possible, use typescript for the app to take advantage of type checking and improved developer experience, but if not possible, javascript is also fine.

- Each module listed in the repo already will be replaced with a corresponding "module" file in the TUI app. These module files will contain the necessary information and commands to manage the corresponding dotfiles.
- Each module should still use stowe to manage the dotfile configurations, but the TUI app will provide a more intuitive way to interact with these configurations.
- The TUI app will also include features such as progress bars, status indicators, and error messages to provide feedback to the user during the installation and management process.

## Pre-requisite installation

- When the app is launched, there should be some pre-requisite checks that ensure that all the things installed by the pre-requisite installation script are present. If any of the pre-requisites are missing, the app should prompt the user to install them before proceeding. If the pre-requisites are missing, the main dashboard should not be displayed, and instead, the user should be guided through the installation process for the missing pre-requisites.
- This "pre-requisite" initial page will explain what the dependencies are, and why they are needed, with links to documentation / repos for each where appropriate. The user should be able to easily install the missing dependencies from this page, and once all dependencies are installed, they can proceed to the main dashboard.

## Layout

- The app should be presented as a single page application in a floating terminal window (or give the floating appearance in the current terminal window) and should fill around 80% of the terminal width and height, with a border around it to give it a distinct appearance.
- The layout should be clean and organized, with clear sections for different functionalities
- The main dashboard should have a sidebar on the left that lists all the available modules (dotfiles) and a main content area on the right that displays the details and options for the selected module.
  - The module list should list each module in its own bordered box, with the name of the module beside an icon for the module (e.g. neovim module would have a neovim icon), a brief description of the app the module installs, and estimated install time and size, with a link to any application websites / repos. An icon should also be present, or background colour changed to indicate which modules have been installed. The module list should also include a search bar at the top to allow users to quickly find specific modules.
- On navigating over a module, the main content area on the right should update to show the details of the selected module.
  - This content area will have multiple tabs for different sections of the module:
    1.  Overview: Summary from any github README file about the module, with links to the relevant documentation and repos.
        - List of dependencies for the module with a column indicating their install method (e.g. npm, cargo, apt, brew, etc.) and a tick or cross icon indicating whether the dependency is already installed on the user's system.
          - From this list, the user can navigate up and down and choose whether to install the dependency directly from the tui app, which will run the appropriate installation command for the dependency (e.g. `npm install -g <dependency>` for npm packages, `cargo install <dependency>` for cargo packages, etc.).
        - Whether the dependency exists should be checked asynchronously and spinner icon used to indicate when the check is in progress.
        - A keyboard shortcut should be available to install the module and all it's dependencies with a single command - A popup to confirm should appear before installing (confirming what will be installed)
        - Any tools to aid in installation (e.g. homebrew, cargo, npm, etc.) should also be included in the dependencies list where required, and installed first when they are needed for the module installation.
    2.  Output window: This tab will show the output of any installation commands run for the module, including progress bars and status indicators to provide feedback to the user during the installation process.
        - Any required input from the user during installation should be displayed as a popup over the main tui app, allowing the user to provide the necessary input without leaving the app or changing tabs.
    3.  Configuration: This tab will provide options for configuring the module, such as selecting specific features to install, choosing between different versions of the software, or customizing the installation process in other ways. The options available in this tab will depend on the specific module and its requirements.

## Controls

The user should have a list of sane keyboard shortcuts that are context aware and displayed in a help menu (e.g. `?` key to open the help menu) to allow for easy navigation and management of the app. For example, keyboard shortcuts could include:

- `j` and `k` to navigate up and down the module list
- `enter` to select a module and view its details in the main content area
- `i` to install the selected module and its dependencies
- `d` to uninstall the selected module
- `s` to open the search bar and quickly find a specific module
- `?` to open the help menu with a list of all available keyboard shortcuts and their functions
- `q` to quit the app
- `tab` to switch between different tabs in the main content area (e.g. Overview, Output window, Configuration)
- `esc` to cancel any ongoing actions or close popups
- `o` to open urls in the default web browser (e.g. for documentation links, module repos, etc.)

Shortcuts should not interfere with normal terminal keybinds (e.g. `ctrl + c` to copy, `ctrl + v` to paste, etc.) and should be designed to be intuitive and easy to remember for users familiar with terminal applications.

## Themes

- The app should use a neon / cyberpunk inspired colour scheme, similar to dracula type colours.

## Other features

- New modules should be easy to add to the app, with a clear structure for how module files should be created and organized within the app's codebase. This will allow for easy expansion of the app's functionality as new dotfiles and configurations are added to the repo, in a similar way to how new modules can be added to the existing stow-based dotfiles repo.
- zsh integration: The app should be designed to work seamlessly with zsh, allowing users to easily manage their zsh configurations and plugins through the TUI interface. This could include features such as automatic detection of existing zsh configurations, easy installation and management of zsh plugins, and integration with popular zsh frameworks like oh-my-zsh or prezto.
- Module installation should be able to occur in parallel, with progress bars and status indicators for each module being installed, to provide feedback to the user during the installation process. This will allow users to efficiently manage their dotfiles without having to wait for each module to be installed sequentially. Conflicting dependencies for parallel running installations should be handled gracefully (i.e. The app should keep track of which dependencies are being installed, what order they need to be installed in, and then coordinate parallel installations accordingly - e.g. if one module requires a dependency that is currently being installed by another module, the app should wait for the first installation to complete before starting the second installation, rather than trying to install both at the same time and potentially causing conflicts or errors).
- After each installation, the installed app should be verified before proceeding by checking for the existence of the app's executable in the system path, and if the verification fails, an error message should be displayed to the user with options to troubleshoot the issue (e.g. check installation logs, retry installation, etc.). Some installations may require refreshing of the terminal environment (e.g. to update the system path after installing a new app), so the app should also include functionality to refresh the terminal environment as needed after installations, without interrupting the user's workflow or requiring them to manually restart their terminal session.

## Guidance

As this is my first tui app, please provide inline documentation in the code to help me learn how to write tui apps in javascript / typescript using the lipgloss library.

- Guidance around adding new modules should also be included, as well as a main develop guide for running and building the app.
