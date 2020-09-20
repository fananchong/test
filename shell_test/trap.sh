#!/bin/bash

# trap -l : can see signal list

on_exit() {
  echo "xxxxxx"
}
trap on_exit EXIT