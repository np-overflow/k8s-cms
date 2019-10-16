#
# k8s-cms
# Contest Web Server Testing 
#

FROM python:3.7.4

# install module requirements
RUN mkdir /project
WORKDIR /project
COPY requirements.txt .
RUN pip install -r requirements.txt

## copy source code
COPY submit.py .

# run submit exposing options as environment variables
ENV CONTEST=test
ENV CONTEST_USERNAME=test
ENV CONTEST_PASSWORD=test
ENV SIMULATE_USERS=1
ENV HIT_AVERAGE=60
ENV HIT_STDDEV=30
ENV SELENIUM_HOST=selenium
ENV SELENIUM_PORT=4444
CMD python3 submit.py \
    --verbose \
    --contest $CONTEST \
    --username $CONTEST_USERNAME \
    --password $CONTEST_PASSWORD \
    --n-users $SIMULATE_USERS \
    --hit-avg $HIT_AVERAGE \
    --hit-stddev $HIT_STDDEV \
    --selenium-port $SELENIUM_PORT \ 
    --timeout $TIMEOUT \
    $SELENIUM_HOST $TARGET_URL
