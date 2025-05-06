#!/bin/bash

cp ../dist/urpc.wasm ./static/urpc/urpc.wasm && \
cp ../dist/wasm_exec.js ./static/urpc/wasm_exec.js && \
cp ./node_modules/web-tree-sitter/tree-sitter.wasm ./static/tree-sitter.wasm && \
cp ./node_modules/curlconverter/dist/tree-sitter-bash.wasm ./static/tree-sitter-bash.wasm
