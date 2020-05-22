#!/usr/bin/env bash

set -e

# Helpful defines for the provisioning process.
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BUILD_DIR="$SCRIPT_DIR/../build"
FORMULA_DIR="$SCRIPT_DIR/provision/formula"

HOMEBREW_REPO="https://github.com/Homebrew/brew"
LINUXBREW_REPO="https://github.com/Linuxbrew/brew"

HOMEBREW_CORE_REPO="https://github.com/Homebrew/homebrew-core"
LINUXBREW_CORE_REPO="https://github.com/Linuxbrew/homebrew-core"

# Set the SHA1 commit hashes for the pinned homebrew Taps.
# Pinning allows determinism for bottle availability, expect to update often.
HOMEBREW_CORE="941ca36839ea354031846d73ad538e1e44e673f4"
LINUXBREW_CORE="f54281a496bb7d3dd2f46b2f3067193d05f5013b"
HOMEBREW_BREW="ac2cbd2137006ebfe84d8584ccdcb5d78c1130d9"
LINUXBREW_BREW="20bcce2c176469cec271b46d523eef1510217436"

# These suffixes are used when building bottle tarballs.
LINUX_BOTTLE_SUFFIX="x86_64_linux"
DARWIN_BOTTLE_SUFFIX="sierra"

# If the world needs to be rebuilt, increase the version
DEPS_VERSION="5"

source "$SCRIPT_DIR/lib.sh"
source "$SCRIPT_DIR/provision_lib.sh"

function platform_linux_main() {
  echo
}

function platform_darwin_main() {
  echo
}

 function platform_posix_main() {
  echo
}

function sysprep(){
  if [[ "$1" = "build" ]]; then
      cd third-party/osquery
    MAKELEVEL_TEMP=$MAKELEVEL
    unset MAKELEVEL
    git branch -D feature/in_proc_extensions
    git checkout feature/in_proc_extensions
    git pull
    #make -j$THREADS deps
    make clean
    make -j$THREADS libosquery
    make -j$THREADS libosquery_additional
    export MAKELEVEL=$MAKELEVEL_TEMP
  else
    log "your $OS does not use a provision script"
  fi
} 

#ToDo: Figure out how to install cloudquery specific dependencies
function main() {

  ACTION=$1

  platform OS
  distro $OS DISTRO
  threads THREADS

  if ! hash sudo 2>/dev/null; then
    echo "Please install sudo in this machine"
    exit 1
  fi

  # Setup the cloudquery dependency directory.
  # One can use a non-build location using CLOUDQUERY_DEPS=/path/to/deps
  if [[ ! -z "$CLOUDQUERY_DEPS" ]]; then
    DEPS_DIR="$CLOUDQUERY_DEPS"
  else
    DEPS_DIR="/usr/local/cloudquery"
  fi

  if [[ "$ACTION" = "clean" ]]; then
    do_sudo rm -rf "$DEPS_DIR"
    return
  fi

  export DEPS_DIR=$DEPS_DIR
  #Link osquery dependency to cloudquery
  do_sudo ln -snf /usr/local/osquery $DEPS_DIR


  # Setup the local ./build/DISTRO cmake build directory.
  if [[ ! -z "$SUDO_USER" ]]; then
    echo "chown -h $SUDO_USER $BUILD_DIR/*"
    chown -h $SUDO_USER:$SUDO_GID "$BUILD_DIR" || true
    if [[ $OS = "linux" ]]; then
      chown -h $SUDO_USER:$SUDO_GID "$BUILD_DIR/linux" || true
    fi
    chown $SUDO_USER:$SUDO_GID "$WORKING_DIR" > /dev/null 2>&1 || true
  fi

  # Provisioning uses either Linux or Home (OS X) brew.
  if [[ $OS = "darwin" ]]; then
    BREW_TYPE="darwin"
  elif [[ $OS = "freebsd" ]]; then
    BREW_TYPE="freebsd"
  else
    BREW_TYPE="linux"
  fi

    
  if [[ "$BREW_TYPE" = "darwin" ]]; then
    platform_darwin_main
  elif [[ "$BREW_TYPE" = "linux" ]]; then
    platform_linux_main
  fi

  cd "$SCRIPT_DIR/../"

  # Additional compilations may occur for Python and Ruby
  export LIBRARY_PATH="$DEPS_DIR/legacy/lib:$DEPS_DIR/lib:$LIBRARY_PATH"
  set_cc clang
  set_cxx clang++

  initialize $OS
  sysprep $ACTION
  
}

check $1 "$2"
main $1 "$2"
