# This is a sample commands.py.  You can add your own commands here.
#
# Please refer to commands_full.py for all the default commands and a complete
# documentation.  Do NOT add them all here, or you may end up with defunct
# commands when upgrading ranger.

# A simple command for demonstration purposes follows.
# -----------------------------------------------------------------------------

from __future__ import absolute_import, division, print_function

# You can import any python module as needed.
import os
import subprocess
import ranger.api
import ranger.core.runner

# You always need to import ranger.api.commands here to get the Command class:
from ranger.api.commands import Command
from ranger.core.loader import CommandLoader

old_hook_init = ranger.api.hook_init


def call(cmds):
    return subprocess.Popen(cmds, stdout=subprocess.PIPE, stderr=subprocess.PIPE)


def git_exists():
    try:
        proc = call(["git", "rev-parse", "--git-dir"])
        proc.communicate()
        return proc.returncode == 0
    except OSError:
        return False


def git_mv_call(fm, flags, fname):
    thisdir = str(fm.thisdir)
    git_mv_command = ["git", "mv", *flags, fname, thisdir]
    # print out the git_mv_command
    fm.notify(" ".join(git_mv_command))
    loader = CommandLoader(git_mv_command, "git:mv")

    def reload_dir():
        thisdir.unload()
        thisdir.load_content()

    loader.signal_bind("after", reload_dir)
    fm.loader.add(loader)


# Any class that is a subclass of "Command" will be integrated into ranger as a
# command.  Try typing ":my_edit<ENTER>" in ranger!
class my_edit(Command):
    # The so-called doc-string of the class will be visible in the built-in
    # help that is accessible by typing "?c" inside ranger.
    """:my_edit <filename>

    A sample command for demonstration purposes that opens a file in an editor.
    """

    # The execute method is called when you run this command in ranger.
    def execute(self):
        # self.arg(1) is the first (space-separated) argument to the function.
        # This way you can write ":my_edit somefilename<ENTER>".
        if self.arg(1):
            # self.rest(1) contains self.arg(1) and everything that follows
            target_filename = self.rest(1)
        else:
            # self.fm is a ranger.core.filemanager.FileManager object and gives
            # you access to internals of ranger.
            # self.fm.thisfile is a ranger.container.file.File object and is a
            # reference to the currently selected file.
            target_filename = self.fm.thisfile.path

        # This is a generic function to print text in ranger.
        self.fm.notify("Let's edit the file " + target_filename + "!")

        # Using bad=True in fm.notify allows you to print error messages:
        if not os.path.exists(target_filename):
            self.fm.notify("The given file does not exist!", bad=True)
            return

        # This executes a function from ranger.core.acitons, a module with a
        # variety of subroutines that can help you construct commands.
        # Check out the source, or run "pydoc ranger.core.actions" for a list.
        self.fm.edit_file(target_filename)

    # The tab method is called when you press tab, and should return a list of
    # suggestions that the user will tab through.
    # tabnum is 1 for <TAB> and -1 for <S-TAB> by default
    def tab(self, tabnum):
        # This is a generic tab-completion function that iterates through the
        # content of the current directory.
        return self._tab_directory_content()


class git_mv(Command):
    """:git_mv

    Perform git mv on the selected files, into the current ranger directory.
    """

    def execute(self):
        if not git_exists():
            self.fm.notify("not in a git directory", bad=True)
            return

        self.fm.notify("In git. Trying to move files", bad=False)

        if self.arg(1):
            flags = self.rest(1)
        else:
            flags = []

        paths = [os.path.basename(f.path) for f in self.fm.thistab.get_selection()]

        for p in paths:
            # self.fm.notify("Moving file " + p + " to folder")
            git_mv_call(self.fm, flags, p)
