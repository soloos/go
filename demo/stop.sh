#!/bin/bash
ps aux|grep gopool|grep -v grep|awk '{print $2}'|xargs kill -9
