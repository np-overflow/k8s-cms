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

# run submit
CMD python3 submit.py
