# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.13.10 as builder

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/middleware2018-PSS/back2_school
WORKDIR $GOPATH/src/github.com/middleware2018-PSS/back2_school

ADD . .
#RUN dep ensure
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# JWT
ADD jwt_keys /home/
ENV JWT_SECRET=/home/jwt_keys/jwtRS256.key
ENV JWT_PUBLIC_KEY=/home/jwt_keys/jwtRS256.key.pub

# Casbin
ADD auth_model.conf /home/
ADD policy.csv /home/
ENV AUTH_MODEL=/home/auth_model.conf
ENV POLICY=/home/policy.csv

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD exec /bin/app
