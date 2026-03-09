# Copyright (c) 2025 IBM Corp.
# All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Alpine-based image with OpenSSL for contract-cli
# This image is built by GoReleaser — the binary is pre-compiled and copied in.
# GoReleaser's dockers_v2 API places binaries in $TARGETPLATFORM/ subdirectories.

FROM alpine:3.23

# Install OpenSSL (required for encryption operations) and ca-certificates
RUN apk add --no-cache \
    openssl \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Copy the pre-built binary from GoReleaser (dockers_v2 context layout)
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/contract-cli /usr/bin/contract-cli

# Copy documentation
COPY LICENSE /usr/share/doc/contract-cli/LICENSE
COPY README.md /usr/share/doc/contract-cli/README.md

# Run as non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["contract-cli"]
