#!/bin/bash

s=145
e=145

echo "./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt"

./crawler -s $s -e $e -t dragonBallZ 2>&1 | tee log.txt
