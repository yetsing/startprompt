#!/usr/bin/env bash

# 运行所有的例子

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# 当前文件和文件夹
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__cwd="$(pwd)"
# 项目根目录
__root="$(cd "$(dirname "${__dir}")" && pwd)"

function main() {
  cd "${__root}"
  echo "============== clockinput =============="
  go run ./examples/clockinput/clockinput.go
  echo "============== complete ================"
  go run ./examples/complete/complete.go
  echo "============== echo ===================="
  go run ./examples/echo/echo.go
  echo "============== multiline ==============="
  go run ./examples/multiline/multiline.go
  echo "============== persistenthistory ======="
  go run ./examples/persistenthistory/persistenthistory.go
  echo "============== prompt =================="
  go run ./examples/prompt/prompt.go
  echo "============== syntaxheightlight ======="
  go run ./examples/syntaxheightlight/syntaxheightlight.go
  cd "${__cwd}"
}

main