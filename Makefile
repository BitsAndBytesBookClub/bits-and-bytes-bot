TAGNAME=bits-and-bytes-bot:latest

.PHONY: docker-build
docker-build:
	@sudo docker build -t $(TAGNAME) .
	@printf "\nYou can run a container with the following command\n\tsudo docker run $(TAGNAME)\n\n"
