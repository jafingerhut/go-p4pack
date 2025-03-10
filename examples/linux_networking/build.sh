#!/bin/bash
#
# Copyright 2022 - Sander Tolsma. All rights reserved
# SPDX-License-Identifier: Apache-2.0

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

docker run --rm -u 1000:1000 -v "${THIS_DIR}":/p4ccode -w /p4ccode ghcr.io/stolsma/p4c:latest /bin/bash -c "p4c-dpdk --arch pna --target dpdk -o linux_networking.spec --p4runtime-files p4info.proto.txt --bf-rt-schema bf-rt.json --context context.json linux_networking.p4"