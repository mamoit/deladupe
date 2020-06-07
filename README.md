# DelADupe

I have a folder where I've thrown stuff into for a looooong time.
It's full of cruft and duplicated data.
Some of this data also lives in my current setup, which is well organized.
But there are still some pearls in there, that I don't want to delete, but need to fish out from the mess of millions of loose files.

This program will take directory A (the old rusty backup), and directory B (the sparkling current setup where everything is neat an organized).
It will list all the files that are in A, and also in B.
After some rigorous testing I hope to be able to actually make it delete the duplicates from dir A.

## Questions

### Why not fdupes
Fdupes is pretty cool, but it deletes all duplications (unless you want to check one by one).
I don't want to delete any duplicates in dir B though.
I already ran fdupes in the messy directory, and that actually reduced the amount of files and space used by a lot, but there's still a lot of unique files in there that are in my current setup

### Don't you prefer to have the files hard linked so you can keep the structure?
No.
The objective is to eventually be able to comb through all the cruft that has been accumulating in the source directory and get rid of it.
Or at least make it manageable.

## TODO
* Only calculate hashes of a file in the target directory if one in the source directory is found with the exact same size.

