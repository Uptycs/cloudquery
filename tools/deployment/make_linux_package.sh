#!/usr/bin/env bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOURCE_DIR="$SCRIPT_DIR/../.."
BUILD_DIR=${BUILD_DIR:="$SOURCE_DIR/build/linux"}

CLOUDQUERY_DEPS="${CLOUDQUERY_DEPS:-/usr/local/cloudquery}"

export PATH="${CLOUDQUERY_DEPS}/bin:$PATH"
source "$SOURCE_DIR/tools/lib.sh"

VERSION=`(cd $SOURCE_DIR; git describe --tags HEAD) || echo 'unknown-version'`
PACKAGE_VERSION=${CLOUDQUERY_BUILD_VERSION:="$VERSION"}
PACKAGE_ARCH="x86_64"
PACKAGE_TYPE=""
PACKAGE_ITERATION=""
DESCRIPTION="cloudquery is an cloudware instrumentation toolchain."
PACKAGE_NAME="cloudquery"
#if [[ $PACKAGE_VERSION == *"-"* ]]; then
#  DESCRIPTION="$DESCRIPTION (unstable/latest version)"
#fi

# Config files
INITD_SRC="$SCRIPT_DIR/cloudqueryd.initd"
INITD_DST="/etc/init.d/cloudqueryd"
SYSTEMD_SERVICE_SRC="$SCRIPT_DIR/cloudqueryd.service"
SYSTEMD_SERVICE_DST="/usr/lib/systemd/system/cloudqueryd.service"
SYSTEMD_SYSCONFIG_SRC="$SCRIPT_DIR/cloudqueryd.sysconfig"
SYSTEMD_SYSCONFIG_DST="/etc/sysconfig/cloudqueryd"
SYSTEMD_SYSCONFIG_DST_DEB="/etc/default/cloudqueryd"
CTL_SRC="$SCRIPT_DIR/cloudqueryctl"
PACKS_SRC="$SOURCE_DIR/packs"
PACKS_DST="/usr/share/cloudquery/packs/"
CLOUDQUERY_POSTINSTALL=${CLOUDQUERY_POSTINSTALL:-""}
CLOUDQUERY_PREUNINSTALL=${CLOUDQUERY_PREUNINSTALL:-"$SCRIPT_DIR/linux_prerm.sh"}
CLOUDQUERY_CONFIG_SRC=${CLOUDQUERY_CONFIG_SRC:-""}
CLOUDQUERY_TLS_CERT_CHAIN_SRC=${CLOUDQUERY_TLS_CERT_CHAIN_SRC:-""}
CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_SRC="${CLOUDQUERY_DEPS}/etc/openssl/cert.pem"
CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_DST="/usr/share/cloudquery/certs/certs.pem"
CLOUDQUERY_EXAMPLE_CONFIG_SRC="$SCRIPT_DIR/cloudquery.example.conf"
CLOUDQUERY_EXAMPLE_CONFIG_DST="/usr/share/cloudquery/cloudquery.example.conf"
CLOUDQUERY_LOG_DIR="/var/log/cloudquery/"
CLOUDQUERY_VAR_DIR="/var/cloudquery"
CLOUDQUERY_ETC_DIR="/etc/cloudquery"

WORKING_DIR=/tmp/cloudquery_packaging
INSTALL_PREFIX=$WORKING_DIR/prefix
DEBUG_PREFIX=$WORKING_DIR/debug

function usage() {
  fatal "Usage: $0 -t deb|rpm -i REVISION -d DEPENDENCY_LIST
    [-u|--preuninst] /path/to/pre-uninstall
    [-p|--postinst] /path/to/post-install
    [-c|--config] /path/to/embedded.config
  This will generate an Linux package with:
  (1) An example config /usr/share/cloudquery/cloudquery.example.conf
  (2) An init.d script /etc/init.d/cloudqueryd
  (3) A systemd service file /usr/lib/systemd/system/cloudqueryd.service and
      a sysconfig file /etc/{default|sysconfig}/cloudqueryd as appropriate
  (4) A default TLS certificate bundle (provided by cURL)
  (5) The cloudquery toolset /usr/bin/cloudquery*"
}

