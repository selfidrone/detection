FROM denismakogon/gocv-runtime:edge
#FROM denismakogon/gocv-build-stage:edge

#FROM alpine:3.7

#RUN echo -e \
#  '@edgunity http://nl.alpinelinux.org/alpine/edge/community\n\
#@edge http://nl.alpinelinux.org/alpine/edge/main\n\
#@testing http://nl.alpinelinux.org/alpine/edge/testing\n\
#@community http://dl-cdn.alpinelinux.org/alpine/edge/community'\
#  >> /etc/apk/repositories
#
#RUN apk update && \
#    apk add --upgrade apk-tools@edge && \
#    apk --no-cache add ca-certificates \
#    libstdc++ libjpeg libtbb@testing libpng jasper-libs tiff openblas libwebp

# Add non root user
RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app
RUN chown app /home/app

# Copy OpenCV Libs
#COPY --from=0 /usr/local/lib64 /usr/local/lib64
#COPY --from=0 /usr/local/include /usr/local/include

WORKDIR /home/app
COPY ./cascades ./cascades
COPY detection.linux .

RUN chown -R app ./cascades
RUN chown -R app ./detection.linux

USER app

#ENV PKG_CONFIG_PATH="/usr/local/lib64/pkgconfig"
#ENV LD_LIBRARY_PATH="/usr/local/lib64"
#ENV CGO_CPPFLAGS="-I/usr/local/include"
#ENV CGO_CXXFLAGS="--std=c++1z"
#ENV CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_face -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_xfeatures2d -lopencv_plot -lopencv_tracking"

ENTRYPOINT "./detection.linux"
