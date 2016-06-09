:mock:template gunter <<GUNTER
#!/bin/bash

echo finn
echo sweet/bubblegum
GUNTER

:mock:template treetrunks <<TREETRUNKS
#!/bin/bash

echo ice/
echo ice/king
TREETRUNKS

:allow-system-command cat

:mock:template guntalina <<GUNTALINA
#!/bin/bash

cat >& 2
GUNTALINA

tests:ensure :bottlebreaker -a --sync
tests:assert-stderr "finn"
tests:assert-stderr "sweet/bubblegum"
tests:assert-stderr "ice/"
tests:assert-stderr "ice/king"
