BUILD=go build
OUT_LINUX=./bin/myweb
OUT_WIN=./bin/myweb_unsigned.exe
SIGNED_WIN_OUT=./bin/myweb.exe
WIN_RESOURCE=./resource.syso
SRV_KEY=./server.key
SRV_PEM=./server.pem
KEY=$(shell cat ${SRV_KEY} | base64)
CERT=$(shell cat ${SRV_PEM} | base64)
WIN_SRC=.
NIC_SRC=*.go
WIN_LDFLAGS=--trimpath --ldflags "-s -w -X main.key=${KEY} -X main.cert=${CERT} -H=windowsgui"
LINUX_LDFLAGS=--trimpath --ldflags "-s -w -X main.key=${KEY} -X main.cert=${CERT}"

depends:
	openssl req -subj '/emailAddress=abuse@acme.com/CN=ACME Company CA/O=ACME Company/C=US' -new -newkey rsa:4096 -days 365 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}

windows32:
	goversioninfo -icon=./resource/icon.ico ./resource/verioninfo.json
	GOOS=windows GOARCH=386 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WIN} ${WIN_SRC}
	osslsigncode sign -certs ${SRV_PEM} -key ${SRV_KEY} -n "Fallingelk" -i http://www.acme.com -in ${OUT_WIN} -out ${SIGNED_WIN_OUT}
	osslsigncode verify -verbose -CAfile ${SRV_PEM} ${SIGNED_WIN_OUT}

windows64:
	goversioninfo -icon=./resource/icon.ico ./resource/verioninfo.json
	GOOS=windows GOARCH=amd64 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WIN} ${WIN_SRC} 
	osslsigncode sign -certs ${SRV_PEM} -key ${SRV_KEY} -n "Fallingelk" -i http://www.acme.com -in ${OUT_WIN} -out ${SIGNED_WIN_OUT}
	osslsigncode verify -verbose -CAfile ${SRV_PEM} ${SIGNED_WIN_OUT}

linux32:
	GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

macos64:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

clean:
	rm -f ${OUT_LINUX} ${OUT_WIN} ${WIN_RESOURCE} ${SIGNED_WIN_OUT} ${WIN_RESOURCE}


