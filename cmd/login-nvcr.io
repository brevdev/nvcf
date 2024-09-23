#!/bin/bash
set -euo pipefail
[ -z "${NGC_API_KEY}" ] && echo "NGC_API_KEY is not set" && exit 1

echo "${NGC_API_KEY}" | docker login nvcr.io -u '$oauthtoken' --password-stdin
