# 1BRC Golang
- Floating arithmetic is slow
    - Always 1 decimal so its 10x 
- Branches are slow
- Keep often used memory close
    - Optimize CPU caches?
    - Use pointers


# Implementations

## Reference
    - 209.50s user 6.39s system 101% cpu 3:32.46 total
## Custom ParseTemperature and print implementation
    - 166.71s user 7.35s system 103% cpu 2:48.49 total
    After implementing ramdisk
    - 156.51s user 8.56s system 102% cpu 2:40.32 total
## Custom batch reader implementation
    - 82.12s user 3.40s system 100% cpu 1:24.99 total
