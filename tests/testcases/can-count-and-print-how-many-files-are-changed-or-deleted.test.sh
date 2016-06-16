:mock:template gunter <<GUNTER
#!/bin/bash

echo finn
echo sweet/bubblegum
GUNTER

:mock:template treetrunks <<TREETRUNKS
#!/bin/bash

echo ice/
echo ice/king
echo ice/crown
TREETRUNKS

:allow-system-command cat

:mock:true guntalina

tests:ensure :bottlebreaker -a --sync
tests:assert-stderr "2 files changed"
tests:assert-stderr "3 files removed"
