#!/bin/bash

if [ $# < 1 ]; then
      usage
      exit 1
fi

while getopts "p:d:" arg #选项后面的冒号表示该选项需要参数
do
    case $arg in
        p)
            projectname=$OPTARG
            ;;
        d)
            directory=$OPTARG
            ;;
        ?)  #当有不认识的选项的时候arg为?
            echo "unknown argument"
            exit 1
            ;;
    esac
done

usage() {
    echo "the script require at least one parameter to work"
    echo "Usage: sh launch.sh"
    echo "-p your project name"
    echo "-d your project path in /vagrant"
}

start() {
    cd /vagrant
    tar -zxf $projectname.tar.gz
    cd $directory
    docker build -t $directory:test .
}

start