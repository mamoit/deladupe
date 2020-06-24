# DelADupe

[![Build Status](https://travis-ci.org/mamoit/deladupe.svg?branch=master)](https://travis-ci.org/mamoit/deladupe)

*THIS IS NOT PRODUCTION READY, DO NOT USE THIS!*

## Motivation
I have a folder where I've thrown stuff into for a looooong time.
It's full of cruft and duplicated data.
Some of this data also lives in my current setup, which is well organized.
But there are still some pearls in there, that I don't want to delete, but need to fish out from the mess of millions of loose files.

This program will take directory A (the old rusty backup), and directory B (the sparkling current setup where everything is neat an organized).
It will list all the files that are in A, and also in B.
After some rigorous testing I hope to be able to actually make it delete the duplicates from dir A.

## Basic operation
DelADupe takes 2 lists of folders:

*keep* - the folders that will be considered for deduplication, but no files will be deleted.
*purge* - the folders that will have files deleted from.

If a file exists multiple times, but only inside *keep* folders, all its instances will be kept.

If a file exists multiple times, but only inside *purge* folders, all instances will be deleted but one (Ordering has not been considered as of the time of writing, so which one is kept should be considered random).

If a file exists multiple times in both *keep* and *purge* folders, then all the ones in *keep* will be kept, and all the ones in *purge* will be... well... purged.

## Data structure
```
Deduper
- lock
- filesBySize:
  123:
  - lock
  - pending
  - filesByHash:
    0123456789abcdef:
    - path
```

## Questions

### Why not fdupes
Fdupes is pretty cool, but it deletes all duplications (unless you want to check one by one), and it is far from parallelized.
I don't want to delete any duplicates in dir B though.
I already ran fdupes in the messy directory, and that actually reduced the amount of files and space used by a lot, but there's still a lot of unique files in there that are in my current setup

### Don't you prefer to have the files hard linked so you can keep the structure?
No.
The objective is to eventually be able to comb through all the cruft that has been accumulating in the source directory and get rid of it.
Or at least make it manageable.
