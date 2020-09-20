#!/bin/bash

# trap -l : can see signal list

on_exit() {
  echo "xxxxxxxx"
}
trap on_exit EXIT