# Prototypes
`kbackup` started life as many backup programs do; as a set of bash scripts. Bash scripts can be quick-and-easy
automation solutions but they have many issues assiocaited with them. The biggest of these is the text-based nature of
the environment.

The pragmatic programmer talks about prototypes. They cover the topic briefly [here](https://www.codingblocks.net/podcast/the-pragmatic-programmer-tracer-bullets-and-prototyping/)
but I recommend reading the [actual book](https://www.amazon.com/Pragmatic-Programmer-journey-mastery-Anniversary/dp/0135957052/ref=sr_1_1?hvadid=580689974848&hvdev=c&hvlocphy=9031945&hvnetw=g&hvqmt=e&hvrand=10297413813783756329&hvtargid=kwd-343418855&hydadcr=16404_13419760&keywords=the+pragmatic+programmer&qid=1707464628&sr=8-1)
.

Our bash script eventually ended up being the prototype and looked something like this:

```bash
#!/usr/bin/env bash
set -eo pipefail

[[ -z $1 ]] && echo "No sources provided" && exit 1
[[ -z $2 ]] && echo "No destination provided" && exit 1

SOURCES="${@:1:$#-1}"
DESTINATION="${@: -1}"

for SOURCE in $SOURCES
do
    TARGET=`basename $SOURCE`
    mkdir -pv $DESTINATION/$TARGET/.kbackup

    CURRENT=`date -Is | head -c 19 | sed "s/:/-/g"`

    if [[ -e $DESTINATION/$TARGET/.kbackup/last ]]
    then
        LAST=`cat $DESTINATION/$TARGET/.kbackup/last`
        echo "> Rotating last backup: $DESTINATION/$TARGET/.kbackup/$LAST"
        mkdir -pv $DESTINATION/$TARGET/.kbackup/$LAST
        ls -A $DESTINATION/$TARGET | grep -v '.kbackup' | while read f; do echo "$DESTINATION/$TARGET/$f"; done \
        | xargs -d '\n' cp -al -t $DESTINATION/$TARGET/.kbackup/$LAST
    else
        echo "> No existing backups"
    fi

    echo "> Backing up: $SOURCE -> $DESTINATION/$TARGET"
    rsync -hav --delete --exclude .kbackup "$SOURCE/" "$DESTINATION/$TARGET" || (
        [[ $LAST ]] &&
        echo "> Cleaning up: $DESTINATION/$TARGET/.kbackup/$LAST" &&
        rm -rf $DESTINATION/$TARGET/.kbackup/$LAST &&
        exit 1
    )

    echo "$CURRENT" > $DESTINATION/$TARGET/.kbackup/last

done
```
