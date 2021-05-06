#!/bin/bash

/race $@ 2> >(./json.sh | tee -a /1.log)
