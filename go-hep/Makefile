all: compile docker push

compile:
	CGO_ENABLED=0 go build -o gohep

docker:
	sudo docker build --no-cache=true --force-rm=true -t dwhitena/diana-go-hep .

push:
	sudo docker push dwhitena/diana-go-hep