function parse_args() {
  while [ "$1" != "" ]; do
    case $1 in
      -t | --type )           shift
                              PACKAGE_TYPE=$1
                              ;;
      -i | --iteration )      shift
                              PACKAGE_ITERATION=$1
                              ;;
      -d | --dependencies )   shift
                              PACKAGE_DEPENDENCIES="${@}"
                              ;;
      -u | --preuninst)       shift
                              CLOUDQUERY_PREUNINSTALL=$1
                              ;;
      -p | --postinst )       shift
                              CLOUDQUERY_POSTINSTALL=$1
                              ;;
      -c | --config )         shift
                              CLOUDQUERY_CONFIG_SRC=$1
                              ;;
      -h | --help )           usage
                              ;;
    esac
    shift
  done
}

function check_parsed_args() {
  if [[ $PACKAGE_TYPE = "" ]] || [[ $PACKAGE_ITERATION = "" ]]; then
    usage
  fi
}

function get_pkg_suffix() {
  if [[ $PACKAGE_TYPE == "deb" ]]; then
    # stay compliant with Debian package naming convention
    echo "_${PACKAGE_VERSION}_${PACKAGE_ITERATION}.amd64.${PACKAGE_TYPE}"
  elif [[ $PACKAGE_TYPE == "rpm" ]]; then
    V=`echo ${PACKAGE_VERSION}|tr '-' '_'`
    echo "-${V}-${PACKAGE_ITERATION}.${PACKAGE_ARCH}.${PACKAGE_TYPE}"
  elif [[ $PACKAGE_TYPE == "pacman" ]]; then
    echo "-${PACKAGE_VERSION}-${PACKAGE_ITERATION}-${PACKAGE_ARCH}.pkg.tar.xz"
  else
    echo "-${PACKAGE_VERSION}_${PACKAGE_ITERATION}_${PACKAGE_ARCH}.tar.gz"
  fi
}

