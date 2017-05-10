#!/usr/bin/env bash

if [ "$(uname)" != "Darwin" ]; then
    echo "This script has not been tested on this platform!  (edit the script to remove this check)"
    exit 1
fi

# bail on errors, treat unset vars as errors, disable filename expansion (globbing), produce an error if any command in a pipeline fails (rather than just the last)
set -euf -o pipefail


# Use absolute paths for commands to prevent alias-related errors
ECHO=/bin/echo
DIRNAME=/usr/bin/dirname
FIND=/usr/bin/find
GREP=/usr/bin/grep

SCRIPT_DIR="$( cd "$( $DIRNAME "${BASH_SOURCE[0]}" )" && pwd )"
: ${SCRIPT_DIR:?}


# Terminal colors!
NORMAL=
RED=
GREEN=
YELLOW=
BOLD=
UNDERLINE=
NOLINE=
TERMINAL_COLORS=0

# Detect support for colors
if [ -t 1 ]; then  # is stdout a terminal?
    TERMINAL_COLORS=$(tput colors)
    if [[ TERMINAL_COLORS && TERMINAL_COLORS -ge 8 ]]; then  # does the terminal support at least 8 colors?
        NORMAL=$(tput sgr0)
        RED=$(tput setaf 1)
        GREEN=$(tput setaf 2)
        YELLOW=$(tput setaf 3)
        BOLD=$(tput bold)
        UNDERLINE=$(tput smul)
        NOLINE=$(tput rmul)
    fi
fi

ERROR=${BOLD}${RED}
WARN=${BOLD}${YELLOW}


printWarn() {
    $ECHO "${WARN}WARNING${NORMAL} - $@" 1>&2;
}

printError() {
    $ECHO "${ERROR}ERROR${NORMAL} - $@" 1>&2;
}


# first arg (optional): message to echo
# second arg (optional): error code to exit with (defaults to 1)
abort() {
    local returnValue=${2:-1}
    if [ "${returnValue}" == "" ]; then
        returnValue=1
    fi

    if [[ "$1" != "" ]]; then
        if [ $returnValue -eq 0 ]; then
            $ECHO "$1"
        else
            printError "$1"
        fi
    fi

    exit $returnValue
}




ROOT_DIR="$SCRIPT_DIR/.."
: ${ROOT_DIR:?}

$ECHO "Generating..."
# clean up any previous failed builds
cd $ROOT_DIR && $FIND . -type f -name "ffjson-inception*.go" -delete || abort "Failed to remove stale ffjson temp files"
cd $ROOT_DIR && $FIND . -type d -name "ffjson-inception*" -delete || abort "Failed to remove stale ffjson temp directories"
# Then normal cleanup
cd $ROOT_DIR && $FIND . -name "*_ffjson.go" -delete || abort "Failed to remove generated ffjson files"
cd $ROOT_DIR && go generate "./..." || abort "'go generate' failed"



# Check for structs that cannot be generated, where ffjson has fallen back to the built-in reflection-based encoder
# (see https://github.com/pquerna/ffjson, tip 4)

# Grep will exit with 1 if no lines are selected (which is success in our case), so we need to temporarily disable
# bail-on-error so that the script will continue running...
set +e
fallback=$(cd $ROOT_DIR && $GREP -Ir "Falling back" . --include="*_ffjson.go")
set -e

# Now check for errors in grep itself (exit codes greater than 1)
if [ $? -gt 1 ]; then
    abort "grep command failed"
fi

if [[ "$fallback" != "" ]]; then
    printWarn "Fallback detected in the following files:"
    $ECHO "$fallback"
    exit 1
fi

$ECHO "${GREEN}${BOLD}SUCCESS${NORMAL}"
