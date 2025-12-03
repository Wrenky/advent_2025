#!/bin/bash
name="$1"

mkdir $name

cp template.go $name/main.go
touch $name/demo
touch $name/input

ls -lash $name/

cd $name
