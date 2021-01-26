#PACKAGE = github.com/eric135/YaNFD-packages/bestroute
STRATEGY = bestroute

.PHONY: all clean

all:
	go build -buildmode=plugin -o ${STRATEGY}.so ${STRATEGY}.go

clean:
	rm ${STRATEGY}.wasm
