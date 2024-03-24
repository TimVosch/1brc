#!/env/bash

mkdir ramdisk
mount -t tmpfs -o size=15G tmpfs ramdisk
cp measurements.txt ramdisk/measurements.txt
