# These variables are used for image building purposes only.
ARG APP

FROM scratch

ENV PORT 80
EXPOSE $PORT

COPY ${APP} /
CMD ["/${APP}"]