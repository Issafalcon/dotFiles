#!/bin/zsh

# Derive JAVA_HOME from `javac`, not `java`: update-alternatives can point
# each at a different JDK (e.g. a JRE-only package registers `java` but has
# no `javac`), so using `java` can land on a JAVA_HOME with no compiler.
if command -v javac >/dev/null 2>&1; then
  export JAVA_HOME="$(dirname "$(dirname "$(readlink -f "$(command -v javac)")")")"
fi
