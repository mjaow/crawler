#!/bin/bash

if [ $# -ne 2 ]
then
    echo "usage: $0 start end"
    exit 1
fi

s=$1
e=$2

echo "./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt"

./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt
