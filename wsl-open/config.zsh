# Adding wsl-open as a browser for Bash for Windows
if [[ $(uname -r) =~ (m|M)icrosoft ]]; then
  if [[ -z $BROWSER ]]; then
    export BROWSER=wsl-open
  else
    export BROWSER=$BROWSER:wsl-open
  fi
fi
