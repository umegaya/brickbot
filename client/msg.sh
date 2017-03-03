NCPATH=`which nc`
if [ "$?" != "0" ]; then
	echo "nc required"
	exit 1
fi

cat <<Payload | nc $2 $3
{"Kind":"PlainMessage","Payload":{"message":"$1"}}
Payload
