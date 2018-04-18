# Helios

This project helps you manage your command line.
Create favourite directories with ease
Use regex to navigate
Autocomplete on your favourite directories
https://unix.stackexchange.com/questions/1800/how-to-specify-a-custom-autocomplete-for-specific-commands
add script to your /bin
hc -s 'executable name' 'type your command here'
use -sb to have helios save the command internally instead of putting a script in your bin.

hc -e /folder/path/
to export your settings
hc -i /dir/path/
to import settings

hc -f nickname /path/to/directory

helios will prioritize nickname directories over those in the current folder when using helios cd

hc -f will show your your favorite directories and their nicknames.

use hc -r for a regex directory search - it will use the first match it finds and navigate there
-f will also favourite this path

hc -h will output a history of nav
hc -h 1 will naviagte to the absolute dir of the last directory navigated to



