#!/bin/bash

s=207
e=218

echo "./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt"

./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt
