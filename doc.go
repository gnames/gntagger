/*
gntagger is a command line application that helps to find/curate scientific
names interactively. For example if there is a monograph about a genus with
hundreds of scientific names, gntagger will find names automatically and then
will let the user to verify found names interactively.

Asciicast: https://asciinema.org/a/wNfIt2TfZiyrAwJZKhuq5DkHV

The user interface of the program consists of 2 panels. The left panel
contains detected scientific names, with a "current name" located in the middle
of the screen and highlighted. The left panel contains the full text, where
the "current name" is highlighted and aligned with the "current name" in the
left panel.

The program is designed to move though the names quickly. Navigate to the
next/previous name in the left panel using Right/Left arrow keys. All names
have an empty annotation at the beginning. Pressing Right Arrow key
automatically "accepts" found name if the annotation is empty. Other keys
allow to annotate the "current name" differently:

* Space: rejects a name with "NotName" annotation

* 'y':   re-accepts mistakenly rejected name with "Accepted" annotation

* 'u':   marks a name as "Uninomial"

* 'g':   marks a name as "Genus"

* 's':  marks a name as "Species"

* 'd':  marks a name as "Doubtful"

* Ctrl-C: saves curation and exits application

* Ctrl-S: saves curations made so far

The program autosaves results of curation. If the program crashes, or exited
the user can continue curation at the last point instead of starting from
scratch.
*/
package gntagger
