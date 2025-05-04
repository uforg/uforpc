#!/bin/bash

WEBTREESITTER=$(find node_modules/.deno \
  -type f \
  -path "*/web-tree-sitter@*/node_modules/web-tree-sitter/tree-sitter.wasm" \
  -print -quit)

cp ../dist/urpc.wasm ./static/urpc/urpc.wasm && \
cp ../dist/wasm_exec.js ./static/urpc/wasm_exec.js && \
cp $WEBTREESITTER ./static/tree-sitter.wasm && \
cp ./node_modules/curlconverter/dist/tree-sitter-bash.wasm ./static/tree-sitter-bash.wasm
