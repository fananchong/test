#!/bin/bash

/race $@ 2> >(./date.sh | tee -a /1.log)
