set -e

output=$(npm exec -c 'nitro-rpc-client address')
echo $output

if echo $output | grep -q '0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE' ; then
    echo "got expected address" && exit 0
else
    echo "did not get expected address" && exit 1
fi