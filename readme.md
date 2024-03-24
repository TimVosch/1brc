# 1BRC Golang
- Branches are slow
- Keep often used memory close
    - Optimize CPU caches?
    - Use pointers
- Or RtL and cut at the same time
- Parse number from LtR use:
    number
    is more? number * 10 + newnumber
    is more? number * 10 + newnumber
- Min Binary Heap

Max of 10000 items with a 100 byte identifier referencing a struct, to be ordered

Sparse set


# Implementations

## Reference
    - 209.50s user 6.39s system 101% cpu 3:32.46 total
## Custom ParseTemperature and print implementation
    - 166.71s user 7.35s system 103% cpu 2:48.49 total
    After implementing ramdisk
    - 156.51s user 8.56s system 102% cpu 2:40.32 total
## Custom batch reader implementation
    - 82.12s user 3.40s system 100% cpu 1:24.99 total
## Sparse Trie (invalid)
    Missing proper sorting

    The Trie uses a sparse set for its children. The idea was that the dense set is sorted and can always be 
    read from left to right, being alphabetically sorted.
    Where the sparse array is used to direct index the next node.

    - 87.14s user 1.63s system 100% cpu 1:28.72 total
## Regular Trie (no sparse set)
    - 60.77s user 1.70s system 100% cpu 1:02.45 total
    - 60.55s user 1.77s system 99% cpu 1:02.34 total # updated Get
    - 46.39s user 1.73s system 100% cpu 48.002 total # With TrieNode pool
    - 48.77s user 1.64s system 99% cpu 50.414 total  # Removed IsValueNode and made Station a pointer
    - 46.35s user 1.65s system 100% cpu 47.983 total # Value Station instead of pointer
