# Dockercraft

## Aspects

2 Sides

- Server Side
- Build side

# start create container

`docker run -p 25565:25565 --mount type=bind,src="$(pwd)"/plugins,dst=/plugins --mount type=bind,src="$(pwd)"/static,dst=/static -i -t -d --name dockercraft-c dockercraft`

## Startup Path

- Check if image exists -> build
- Check if container exists -> create
- Start container
- Check if all plugins built?
- Build plugins if not
