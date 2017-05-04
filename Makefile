youtuber:
	GOOS=linux GOARCH=arm GOARM=5 go build -v
deploy: youtuber
	ssh osmc@yanzay.com -p 2222 "sudo sv stop youtuber"
	scp -P 2222 youtuber osmc@yanzay.com:~
	ssh osmc@yanzay.com -p 2222 "sudo sv start youtuber"
local: youtuber
	ssh osmc@osmc "sudo sv stop youtuber"
	scp youtuber osmc@osmc:~
	ssh osmc@osmc "sudo sv start youtuber"

