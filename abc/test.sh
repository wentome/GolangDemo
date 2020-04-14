#! /bin/bash

files=$(ls abc | grep -v $1)
echo $files
