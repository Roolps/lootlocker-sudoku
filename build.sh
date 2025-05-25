# lol ignore this line - hot patch so I can build it :,)
chown -R root:root /var/lib/pterodactyl/volumes/d7b6a87e-5f25-42dd-a983-7d1781c20788/public/
cd ./public
npm run build
cd ../
go build main.go