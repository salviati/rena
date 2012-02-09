# Rename Archive.
A tool crafted to rename files in anime/manga archives/dirs.
Also can be used to list missing/duplicate episodes.
  
To see what //would// be done
  `rena -n -r -N 3 -t 'Naruto_Shippuuden_{N}' /path/to/Naruto/Shippuden/dir`
Make sure you're okay with this list. Then
  `rena -y -r -N 3 -t 'Naruto_Shippuuden_{N}' /path/to/Naruto/Shippuden/dir`
  (you can omit `-y` if you want to confirm each step)
Just looking for duplicates? Or missing volumes in your manga collection?
  `rena -r -n /path/to/Naruto/manga/dir`

## How Does It Work?
  - Drop the predefined list of strings from the filename (such as "480p" or 8-byte hexadecimal numbers in the given filename (most likely CRC32))
  - Pick the list of numbers in the filename
  - Assume that the first number is the episode number
  - Include the following sequential numbers (if any)
  - Output the new name using the given template
To demonstrate, the filename "[DB]_Naruto_Shippuuden_99_-_100_-_101_720p.avi[CFA31234]" will yield the 8-byte hex number CFA31234 and 720p, which are discarded readily. Second step gives 99, 100, 101. First number is 99. Next number 100 is sequential, so it's in. 101 is sequential again, so it's in as well.
Assuming the template was 'Naruto Shippuuden {N}' along with options `-s '-' -N 3`, the output name is "Naruto Shippuuden 099-100-101.avi"

There are cases this algorithm would fail. "FMA2_37.avi", for instance, would yield 2 as the episode number, and 37 would be discarded. To remedy such situations, there's crop option. -C regexp will crop what matches regexp beforehand. Adding -C 'FMA2' would solve the problem in our case.
Sometimes, a troublesome checksum appears before the episode name, such as "[12345678]Nyan Koi 01.avi". We, then, need to drop the hexadecimal number enclosed with [] to handle the problem in general, which can be achieved using -C '^\[[0-9A-Za-z]+\]'.

You see that there are various exceptions; and you may need to rename your archives in several steps, starting from these exceptions (once exceptions are handled, they're likely to be handled properly in a "normal" rename).

## License
GNU General Public License, version 3 or later.
See the COPYING file for details.
