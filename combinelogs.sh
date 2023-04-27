cd artifacts
for d in */ ; do
    echo "$d"
    cat $d/*.log | pino-pretty -S -o "{engine} {To} < {From}" > $d/combined.tmp
done
