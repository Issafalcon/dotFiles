#!/bin/zsh

# Derive JAVA_HOME from the active `java` binary (works across JDK versions)
if command -v java >/dev/null 2>&1; then
  export JAVA_HOME="$(dirname "$(dirname "$(readlink -f "$(command -v java)")")")"
fi
