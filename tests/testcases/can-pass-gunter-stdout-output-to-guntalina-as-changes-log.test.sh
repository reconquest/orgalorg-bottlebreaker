:mock:template gunter <<GUNTER
#!/bin/bash

echo finn
echo sweet/bubblegum
GUNTER

:allow-system-command cat

:mock:template guntalina <<GUNTALINA
#!/bin/bash

echo guntalina args: "\${@}" >& 2

cat >& 2
GUNTALINA

:mock:true treetrunks

tests:ensure :bottlebreaker -a --sync
tests:assert-stderr "guntalina args: -s /dev/stdin"
tests:assert-stderr "finn"
tests:assert-stderr "sweet/bubblegum"
