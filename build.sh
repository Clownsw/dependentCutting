
if [ $1 == "windows" ];then
	go build -ldflags="-s -w" -gcflags="-B"
fi

if [ $1 == "linux" ];then
        GOOS=linux go build -ldflags="-s -w" -gcflags="-B"
fi

if [ $1 == "darwin" ];then
	GOOS=darwin go build -ldflags="-s -w" -gcflags="-B"
fi
