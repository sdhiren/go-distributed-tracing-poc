#!/bin/bash

# ./api1/main
# ./api2/main
# ./api3/main
# ./api4/main

gnome-terminal -- bash -c "./api1/main; exec bash"
gnome-terminal -- bash -c "./api2/main; exec bash"
gnome-terminal -- bash -c "./api3/main; exec bash"
gnome-terminal -- bash -c "./api4/main; exec bash"
