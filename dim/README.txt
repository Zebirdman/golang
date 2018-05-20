Docker Image Manager:

Simple command line ap for adding prefix's to all your images in one hit.
Usefull when retagging for use with a private registry for example. can
add, replace and delete prefix's from all images. This tool will always
work on all images on your repo except for the delete prefix option which
targets a specific string.

to run this on your local machine simply build and it will be compatoble
with your docker API version. I have tested this with versions as old
as 1.24 up to latest 1.37.

Add binary to /usr/bin/here for example and use it like a commad line
app as normal. Go dim --help to see help page.

Author: Ben Futterleib
