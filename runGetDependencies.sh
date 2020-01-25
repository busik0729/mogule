#!/bin/bash
file="package.txt"
for var in $(cat $file)
do
go get $var
done