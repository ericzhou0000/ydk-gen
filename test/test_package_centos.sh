#!/bin/bash
#  ----------------------------------------------------------------
# Copyright 2016 Cisco Systems
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------
#
# Script for running ydk CI on travis-ci.org
#
# ------------------------------------------------------------------

# Terminal colors
RED="\033[0;31m"
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

######################################################################
# Utility functions
######################################################################

function print_msg {
    echo -e "${MSG_COLOR}*** $(date): test_package_centos.sh | $@ ${NOCOLOR}"
}

function init_confd {
    cd $1
    print_msg "Initializing confd in $(pwd)"
    source $YDKGEN_HOME/../confd/confdrc
    make stop > /dev/null
    make clean > /dev/null
    make all > /dev/null
    make start
    cd -
}

########################## EXECUTION STARTS HERE #############################
######################################
# Parse args
######################################
PYTHON_VERSION=""

args=$(getopt p:d $*)
set -- $args
PYTHON_VERSION=${2}

PYTHON_BIN=python${PYTHON_VERSION}

if [[ ${PYTHON_VERSION} = *"2"* ]]; then
    PIP_BIN=pip
elif [[ ${PYTHON_VERSION} = *"3.5"* ]]; then
    PIP_BIN=pip3
else
    PIP_BIN=pip${PYTHON_VERSION}
fi

print_msg "Using ${PYTHON_BIN} & ${PIP_BIN}"

######################################
# Set up env
######################################
export YDKGEN_HOME="$(pwd)"

init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/ydktest

######################################
# Install
######################################
print_msg "Installing libydk"
yum install -y libydk*.rpm

print_msg "Installing ydk-py"
cd sdk/python/core
python setup.py sdist
${PIP_BIN} install -v dist/ydk*.tar.gz

######################################
# Test
######################################
print_msg "Running basic python test"
${PYTHON_BIN} -c 'import os;from ydk.path import Repository; from ydk.providers import NetconfServiceProvider; repo=Repository(os.getcwd()); p=NetconfServiceProvider("127.0.0.1","admin","admin",12022,repo=repo)'
status=$?
if [ $status -ne 0 ]; then
    MSG_COLOR=$RED
    print_msg "Failed execute basic Python test"
    exit $status
fi
