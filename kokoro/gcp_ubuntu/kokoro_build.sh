#!/bin/bash

# Fail on any error.
set -e

# Install Bazel.
use_bazel.sh 4.2.1
command -v bazel
bazel version

# Code under repo is checked out to ${KOKORO_ARTIFACTS_DIR}/git.
# The final directory name in this path is determined by the scm name specified
# in the job configuration.
cd "${KOKORO_ARTIFACTS_DIR}/git/localtoast"
bazel build ...
