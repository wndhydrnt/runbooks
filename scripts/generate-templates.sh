#! /usr/bin/env bash

set -e

getRunbook=$(cat pkg/ui/getRunbook.html.mustache)
layout=$(cat pkg/ui/layout.html.mustache)
listRunbooks=$(cat pkg/ui/listRunbooks.html.mustache)

cat << EOF > pkg/ui/templates.go
package ui

var getRunbookTplString = \`${getRunbook}
\`

var layoutTplString = \`${layout}
\`

var listRunbooksTplString = \`${listRunbooks}
\`
EOF