function main() {
  parse_args $@
  check_parsed_args

  platform OS
  distro $OS DISTRO

  OUTPUT_PKG_PATH=`readlink --canonicalize "$BUILD_DIR"`/$PACKAGE_NAME$(get_pkg_suffix)

  rm -rf $WORKING_DIR
  rm -f $OUTPUT_PKG_PATH
  mkdir -p $INSTALL_PREFIX

  log "copying cloudquery binaries"
  BINARY_INSTALL_DIR="$INSTALL_PREFIX/usr/bin/"
  mkdir -p $BINARY_INSTALL_DIR
  cp "$BUILD_DIR/cloudquery/cloudqueryd" $BINARY_INSTALL_DIR
  ln -s cloudqueryd $BINARY_INSTALL_DIR/cloudqueryi
  strip $BINARY_INSTALL_DIR/*
  cp "$CTL_SRC" $BINARY_INSTALL_DIR

  # Create the prefix log dir and copy source configs
  log "copying cloudquery configurations"
  mkdir -p $INSTALL_PREFIX/$CLOUDQUERY_VAR_DIR
  mkdir -p $INSTALL_PREFIX/$CLOUDQUERY_LOG_DIR
  mkdir -p $INSTALL_PREFIX/$CLOUDQUERY_ETC_DIR
  mkdir -p $INSTALL_PREFIX/$PACKS_DST
  mkdir -p `dirname $INSTALL_PREFIX$CLOUDQUERY_EXAMPLE_CONFIG_DST`
  cp $CLOUDQUERY_EXAMPLE_CONFIG_SRC $INSTALL_PREFIX$CLOUDQUERY_EXAMPLE_CONFIG_DST
  cp $PACKS_SRC/* $INSTALL_PREFIX/$PACKS_DST

  if [[ $CLOUDQUERY_CONFIG_SRC != "" ]] && [[ -f $CLOUDQUERY_CONFIG_SRC ]]; then
    log "config setup"
    cp $CLOUDQUERY_CONFIG_SRC $INSTALL_PREFIX/$CLOUDQUERY_ETC_DIR/cloudquery.conf
  fi

  if [[ $CLOUDQUERY_TLS_CERT_CHAIN_SRC != "" ]] && [[ -f $CLOUDQUERY_TLS_CERT_CHAIN_SRC ]]; then
    log "custom tls server certs file setup"
    cp $CLOUDQUERY_TLS_CERT_CHAIN_SRC $INSTALL_PREFIX/$CLOUDQUERY_ETC_DIR/tls-server-certs.pem
  fi

  if [[ $CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_SRC != "" ]] && [[ -f $CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_SRC ]]; then
    log "built-in tls server certs file setup"
    mkdir -p `dirname $INSTALL_PREFIX/$CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_DST`
    cp $CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_SRC $INSTALL_PREFIX/$CLOUDQUERY_TLS_CERT_CHAIN_BUILTIN_DST
  fi

  if [[ $PACKAGE_TYPE = "deb" ]]; then
    #Change config path to Ubuntu default
    SYSTEMD_SYSCONFIG_DST=$SYSTEMD_SYSCONFIG_DST_DEB
  fi

  log "copying cloudquery init scripts"
  mkdir -p `dirname $INSTALL_PREFIX$INITD_DST`
  mkdir -p `dirname $INSTALL_PREFIX$SYSTEMD_SERVICE_DST`
  mkdir -p `dirname $INSTALL_PREFIX$SYSTEMD_SYSCONFIG_DST`
  cp $INITD_SRC $INSTALL_PREFIX$INITD_DST
  cp $SYSTEMD_SERVICE_SRC $INSTALL_PREFIX$SYSTEMD_SERVICE_DST
  cp $SYSTEMD_SYSCONFIG_SRC $INSTALL_PREFIX$SYSTEMD_SYSCONFIG_DST

  if [[ $PACKAGE_TYPE = "deb" ]]; then
    #Change config path in service unit
    sed -i 's/sysconfig/default/g' $INSTALL_PREFIX$SYSTEMD_SERVICE_DST
  fi

  log "creating $PACKAGE_TYPE package"
  IFS=',' read -a deps <<< "$PACKAGE_DEPENDENCIES"
  PACKAGE_DEPENDENCIES=
  for element in "${deps[@]}"
  do
    element=`echo $element | sed 's/ *$//'`
    PACKAGE_DEPENDENCIES="$PACKAGE_DEPENDENCIES -d \"$element\""
  done

  # Let callers provide their own fpm if desired
  FPM=${FPM:="fpm"}

  if [[ $CLOUDQUERY_POSTINSTALL = "" ]]; then
    if [[ $PACKAGE_TYPE == "rpm" ]]; then
      CLOUDQUERY_POSTINSTALL="$SCRIPT_DIR/rpm_postinstall.sh"
    else
      CLOUDQUERY_POSTINSTALL="$SCRIPT_DIR/deb_postinstall.sh"
    fi
  fi

  POSTINST_CMD=""
  if [[ $CLOUDQUERY_POSTINSTALL != "" ]] && [[ -f $CLOUDQUERY_POSTINSTALL ]]; then
    POSTINST_CMD="--after-install $CLOUDQUERY_POSTINSTALL"
  fi

  PREUNINST_CMD=""
  if [[ $CLOUDQUERY_PREUNINSTALL != "" ]] && [[ -f $CLOUDQUERY_PREUNINSTALL ]]; then
    PREUNINST_CMD="--before-remove $CLOUDQUERY_PREUNINSTALL"
  fi

  # Change directory modes
  find $INSTALL_PREFIX/ -type d | xargs chmod 755

  EPILOG="--url https://www.uptycs.com \
    -m contact@uptycs.com              \
    --vendor 'Uptycs Inc'              \
    --license BSD                      \
    --description \"$DESCRIPTION\""

  CMD="$FPM -s dir -t $PACKAGE_TYPE \
    -n $PACKAGE_NAME -v $PACKAGE_VERSION \
    --iteration $PACKAGE_ITERATION \
    -a $PACKAGE_ARCH               \
    --log error                    \
    --config-files $INITD_DST      \
    --config-files $SYSTEMD_SYSCONFIG_DST \
    $PREUNINST_CMD                 \
    $POSTINST_CMD                  \
    $PACKAGE_DEPENDENCIES          \
    -p $OUTPUT_PKG_PATH            \
    $EPILOG \"$INSTALL_PREFIX/=/\""
  eval "$CMD"
  log "package created at $OUTPUT_PKG_PATH"

  # Generate debug packages for Linux or CentOS
  BUILD_DEBUG_PKG=false
  if [[ $PACKAGE_TYPE = "deb" ]]; then
    BUILD_DEBUG_PKG=true
    PACKAGE_DEBUG_NAME="$PACKAGE_NAME-dbg"
    PACKAGE_DEBUG_DEPENDENCIES="cloudquery (= $PACKAGE_VERSION-$PACKAGE_ITERATION)"

    # Debian only needs the non-stripped binaries.
    BINARY_DEBUG_DIR=$DEBUG_PREFIX/usr/lib/debug/usr/bin
    mkdir -p $BINARY_DEBUG_DIR
    cp "$BUILD_DIR/cloudquery/cloudqueryd" $BINARY_DEBUG_DIR
    strip --only-keep-debug "$BINARY_DEBUG_DIR/cloudqueryd"
    ln -s cloudqueryd $BINARY_DEBUG_DIR/cloudqueryi
  elif [[ $PACKAGE_TYPE = "rpm" ]]; then
    BUILD_DEBUG_PKG=true
    PACKAGE_DEBUG_NAME="$PACKAGE_NAME-debuginfo"
    PACKAGE_DEBUG_DEPENDENCIES="cloudquery = $PACKAGE_VERSION"

    # Create Build-ID links for executables and Dwarfs.
    BUILD_ID=`readelf -n "$BUILD_DIR/cloudquery/cloudqueryd" | grep "Build ID" | awk '{print $3}'`
    if [[ ! "$BUILD_ID" = "" ]]; then
      BUILDLINK_DEBUG_DIR=$DEBUG_PREFIX/usr/lib/debug/.build-id/${BUILD_ID:0:2}
      mkdir -p $BUILDLINK_DEBUG_DIR
      ln -sf ../../../../bin/cloudqueryd $BUILDLINK_DEBUG_DIR/${BUILD_ID:2}
      ln -sf ../../bin/cloudqueryd.debug $BUILDLINK_DEBUG_DIR/${BUILD_ID:2}.debug
    fi

    # Install the non-stripped binaries.
    BINARY_DEBUG_DIR=$DEBUG_PREFIX/usr/lib/debug/usr/bin/
    mkdir -p $BINARY_DEBUG_DIR
    cp "$BUILD_DIR/cloudquery/cloudqueryd" "$BINARY_DEBUG_DIR/cloudqueryd.debug"
    strip --only-keep-debug "$BINARY_DEBUG_DIR/cloudqueryd.debug"
    ln -s cloudqueryd "$BINARY_DEBUG_DIR/cloudqueryi.debug"

    # Finally install the source.
    SOURCE_DEBUG_DIR=$DEBUG_PREFIX/usr/src/debug/cloudquery-$PACKAGE_VERSION
    BUILD_DIR=`readlink --canonicalize "$BUILD_DIR"`
    SOURCE_DIR=`readlink --canonicalize "$SOURCE_DIR"`
    for file in `"$SCRIPT_DIR/getfiles.py" --build "$BUILD_DIR/" --base "$SOURCE_DIR/"`
    do
      mkdir -p `dirname "$SOURCE_DEBUG_DIR/$file"`
      cp "$SOURCE_DIR/$file" "$SOURCE_DEBUG_DIR/$file"
    done
  fi

  PACKAGE_DEBUG_DEPENDENCIES=`echo "$PACKAGE_DEBUG_DEPENDENCIES"|tr '-' '_'`
  OUTPUT_DEBUG_PKG_PATH=`readlink --canonicalize "$BUILD_DIR"`/$PACKAGE_DEBUG_NAME$(get_pkg_suffix)
  if [[ "$BUILD_DEBUG_PKG" = "true" ]]; then
    rm -f $OUTPUT_DEBUG_PKG_PATH
    CMD="$FPM -s dir -t $PACKAGE_TYPE            \
      -n $PACKAGE_DEBUG_NAME -v $PACKAGE_VERSION \
      -a $PACKAGE_ARCH                           \
      --iteration $PACKAGE_ITERATION             \
      --log error                                \
      -d \"$PACKAGE_DEBUG_DEPENDENCIES\"         \
      -p $OUTPUT_DEBUG_PKG_PATH                  \
      $EPILOG \"$DEBUG_PREFIX/=/\""
    eval "$CMD"
    log "debug created at $OUTPUT_DEBUG_PKG_PATH"
  fi
}

main $@
