dirs=$(find . -type d -print | grep -v "^./docs/" | grep 'apiv[1-9]/.*$')

if [ "$dirs" != "" ] ; then
    echo "FAILED: Cannot have dirs under apivX directories"
    echo "Reason:"
    echo $dirs
    echo ""
    exit 1
fi
