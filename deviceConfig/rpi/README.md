This is the directory for setting up mural on a raspberry pi

### Steps for setting up a rpi

run "make multi" from project root to make the image and push to docker
Ensure that there is a armv7 image available for docker rpi

Setup a raspberry pi with docker https://phoenixnap.com/kb/docker-on-raspberry-pi
It has an install script + some steps on how to run docker

Pull the container from docker
"docker pull waduphaitian/mural_dev:multi"

Setup a directory to contain the images and mural info
folder
  |                                            |
   software.json (contains mural information)   artifacts(folder containing files to distribute)

From rpi folder root
"docker run -v {Path to host folder}:containerFiles -p {Host port to use}:42069 --restart unless-stopped --name digi --detach waduphaitian/mural_dev:multi"
