#!/bin/sh

set -e

exec reflex --regex='\.(go|html|prompt)$' --start-service --decoration=none -- "$@"