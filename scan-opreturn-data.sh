#!/bin/bash

for block in output/op-returns/*/
do
    for f in $block/*.dat
    do
        file $f | grep -v ': data'
    done
done