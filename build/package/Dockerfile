# build from scratch image
FROM scratch
# set label to github repository
LABEL github-repository="https://github.com/mmichaelb/gosharexserver"
# set working directory
WORKDIR /app/
# add compiled binary to root path
ADD ./gosharexserver .
# start application located at /app/gosharexserver
ENTRYPOINT ["./gosharexserver"]
